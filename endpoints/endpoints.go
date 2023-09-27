package endpoints

import (
	"go-htmx/database"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Header struct {
    Title string
}

type Nav struct {
	CurrentPage string
}

type Userdata struct {
	Username string
	Name 	 string
	Surname  string
	Email 	 string
}

type Page struct {
	Header
	Nav
	Userdata 
}

func LoadUserDataFromCookie(c echo.Context) (*Userdata) {
	cookie, err := c.Cookie("user")
	if err != nil {
		log.Printf("Could not get cookie from request: %+v", err)
		return nil
	}

	user, err := database.GetUser(cookie.Value)

	return &Userdata{
		Username: user.Username,
		Name: user.Name,
		Surname: user.Surname,
		Email: user.Email,
	}
}

type TemplateRenderer struct {
	tmpl *template.Template
}

func NewTemplateRenderer(tmpls *template.Template) TemplateRenderer {
	return TemplateRenderer{tmpls}
}

func (t TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

func HandleIndex(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/login")
}

func HandleHome(c echo.Context) error {
	return c.Render(200, "home.html", Page{
		Header: Header{
			Title: "TODO - Home",
		},
		Nav: Nav{
			CurrentPage: "home",
		},
	})
}

func HandleSettings(c echo.Context) error {
	return c.Render(200, "settings.html", Page{
		Header: Header{
			Title: "TODO - Settings",
		},
		Nav: Nav{
			CurrentPage: "settings",
		},
		Userdata: *LoadUserDataFromCookie(c),
	})
}

func HandleHelp(c echo.Context) error {
	return c.Render(200, "help.html", Page{
		Header: Header{
			Title: "TODO - Help",
		},
		Nav: Nav{
			CurrentPage: "help",
		},
	})
}

