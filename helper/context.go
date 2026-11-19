package helper

import (
	"github.com/gofiber/fiber/v2"

	"github.com/vyerru/lapor-kampus/app/model"
)

const LocalsAuthUser = "authUser"

// CurrentUser mengambil identitas yang sudah ditempelkan middleware RequireAuth.
func CurrentUser(c *fiber.Ctx) (model.AuthUser, bool) {
	user, ok := c.Locals(LocalsAuthUser).(model.AuthUser)
	return user, ok
}
