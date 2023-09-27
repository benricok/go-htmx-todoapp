package main

import (
	"fmt"
	"go-htmx/auth"
	"go-htmx/database"
	"go-htmx/endpoints"
	"go-htmx/user"
	"html/template"
	"log"
	"os"

	// trunk-ignore(trufflehog/SQLServer)
	"strconv"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func GetDBcredentials() string {
	var (
		db_host      string = os.Getenv("POSTGRES_HOST")
		db_port, err        = strconv.Atoi(os.Getenv("POSTGRES_PORT"))
		db_user      string = os.Getenv("POSTGRES_USER")
		db_password  string = os.Getenv("POSTGRES_PASSWORD")
		db_name      string = os.Getenv("POSTGRES_DB")
	)

	if err != nil {
		log.Fatalf("Could not determine DB URL: %+v", err)
	}

	return fmt.Sprintf("host=%s port=%d user=%s "+
		// trunk-ignore(trufflehog/SQLServer)
		"password=%s dbname=%s sslmode=disable",
		db_host, db_port, db_user, db_password, db_name)
}

func main() {
	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		db_url = GetDBcredentials()
	}

	err := database.NewDatabase(db_url)
	if err != nil {
		log.Fatalf("Could not init db: %+v", err)
	}

	err = database.AddUser(user.LoadTestUser())
	if err != nil {
		log.Fatalf("Could not add user: %+v", err)
	}

	user, err := database.GetUser("test")
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
			NewClaimsFunc: auth.Claim,
			SigningKey:    []byte(auth.GetJWTSecret()),
			TokenLookup:   "cookie:access-token",
			ErrorHandler:  auth.JWTErrorChecker,
		}))

		app.GET("/home", endpoints.HandleHome)
		app.GET("/settings", endpoints.HandleSettings)
		app.GET("/help", endpoints.HandleHelp)
	}

	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}
