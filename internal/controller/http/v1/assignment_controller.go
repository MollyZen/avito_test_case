package v1

import (
	"avito_test_case/internal/dto"
	"avito_test_case/internal/service"
	"avito_test_case/pkg/logger"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type assignmentController struct {
	s *service.PostgresAssignmentService
	l logger.Logger
}

func NewAssignmentController(handler *chi.Mux, s *service.PostgresAssignmentService, l logger.Logger) {
	c := &assignmentController{s, l}
	_ = handler.Route("/assignment", func(r chi.Router) {
		r.Put("/", c.ChangeAssignment)
		r.Get("/history", c.GetHistory)
	})
}

func (c *assignmentController) ChangeAssignment(w http.ResponseWriter, r *http.Request) {
	var ch dto.UserSegmentChange
	if err := json.NewDecoder(r.Body).Decode(&ch); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if err := validator.New().Struct(ch); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	if err := c.s.Assign(context.TODO(), ch.UserID, ch.SegmentAdd, ch.SegmentRemove); err != nil {
		return
	}
}

func (c *assignmentController) GetHistory(w http.ResponseWriter, r *http.Request) {

}
