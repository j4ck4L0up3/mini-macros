package handlers

import (
	b64 "encoding/base64"
	"fmt"
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
	"strings"
)

type GetDashboardHandler struct {
	sessionCookieName string
	sessionStore      store.SessionStore
}

type GetDashboardHandlerParams struct {
	SessionCookieName string
	SessionStore      store.SessionStore
}

func NewGetDashboardHandler(params GetDashboardHandlerParams) *GetDashboardHandler {
	return &GetDashboardHandler{
		sessionCookieName: params.SessionCookieName,
		sessionStore:      params.SessionStore,
	}
}

func (h *GetDashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	currCookie, err := r.Cookie(h.sessionCookieName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
		fmt.Printf(
			"Could not extract user from session: %v\nsessionID: %v, userID: %v, user: %v",
			err,
			sessionID,
			userID,
			user,
		)
	}

	// TODO: pass in macros from user ID
	c := templates.Dashboard(user.FirstName)
	err = templates.Layout(c, "Mini Macros").Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error rendering dashboard: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
