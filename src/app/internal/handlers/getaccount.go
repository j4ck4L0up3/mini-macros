package handlers

import (
	b64 "encoding/base64"
	"fmt"
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
	"strings"
)

type GetAccountHandler struct {
	sessionStore      store.SessionStore
	sessionCookieName string
}

type GetAccountHandlerParams struct {
	SessionStore      store.SessionStore
	SessionCookieName string
}

func NewGetAccountHandler(params GetAccountHandlerParams) *GetAccountHandler {
	return &GetAccountHandler{
		sessionStore:      params.SessionStore,
		sessionCookieName: params.SessionCookieName,
	}
}

func (h *GetAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	currCookie, err := r.Cookie(h.sessionCookieName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	valueBytes, err := b64.StdEncoding.DecodeString(currCookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Could not decode cookie value: %v , error: %v\n", currCookie.Value, err)
		return
	}
	splitVals := strings.Split(string(valueBytes), ":")

	sessionID := splitVals[0]
	userID := splitVals[1]

	user, err := h.sessionStore.GetUserFromSession(sessionID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Could not retrieve user from session: %v\n", err)
		return
	}

	c := templates.Account(user.Email, user.FirstName, user.LastName)
	err = templates.Layout(c, "Mini Macros").Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error rendering page template: %v\n", err)
	}
}
