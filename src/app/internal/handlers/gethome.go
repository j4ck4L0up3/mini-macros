package handlers

import (
	"fmt"
	"goth/internal/middleware"
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(middleware.UserKey).(*store.User)

	if !ok {
		c := templates.GuestIndex()

		err := templates.Layout(c, "Mini Macros").Render(r.Context(), w)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("Error rendering template: %v\n", err)
			return
		}

		return
	}

	c := templates.Index(user.FirstName)
	err := templates.Layout(c, "Mini Macros").Render(r.Context(), w)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error rendering template: %v\n", err)
		return
	}
}
