package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/esgj/gochat/engine"
	"github.com/esgj/gochat/model"
)

func main() {
	intentsData, err := os.ReadFile("./chat/intents.json")

	if err != nil {
		log.Fatal(err)
	}

	intentClasses, err := os.ReadFile("./chat/classify_intents.json")

	if err != nil {
		log.Fatal(err)
	}

	intents := make([]model.Intent, 0)
	classes := make([]model.IntentClass, 0)

	if err := json.Unmarshal(intentsData, &intents); err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(intentClasses, &classes); err != nil {
		log.Fatal(err)
	}

	engine.Run(intents, classes)
}
