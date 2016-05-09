package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/jeffbmartinez/userauth/handler"
)

func init() {
	viper.BindEnv("serviceHost", "USERAUTH_SERVICE_HOST")
	viper.BindEnv("servicePort", "USERAUTH_SERVICE_PORT")

	viper.SetDefault("serviceHost", "localhost")

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	userauthServiceHost := getUserauthListenHost()
	userauthServicePort := getUserauthListenPort()

	r := mux.NewRouter()

	r.HandleFunc("/login/google", handler.LoginGoogle)
	r.HandleFunc("/logout", handler.Logout)
	r.HandleFunc("/verify/session", handler.VerifySession)

	http.Handle("/", r)

	listenDomain := fmt.Sprintf("%s:%s", userauthServiceHost, userauthServicePort)
	log.WithFields(log.Fields{
		"host": userauthServiceHost,
		"port": userauthServicePort,
	}).Info("userauth service is starting")

	err := http.ListenAndServe(listenDomain, nil)
	log.WithError(err).Fatal("Problem running server")
}

func getUserauthListenHost() string {
	return viper.GetString("serviceHost")
}

func getUserauthListenPort() string {
	userauthServicePort := viper.GetString("servicePort")
	if userauthServicePort == "" {
		log.Fatal("Listen port configuration is not set. It is required.")
	}

	return userauthServicePort
}
