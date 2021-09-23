package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"callback"
	"home"
	"login"
	"logout"
	"middlewares"
	"user"

	"app"
)

var r *mux.Router

func setupRouter() {
	if r == nil {
		r = mux.NewRouter()

		r.HandleFunc("/", home.HomeHandler)
		r.HandleFunc("/login", login.LoginHandler)
		r.HandleFunc("/logout", logout.LogoutHandler)
		r.HandleFunc("/callback", callback.CallbackHandler)
		r.Handle("/user", negroni.New(
			negroni.HandlerFunc(middlewares.IsAuthenticated),
			negroni.Wrap(http.HandlerFunc(user.UserHandler)),
		))
		r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))

		http.Handle("/", r)
	}
}

func startServer(errch chan error) *http.Server {
	srv := &http.Server{
		Addr:         ":3000",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Print("HTTP Server listening on :3000")

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errch <- err
		}
	}()

	return srv
}

func startHttpsServer(errch chan error) *http.Server {

	if len(app.Auth0LocalDomain) > 0 {

		// m := &autocert.Manager{
		// 	Cache:      autocert.DirCache("secret-dir"),
		// 	Prompt:     autocert.AcceptTOS,
		// 	HostPolicy: autocert.HostWhitelist(app.Auth0LocalDomain),
		// }
		// srv := &http.Server{
		// 	Addr:         ":4443",
		// 	TLSConfig:    m.TLSConfig(),
		// 	Handler:      r,
		// 	ReadTimeout:  5 * time.Second,
		// 	WriteTimeout: 10 * time.Second,
		// 	IdleTimeout:  120 * time.Second,
		// }

		// log.Print("HTTPS Server listening on :4443")

		// go func() {
		// 	if err := srv.ListenAndServeTLS("", ""); err != nil {
		// 		errch <- err
		// 	}
		// }()

		srv := &http.Server{
			Addr: ":4443",
			// TLSConfig:    m.TLSConfig(),
			Handler:      r,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}

		log.Print("HTTPS Server listening on :4443")

		go func() {
			if err := srv.ListenAndServeTLS("server.crt", "server.key"); err != nil {
				errch <- err
			}
		}()

		return srv
	} else {
		return nil
	}

}

func RunServers() {
	setupRouter()

	errs := make(chan error)

	httpSrv := startServer(errs)
	httpsSrv := startHttpsServer(errs)

	// Wait for an interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	app.Log.Infof("Servers up")
sigloop:
	for true {
		select {
		case sig := <-c:
			app.Log.Warnf("sig INT: %+v", sig)
			break sigloop
		case err := <-errs:
			app.Log.Errorf("err -> %s", err.Error())
		}

	}

	// Attempt a graceful shutdown
	app.Log.Infof("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	httpSrv.Shutdown(ctx)
	if httpsSrv != nil {
		httpsSrv.Shutdown(ctx)
	}
	app.Log.Infof("Shutdown")
}
