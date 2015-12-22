// Parser.go specifies the common interface for service discovery engines
// You'll notice this does not pass a config state in, currently any externaly visible scope is just global. Mostly this stuff is resource connection locators or client configuration. Too bad the resources can't all be adapted into URIs.

package parser

import (
	"fmt"
	"github.com/discoteq/discoteq-go/service"
	"github.com/discoteq/discoteq-go/chef"
)


// ServiceFromRaw passes a the data from a service config to the appropriate engine
func Parse(name string, raw map[string]interface{}) (discoteq.Service, error) {
	engine, ok := raw["engine"].(string)
    if !ok {
        return nil, fmt.Errorf("No engine provided for service:%q.", name)
        // return nil, fmt.Errorf("No engine provided for service: who knows.")
    }

	if engine == "chef" {
		return chef.ServiceFromRaw(name, raw), nil
	} else {
		return nil, fmt.Errorf("Unknown engine:%s for service:%s. Supported engines: chef.", name, engine)
	}
}
