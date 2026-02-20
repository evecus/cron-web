package main

import (
	"bufio"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ── Data types ────────────────────────────────────────────────────────────────

type CronJob struct {
	ID       string `json:"id"`
	Schedule string `json:"schedule"`
	Command  string `json:"command"`
	Comment  string `json:"comment"`
	Enabled  bool   `json:"enabled"`
	Raw      string `json:"raw"`
	// Log-related fields (populated from meta, not from crontab line)
	SaveLog bool   `json:"saveLog"`
	LogDir  string `json:"logDir"`
	RealCmd string `json:"realCmd"` // original command before wrapper
}

type AddJobRequest struct {
	Mode          string `json:"mode"`
	Days          string `json:"days"`
	Weekday       string `json:"weekday"`
	MonthDay      string `json:"monthDay"`
	Month         string `json:"month"`
	Hour          string `json:"hour"`
	Minute        string `json:"minute"`
	Command       string `json:"command"`
	ScriptPath    string `json:"scriptPath"`
	ScriptContent string `json:"scriptContent"`
	Comment       string `json:"comment"`
	CustomCron    string `json:"customCron"`
	SaveLog       bool   `json:"saveLog"`
}

type EditJobRequest struct {
	ID            string `json:"id"`
	Mode          string `json:"mode"`
	Days          string `json:"days"`
	Weekday       string `json:"weekday"`
	MonthDay      string `json:"monthDay"`
	Month         string `json:"month"`
	Hour          string `json:"hour"`
	Minute        string `json:"minute"`
	Command       string `json:"command"`
	ScriptPath    string `json:"scriptPath"`
	ScriptContent string `json:"scriptContent"`
	Comment       string `json:"comment"`
	CustomCron    string `json:"customCron"`
	SaveLog       bool   `json:"saveLog"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type LogEntry struct {
	Filename  string `json:"filename"`
	CreatedAt string `json:"createdAt"`
	Size      int64  `json:"size"`
}

// LogMeta stores logging config for each cron job, keyed by a stable job key.
// Persisted to disk as cronpanel-logs/meta.json.
type JobLogMeta struct {
	SaveLog     bool   `json:"saveLog"`
	LogDir      string `json:"logDir"`
	WrapperPath string `json:"wrapperPath"` // path of the cplog_*.sh script
	RealCmd     string `json:"realCmd"`     // original command before wrapping
}

type MetaStore struct {
	mu   sync.Mutex
	path string
	data map[string]*JobLogMeta // key: jobKey()
}

func newMetaStore(path string) *MetaStore {
	ms := &MetaStore{path: path, data: make(map[string]*JobLogMeta)}
	ms.load()
	return ms
}

func (ms *MetaStore) load() {
	b, err := os.ReadFile(ms.path)
	if err != nil {
		return
	}
	json.Unmarshal(b, &ms.data)
}

func (ms *MetaStore) save() {
	b, _ := json.MarshalIndent(ms.data, "", "  ")
	os.WriteFile(ms.path, b, 0644)
}

func (ms *MetaStore) Get(key string) *JobLogMeta {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return ms.data[key]
}

func (ms *MetaStore) Set(key string, m *JobLogMeta) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.data[key] = m
	ms.save()
}

func (ms *MetaStore) Delete(key string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	delete(ms.data, key)
	ms.save()
}

// ── Session store ─────────────────────────────────────────────────────────────

type SessionStore struct {
	mu       sync.Mutex
	sessions map[string]time.Time
}

func newSessionStore() *SessionStore {
	return &SessionStore{sessions: make(map[string]time.Time)}
}

func (s *SessionStore) Create() string {
	b := make([]byte, 16)
	rand.Read(b)
	token := hex.EncodeToString(b)
	s.mu.Lock()
	s.sessions[token] = time.Now().Add(24 * time.Hour)
	s.mu.Unlock()
	return token
}

func (s *SessionStore) Valid(token string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	exp, ok := s.sessions[token]
	if !ok {
		return false
	}
	if time.Now().After(exp) {
		delete(s.sessions, token)
		return false
	}
	return true
}

func (s *SessionStore) Delete(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

// ── Globals ───────────────────────────────────────────────────────────────────

var (
	scriptDir   string
	logBaseDir  string
	Version     = "dev"
	authUser    string
	authPass    string
	authEnabled bool
	sessions    = newSessionStore()
	metaStore   *MetaStore
)

// jobKey produces a stable identifier for a cron entry from its schedule+command.
func jobKey(schedule, command string) string {
	h := md5.Sum([]byte(schedule + "|" + command))
	return hex.EncodeToString(h[:])
}

// logDirName builds a human-friendly log directory name.
func logDirName(comment, realCmd string) string {
	h := md5.Sum([]byte(realCmd))
	slug := hex.EncodeToString(h[:])[:10]
	if comment != "" {
		safe := strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
				(r >= '0' && r <= '9') || r == '-' || r == '_' {
				return r
			}
			return '_'
		}, comment)
		if len(safe) > 24 {
			safe = safe[:24]
		}
		return filepath.Join(logBaseDir, safe+"_"+slug)
	}
	return filepath.Join(logBaseDir, slug)
}

// createWrapperScript writes a bash script that runs realCmd and appends
// stdout+stderr to a timestamped file inside logDir.
// Returns the absolute path of the script.
func createWrapperScript(realCmd, logDir string) (string, error) {
	// The script uses $() expansion for the filename — no quoting issues because
	// all paths are in the script body, not on the crontab line.
	content := fmt.Sprintf(`#!/bin/bash
LOGDIR=%s
mkdir -p "$LOGDIR"
LOGFILE="$LOGDIR/$(date +%%Y%%m%%d_%%H%%M%%S).log"
%s >> "$LOGFILE" 2>&1
`, logDir, realCmd)

	fname := fmt.Sprintf("cplog_%d.sh", time.Now().UnixNano())
	path := filepath.Join(scriptDir, fname)
	if err := os.WriteFile(path, []byte(content), 0755); err != nil {
		return "", err
	}
	return path, nil
}

// ── main ──────────────────────────────────────────────────────────────────────

func main() {
	portFlag := flag.String("port", "", "Port to listen on (default 8899)")
	authFlag := flag.String("auth", "", "Authentication in user:password format")
	flag.Parse()

	port := "8899"
	if *portFlag != "" {
		port = *portFlag
	} else if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	if *authFlag != "" {
		parts := strings.SplitN(*authFlag, ":", 2)
		if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
			authUser = parts[0]
			authPass = parts[1]
			authEnabled = true
			fmt.Printf("Auth enabled for user: %s\n", authUser)
		} else {
			fmt.Println("Warning: invalid --auth format, expected user:password. Auth disabled.")
		}
	}

	fmt.Printf("CronPanel %s\n", Version)
	exePath, err := os.Executable()
	if err != nil {
		exePath = "."
	}
	scriptDir = filepath.Join(filepath.Dir(exePath), "cronpanel-scripts")
	os.MkdirAll(scriptDir, 0755)
	logBaseDir = filepath.Join(filepath.Dir(exePath), "cronpanel-logs")
	os.MkdirAll(logBaseDir, 0755)
	metaStore = newMetaStore(filepath.Join(logBaseDir, "meta.json"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/auth/login", handleLogin)
	mux.HandleFunc("/api/auth/logout", handleLogout)
	mux.HandleFunc("/api/auth/check", handleAuthCheck)
	mux.HandleFunc("/api/jobs/read-script", authMiddleware(handleReadScript))
	mux.HandleFunc("/api/jobs", authMiddleware(handleJobs))
	mux.HandleFunc("/api/jobs/add", authMiddleware(handleAddJob))
	mux.HandleFunc("/api/jobs/edit", authMiddleware(handleEditJob))
	mux.HandleFunc("/api/jobs/delete", authMiddleware(handleDeleteJob))
	mux.HandleFunc("/api/jobs/toggle", authMiddleware(handleToggleJob))
	mux.HandleFunc("/api/jobs/logs", authMiddleware(handleListLogs))
	mux.HandleFunc("/api/jobs/logs/content", authMiddleware(handleLogContent))
	mux.HandleFunc("/api/jobs/logs/delete", authMiddleware(handleDeleteLog))

	fmt.Printf("CronPanel running at http://0.0.0.0:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// ── Auth middleware ───────────────────────────────────────────────────────────

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if !authEnabled {
			next(w, r)
			return
		}
		token := getSessionToken(r)
		if !sessions.Valid(token) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Response{Success: false, Message: "unauthorized"})
			return
		}
		next(w, r)
	}
}

func getSessionToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	c, err := r.Cookie("cp_session")
	if err == nil {
		return c.Value
	}
	return ""
}

// ── Auth handlers ─────────────────────────────────────────────────────────────

func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !authEnabled {
		json.NewEncoder(w).Encode(Response{Success: true, Message: "no_auth"})
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "invalid request"})
		return
	}
	if req.Username == authUser && req.Password == authPass {
		token := sessions.Create()
		http.SetCookie(w, &http.Cookie{
			Name: "cp_session", Value: token, Path: "/",
			MaxAge: 86400, HttpOnly: true, SameSite: http.SameSiteStrictMode,
		})
		json.NewEncoder(w).Encode(Response{Success: true, Data: token})
	} else {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "invalid credentials"})
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := getSessionToken(r)
	if token != "" {
		sessions.Delete(token)
	}
	http.SetCookie(w, &http.Cookie{Name: "cp_session", Value: "", Path: "/", MaxAge: -1})
	json.NewEncoder(w).Encode(Response{Success: true})
}

func handleAuthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type authInfo struct {
		Required bool `json:"required"`
		LoggedIn bool `json:"loggedIn"`
	}
	info := authInfo{Required: authEnabled}
	if authEnabled {
		info.LoggedIn = sessions.Valid(getSessionToken(r))
	} else {
		info.LoggedIn = true
	}
	json.NewEncoder(w).Encode(Response{Success: true, Data: info})
}

// ── Crontab helpers ───────────────────────────────────────────────────────────

func getCrontab() ([]CronJob, error) {
	cmd := exec.Command("crontab", "-l")
	out, err := cmd.Output()
	if err != nil {
		return []CronJob{}, nil
	}

	var jobs []CronJob
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	id := 0
	for scanner.Scan() {
		line := scanner.Text()
		raw := line
		comment := ""
		enabled := true

		if strings.HasPrefix(line, "#!cm:") || strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, "#CM_DISABLED:") {
			line = strings.TrimPrefix(line, "#CM_DISABLED:")
			enabled = false
		} else if strings.HasPrefix(line, "#") {
			continue
		}

		// Extract optional comment (## suffix written by us)
		if idx := strings.Index(line, " ##"); idx != -1 {
			comment = strings.TrimSpace(line[idx+3:])
			line = strings.TrimSpace(line[:idx])
		}

		parts := strings.Fields(line)
		if len(parts) < 6 {
			continue
		}
		schedule := strings.Join(parts[:5], " ")
		command := strings.Join(parts[5:], " ")

		// Look up log meta for this job
		key := jobKey(schedule, command)
		job := CronJob{
			ID: strconv.Itoa(id), Schedule: schedule,
			Command: command, Comment: comment, Enabled: enabled, Raw: raw,
		}
		if m := metaStore.Get(key); m != nil && m.SaveLog {
			job.SaveLog = true
			job.LogDir = m.LogDir
			job.RealCmd = m.RealCmd
		}

		jobs = append(jobs, job)
		id++
	}
	return jobs, nil
}

func writeCrontab(jobs []CronJob) error {
	var lines []string
	lines = append(lines, "#!cm:managed by crontab-manager")
	for _, job := range jobs {
		line := job.Schedule + " " + job.Command
		if job.Comment != "" {
			line += " ##" + job.Comment
		}
		if !job.Enabled {
			line = "#CM_DISABLED:" + line
		}
		lines = append(lines, line)
	}
	content := strings.Join(lines, "\n") + "\n"
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(content)
	return cmd.Run()
}

func buildSchedule(req AddJobRequest) string {
	min := req.Minute
	if min == "" {
		min = "0"
	}
	hour := req.Hour
	if hour == "" {
		hour = "0"
	}
	switch req.Mode {
	case "interval":
		n, _ := strconv.Atoi(req.Days)
		if n <= 0 {
			n = 1
		}
		return fmt.Sprintf("%s %s */%d * *", min, hour, n)
	case "weekly":
		wd := req.Weekday
		if wd == "" {
			wd = "0"
		}
		return fmt.Sprintf("%s %s * * %s", min, hour, wd)
	case "monthly":
		md := req.MonthDay
		if md == "" {
			md = "1"
		}
		month := req.Month
		if month == "" {
			month = "*"
		}
		return fmt.Sprintf("%s %s %s %s *", min, hour, md, month)
	case "daily":
		return fmt.Sprintf("%s %s * * *", min, hour)
	case "custom":
		return req.CustomCron
	default:
		return fmt.Sprintf("%s %s * * *", min, hour)
	}
}

func buildScheduleFromEdit(req EditJobRequest) string {
	return buildSchedule(AddJobRequest{
		Mode: req.Mode, Days: req.Days, Weekday: req.Weekday,
		MonthDay: req.MonthDay, Month: req.Month,
		Hour: req.Hour, Minute: req.Minute, CustomCron: req.CustomCron,
	})
}

// resolveCommand turns form inputs into the actual crontab command string.
// For script-content jobs it saves the script to disk.
func resolveCommand(cmd, scriptPath, scriptContent string) (string, error) {
	if scriptContent != "" {
		fname := fmt.Sprintf("script_%d.sh", time.Now().UnixNano())
		fpath := filepath.Join(scriptDir, fname)
		if err := os.WriteFile(fpath, []byte(scriptContent), 0755); err != nil {
			return "", err
		}
		return "/bin/bash " + fpath, nil
	}
	if scriptPath != "" {
		return "/bin/bash " + scriptPath, nil
	}
	return cmd, nil
}

// ── Job handlers ──────────────────────────────────────────────────────────────

func handleReadScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "invalid request"})
		return
	}
	fpath := strings.TrimSpace(strings.TrimPrefix(req.Path, "/bin/bash "))
	absPath, err := filepath.Abs(fpath)
	if err != nil || !strings.HasPrefix(absPath, scriptDir) {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "access denied"})
		return
	}
	content, err := os.ReadFile(absPath)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "file not found"})
		return
	}
	json.NewEncoder(w).Encode(Response{Success: true, Data: string(content)})
}

func handleJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := getCrontab()
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Success: true, Data: jobs})
}

func handleAddJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "method not allowed"})
		return
	}
	var req AddJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}

	realCmd, err := resolveCommand(req.Command, req.ScriptPath, req.ScriptContent)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "script error: " + err.Error()})
		return
	}
	if realCmd == "" {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "command is required"})
		return
	}

	schedule := buildSchedule(req)

	// The command written to crontab: wrapper script if logging, else realCmd directly
	cronCmd := realCmd
	var logMeta *JobLogMeta

	if req.SaveLog {
		logDir := logDirName(req.Comment, realCmd)
		os.MkdirAll(logDir, 0755)
		wrapperPath, werr := createWrapperScript(realCmd, logDir)
		if werr != nil {
			json.NewEncoder(w).Encode(Response{Success: false, Message: "wrapper error: " + werr.Error()})
			return
		}
		cronCmd = "/bin/bash " + wrapperPath
		logMeta = &JobLogMeta{
			SaveLog:     true,
			LogDir:      logDir,
			WrapperPath: wrapperPath,
			RealCmd:     realCmd,
		}
	}

	jobs, err := getCrontab()
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}
	newJob := CronJob{
		ID: strconv.Itoa(len(jobs)), Schedule: schedule, Command: cronCmd,
		Comment: req.Comment, Enabled: true,
		SaveLog: req.SaveLog, RealCmd: realCmd,
	}
	if logMeta != nil {
		newJob.LogDir = logMeta.LogDir
	}
	jobs = append(jobs, newJob)

	if err := writeCrontab(jobs); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "crontab error: " + err.Error()})
		return
	}

	// Persist meta AFTER writing crontab (key uses final cronCmd)
	key := jobKey(schedule, cronCmd)
	if logMeta != nil {
		metaStore.Set(key, logMeta)
	}

	json.NewEncoder(w).Encode(Response{Success: true, Message: "job added", Data: newJob})
}

func handleEditJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "method not allowed"})
		return
	}
	var req EditJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}

	realCmd, err := resolveCommand(req.Command, req.ScriptPath, req.ScriptContent)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "script error: " + err.Error()})
		return
	}
	if realCmd == "" {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "command is required"})
		return
	}

	schedule := buildScheduleFromEdit(req)

	jobs, err := getCrontab()
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}

	// Find the old job to get its old key (so we can clean up old meta)
	var oldJob *CronJob
	for i := range jobs {
		if jobs[i].ID == req.ID {
			oldJob = &jobs[i]
			break
		}
	}
	if oldJob == nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "job not found"})
		return
	}
	oldKey := jobKey(oldJob.Schedule, oldJob.Command)
	oldMeta := metaStore.Get(oldKey)

	// If logging was previously on, remove old wrapper script
	if oldMeta != nil && oldMeta.WrapperPath != "" {
		os.Remove(oldMeta.WrapperPath)
	}

	// Build new command and meta
	cronCmd := realCmd
	var newMeta *JobLogMeta

	if req.SaveLog {
		logDir := logDirName(req.Comment, realCmd)
		os.MkdirAll(logDir, 0755)
		wrapperPath, werr := createWrapperScript(realCmd, logDir)
		if werr != nil {
			json.NewEncoder(w).Encode(Response{Success: false, Message: "wrapper error: " + werr.Error()})
			return
		}
		cronCmd = "/bin/bash " + wrapperPath
		newMeta = &JobLogMeta{
			SaveLog:     true,
			LogDir:      logDir,
			WrapperPath: wrapperPath,
			RealCmd:     realCmd,
		}
	}

	// Update the job in the list
	for i := range jobs {
		if jobs[i].ID == req.ID {
			jobs[i].Schedule = schedule
			jobs[i].Command = cronCmd
			jobs[i].Comment = req.Comment
			break
		}
	}

	if err := writeCrontab(jobs); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "crontab error: " + err.Error()})
		return
	}

	// Update meta store: remove old key, add new key
	metaStore.Delete(oldKey)
	newKey := jobKey(schedule, cronCmd)
	if newMeta != nil {
		metaStore.Set(newKey, newMeta)
	}

	json.NewEncoder(w).Encode(Response{Success: true, Message: "job updated"})
}

func handleDeleteJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "method not allowed"})
		return
	}
	var req struct {
		ID string `json:"id"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	jobs, err := getCrontab()
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}

	var newJobs []CronJob
	for _, j := range jobs {
		if j.ID == req.ID {
			// Clean up wrapper script if any
			key := jobKey(j.Schedule, j.Command)
			if m := metaStore.Get(key); m != nil {
				if m.WrapperPath != "" {
					os.Remove(m.WrapperPath)
				}
				metaStore.Delete(key)
			}
		} else {
			newJobs = append(newJobs, j)
		}
	}

	if err := writeCrontab(newJobs); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Success: true, Message: "job deleted"})
}

