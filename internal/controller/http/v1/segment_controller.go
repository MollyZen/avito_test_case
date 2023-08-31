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
	"time"
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

// @Summary CreateNewSegment
// @Tags segment
// @Description Creates new segment. If it already exists - sets active state to true
// @ID create-segment
// @Accept json
// @Produce json
// @Param input body dto.Segment true "Segment slug with percent of people to assign this segment to"
// @Success 200
// @Router /segment [put]
func (s *segmentController) Create(w http.ResponseWriter, r *http.Request) {
	var seg dto.Segment
	if err := json.NewDecoder(r.Body).Decode(&seg); err != nil {
		http.Error(w, http.StatusText(500), 500)
		s.l.Warn(err)
		return
	}
	if err := validator.New().Struct(seg); err != nil {
		http.Error(w, http.StatusText(400), 400)
		s.l.Warn(err)
		return
	}
	if len(seg.UntilDate) > 0 {
		_, err := time.Parse(time.RFC3339, seg.UntilDate)
		if err != nil {
			s.l.Warn(err)
			http.Error(w, http.StatusText(400), 400)
			return
		}
	}
	if _, err := s.s.Create(context.TODO(), seg); err != nil {
		http.Error(w, http.StatusText(500), 500)
		s.l.Warn(err)
		return
	}
}

// @Summary DeleteSegment
// @Tags segment
// @Description Deletes a segment. All existing assignments with this segment will be deleted with it
// @ID delete-segment
// @Accept json
// @Produce json
// @Param input body dto.SegmentToDelete true "Segment slug"
// @Success 200
// @Router /segment [delete]
func (s *segmentController) Delete(w http.ResponseWriter, r *http.Request) {
	var seg dto.SegmentToDelete
	if err := json.NewDecoder(r.Body).Decode(&seg); err != nil {
		http.Error(w, http.StatusText(500), 500)
		s.l.Warn(err)
		return
	}
	if err := validator.New().Struct(seg); err != nil {
		http.Error(w, http.StatusText(400), 400)
		s.l.Warn(err)
		return
	}
	if _, err := s.s.Delete(context.TODO(), datastruct.Segment{
		Slug: seg.Slug,
	}); err != nil {
		http.Error(w, http.StatusText(400), 400)
		s.l.Warn(err)
		return
	}
}
