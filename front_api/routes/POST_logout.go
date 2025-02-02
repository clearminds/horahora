package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Route: POST /logout
// Accepts no parameters
// response: 200 if ok
func (r RouteHandler) handleLogout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = ""

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, nil)
}
