package handlers

import (
	b64 "encoding/base64"
	"fmt"
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
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
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf(
			"Could not extract user from session: %v\n",
			err,
		)
		return
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
