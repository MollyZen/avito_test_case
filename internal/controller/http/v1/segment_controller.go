package v1

import (
	"avito_test_case/internal/datastruct"
	"avito_test_case/internal/dto"
	"avito_test_case/internal/service"
	"avito_test_case/pkg/logger"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type segmentController struct {
	s *service.SegmentService
	l logger.Logger
}

func NewSegmentController(handler *chi.Mux, s *service.SegmentService, l logger.Logger) {
	c := &segmentController{s, l}
	_ = handler.Route("/segment", func(r chi.Router) {
		r.Put("/", c.Create)
		r.Delete("/", c.Delete)
	})
}

func (s *segmentController) Create(w http.ResponseWriter, r *http.Request) {
	var seg dto.Segment
	if err := json.NewDecoder(r.Body).Decode(&seg); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	if err := validator.New().Struct(seg); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	if _, err := s.s.Create(context.TODO(), datastruct.Segment{
		Slug: seg.Slug,
	}); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
}

func (s *segmentController) Delete(w http.ResponseWriter, r *http.Request) {
	var seg dto.SegmentToDelete
	if err := json.NewDecoder(r.Body).Decode(&seg); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	if err := validator.New().Struct(seg); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	if _, err := s.s.Delete(context.TODO(), datastruct.Segment{
		Slug: seg.Slug,
	}); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
}
