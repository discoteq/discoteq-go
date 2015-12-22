package consul

import (
	"fmt"
	"log"
	"sort"
	"strings"

	consulApi "github.com/hashicorp/consul/api"

	"github.com/discoteq/discoteq-go/service"
)

// consul.Service: query object generated from config
type Service struct {
	name        string
	Service     string
	Datacenters []string
	OnlyPassing bool
	Tags        []string
	Attrs       map[string]string
}

func (s *Service) Name() string {
	return s.name
}

func ServiceFromRaw(name string, raw map[string]interface{}) *Service {
	service := new(Service)
	service.name = name
	service.Service, _ = raw["service"].(string)
	return service
}

func (s *Service) HostRecordList() discoteq.ServiceHostRecordList {
	log.Print("[DEBUG] Entering HostRecordList()")
	c, err := consulApi.NewClient(consulApi.DefaultConfig())
	// Get a new client
	if err != nil {
		log.Fatalf("Could not create consul client, error: %s", err)
	}

	// request service data
	service := s.Service
	log.Print("[DEBUG] Searching for service: ", service)

	tag := ""
	q := new(consulApi.QueryOptions)

	svcs, _, err := c.Catalog().Service(service, tag, q)

	if err != nil {
		log.Fatalf("Could not query catalog for service:\"%s\", tag:\"%s\",  q:\"%s\", error: %s", service, tag, q, err)
	}
	log.Print("[DEBUG] Searching results: ", svcs)

	return s.hostRecordListFromResults(svcs)
}

func (s *Service) hostRecordListFromResults(svcs []*consulApi.CatalogService) discoteq.ServiceHostRecordList {
	discoveredService := make(discoteq.ServiceHostRecordList, 0)

	for _, svc := range svcs {
		attrs := make(discoteq.ServiceHostRecord)
		attrs["name"] = svc.ServiceName
		attrs["hostname"] = svc.ServiceAddress
		attrs["node"] = svc.Node
		// attrs["node_address"] = svc.Address
		attrs["port"] = fmt.Sprintf("%d", svc.ServicePort)
		attrs["tags"] = strings.Join(svc.ServiceTags, ",") // TODO: toJSON?
		attrs["id"] = svc.ServiceID

		discoveredService = append(discoveredService, attrs)
	}

	sort.Sort(discoveredService)

	return discoveredService
}
