package main

import (
	"os"
	"log"
	"path/filepath"
	"github.com/jchesterpivotal/knative-service-resource/pkg/in"
	"encoding/json"
	"github.com/jchesterpivotal/knative-service-resource/pkg"
	"github.com/jchesterpivotal/knative-service-resource/pkg/config"
)

func main() {
	var input *config.InRequest
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

	inner := in.NewInner(client, &input.Source, &input.Version)
	output, svc, _, err := inner.In()
	if err != nil {
		log.Printf("failed to get information from Knative: %s", err)
		os.Exit(1)
		return
	}

	inDir := os.Args[1]

	svcJson, err := os.Create(filepath.Join(inDir, "service.json"))
	if err != nil {
		log.Printf("failed to create service.json: %s\n", err)
		os.Exit(1)
		return
	}
	defer svcJson.Close()

	json.NewEncoder(svcJson).Encode(svc)
	json.NewEncoder(os.Stdout).Encode(output)
}
