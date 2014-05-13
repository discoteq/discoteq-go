// Parse configuration from command-line flags and json file
package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	chefClient "github.com/marpaia/chef-golang"
)

const ()

var (
	ConfigPath             = flag.String("c", "/etc/discoteq-chef.json", "config file path")
	defaultServerUrl       = "http://localhost:4545"
	ServerUrl              = flag.String("s", "", "chef server URL, default: http://localhost:4545")
	ClientUsername         = flag.String("u", "", "API client username, default: node.fqdn")
	defaultKeyPath         = "/etc/chef/client.pem"
	ClientKeyPath          = flag.String("k", "", "API client key file path, default: /etc/chef/client.pem")
	defaultChefEnvironment = "_default"
	ChefEnvironment        = flag.String("E", "", "environment query scope, default: node.chef_environment || _default")
	Services               map[string]map[string]interface{}
)

type Config struct {
	ServerUrl       string                            `json:server-url`
	ClientUsername  string                            `json:user`
	ClientKeyPath   string                            `json:key`
	ChefEnvironment string                            `json:environment`
	Services        map[string]map[string]interface{} `json:"services"`
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
	flag.Parse()

	// load conf file
	fileConfig, err := UnmarshalFile(*ConfigPath)
	if err != nil {
		// fileConfig is /required/ for Services config, cannot function w/o it
		log.Fatalf("Could not parse config file %s: %s", ConfigPath, err)
	}

	// Set global values, prefering flag value, then file, then default.
	if *ServerUrl == "" {
		if fileConfig.ServerUrl == "" {
			ServerUrl = &fileConfig.ServerUrl
		} else {
			ServerUrl = &defaultServerUrl
		}
	}

	if *ClientUsername == "" {
		if fileConfig.ClientUsername == "" {
			ClientUsername = &fileConfig.ClientUsername
		} else {
			// Default client-username is FQDN if not specified in flags or file
			hostname, err := os.Hostname()
			if err != nil {
				// ClientUsername is /required/ to access chef, cannot function w/o it
				log.Fatal("Could not find hostname.")
			}
			ClientUsername = &hostname
		}
	}

	if *ClientKeyPath == "" {
		if fileConfig.ClientKeyPath == "" {
			ClientKeyPath = &fileConfig.ClientKeyPath
		} else {
			ClientKeyPath = &defaultKeyPath
		}
	}

	if *ChefEnvironment == "" {
		if fileConfig.ChefEnvironment != "" {
			ChefEnvironment = &fileConfig.ChefEnvironment
		} else {
			c, err := chefClient.Connect()
			if err != nil {
				log.Fatalf("Could not retrieve node from chef:", err)
			}
			c.SSLNoVerify = true

			node, _, _ := c.GetNode(*ClientUsername)
			// it doesn't matter if response is !ok or err, node.Environment == "" in each case
			if node.Environment != "" {
				ChefEnvironment = &node.Environment
			} else {
				log.Print("Couldn't find node to set ChefEnvironment, using _default")
				ChefEnvironment = &defaultChefEnvironment
			}
		}
	}

	// services must be set from file
	Services = fileConfig.Services
}
