package accesslog

import (
	"megaease/access-log-go/accesslog/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (m *AccessLogMiddleware) GetGinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.checkSkip(c.Request) {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		log := api.NewAccessLog()
		log.SetService(m.serviceName, m.hostname, "")
		log.SetRequest(c.Request, c.FullPath(), c.ClientIP())
		log.SetResponse(c.Writer.Status(), int64(c.Writer.Size()), time.Since(start).Milliseconds())

		err := m.backend.Send(log)
		if err != nil {
			logrus.Errorf("send access log failed: %v", err)
		}
	}
}
