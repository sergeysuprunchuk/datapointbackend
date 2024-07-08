package handler

import (
	"datapointbackend/internal/entity"
	"datapointbackend/internal/service"
	"github.com/gofiber/fiber/v2"
)

type queryHandler struct {
	qs *service.QueryService
}

func newQueryHandler(app *fiber.App, qs *service.QueryService) {
	h := queryHandler{qs: qs}
	group := app.Group("/queries")
	group.Post("/execute", h.execute)
}

// @tags	запросы
// @param	query	body	entity.Query	true	"запрос"
// @router	/queries/execute [post]
func (h *queryHandler) execute(ctx *fiber.Ctx) error {
	var query entity.Query

	err := ctx.BodyParser(&query)
	if err != nil {
		return err
	}

	return ctx.JSON(h.qs.Execute(ctx.Context(), query))
}
