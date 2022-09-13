package main

import (
	_ "embed"
	"log"

	"gopkg.in/yaml.v3"
)

//go:embed resource.yaml
var resourcesData []byte
var Resource ResourceType

type ResourceType struct {
	Short string
	Long  string
}

func init() {
	err := yaml.Unmarshal(resourcesData, &Resource)
	if err != nil {
		// should not occur...
		log.Fatalf("error: %v", err)
	}
}
