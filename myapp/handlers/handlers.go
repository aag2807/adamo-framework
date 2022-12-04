package handlers

import (
	"net/http"

	"github.com/aag2807/adamo-framework"
)

type Handlers struct {
	App *adamo.Adamo
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "home", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}
