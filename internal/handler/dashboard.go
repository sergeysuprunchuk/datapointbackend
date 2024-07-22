package handler

import (
	"datapointbackend/internal/entity"
	"datapointbackend/internal/service"
	"github.com/gofiber/fiber/v2"
)

type dashboardHandler struct {
	ds *service.DashboardService
}

func newDashboardHandler(app *fiber.App, ds *service.DashboardService) {
	h := dashboardHandler{ds: ds}
	g := app.Group("/dashboards")
	g.Get("/", h.getAll)
	g.Get("/:id", h.getOne)
	g.Post("/", h.create)
	g.Delete("/:id", h.delete)
}

func (h *dashboardHandler) getAll(ctx *fiber.Ctx) error {
	sl, err := h.ds.GetAll(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(sl)
}

func (h *dashboardHandler) getOne(ctx *fiber.Ctx) error {
	return nil
}

func (h *dashboardHandler) delete(ctx *fiber.Ctx) error {
	return nil
}

func (h *dashboardHandler) create(ctx *fiber.Ctx) error {
	var dashboard entity.Dashboard

	err := ctx.BodyParser(&dashboard)
	if err != nil {
		return err
	}

	if dashboard.Id, err = h.ds.Create(ctx.Context(), dashboard); err != nil {
		return err
	}

	return ctx.SendString(dashboard.Id)
}
