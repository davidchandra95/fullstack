package controllers

import (
	"github.com/davidchandra95/fullstack/api/responses"
	"net/http"
)

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome..")
}