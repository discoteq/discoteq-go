package chef

import (
	"fmt"
	"log"
	"sort"
	"strings"

	chefClient "github.com/go-chef/chef"

	"github.com/discoteq/discoteq-go/chef/config"
	"github.com/discoteq/discoteq-go/common"
)

var (
	defaultAttrsFqdn = map[string]string{"hostname": "fqdn"}
)

// chef.Service: query object generated from config
type Service struct {
	Name  string
	Query string
	Attrs map[string]string
}

// Node attribute map
type ChefNodeMap map[string]interface{}

func (s *Service) FullQuery() string {
	return s.Query
}
func ServiceFromRaw(name string, raw map[string]interface{}) *Service {
	service := new(Service)
	service.Name = name

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
					log.Printf("[WARN] Could not cast attribute into string: %v", attrs[k])
				}
			}
		} else {
			log.Printf("[WARN] Could not cast attributes into map[string]interface{}: %v", raw["attrs"])
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
	discoveredService := make(discoteq.ServiceHostRecordList, 0)

	for _, node := range searchResults.Rows {

		attrs := make(discoteq.ServiceHostRecord)
		nodeMap, _ := node.(ChefNodeMap)
		mergedNodeMap := mergeNodeAttrs(nodeMap)
		requestedAttrs := s.Attrs
		for k, v := range requestedAttrs {
			attrs[k] = getAttr(mergedNodeMap, v).(string)
		}

		discoveredService = append(discoveredService, attrs)
	}

	sort.Sort(discoveredService)

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
		current, _ = current[seg].(ChefNodeMap)
	}
	return result
}

// take a node with default, normal and automatic attributes
// and return a single merged map of the highest precedence values
func mergeNodeAttrs(node ChefNodeMap) ChefNodeMap {
	// default is a keyword, dfault will have to do
	dfault, _ := node["default"].(ChefNodeMap)
	normal, _ := node["normal"].(ChefNodeMap)
	automatic, _ := node["automatic"].(ChefNodeMap)
	// merge together attributes with automatic at highest precedence,
	// followed by normal, followed by default
	result := mergeAttrMap(mergeAttrMap(dfault, normal), automatic)
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
