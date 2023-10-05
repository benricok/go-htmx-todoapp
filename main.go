package main

import (
	"fmt"
	"go-htmx/internal/auth"
	"go-htmx/internal/db"
	"go-htmx/internal/endpoints"
	"go-htmx/internal/user"
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

	err := db.NewDatabase(db_url)
	if err != nil {
		log.Fatalf("Could not init db: %+v", err)
	}

	u, _ := db.GetUser("test")
	if u.Username == "" {
		err = db.AddUser(user.LoadTestUser())
		if err != nil {
			log.Fatalf("Could not add user: %+v", err)
		}
	}
	print(u.Username)

	tmpl, err := template.ParseFiles(
		"./web/public/login.html",
		"./web/public/header.html",
		"./web/public/nav.html",
		"./web/public/home.html",
		"./web/public/help.html",
		"./web/public/settings.html",
		"./web/public/usermng.html",
		"./web/public/signup.html",
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
		return c.File("./web/public/css/style.css")
	})
	e.GET("/login", endpoints.HandleLoginForm)
	e.GET("/signup", endpoints.HandleSignup)
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
