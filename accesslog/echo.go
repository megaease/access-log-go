package accesslog

import (
	"megaease/access-log-go/accesslog/api"
	"time"

	"github.com/labstack/echo/v4"
)

// GetEchoMiddleWare returns the Echo middleware for access log.
func (m *AccessLogMiddleware) GetEchoMiddleWare() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if m.checkSkip(ctx.Request()) {
				return next(ctx)
			}

			start := time.Now()
			err := next(ctx)

			log := api.NewAccessLog(m.serviceName)
			log.SetRequest(ctx.Request(), ctx.Path(), ctx.RealIP())
			log.SetResponse(ctx.Response().Status, ctx.Response().Size, time.Since(start).Milliseconds())
			m.backend.Send(log)
			return err
		}
	}
}
