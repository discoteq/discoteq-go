// Generate a Service Map from a chef server
//
//discoteq-chef generates a service map from chef data

package main

import (
	"log"
	"os"

	"github.com/hashicorp/logutils"

	"github.com/discoteq/discoteq-go/config"
	"github.com/discoteq/discoteq-go/service"
	"github.com/discoteq/discoteq-go/service/parser"
)

func main() {
	log.SetFlags(0) // disable time in output
	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{
			"DEBUG",
			"INFO",
			"NOTICE",
			"WARN",
			"ERROR",
			"CRITICAL",
			"ALERT",
			"PANIC",
		},
		MinLevel: "WARN",
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
	config.Parse()

	discoveredServices := make(discoteq.ServiceMap)
	for name, svc := range config.Services {
		service, err := parser.Parse(name, svc)
		if err != nil {
			log.Fatalf("Could not handle service: %s. Error: %v", name, err)
		}

		hostRecordList := service.HostRecordList()
		discoveredServices[name] = hostRecordList
	}

	json, err := discoveredServices.Marshal()
	if err != nil {
		log.Fatal("Could not generate JSON from Service Map")
	}
	os.Stdout.Write(json)
}
