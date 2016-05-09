package safecookie

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/securecookie"
	"github.com/spf13/viper"
)

var safeCookie *securecookie.SecureCookie

func init() {
	viper.BindEnv("googleClientID", "USERAUTH_GOOGLE_CLIENT_ID")
	viper.BindEnv("secureCookieHashKey", "USERAUTH_SECURE_COOKIE_HASH_KEY")
	viper.BindEnv("secureCookieBlockKey", "USERAUTH_SECURE_COOKIE_BLOCK_KEY")

	hashKey := viper.GetString("secureCookieHashKey")
	if hashKey == "" {
		log.Fatal("Secure cookie hash key has not been configured. User authentication is not possible.")
	}

	blockKey := viper.GetString("secureCookieBlockKey")
	if blockKey == "" {
		log.Fatal("Secure cookie block key has not been configured. User authentication is unsafe.")
	}

	safeCookie = securecookie.New([]byte(hashKey), []byte(blockKey))
}

/*
Get returns an initialized instance of gorilla mux's securecookie.
This makes it easier to share the same secure cookie instance.
*/
func Get() *securecookie.SecureCookie {
	return safeCookie
}
