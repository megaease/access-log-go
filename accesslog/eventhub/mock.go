package eventhub

import "megaease/access-log-go/accesslog/api"

type (
	// eventHubMock is the mock event hub.
	EventHubMock struct {
		Logs []*api.AccessLog
	}
)

// Send sends the access log.
func (m *EventHubMock) Send(log *api.AccessLog) error {
	if m.Logs == nil {
		m.Logs = []*api.AccessLog{}
	}
	m.Logs = append(m.Logs, log)
	return nil
}

// Close closes the mock event hub.
func (m *EventHubMock) Close() {
}
