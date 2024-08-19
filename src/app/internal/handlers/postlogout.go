package handlers

import (
	b64 "encoding/base64"
	"fmt"
	"goth/internal/store"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PostLogoutHandler struct {
	sessionCookieName string
	userStore         store.UserStore
}

type PostLogoutHandlerParams struct {
	SessionCookieName string
	UserStore         store.UserStore
}

func NewPostLogoutHandler(params PostLogoutHandlerParams) *PostLogoutHandler {
	return &PostLogoutHandler{
		sessionCookieName: params.SessionCookieName,
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

	valueBytes, err := b64.StdEncoding.DecodeString(currCookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Could not decode cookie value: %v , error: %v", currCookie.Value, err)
		return
	}
	splitVals := strings.Split(string(valueBytes), ":")

	userIDStr := splitVals[1]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Could not convert userID into integer: %v", err)
		return
	}

	err = h.userStore.SetInactive(uint(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("User could not be set inactive: %v", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    h.sessionCookieName,
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
