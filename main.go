package main

import (
	"bufio"
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
	"strconv"
	"strings"
	"sync"
	"time"
)

type CronJob struct {
	ID       string `json:"id"`
	Schedule string `json:"schedule"`
	Command  string `json:"command"`
	Comment  string `json:"comment"`
	Enabled  bool   `json:"enabled"`
	Raw      string `json:"raw"`
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
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Session store: token -> expiry
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

var (
	scriptDir  string
	Version    = "dev"
	authUser   string
	authPass   string
	authEnabled bool
	sessions   = newSessionStore()
)

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
	scriptDir = filepath.Join(os.TempDir(), "crontab-manager-scripts")
	os.MkdirAll(scriptDir, 0755)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/auth/login", handleLogin)
	mux.HandleFunc("/api/auth/logout", handleLogout)
	mux.HandleFunc("/api/auth/check", handleAuthCheck)
	mux.HandleFunc("/api/jobs", authMiddleware(handleJobs))
	mux.HandleFunc("/api/jobs/add", authMiddleware(handleAddJob))
	mux.HandleFunc("/api/jobs/edit", authMiddleware(handleEditJob))
	mux.HandleFunc("/api/jobs/delete", authMiddleware(handleDeleteJob))
	mux.HandleFunc("/api/jobs/toggle", authMiddleware(handleToggleJob))

	fmt.Printf("CronPanel running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// authMiddleware wraps handlers that require authentication
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
	// Check Authorization header: Bearer <token>
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	// Check cookie fallback
	c, err := r.Cookie("cp_session")
	if err == nil {
		return c.Value
	}
	return ""
}

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
			Name:     "cp_session",
			Value:    token,
			Path:     "/",
			MaxAge:   86400,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
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
		token := getSessionToken(r)
		info.LoggedIn = sessions.Valid(token)
	} else {
		info.LoggedIn = true
	}
	json.NewEncoder(w).Encode(Response{Success: true, Data: info})
}

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

		if strings.HasPrefix(line, "#!cm:") {
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, "#CM_DISABLED:") {
			line = strings.TrimPrefix(line, "#CM_DISABLED:")
			enabled = false
		} else if strings.HasPrefix(line, "#") {
			continue
		}

		if idx := strings.Index(line, " #"); idx != -1 {
			comment = strings.TrimSpace(line[idx+2:])
			line = strings.TrimSpace(line[:idx])
		}

		parts := strings.Fields(line)
		if len(parts) < 6 {
			continue
		}
		schedule := strings.Join(parts[:5], " ")
		command := strings.Join(parts[5:], " ")

		jobs = append(jobs, CronJob{
			ID: strconv.Itoa(id), Schedule: schedule,
			Command: command, Comment: comment, Enabled: enabled, Raw: raw,
		})
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
			line += " #" + job.Comment
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
	if min == "" { min = "0" }
	hour := req.Hour
	if hour == "" { hour = "0" }
	switch req.Mode {
	case "interval":
		n, _ := strconv.Atoi(req.Days)
		if n <= 0 { n = 1 }
		return fmt.Sprintf("%s %s */%d * *", min, hour, n)
	case "weekly":
		wd := req.Weekday
		if wd == "" { wd = "0" }
		return fmt.Sprintf("%s %s * * %s", min, hour, wd)
	case "monthly":
		md := req.MonthDay
		if md == "" { md = "1" }
		month := req.Month
		if month == "" { month = "*" }
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

func resolveCommand(req AddJobRequest) (string, error) {
	command := req.Command
	if req.ScriptContent != "" {
		fname := fmt.Sprintf("script_%d.sh", time.Now().UnixNano())
		fpath := filepath.Join(scriptDir, fname)
		if err := os.WriteFile(fpath, []byte(req.ScriptContent), 0755); err != nil {
			return "", err
		}
		command = "/bin/bash " + fpath
	} else if req.ScriptPath != "" {
		command = "/bin/bash " + req.ScriptPath
	}
	return command, nil
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
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Method not allowed"})
		return
	}
	var req AddJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}
	command, err := resolveCommand(req)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Failed to write script: " + err.Error()})
		return
	}
	if command == "" {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Command is required"})
		return
	}
	schedule := buildSchedule(req)
	jobs, err := getCrontab()
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}
	newJob := CronJob{ID: strconv.Itoa(len(jobs)), Schedule: schedule, Command: command, Comment: req.Comment, Enabled: true}
	jobs = append(jobs, newJob)
	if err := writeCrontab(jobs); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Failed to write crontab: " + err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Success: true, Message: "Job added successfully", Data: newJob})
}

func handleEditJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Method not allowed"})
		return
	}
	var req EditJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}
	addReq := AddJobRequest{Command: req.Command, ScriptPath: req.ScriptPath, ScriptContent: req.ScriptContent}
	command, err := resolveCommand(addReq)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Failed to write script: " + err.Error()})
		return
	}
	if command == "" {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Command is required"})
		return
	}
	schedule := buildScheduleFromEdit(req)
	jobs, err := getCrontab()
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}
	found := false
	for i, j := range jobs {
		if j.ID == req.ID {
			jobs[i].Schedule = schedule
			jobs[i].Command = command
			jobs[i].Comment = req.Comment
			found = true
			break
		}
	}
	if !found {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Job not found"})
		return
	}
	if err := writeCrontab(jobs); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Failed to write crontab: " + err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Success: true, Message: "Job updated successfully"})
}

func handleDeleteJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Method not allowed"})
		return
	}
	var req struct{ ID string `json:"id"` }
	json.NewDecoder(r.Body).Decode(&req)
	jobs, err := getCrontab()
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}
	var newJobs []CronJob
	for _, j := range jobs {
		if j.ID != req.ID { newJobs = append(newJobs, j) }
	}
	if err := writeCrontab(newJobs); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Success: true, Message: "Job deleted"})
}

func handleToggleJob(w http.ResponseWriter, r *http.Request) {
	var req struct{ ID string `json:"id"` }
	json.NewDecoder(r.Body).Decode(&req)
	jobs, err := getCrontab()
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}
	for i, j := range jobs {
		if j.ID == req.ID { jobs[i].Enabled = !j.Enabled }
	}
	if err := writeCrontab(jobs); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Success: true, Message: "Toggled"})
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}
