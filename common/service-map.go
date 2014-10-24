package discoteq

import (
	"encoding/json"
)

type ServiceHostRecord map[string]string
type ServiceHostRecordList []ServiceHostRecord
type ServiceMap map[string]ServiceHostRecordList

// Sorting should occur by hostname attribute
func (a ServiceHostRecordList) Len() int           { return len(a) }
func (a ServiceHostRecordList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ServiceHostRecordList) Less(i, j int) bool { return a[i]["hostname"] < a[j]["hostname"] }

func (sm *ServiceMap) Marshal() ([]byte, error) {
	return json.MarshalIndent(sm, "", "  ")
}
