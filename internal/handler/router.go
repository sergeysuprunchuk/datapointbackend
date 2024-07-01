package handler

import (
	_ "datapointbackend/docs"
	"datapointbackend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func NewRouter(app *fiber.App, ss *service.SourceService, qs *service.QueryService) {
	app.Get("/swagger/*", swagger.HandlerDefault)
	newSourceHandler(app, ss)
	newQueryHandler(app, qs)
}
