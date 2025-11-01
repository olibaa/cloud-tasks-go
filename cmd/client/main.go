package main

import (
	"context"
	"fmt"
	"log"
	"os"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

func createTasks(numTasks int) error {
	ctx := context.Background()
	
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// 環境変数からGCPプロジェクト設定を取得
	projectID := os.Getenv("GCP_PROJECT_ID")
	location := os.Getenv("GCP_LOCATION")
	queueName := os.Getenv("GCP_QUEUE_NAME")
	endpointURL := os.Getenv("TASK_ENDPOINT_URL")
	sleepDuration := os.Getenv("TASK_SLEEP_DURATION")
	
	if projectID == "" || location == "" || queueName == "" || endpointURL == "" {
		return fmt.Errorf("required environment variables not set: GCP_PROJECT_ID, GCP_LOCATION, GCP_QUEUE_NAME, TASK_ENDPOINT_URL")
	}
	
	if sleepDuration == "" {
		sleepDuration = "10"
	}
	
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, location, queueName)
	taskURL := fmt.Sprintf("%s/task?sleep=%s", endpointURL, sleepDuration)
	
	for i := 0; i < numTasks; i++ {
		task := &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					HttpMethod: cloudtaskspb.HttpMethod_POST,
					Url:        taskURL,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: []byte(fmt.Sprintf(`{"user_id": "user_%d", "task_number": %d}`, i+1, i+1)),
				},
			},
		}

		req := &cloudtaskspb.CreateTaskRequest{
			Parent: queuePath,
			Task:   task,
		}

		createdTask, err := client.CreateTask(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to create task %d: %v", i, err)
		}
		
		log.Printf("Created task %d: %s", i, createdTask.Name)
	}
	
	return nil
}

func main() {
	numTasks := 10
	if len(os.Args) > 1 {
		if n, err := fmt.Sscanf(os.Args[1], "%d", &numTasks); err != nil || n != 1 {
			log.Fatal("Usage: go run client.go [number_of_tasks]")
		}
	}
	
	log.Printf("Creating %d tasks...", numTasks)
	if err := createTasks(numTasks); err != nil {
		log.Fatal(err)
	}
	log.Println("All tasks created successfully")
}