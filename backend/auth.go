package main

import (
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func loginHandler(e *core.RequestEvent, app *pocketbase.PocketBase) error {
	//get credentials from request
	email := e.Request.FormValue("email")
	password := e.Request.FormValue("password")

	user, err := app.FindAuthRecordByEmail("users", email)
	if err != nil {
		return e.HTML(200, "User does not exist")
	}

	//Check if credentials are correct and return token
	isAuthenticated := user.ValidatePassword(password)
	if isAuthenticated {
		token, err := user.NewAuthToken()
		if err != nil {
			return e.HTML(200, "Error generating token, please try again.")
		}
		token_cookie := http.Cookie{
			Name:     "explore_token",
			Value:    token,
			Secure:   true,
			HttpOnly: true,
			Expires:  time.Now().Add(30 * time.Minute),
		}
		e.SetCookie(&token_cookie)
		id_cookie := http.Cookie{
			Name:     "id",
			Value:    user.Id,
			Secure:   true,
			HttpOnly: true,
			Expires:  time.Now().Add(30 * time.Minute),
		}
		e.SetCookie(&id_cookie)
		e.Response.Header().Set("HX-Redirect", "/")
	}

	loginError := "Login Failed, Please try again"
	return e.HTML(200, loginError)
}
