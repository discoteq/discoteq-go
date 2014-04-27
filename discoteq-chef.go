package main

import (
	"os"
	"fmt"
	"flag"
	"log"
//	"github.com/marpaia/chef-golang"
	"./chef"
//	"chef"
)

var configPath string

func init() {
	const (
		defaultConfigPath = "/etc/discoteq-chef.json"
		usage         = "config file"
	)
	flag.StringVar(&configPath, "config", defaultConfigPath, usage)
	flag.StringVar(&configPath, "c", defaultConfigPath, usage+" (shorthand)")
}


func main() {
	flag.Parse()

	// load conf
	config, err := chef.UnmarshalFile(configPath)
	if err != nil {
		log.Fatalf("could not parse config file %s: 5s", configPath, err)
	}

	// discoveredServices = new Services()
	for name, svc := range config.Services {
		fmt.Printf("%s: %s\n", name, svc)
		// request service data
		// discoveredService = new Service(svc.Name())
		// nodes, err := chef.Search("nodes", svc.Query())
		// for node, i := range nodes {
			// attrs := svc.ExtractAttrs(node)
			// discoveredService.Append(attrs)
		// }
		// discoveredServices.Append(discoveredService)
	}
	
	// render services data
	json, err := config.Marshal()
	if err != nil {
		log.Fatal("could not generate json")
	}
	output :=  fmt.Sprintf("%s: %s\n", configPath, json)
	os.Stdout.Write( []byte(output))
}

