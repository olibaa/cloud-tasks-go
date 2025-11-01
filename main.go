package helloworld

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("TaskHandler", handleTask)
}

type TaskRequest struct {
	UserID     string `json:"user_id"`
	TaskNumber int    `json:"task_number"`
}

type TaskResponse struct {
	Status   string `json:"status"`
	TaskID   string `json:"task_id"`
	Duration string `json:"duration"`
	Message  string `json:"message"`
	UserID   string `json:"user_id,omitempty"`
}

// handleTask is an HTTP Cloud Function for processing tasks
func handleTask(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// リクエスト内容をログ出力
	log.Printf("=== Task Request ===")
	log.Printf("Method: %s", r.Method)
	log.Printf("URL: %s", r.URL.String())
	log.Printf("Headers:")
	for key, values := range r.Header {
		for _, value := range values {
			log.Printf("  %s: %s", key, value)
		}
	}

	// リクエストボディをログ出力
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	log.Printf("Body: %s", string(body))

	// JSONデータのパース
	var taskReq TaskRequest
	if len(body) > 0 {
		if err := json.Unmarshal(body, &taskReq); err != nil {
			log.Printf("Error parsing JSON: %v", err)
		}
	}

	// デフォルトのsleep時間を環境変数から取得
	defaultSleep := os.Getenv("TASK_SLEEP_DURATION")
	if defaultSleep == "" {
		defaultSleep = "10" // 10秒
	}
	
	defaultDuration, err := strconv.Atoi(defaultSleep)
	if err != nil {
		log.Printf("Error converting default sleep duration: %v", err)
		defaultDuration = 10
	}
	sleepDuration := time.Duration(defaultDuration) * time.Second

	// クエリパラメータでsleep時間を上書き可能
	if sleepParam := r.URL.Query().Get("sleep"); sleepParam != "" {
		if duration, err := strconv.Atoi(sleepParam); err == nil {
			sleepDuration = time.Duration(duration) * time.Second
		}
	}

	taskID := r.Header.Get("X-CloudTasks-TaskName")
	if taskID == "" {
		taskID = "unknown"
	}

	log.Printf("[%s] Task started, will sleep for %v", taskID, sleepDuration)
	if taskReq.UserID != "" {
		log.Printf("[%s] Processing for user: %s", taskID, taskReq.UserID)
	}

	time.Sleep(sleepDuration)

	elapsed := time.Since(start)
	log.Printf("[%s] Task completed in %v", taskID, elapsed)

	response := TaskResponse{
		Status:   "completed",
		TaskID:   taskID,
		Duration: elapsed.String(),
		Message:  fmt.Sprintf("Task processed with %v sleep", sleepDuration),
		UserID:   taskReq.UserID,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
