package main

import (
	"fmt"
	"go-htmx/auth"
	"go-htmx/database"
	"go-htmx/endpoints"
	"html/template"
	"log"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	db_host     = "postgres_db"
	db_port     = 5432
	db_user     = "todoapp"
	db_password = "todoapp123"
	db_name   = "todo"
  )

func main() {
	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		db_url = fmt.Sprintf("host=%s port=%d user=%s "+
    	"password=%s dbname=%s sslmode=disable",
    	db_host, db_port, db_user, db_password, db_name)
	}

	err := database.NewDatabase(db_url)
	if err != nil {
		log.Fatalf("Could not init db: %+v", err)
	}

	//err = database.AddUser(user.LoadTestUser())
	//if err != nil {
	//	log.Fatalf("Could not add user: %+v", err)
	//}
	
	user, err:= database.GetUser("test")
	if err != nil {
		log.Fatalf("could not get user: %+v", err)
	}
	print(user.Username)


	tmpl, err := template.ParseFiles(
		"./public/login.html",
		"./public/header.html",
		"./public/nav.html",
		"./public/home.html",
		"./public/help.html",
		"./public/settings.html",
		"./public/usermng.html",
	)

	if err != nil {
		log.Fatalf("Could not initialise templates: %+v", err)
	}

    e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Renderer = endpoints.NewTemplateRenderer(tmpl)

	e.GET("/", endpoints.HandleIndex)
	e.GET("/css/style.css", func(c echo.Context) error { 
		return c.File("./public/css/style.css")
	})
	e.GET("/login", endpoints.HandleLoginForm)
	e.POST("/login", endpoints.Login)
	e.GET("/logout", endpoints.Logout)

	app := e.Group("/app")
	{
		app.Use(echojwt.WithConfig(echojwt.Config{
			NewClaimsFunc:	auth.Claim,
			SigningKey: 	[]byte(auth.GetJWTSecret()),
			TokenLookup: 	"cookie:access-token",
			ErrorHandler:	auth.JWTErrorChecker,
		}))

		app.GET("/home", endpoints.HandleHome)
		app.GET("/settings", endpoints.HandleSettings)
		app.GET("/help", endpoints.HandleHelp)
	}

    e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}