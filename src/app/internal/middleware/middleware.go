package middleware

import (
	"context"
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"goth/internal/store"
	"log"
	"net/http"
)

type key string

var NonceKey key = "nonces"

type Nonces struct {
	Htmx            string
	ResponseTargets string
	Tw              string
	HtmxCSSHash     string
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func CSPMiddleware(next http.Handler) http.Handler {
	// To use the same nonces in all responses, move the Nonces
	// struct creation to here, outside the handler.

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new Nonces struct for every request when here.
		// move to outside the handler to use the same nonces in all responses
		nonceSet := Nonces{
			Htmx:            generateRandomString(16),
			ResponseTargets: generateRandomString(16),
			Tw:              generateRandomString(16),
			HtmxCSSHash:     "sha256-pgn1TCGZX6O77zDvy0oTODMOxemn0oj0LeCnQTRj7Kg=",
		}

		// set nonces in context
		ctx := context.WithValue(r.Context(), NonceKey, nonceSet)
		// insert the nonces into the content security policy header
		cspHeader := fmt.Sprintf(
			"default-src 'self'; script-src 'nonce-%s' 'nonce-%s' ; style-src 'nonce-%s' '%s';",
			nonceSet.Htmx,
			nonceSet.ResponseTargets,
			nonceSet.Tw,
			nonceSet.HtmxCSSHash,
		)
		w.Header().Set("Content-Security-Policy", cspHeader)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TextHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// get the Nonce from the context, it is a struct called Nonces,
// so we can get the nonce we need by the key, i.e. HtmxNonce
func GetNonces(ctx context.Context) Nonces {
	nonceSet := ctx.Value(NonceKey)
	if nonceSet == nil {
		log.Fatal("error getting nonce set - is nil")
	}

	nonces, ok := nonceSet.(Nonces)

	if !ok {
		log.Fatal("error getting nonce set - not ok")
	}

	return nonces
}

func GetHtmxNonce(ctx context.Context) string {
	nonceSet := GetNonces(ctx)

	return nonceSet.Htmx
}

func GetResponseTargetsNonce(ctx context.Context) string {
	nonceSet := GetNonces(ctx)
	return nonceSet.ResponseTargets
}

func GetTwNonce(ctx context.Context) string {
	nonceSet := GetNonces(ctx)
	return nonceSet.Tw
}

type AuthMiddleware struct {
	userStore         store.UserStore
	sessionStore      store.SessionStore
	sessionCookieName string
}

func NewAuthMiddleware(
	userStore store.UserStore,
	sessionStore store.SessionStore,
	sessionCookieName string,
) *AuthMiddleware {
	return &AuthMiddleware{
		userStore:         userStore,
		sessionStore:      sessionStore,
		sessionCookieName: sessionCookieName,
	}
}

type UserContextKey string

var UserKey UserContextKey = "user"

func (m *AuthMiddleware) AddUserToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sessionCookie, err := r.Cookie(m.sessionCookieName)

		if err != nil {
			fmt.Println("error getting session cookie", err)
			next.ServeHTTP(w, r)
			return
		}

		decodedValue, err := b64.StdEncoding.DecodeString(sessionCookie.Value)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		sessionID := string(decodedValue)

		fmt.Println("sessionID", sessionID)

		user, err := m.sessionStore.GetUserFromSession(sessionID)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, user)

		err = m.userStore.SetIsActive(user.ID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(ctx context.Context) *store.User {
	user := ctx.Value(UserKey)
	if user == nil {
		return nil
	}

	return user.(*store.User)
}
