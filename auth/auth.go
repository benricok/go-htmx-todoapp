package auth

import (
	"go-htmx/user"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	accessTokenCookieName = "access-token"
	// TODO: Get from env variable
	jwtSecretKey          = "very-secret-keyy"
)

func GetJWTSecret() string {
	return jwtSecretKey
}

func GetAccessTokenCookieName() string {
	return accessTokenCookieName
}

// will be encoed to a JWT
type Claims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims // embedded type
}

func Claim(c echo.Context) jwt.Claims {
	return &Claims{}
}

func GenerateTokensAndSetCookies(user *user.User, c echo.Context) error {
	accessToken, exp, err := GenerateAccessToken(user)

	if err != nil {
		return err
	}

	setTokenCookie(accessTokenCookieName, accessToken, exp, c)
	setUserCookie(user, exp, c)

	return nil
}

func GenerateAccessToken(user *user.User) (string, time.Time, error) {
	// Declare expiration time of token
	// TODO: add to config toml file
	// Default 1hr
	expirationTime := time.Now().Add(1 * time.Hour)
	
	secret := []byte(GetJWTSecret())

	// Create JWT Claims
	// Include username and expiry time
	claims := &Claims {
		Name: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt:  jwt.NewNumericDate(expirationTime),
			IssuedAt:   jwt.NewNumericDate(time.Now()),
			NotBefore:  jwt.NewNumericDate(time.Now()),
			Issuer: 	"todo-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenStr, expirationTime, nil
}

func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

func setUserCookie(user *user.User, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.Username
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}

// JWTErrorChecker will be executed when user try to access a protected path.
func JWTErrorChecker(c echo.Context, err error) error {
    // Redirects to the signIn form.
	return c.Redirect(http.StatusMovedPermanently, "/login")
}