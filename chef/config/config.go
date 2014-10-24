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
  log.Print("[DEBUG] Parsing flags")
	flag.Parse()

  log.Print("[DEBUG] Loading config file")
	fileConfig, err := UnmarshalFile(*ConfigPath)
	if err != nil {
		// fileConfig is /required/ for Services config, cannot function w/o it
		log.Fatalf("Could not parse config file %s: %s", ConfigPath, err)
	}

	// Set global values, prefering flag value, then file, then default.
   log.Print("[DEBUG] Selecting ServerUrl")
	if *ServerUrl == "" {
		if fileConfig.ServerUrl != "" {
			ServerUrl = &fileConfig.ServerUrl
		} else {
			ServerUrl = &defaultServerUrl
		}
	}
  log.Print("[DEBUG] ServerUrl: ", *ServerUrl)

  log.Print("[DEBUG] Selecting ClientUsername")
	if *ClientUsername == "" {
		if fileConfig.ClientUsername != "" {
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
  log.Print("[DEBUG] ClientUsername: ", ClientUsername)

  log.Print("[DEBUG] Selecting ClientKeyPath")
	if *ClientKeyPath == "" {
		if fileConfig.ClientKeyPath != "" {
			ClientKeyPath = &fileConfig.ClientKeyPath
		} else {
			ClientKeyPath = &defaultKeyPath
		}
	}
  log.Print("[DEBUG] ClientKeyPath: ", ClientKeyPath)

  log.Print("[DEBUG] Selecting ChefEnvironment")
	if *ChefEnvironment == "" {
		if fileConfig.ChefEnvironment != "" {
			ChefEnvironment = &fileConfig.ChefEnvironment
		} else {
      c := ChefClient()

			node, _ := c.Nodes.Get(*ClientUsername)
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
  log.Print("[DEBUG] Reading key:", *ClientKeyPath)
	key, err := ioutil.ReadFile(*ClientKeyPath)
	if err != nil {
    log.Fatal("Couldn't read key:", *ClientKeyPath, err)
	}

  config := chefClient.Config{
    Name: *ClientUsername,
    Key: string(key),
    BaseURL: *ServerUrl,
    SkipSSL: true,
  }

	c, err := chefClient.NewClient(&config)
	if err != nil {
		log.Fatalf("Could not connect to Chef:", err)
	}
	return c
}
