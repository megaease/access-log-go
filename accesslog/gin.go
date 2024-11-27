package accesslog

import (
	"megaease/access-log-go/accesslog/api"
	"megaease/access-log-go/accesslog/utils/fasttime"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetGinMiddleware returns the Gin middleware for access log.
func (m *AccessLogMiddleware) GetGinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.checkSkip(c.Request) {
			c.Next()
			return
		}

		start := fasttime.Now()
		c.Next()

		log := api.NewAccessLog(m.serviceName)
		log.SetRequest(c.Request, c.FullPath(), c.ClientIP())
		log.SetResponse(c.Writer.Status(), int64(c.Writer.Size()), fasttime.Since(start).Milliseconds())

		err := m.backend.Send(log)
		if err != nil {
			logrus.Errorf("send access log failed: %v", err)
		}
	}
}
