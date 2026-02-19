package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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

var (
	scriptDir string
	Version   = "dev"
)

func main() {
	port := "8899"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	fmt.Printf("CronPanel %s\n", Version)
	scriptDir = filepath.Join(os.TempDir(), "crontab-manager-scripts")
	os.MkdirAll(scriptDir, 0755)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/jobs", handleJobs)
	mux.HandleFunc("/api/jobs/add", handleAddJob)
	mux.HandleFunc("/api/jobs/edit", handleEditJob)
	mux.HandleFunc("/api/jobs/delete", handleDeleteJob)
	mux.HandleFunc("/api/jobs/toggle", handleToggleJob)

	fmt.Printf("CronPanel running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
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

		// Skip manager header
		if strings.HasPrefix(line, "#!cm:") {
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Handle disabled jobs — MUST check before generic # skip
		if strings.HasPrefix(line, "#CM_DISABLED:") {
			line = strings.TrimPrefix(line, "#CM_DISABLED:")
			enabled = false
		} else if strings.HasPrefix(line, "#") {
			// Regular comment, skip
			continue
		}

		// Extract inline comment
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
			ID:       strconv.Itoa(id),
			Schedule: schedule,
			Command:  command,
			Comment:  comment,
			Enabled:  enabled,
			Raw:      raw,
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
	add := AddJobRequest{
		Mode: req.Mode, Days: req.Days, Weekday: req.Weekday,
		MonthDay: req.MonthDay, Month: req.Month,
		Hour: req.Hour, Minute: req.Minute, CustomCron: req.CustomCron,
	}
	return buildSchedule(add)
}

func handleJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jobs, err := getCrontab()
	if err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Success: true, Data: jobs})
}

func handleAddJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Method not allowed"})
		return
	}

	var req AddJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}

	command := req.Command
	if req.ScriptContent != "" {
		fname := fmt.Sprintf("script_%d.sh", time.Now().UnixNano())
		fpath := filepath.Join(scriptDir, fname)
		if err := os.WriteFile(fpath, []byte(req.ScriptContent), 0755); err != nil {
			json.NewEncoder(w).Encode(Response{Success: false, Message: "Failed to write script: " + err.Error()})
			return
		}
		command = "/bin/bash " + fpath
	} else if req.ScriptPath != "" {
		command = "/bin/bash " + req.ScriptPath
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

	newJob := CronJob{
		ID:       strconv.Itoa(len(jobs)),
		Schedule: schedule,
		Command:  command,
		Comment:  req.Comment,
		Enabled:  true,
	}
	jobs = append(jobs, newJob)

	if err := writeCrontab(jobs); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Failed to write crontab: " + err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Success: true, Message: "Job added successfully", Data: newJob})
}

func handleEditJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Method not allowed"})
		return
	}

	var req EditJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}

	command := req.Command
	if req.ScriptContent != "" {
		fname := fmt.Sprintf("script_%d.sh", time.Now().UnixNano())
		fpath := filepath.Join(scriptDir, fname)
		if err := os.WriteFile(fpath, []byte(req.ScriptContent), 0755); err != nil {
			json.NewEncoder(w).Encode(Response{Success: false, Message: "Failed to write script: " + err.Error()})
			return
		}
		command = "/bin/bash " + fpath
	} else if req.ScriptPath != "" {
		command = "/bin/bash " + req.ScriptPath
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
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(Response{Success: false, Message: "Method not allowed"})
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
		if j.ID != req.ID {
			newJobs = append(newJobs, j)
		}
	}

	if err := writeCrontab(newJobs); err != nil {
		json.NewEncoder(w).Encode(Response{Success: false, Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Success: true, Message: "Job deleted"})
}

func handleToggleJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	json.NewEncoder(w).Encode(Response{Success: true, Message: "Toggled"})
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}
