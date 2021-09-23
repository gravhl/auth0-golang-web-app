package login

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"app"
	"auth"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate random state
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		// if strings.HasPrefix(err.Error(), "no such file") {
		// 	session, err := app.Store.New()
		// } else {
		if session != nil && !session.IsNew {
			if session == nil {
				app.Log.Errorf("Session is nil")
			}
			app.Log.Errorf("Store: %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			app.Log.Infof("Store: New Session")
		}
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		app.Log.Errorf("Session: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authenticator, err := auth.NewAuthenticator()
	if err != nil {
		app.Log.Errorf("Auth: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}
