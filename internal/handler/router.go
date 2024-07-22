package handler

import (
	_ "datapointbackend/docs"
	"datapointbackend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func NewRouter(
	app *fiber.App,
	ss *service.SourceService,
	qs *service.QueryService,
	ws *service.WidgetService,
	ds *service.DashboardService,
) {
	app.Get("/swagger/*", swagger.HandlerDefault)
	newSourceHandler(app, ss)
	newQueryHandler(app, qs)
	newWidgetHandler(app, ws)
	newDashboardHandler(app, ds)
}
