package main

import (
	"log"
	"os"

	"github.com/ophum/redeploy-coop/pkg"
	"gopkg.in/yaml.v2"
)

var (
	config *pkg.Config
)

func init() {
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	config = &pkg.Config{}
	if err := yaml.NewDecoder(f).Decode(config); err != nil {
		log.Fatal(err)
	}
}
func main() {
	agent, err := pkg.NewCoopAgent(config)
	if err != nil {
		log.Fatal(err)
	}

	agent.Run()
}
