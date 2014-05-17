// Generate a Service Map from a chef server
//
//discoteq-chef generates a service map from chef data

package main

import (
	"log"
	"os"

	"github.com/josephholsten/discoteq/chef"
	"github.com/josephholsten/discoteq/common"
	"github.com/josephholsten/discoteq/chef/config"
)

func main() {
	log.SetFlags(0) // disable time in output
	config.Parse()

	discoveredServices := make(discoteq.ServiceMap)
	for name, svc := range config.Services {
		service := chef.ServiceFromRaw(name, svc)
		discoveredServices[service.Name] = service.HostRecordList()
	}

	json, err := discoveredServices.Marshal()
	if err != nil {
		log.Fatal("Could not generate JSON from Service Map")
	}
	os.Stdout.Write(json)
}
