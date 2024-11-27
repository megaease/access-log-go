package accesslog

import (
	"megaease/access-log-go/accesslog/api"
	"time"

	"github.com/labstack/echo/v4"
)

func (m *AccessLogMiddleware) GetEchoMiddleWare() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if m.checkSkip(ctx.Request()) {
				return next(ctx)
			}

			start := time.Now()
			err := next(ctx)

			log := api.NewAccessLog()
			log.SetService(m.serviceName, m.hostname, "")
			log.SetRequest(ctx.Request(), ctx.Path(), ctx.RealIP())
			log.SetResponse(ctx.Response().Status, ctx.Response().Size, time.Since(start).Milliseconds())
			m.backend.Send(log)
			return err
		}
	}
}
