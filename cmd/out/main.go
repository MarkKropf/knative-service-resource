package main

import (
	"os"
	"log"
	"github.com/jchesterpivotal/knative-service-resource/pkg/config"
	"encoding/json"
	"github.com/jchesterpivotal/knative-service-resource/pkg"
	"github.com/jchesterpivotal/knative-service-resource/pkg/out"
)

func main() {
	var input *config.InRequest
	err := json.NewDecoder(os.Stdin).Decode(&input)
	if err != nil {
		log.Fatalf("failed to parse input JSON: %s", err)
	}

	client, err := clients.NewClients(&input.Source, "check")
	if err != nil {
		log.Fatalf("failed to create clients: %s", err)
	}

	outer := out.NewOuter(client, &input.Source, &input.Params)

}
