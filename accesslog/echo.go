package accesslog

import (
	"github.com/megaease/access-log-go/accesslog/api"
	"github.com/megaease/access-log-go/accesslog/utils/fasttime"

	"github.com/labstack/echo/v4"
)

// GetEchoMiddleWare returns the Echo middleware for access log.
func (m *AccessLogMiddleware) GetEchoMiddleWare() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if m.hostIP == "" {
				m.setHostIP(ctx.Request())
			}
			if m.checkSkip(ctx.Request()) {
				return next(ctx)
			}

			start := fasttime.Now()
			err := next(ctx)

			log := api.NewAccessLog(m.serviceName, m.hostName, m.tenantID)
			log.SetRequest(ctx.Request(), ctx.Path(), ctx.RealIP(), m.hostIP)
			log.SetResponse(ctx.Response().Status, ctx.Response().Size, fasttime.Since(start).Milliseconds())
			m.backend.Send(log)
			return err
		}
	}
}
