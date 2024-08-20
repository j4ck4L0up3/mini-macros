package handlers

import (
	b64 "encoding/base64"
	"fmt"
	"goth/internal/store"
	"net/http"
	"time"
)

type PostLogoutHandler struct {
	sessionCookieName string
	sessionStore      store.SessionStore
	userStore         store.UserStore
}

type PostLogoutHandlerParams struct {
	SessionCookieName string
	SessionStore      store.SessionStore
	UserStore         store.UserStore
}

func NewPostLogoutHandler(params PostLogoutHandlerParams) *PostLogoutHandler {
	return &PostLogoutHandler{
		sessionCookieName: params.SessionCookieName,
		sessionStore:      params.SessionStore,
		userStore:         params.UserStore,
	}
}

func (h *PostLogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	currCookie, err := r.Cookie(h.sessionCookieName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Cookie not found: %v", err)
		return
	}

	sessionBytes, err := b64.StdEncoding.DecodeString(currCookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Could not decode cookie value: %v , error: %v\n", currCookie.Value, err)
		return
	}
	sessionID := string(sessionBytes)
	user, err := h.sessionStore.GetUserFromSession(sessionID)

	err = h.userStore.SetInactive(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("User could not be set inactive: %v\n", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    h.sessionCookieName,
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
		Path:    "/",
	})

	err = h.sessionStore.DeleteSession(sessionID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error deleting session: %v\n", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
