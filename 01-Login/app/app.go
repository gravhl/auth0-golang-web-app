package app

import (
	"encoding/base64"
	"encoding/gob"
	"os"
	"path"
	"runtime"

	"fmt"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

var (
	Store             *sessions.FilesystemStore
	Log               *logrus.Logger
	Auth0Domain       string
	Auth0LocalDomain  string
	Auth0ClientID     string
	Auth0ClientSecret string
	Auth0CallbackURL  string
)

// CheckForEnv will check and fill in mandatory globals from env vars
func CheckForEnv() (err error) {
	envvar := os.Getenv("AUTH0_DOMAIN")
	Auth0Domain = "NOT SET"

	if len(envvar) > 0 {
		bytes, err := base64.StdEncoding.DecodeString(envvar)
		if err != nil {
			Log.Errorf("AUTH0_DOMAIN is not propely encoded base64. Bailing")
			return err
		} else {
			Auth0Domain = string(bytes)
		}
	} else {
		Log.Errorf("AUTH0_DOMAIN is not set.")
		err = fmt.Errorf("AUTH0_DOMAIN not set")
	}

	envvar = os.Getenv("AUTH0_LOCAL_DOMAIN")
	//	Auth0LocalDomain = ""

	if len(envvar) > 0 {
		bytes, err := base64.StdEncoding.DecodeString(envvar)
		if err != nil {
			Log.Errorf("AUTH0_DOMAIN is not propely encoded base64. Bailing")
			return err
		} else {
			Auth0LocalDomain = string(bytes)
			Log.Infof("ClientID: \"%s\"", Auth0LocalDomain)
		}
	} else {
		Log.Warnf("AUTH0_LOCAL_DOMAIN is not set.")
		//		err = fmt.Errorf("AUTH0_LOCAL_DOMAIN not set")
	}

	envvar = os.Getenv("AUTH0_CLIENT_ID")
	Auth0ClientID = "NOT SET"

	if len(envvar) < 1 {
		Log.Errorf("AUTH0_CLIENT_ID is not set.")
		err = fmt.Errorf("AUTH0_CLIENT_ID not set")
	} else {
		Log.Infof("ClientID: \"%s\"", envvar)
		Auth0ClientID = envvar
		// Auth0ClientID = base64.StdEncoding.EncodeToString([]byte(envvar))
		// Log.Infof("ClientID base64: \"%s\"", Auth0ClientID)
	}

	envvar = os.Getenv("AUTH0_CLIENT_SECRET")
	Auth0ClientSecret = "NOT SET"

	if len(envvar) > 0 {
		bytes, err := base64.StdEncoding.DecodeString(envvar)
		if err != nil {
			Log.Errorf("AUTH0_CLIENT_SECRET is not propely encoded base64.")
		} else {
			Auth0ClientSecret = string(bytes)
		}
	} else {
		Log.Errorf("AUTH0_CLIENT_SECRET is not set.")
		err = fmt.Errorf("AUTH0_CLIENT_SECRET not set")
	}

	envvar = os.Getenv("AUTH0_CALLBACK_URL")
	Auth0CallbackURL = "NOT SET"

	if len(envvar) > 0 {
		bytes, err := base64.StdEncoding.DecodeString(envvar)
		if err != nil {
			Log.Errorf("AUTH0_CALLBACK_URL is not propely encoded base64.")
		} else {
			Auth0CallbackURL = string(bytes)
		}
	} else {
		Log.Errorf("AUTH0_CALLBACK_URL is not set.")
		err = fmt.Errorf("AUTH0_CALLBACK_URL not set")
	}

	return
}

func Init() error {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Print(err.Error())
	// 	return err
	// }
	Log = logrus.New()
	Log.SetReportCaller(true)
	Log.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	err := CheckForEnv()
	if err != nil {
		Log.Warnf("NOTE: some environment vars not set correctly. App may fail.")
	}

	Store = sessions.NewFilesystemStore("", []byte("something-very-secret"))
	gob.Register(map[string]interface{}{})
	return nil
}
