package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/jeffbmartinez/userauth/handler"
)

const (
	userauthServiceHostEnv = "USERAUTH_SERVICE_HOST"
	userauthServicePortEnv = "USERAUTH_SERVICE_PORT"
)

func main() {
	userauthServiceHost := getUserauthListenHost()
	userauthServicePort := getUserauthListenPort()

	r := mux.NewRouter()

	r.HandleFunc("/verify/token/google", handler.VerifyGoogleIDToken)

	http.Handle("/", r)

	listenDomain := fmt.Sprintf("%s:%s", userauthServiceHost, userauthServicePort)
	log.Printf("Listening for connections on %s.\nConfigure by setting %s and %s environment variables.",
		listenDomain, userauthServiceHostEnv, userauthServicePortEnv)

	err := http.ListenAndServe(listenDomain, nil)
	log.Fatalf("Problem running server (%v)\n", err)
}

func getUserauthListenHost() string {
	return os.Getenv(userauthServiceHostEnv)
}

func getUserauthListenPort() string {
	userauthServicePort := os.Getenv(userauthServicePortEnv)
	if userauthServicePort == "" {
		log.Fatalf("Environment variable (%v) required and not found\n", userauthServicePortEnv)
	}

	return userauthServicePort
}
