package shared

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	// Store is the cookie store
	Store *sessions.CookieStore
	// Name stores the session name
	Name string
)

// Session stores session level information
type Session struct {
	Options sessions.Options `json:"options"` // http://www.gorillatoolkit.org/pkg/sessions#Options
	Name    string           `json:"name"`    // http://www.gorillatoolkit.org/pkg/sessions#CookieStore.Get
	// used to authenticate cookie using HMAC
	HashKey []byte `json:"hashkey"` // http://www.gorillatoolkit.org/pkg/securecookie#New
	// used to encrypt cookies
	BlockKey []byte `json:"blockkey"` // http://www.gorillatoolkit.org/pkg/securecookie#New
}

// Configure configures the session cookie store
func Configure(s Session) {

	// Store = sessions.NewCookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	Store = sessions.NewCookieStore(s.HashKey, s.BlockKey)
	Store.Options = &s.Options
	Name = s.Name

}

// NewSession returns a new session, never returns an error
func NewSession(r *http.Request) *sessions.Session {
	session, err := Store.Get(r, Name)
	if err != nil {
		log.Println(err)
	}

	return session
}

// Empty deletes all the current session values
func Empty(sess *sessions.Session) {

	// Clear out all stored values in the cookie
	for k := range sess.Values {
		delete(sess.Values, k)
	}

}
