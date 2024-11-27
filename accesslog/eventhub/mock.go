package eventhub

import "megaease/access-log-go/accesslog/api"

type (
	// eventHubMock is the mock event hub.
	eventHubMock struct {
		logs []*api.AccessLog
	}
)

// newEventHubMock creates a new mock event hub.
func newEventHubMock() (*eventHubMock, error) {
	return &eventHubMock{
		logs: []*api.AccessLog{},
	}, nil
}

// Send sends the access log.
func (m *eventHubMock) Send(log *api.AccessLog) error {
	m.logs = append(m.logs, log)
	return nil
}

// Close closes the mock event hub.
func (m *eventHubMock) Close() {
}
