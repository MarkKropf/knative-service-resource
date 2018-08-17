package main

import (
	"os"
		"encoding/json"
	"github.com/jchesterpivotal/knative-service-resource/pkg/check"
	"github.com/jchesterpivotal/knative-service-resource/pkg"
	"fmt"
	"log"
	"github.com/jchesterpivotal/knative-service-resource/pkg/config"
)

func main() {
	var input *config.CheckRequest
	err := json.NewDecoder(os.Stdin).Decode(&input)
	if err != nil {
		log.Printf("failed to parse input JSON: %s", err)
		os.Exit(1)
		return
	}

	client, err := clients.NewClients(&input.Source, "check")
	if err != nil {
		log.Printf("failed to create clients: %s", err)
		os.Exit(1)
		return
	}

	checker := check.NewChecker(client, &input.Source, &input.Version)

	checkResult, err := checker.Check()
	if err != nil {
		fmt.Printf("failed to perform check operation: %s", err)
		os.Exit(1)
		return
	}

	json.NewEncoder(os.Stdout).Encode(checkResult)
}
