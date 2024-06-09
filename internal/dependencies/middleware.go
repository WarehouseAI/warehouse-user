package dependencies

import "github.com/warehouse/user-service/internal/handler/middlewares"

func (d *dependencies) HandlerMiddleware() middlewares.Middleware {
	if d.handlerMiddleware == nil {
		d.handlerMiddleware = middlewares.NewMiddleware(
			d.log,
			d.cfg.Timeouts,
			d.AuthAdapter(),
		)
	}

	return d.handlerMiddleware
}
