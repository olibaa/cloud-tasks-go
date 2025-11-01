package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/task", handleTask)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port
	log.Printf("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func handleTask(c *gin.Context) {
	start := time.Now()

	// リクエストボディをログ出力
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	} else {
		log.Printf("Body: %s", string(body))
	}

	// デフォルトのsleep時間を環境変数から取得
	defaultSleep := os.Getenv("TASK_SLEEP_DURATION")
	if defaultSleep == "" {
		defaultSleep = "300" // 5分
	}
	defaultDuration, err := strconv.Atoi(defaultSleep)
	if err != nil {
		log.Printf("Error converting default sleep duration: %v", err)
	}
	sleepDuration := time.Duration(defaultDuration) * time.Second

	// クエリパラメータでsleep時間を上書き可能
	if sleepParam := c.Query("sleep"); sleepParam != "" {
		if duration, err := strconv.Atoi(sleepParam); err == nil {
			sleepDuration = time.Duration(duration) * time.Second
		}
	}

	taskID := c.GetHeader("X-CloudTasks-TaskName")
	if taskID == "" {
		taskID = "unknown"
	}

	log.Printf("[%s] Task started, will sleep for %v", taskID, sleepDuration)

	time.Sleep(sleepDuration)

	elapsed := time.Since(start)
	log.Printf("[%s] Task completed in %v", taskID, elapsed)

	c.JSON(http.StatusOK, gin.H{
		"status":   "completed",
		"task_id":  taskID,
		"duration": elapsed.String(),
		"message":  fmt.Sprintf("Task processed with %v sleep", sleepDuration),
	})
}
