package main

import (
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/ophum/redeploy-coop/pkg"
	"gopkg.in/yaml.v2"
)

var (
	config *pkg.Config
)

func init() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://81bcc4fb3bc74e629a6229b0846de171@o536371.ingest.sentry.io/5654899",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	f, err := os.Open("config.yaml")
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}
	defer f.Close()

	config = &pkg.Config{}
	if err := yaml.NewDecoder(f).Decode(config); err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}
}
func main() {
	defer sentry.Flush(2 * time.Second)
	agent, err := pkg.NewCoopAgent(config)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	agent.Run()
}
