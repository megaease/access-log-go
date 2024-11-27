package accesslog

import "megaease/access-log-go/accesslog/eventhub"

type (
	Config struct {
		Backend eventhub.Config
	}
)
