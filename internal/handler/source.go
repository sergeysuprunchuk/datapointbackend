package handler

import (
	"datapointbackend/internal/entity"
	"datapointbackend/internal/service"
	"github.com/gofiber/fiber/v2"
)

type sourceHandler struct {
	ss *service.SourceService
}

func newSourceHandler(app *fiber.App, ss *service.SourceService) {
	h := sourceHandler{ss: ss}
	g := app.Group("/sources")
	g.Get("/", h.getAll)
	g.Get("/drivers", h.getDrivers)
	g.Get("/:id", h.getOne)
	g.Post("/", h.create)
	g.Patch("/", h.edit)
	g.Delete("/:id", h.delete)
}

// @tags		источники
// @success	200	{array}	entity.Source
// @router		/sources [get]
func (h *sourceHandler) getAll(ctx *fiber.Ctx) error {
	sl, err := h.ss.GetAll(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(sl)
}

// @tags		источники
// @param		id	path		string	true	"идентификатор источника"
// @success	200	{object}	entity.Source
// @router		/sources/{id} [get]
func (h *sourceHandler) getOne(ctx *fiber.Ctx) error {
	source, err := h.ss.GetOne(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}

	return ctx.JSON(source)
}

// @tags	источники
// @param	source	body	entity.Source	true	"источник"
// @router	/sources [patch]
func (h *sourceHandler) edit(ctx *fiber.Ctx) error {
	var source entity.Source

	err := ctx.BodyParser(&source)
	if err != nil {
		return err
	}

	if err = h.ss.Edit(ctx.Context(), source); err != nil {
		return err
	}

	return nil
}

// @tags	источники
// @param	id	path	string	true	"идентификатор источника"
// @router	/sources/{id} [delete]
func (h *sourceHandler) delete(ctx *fiber.Ctx) error {
	err := h.ss.Delete(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}
	return err
}

// @tags	источники
// @param	source	body	entity.Source	true	"источник"
// @router	/sources [post]
func (h *sourceHandler) create(ctx *fiber.Ctx) error {
	var source entity.Source

	err := ctx.BodyParser(&source)
	if err != nil {
		return err
	}

	if source.Id, err = h.ss.Create(ctx.Context(), source); err != nil {
		return err
	}

	return ctx.SendString(source.Id)
}

// @tags		источники
// @success	200	{array}	string
// @router		/sources/drivers [get]
func (h *sourceHandler) getDrivers(ctx *fiber.Ctx) error {
	return ctx.JSON(h.ss.GetDrivers())
}