func handleToggleJob(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	jobs, err := getCrontab()
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}
	for i, j := range jobs {
		if j.ID == req.ID {
			jobs[i].Enabled = !j.Enabled
		}
	}
	if err := writeCrontab(jobs); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Success: true, Message: "toggled"})
}

// ── Log handlers ──────────────────────────────────────────────────────────────

func handleListLogs(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LogDir string `json:"logDir"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "invalid request"})
		return
	}
	absDir, err := filepath.Abs(req.LogDir)
	if err != nil || !strings.HasPrefix(absDir, logBaseDir) {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "access denied"})
		return
	}
	entries, err := os.ReadDir(absDir)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: true, Data: []LogEntry{}})
		return
	}
	var logs []LogEntry
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".log") {
			continue
		}
		info, _ := e.Info()
		ts := strings.TrimSuffix(e.Name(), ".log")
		t, perr := time.ParseInLocation("20060102_150405", ts, time.Local)
		createdAt := e.Name()
		if perr == nil {
			createdAt = t.Format("2006-01-02 15:04:05")
		}
		logs = append(logs, LogEntry{Filename: e.Name(), CreatedAt: createdAt, Size: info.Size()})
	}
	sort.Slice(logs, func(i, j int) bool { return logs[i].Filename > logs[j].Filename })
	json.NewEncoder(w).Encode(Response{Success: true, Data: logs})
}

func handleLogContent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LogDir   string `json:"logDir"`
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "invalid request"})
		return
	}
	absDir, err := filepath.Abs(req.LogDir)
	if err != nil || !strings.HasPrefix(absDir, logBaseDir) {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "access denied"})
		return
	}
	if strings.ContainsAny(req.Filename, "/\\") {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "invalid filename"})
		return
	}
	content, err := os.ReadFile(filepath.Join(absDir, req.Filename))
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "file not found"})
		return
	}
	const maxSize = 512 * 1024
	if len(content) > maxSize {
		content = append([]byte("...(showing last 512KB)...\n"), content[len(content)-maxSize:]...)
	}
	json.NewEncoder(w).Encode(Response{Success: true, Data: string(content)})
}

func handleDeleteLog(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LogDir   string `json:"logDir"`
		Filename string `json:"filename"` // empty = clear all
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "invalid request"})
		return
	}
	absDir, err := filepath.Abs(req.LogDir)
	if err != nil || !strings.HasPrefix(absDir, logBaseDir) {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "access denied"})
		return
	}
	if req.Filename == "" {
		entries, _ := os.ReadDir(absDir)
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".log") {
				os.Remove(filepath.Join(absDir, e.Name()))
			}
		}
	} else {
		if strings.ContainsAny(req.Filename, "/\\") {
			json.NewEncoder(w).Encode(Response{Success: false, Message: "invalid filename"})
			return
		}
		os.Remove(filepath.Join(absDir, req.Filename))
	}
	json.NewEncoder(w).Encode(Response{Success: true, Message: "deleted"})
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}
