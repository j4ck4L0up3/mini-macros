package handlers

import (
	b64 "encoding/base64"
	"fmt"
	"goth/internal/hash"
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
	"time"
)

type PostLoginHandler struct {
	userStore         store.UserStore
	sessionStore      store.SessionStore
	passwordhash      hash.PasswordHash
	sessionCookieName string
}

type PostLoginHandlerParams struct {
	UserStore         store.UserStore
	SessionStore      store.SessionStore
	PasswordHash      hash.PasswordHash
	SessionCookieName string
}

func NewPostLoginHandler(params PostLoginHandlerParams) *PostLoginHandler {
	return &PostLoginHandler{
		userStore:         params.UserStore,
		sessionStore:      params.SessionStore,
		passwordhash:      params.PasswordHash,
		sessionCookieName: params.SessionCookieName,
	}
}

func (h *PostLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.userStore.GetUser(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.LoginError()
		c.Render(r.Context(), w)
		return
	}

	lockoutDiff := user.LockoutDuration.Sub(time.Now())
	noTime, err := time.ParseDuration("0m")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if lockoutDiff > noTime && user.LockedOut == true {
		w.WriteHeader(http.StatusLocked)
		c := templates.LockoutError()
		c.Render(r.Context(), w)
		return
	}

	passwordIsValid, err := h.passwordhash.ComparePasswordAndHash(password, user.Password)

	if err != nil || !passwordIsValid {
		h.userStore.IncrementLoginAttempts(user)
		w.WriteHeader(http.StatusUnauthorized)
		c := templates.LoginError()
		c.Render(r.Context(), w)
		return
	}

	h.userStore.ResetLoginAttempts(user)
	session, err := h.sessionStore.CreateSession(&store.Session{
		UserID: user.ID,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.userStore.SetIsActive(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("User could not be set active: %v\n", err)
		return
	}

	cookieValue := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s", session.SessionID)))
	expiration := time.Now().Add(3 * time.Hour)
	cookie := http.Cookie{
		Name:     h.sessionCookieName,
		Value:    cookieValue,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)

	w.Header().Set("HX-Redirect", "/dashboard")
	w.WriteHeader(http.StatusAccepted)
}
