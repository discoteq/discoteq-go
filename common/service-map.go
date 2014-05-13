package discoteq

import (
	"encoding/json"
)

type ServiceHostRecord map[string]interface{}
type ServiceHostRecordList []ServiceHostRecord
type ServiceMap map[string]ServiceHostRecordList

func (sm *ServiceMap) Marshal() ([]byte, error) {
	return json.MarshalIndent(sm, "", "  ")
}
