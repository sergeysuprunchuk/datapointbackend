package app

import (
	"datapointbackend/config"
	"datapointbackend/internal/handler"
	"datapointbackend/internal/repository"
	"datapointbackend/internal/service"
	"datapointbackend/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jsoniter "github.com/json-iterator/go"
)

func Run(cfg *config.Config) error {
	app := fiber.New(fiber.Config{
		JSONEncoder: jsoniter.Marshal,
		JSONDecoder: jsoniter.Unmarshal,
	})

	app.Use(cors.New(cors.Config{
		Next: func(c *fiber.Ctx) bool {
			return false
		},
	}))

	db, err := database.New(cfg.Database)
	if err != nil {
		return err
	}

	var (
		sr = repository.NewSourceRepository(db)
		wr = repository.NewWidgetRepository(db)
		dr = repository.NewDashboardRepository(db)
	)

	var (
		ss = service.NewSourceService(sr)
		qs = service.NewQueryService(ss)
		ws = service.NewWidgetService(wr)
		ds = service.NewDashboardService(dr)
	)

	handler.NewRouter(app, ss, qs, ws, ds)

	return app.Listen(cfg.Http.Addr)
}
