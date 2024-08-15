package handlers

import (
	b64 "encoding/base64"
	"fmt"
	"goth/internal/store"
	"net/http"
	"strings"
	"time"
)

type PostLogoutHandler struct {
	sessionCookieName string
	sessionStore      store.SessionStore
}

type PostLogoutHandlerParams struct {
	SessionCookieName string
	SessionStore      store.SessionStore
}

func NewPostLogoutHandler(params PostLogoutHandlerParams) *PostLogoutHandler {
	return &PostLogoutHandler{
		sessionCookieName: params.SessionCookieName,
		sessionStore:      params.SessionStore,
	}
}

func (h *PostLogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	currCookie, cookieErr := r.Cookie(h.sessionCookieName)
	if cookieErr != nil {
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		fmt.Errorf("Cookie not found: %v", cookieErr)
		return
	}

	valueStr := fmt.Sprint(b64.StdEncoding.DecodeString(currCookie.Value))
	splitVals := strings.Split(valueStr, ":")

	sessionID := splitVals[0]
	userID := splitVals[1]

	user, userErr := h.sessionStore.GetUserFromSession(sessionID, userID)
	if userErr != nil {
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		fmt.Errorf(
			"User could not be retrieved from session: %v, sessionID: %v, userID: %v",
			userErr,
			sessionID,
			userID,
		)
		return
	}

	user.Active = false

	http.SetCookie(w, &http.Cookie{
		Name:    h.sessionCookieName,
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
