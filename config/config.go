// Parse configuration from command-line flags and json file
package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	chefClient "github.com/go-chef/chef"
)

var (
	ConfigPath             = flag.String("c", "/etc/discoteq.json", "config file path")
	defaultChefServerUrl   = "http://localhost:4545"
	ChefServerUrl          = flag.String("s", "", "chef server URL, default: http://localhost:4545")
	ChefClientUsername     = flag.String("u", "", "API client username, default: node.fqdn")
	defaultChefKeyPath     = "/etc/chef/client.pem"
	ChefClientKeyPath      = flag.String("k", "", "API client key file path, default: /etc/chef/client.pem")
	defaultChefEnvironment = "_default"
	ChefEnvironment        = flag.String("E", "", "environment query scope, default: node.chef_environment || _default")
	Services               map[string]map[string]interface{}
)

type Config struct {
	ChefServerUrl      string                            `json:chef-server-url`
	ChefClientUsername string                            `json:chef-user`
	ChefClientKeyPath  string                            `json:chef-key`
	ChefEnvironment    string                            `json:chef-environment`
	Services           map[string]map[string]interface{} `json:services`
	//	Services map[string]Service  `json:"services"`
}

func Unmarshal(data []byte) (*Config, error) {
	var config *Config
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func UnmarshalFile(filename string) (*Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config, err := Unmarshal(file)
	return config, err
}

func (config *Config) Marshal() ([]byte, error) {
	json, err := json.MarshalIndent(config, "", "  ")
	return json, err
}

func Parse() {
	log.Print("[DEBUG] Parsing flags")
	flag.Parse()

	log.Print("[DEBUG] Loading config file")
	fileConfig, err := UnmarshalFile(*ConfigPath)
	if err != nil {
		// fileConfig is /required/ for Services config, cannot function w/o it
		log.Fatalf("Could not parse config file %s: %s", ConfigPath, err)
	}

	// Set global values, prefering flag value, then file, then default.
	log.Print("[DEBUG] Selecting ChefServerUrl")
	if *ChefServerUrl == "" {
		if fileConfig.ChefServerUrl != "" {
			ChefServerUrl = &fileConfig.ChefServerUrl
		} else {
			ChefServerUrl = &defaultChefServerUrl
		}
	}
	log.Print("[DEBUG] ChefServerUrl: ", *ChefServerUrl)

	log.Print("[DEBUG] Selecting ChefClientUsername")
	if *ChefClientUsername == "" {
		if fileConfig.ChefClientUsername != "" {
			ChefClientUsername = &fileConfig.ChefClientUsername
		} else {
			// Default client-username is FQDN if not specified in flags or file
			hostname, err := os.Hostname()
			if err != nil {
				// ClientUsername is /required/ to access chef, cannot function w/o it
				log.Fatal("Could not find hostname.")
			}
			ChefClientUsername = &hostname
		}
	}
	log.Print("[DEBUG] ChefClientUsername: ", ChefClientUsername)

	log.Print("[DEBUG] Selecting ChefClientKeyPath")
	if *ChefClientKeyPath == "" {
		if fileConfig.ChefClientKeyPath != "" {
			ChefClientKeyPath = &fileConfig.ChefClientKeyPath
		} else {
			ChefClientKeyPath = &defaultChefKeyPath
		}
	}
	log.Print("[DEBUG] ChefClientKeyPath: ", ChefClientKeyPath)

	log.Print("[DEBUG] Selecting ChefEnvironment")
	if *ChefEnvironment == "" {
		if fileConfig.ChefEnvironment != "" {
			ChefEnvironment = &fileConfig.ChefEnvironment
		} else {
			c := ChefClient()

			node, _ := c.Nodes.Get(*ChefClientUsername)
			// it doesn't matter if response is !ok or err, node.Environment == "" in each case
			if node.Environment != "" {
				ChefEnvironment = &node.Environment
			} else {
				log.Print("[WARN] Couldn't find node to set ChefEnvironment, using _default")
				ChefEnvironment = &defaultChefEnvironment
			}
		}
	}
	log.Print("[DEBUG] ChefEnvironment: ", ChefEnvironment)

	// services must be set from file
	Services = fileConfig.Services
}

func ChefClient() *chefClient.Client {
	// read a client key
	log.Print("[DEBUG] Reading key:", *ChefClientKeyPath)
	key, err := ioutil.ReadFile(*ChefClientKeyPath)
	if err != nil {
		log.Fatal("Couldn't read key:", *ChefClientKeyPath, err)
	}

	config := chefClient.Config{
		Name:    *ChefClientUsername,
		Key:     string(key),
		BaseURL: *ChefServerUrl,
		SkipSSL: true,
	}

	c, err := chefClient.NewClient(&config)
	if err != nil {
		log.Fatalf("Could not connect to Chef:", err)
	}
	return c
}
