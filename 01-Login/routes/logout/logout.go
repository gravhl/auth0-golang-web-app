package logout

import (
	"net/http"
	"net/url"

	"app"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	domain := app.Auth0Domain

	logoutUrl, err := url.Parse("https://" + domain)

	if err != nil {
		app.Log.Errorf("url.Parse: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logoutUrl.Path += "/v2/logout"
	parameters := url.Values{}

	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + r.Host)
	if err != nil {
		app.Log.Errorf("url.Parse: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", app.Auth0ClientID)
	logoutUrl.RawQuery = parameters.Encode()

	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}
