module 01-Login

go 1.12

require (
	app v0.0.0
	auth v0.0.0
	callback v0.0.0
	github.com/codegangsta/negroni v1.0.0
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	home v0.0.0
	login v0.0.0
	logout v0.0.0
	middlewares v0.0.0
	templates v0.0.0
	user v0.0.0
)

replace app => ./app

replace auth => ./auth

replace callback => ./routes/callback

replace home => ./routes/home

replace login => ./routes/login

replace logout => ./routes/logout

replace middlewares => ./routes/middlewares

replace user => ./routes/user

replace templates => ./routes/templates
