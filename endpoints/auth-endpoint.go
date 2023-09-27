package endpoints

import (
	"go-htmx/auth"
	"go-htmx/user"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func HandleLoginForm(c echo.Context) error {
	return c.Render(200, "login.html", Page{
		Header: Header{
			Title: "TODO - Login",
		},
	})
}

func Login(c echo.Context) error {
	storedUser := user.LoadTestUser()

	u := new(user.User)

	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(u.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Password is incorrect")
	}
	
	if err := auth.GenerateTokensAndSetCookies(storedUser, c); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Token is incorrect")
	}

	return c.Redirect(http.StatusMovedPermanently, "/app/home")
}

func Logout(c echo.Context) error {
	usercookie := new(http.Cookie)
	usercookie.Name = "user"
	usercookie.MaxAge = -1
	usercookie.Path = "/"
	c.SetCookie(usercookie)

	cookie := new(http.Cookie)
	cookie.Name = auth.GetAccessTokenCookieName()
	cookie.MaxAge = -1
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.Redirect(http.StatusMovedPermanently, "/login")
}