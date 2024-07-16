package handler

import (
	"datapointbackend/internal/entity"
	"datapointbackend/internal/service"
	"github.com/gofiber/fiber/v2"
)

type widgetHandler struct {
	ws *service.WidgetService
}

func newWidgetHandler(app *fiber.App, ws *service.WidgetService) {
	h := widgetHandler{ws: ws}
	g := app.Group("/widgets")
	g.Get("/", h.getAll)
	g.Get("/:id", h.getOne)
	g.Delete("/:id", h.delete)
	g.Post("/", h.create)
	g.Patch("/", h.edit)
}

// @tags		виджеты
// @success	200	{array}	entity.Widget
// @router		/widgets [get]
func (h *widgetHandler) getAll(ctx *fiber.Ctx) error {
	sl, err := h.ws.GetAll(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(sl)
}

// @tags		виджеты
// @param		id	path		string	true	"идентификатор виджета"
// @success	200	{object}	entity.Widget
// @router		/widgets/{id} [get]
func (h *widgetHandler) getOne(ctx *fiber.Ctx) error {
	widget, err := h.ws.GetOne(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}
	return ctx.JSON(widget)
}

// @tags	виджеты
// @param	id	path	string	true	"идентификатор виджета"
// @router	/widgets/{id} [delete]
func (h *widgetHandler) delete(ctx *fiber.Ctx) error {
	return h.ws.Delete(ctx.Context(), ctx.Params("id"))
}

// @tags	виджеты
// @param	widget	body	entity.Widget	true	"виджет"
// @router	/widgets [post]
func (h *widgetHandler) create(ctx *fiber.Ctx) error {
	var widget entity.Widget

	err := ctx.BodyParser(&widget)
	if err != nil {
		return err
	}

	if widget.Id, err = h.ws.Create(ctx.Context(), widget); err != nil {
		return err
	}

	return ctx.SendString(widget.Id)
}

func (h *widgetHandler) edit(ctx *fiber.Ctx) error {
	return nil
}
