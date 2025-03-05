package main

import (
	"log"

	"github.com/Yorshik/final_task_sprint_1/internal/server/agent"
	"github.com/Yorshik/final_task_sprint_1/internal/server/orchestrator"
)

func main() {
	go func() {
		log.Println("Starting orchestrator...")
		orchestrator.StartServer()
	}()

	log.Println("Starting agent...")
	agent.StartAgent()
}
