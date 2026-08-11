package vtr

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type JsonIntent struct {
	Name              string   `json:"name"`
	Keyphrases        []string `json:"keyphrases"`
	RequireExactMatch bool     `json:"requiresexact"`
}

var IntentList []JsonIntent

func loadIntents() {
	file, err := os.ReadFile("/anki/data/cloudswitch/en-US/en-US.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &IntentList)
	if err != nil {
		fmt.Println(err)
	}
}
