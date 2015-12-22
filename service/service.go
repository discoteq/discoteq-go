// Service.go specifies the common interface for service discovery engines
// You'll notice this does not pass a config state in, currently any externaly visible scope is just global. Mostly this stuff is resource connection locators or client configuration. Too bad the resources can't all be adapted into URIs.

package discoteq

// Service can list its hosts
type Service interface {
	Name() string
	HostRecordList() ServiceHostRecordList
}
