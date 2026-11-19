package route

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	"github.com/vyerru/lapor-kampus/helper"
)

// Dependencies adalah kumpulan hal yang dibutuhkan handler.
// route tidak mengimpor config (config mengimpor route), jadi nilai seperti
// nama aplikasi dan versi diteruskan dari main.
type Dependencies struct {
	Pool        *pgxpool.Pool
	Mongo       *mongo.Client
	JWT         *helper.JWTManager
	Permissions *helper.PermissionSet

	AppName string
	Env     string
	Version string
	Commit  string
}

func Register(app *fiber.App, deps Dependencies) {
	api := app.Group("/api/v1")

	api.Get("/health", healthCheck(deps))
	api.Get("/readyz", readinessCheck(deps))

	// Titik tumbuh endpoint berikutnya:
	//   auth      := api.Group("/auth", middleware.RequireJSON)
	//   students  := api.Group("/students", middleware.RequireJSON, middleware.RequireAuth(deps.JWT))
}

// healthCheck hanya menjawab "proses ini hidup" — tidak menyentuh dependency.
func healthCheck(deps Dependencies) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return helper.Success(c, fiber.StatusOK, "aplikasi sehat", fiber.Map{
			"service": deps.AppName,
			"env":     deps.Env,
			"version": deps.Version,
			"commit":  deps.Commit,
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
	}
}

// readinessCheck menjawab "apakah dependency siap melayani trafik".
func readinessCheck(deps Dependencies) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()

		if deps.Pool != nil {
			if err := deps.Pool.Ping(ctx); err != nil {
				return helper.ServiceUnavailable("postgres belum siap")
			}
		}
		if deps.Mongo != nil {
			if err := deps.Mongo.Ping(ctx, readpref.Primary()); err != nil {
				return helper.ServiceUnavailable("mongodb belum siap")
			}
		}

		return helper.Success(c, fiber.StatusOK, "semua dependency siap", fiber.Map{
			"postgres": deps.Pool != nil,
			"mongodb":  deps.Mongo != nil,
			"time":     time.Now().UTC().Format(time.RFC3339),
		})
	}
}
