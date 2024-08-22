package handlers

import (
	b64 "encoding/base64"
	"fmt"
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
)

type DeleteAccountHandler struct {
	userStore         store.UserStore
	sessionStore      store.SessionStore
	sessionCookieName string
}

type DeleteAccountHandlerParams struct {
	UserStore         store.UserStore
	SessionStore      store.SessionStore
	SessionCookieName string
}

func NewDeleteAccountHandler(params DeleteAccountHandlerParams) *DeleteAccountHandler {
	return &DeleteAccountHandler{
		userStore:         params.UserStore,
		sessionStore:      params.SessionStore,
		sessionCookieName: params.SessionCookieName,
	}
}

func (h *DeleteAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	currCookie, err := r.Cookie(h.sessionCookieName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.DeleteAccountError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not retrieve session cookie: %v\n", err)
		return
	}

	sessionBytes, err := b64.RawStdEncoding.DecodeString(currCookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.DeleteAccountError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not decode session cookie: %v\n", err)
		return
	}

	sessionID := string(sessionBytes)

	user, err := h.sessionStore.GetUserFromSession(sessionID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.DeleteAccountError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not retrieve user from session: %v\n", err)
		return
	}

	defer func() {
		err = h.userStore.DeleteUser(user.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("Could not delete user: %v\n", err)
			return
		}
	}()

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusPermanentRedirect)
	return
}
