package chef

import (
	"fmt"
	"log"
	"sort"
	"strings"

	chefClient "github.com/go-chef/chef"

	"github.com/discoteq/discoteq-go/config"
	"github.com/discoteq/discoteq-go/service"
)

var (
	defaultAttrsFqdn = map[string]string{"hostname": "fqdn"}
)

// chef.Service: query object generated from config
type Service struct {
	name  string
	Query string
	Attrs map[string]string
}

func (s *Service) Name() string {
	return s.name
}

// Node attribute map
type ChefNodeMap map[string]interface{}

func (s *Service) FullQuery() string {
	return s.Query
}
func ServiceFromRaw(name string, raw map[string]interface{}) *Service {
	service := new(Service)
	service.name = name

	// build query
	role, _ := raw["role"].(string)
	tag, _ := raw["tag"].(string)
	raw_query, _ := raw["query"].(string)
	if (raw_query != "" && (role != "" || tag != "")) || (role != "" && tag != "") {
		log.Fatalf("Service %s sets more than of of the query, role and tag attributes. Please only define a single one.", name)
	}

	include_chef_environment, ok := raw["include_chef_environment"].(bool)
	if !ok {
		// default to true
		include_chef_environment = true
	}

	query := ""

	if role != "" {
		query += fmt.Sprintf("role:%s", role)
	}
	if tag != "" {
		query += fmt.Sprintf("tag:%s", tag)
	}
	if raw_query != "" {
		query += raw_query
	}
	if include_chef_environment {
		query += fmt.Sprintf(" AND chef_environment:%s", *config.ChefEnvironment)
	}
	service.Query = query

	service.Attrs = make(map[string]string)
	defaultAttrs := defaultAttrsFqdn
	for k, v := range defaultAttrs {
		service.Attrs[k] = v
	}
	if raw["attrs"] != nil {
		attrs, ok := raw["attrs"].(map[string]interface{})
		if ok {
			for k := range attrs {
				v, ok := attrs[k].(string)
				if ok {
					service.Attrs[k] = v
				} else {
					log.Printf("[WARN] Could not cast attribute into string: %q", attrs[k])
				}
			}
		} else {
			log.Printf("[WARN] Could not cast attributes into map[string]interface{}: %q", raw["attrs"])
		}
	}

	return service
}

func (s *Service) HostRecordList() discoteq.ServiceHostRecordList {
	log.Print("[DEBUG] Entering HostRecordList()")
	c := config.ChefClient()
	// request service data
	query := s.FullQuery()
	log.Print("[DEBUG] Searching with query: ", query)
	searchResults, err := c.Search.Exec("node", query)
	if err != nil {
		log.Fatalf("Could not search for nodes with query:\"%s\", error: %s", query, err)
	}
	log.Print("[DEBUG] Searching results: ", searchResults)
	return s.hostRecordListFromResults(&searchResults)
}

func (s *Service) hostRecordListFromResults(searchResults *chefClient.SearchResult) discoteq.ServiceHostRecordList {
	log.Print("[DEBUG] Entering hostRecordListFromResults")
	discoveredService := make(discoteq.ServiceHostRecordList, 0)

	for _, node := range searchResults.Rows {
		log.Printf("[DEBUG] node:%q", node)

		attrs := make(discoteq.ServiceHostRecord)
		nodeMap, ok := node.(map[string]interface{})
		if !ok {
			log.Printf("[DEBUG] node could not be cast to map[string]interface{}: %q", node)
		}

		log.Printf("[DEBUG] nodeMap:%q", nodeMap)
		mergedNodeMap := mergeNodeAttrs(nodeMap)
		log.Printf("[DEBUG] mergedNodeMap:%q", mergedNodeMap)
		requestedAttrs := s.Attrs
		for k, v := range requestedAttrs {
			log.Printf("[DEBUG] requested attr k:%q v:%q", k, v)
			switch val := getAttr(mergedNodeMap, v).(type) {
			default:
				log.Printf("[WARN] Could not identify type of attr: %q", val)
				attrs[k] = ""
			case string:
				attrs[k] = val
			case float64:
				attrs[k] = fmt.Sprintf("%v", val)
			}
		}

		discoveredService = append(discoveredService, attrs)
	}

	sort.Sort(discoveredService)

	log.Print("[DEBUG] discoveredService: ", discoveredService)
	return discoveredService
}

// query a node attribute map using a query string with simplified
// syntax, so that foo.bar.baz is equivalent to node["foo"]["bar"]["baz"],
// and returning nil in the event of invalid access
func getAttr(node ChefNodeMap, query string) interface{} {
	segments := strings.Split(query, ".")
	current := node
	var result interface{}
	for _, seg := range segments {
		result = current[seg]
		// descent into empty map doesn't matter, it
		// correctly returns null regardless
		next, _ := current[seg].(map[string]interface{})
		current = next
	}
	return result
}

// take a node with default, normal and automatic attributes
// and return a single merged map of the highest precedence values
func mergeNodeAttrs(node ChefNodeMap) ChefNodeMap {
	log.Print("[DEBUG] mergeNodeAttrs()")
	// default is a keyword, dfault will have to do
	dfault, _ := node["default"].(map[string]interface{})
	log.Printf("[DEBUG] default:%q", dfault)
	normal, _ := node["normal"].(map[string]interface{})
	log.Printf("[DEBUG] normal:%q", normal)
	automatic, _ := node["automatic"].(map[string]interface{})
	log.Printf("[DEBUG] automatic:%q", automatic)
	// merge together attributes with automatic at highest precedence,
	// followed by normal, followed by default
	result := mergeAttrMap(mergeAttrMap(dfault, normal), automatic)
	log.Printf("[DEBUG] result:%q", result)
	return result
}

// merge n into m, preferring values from n
func mergeAttrMap(m, n map[string]interface{}) map[string]interface{} {
	result := m
	for k := range n {
		mmap, mok := m[k].(map[string]interface{})
		nmap, nok := n[k].(map[string]interface{})
		if mok && nok {
			result[k] = mergeAttrMap(mmap, nmap)
		} else {
			result[k] = n[k]
		}
	}
	return result
}
