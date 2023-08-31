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
	"time"
)

type assignmentController struct {
	s *service.PostgresAssignmentService
	l logger.Logger
}

func NewAssignmentController(handler *chi.Mux, s *service.PostgresAssignmentService, l logger.Logger) {
	c := &assignmentController{s, l}
	_ = handler.Route("/assignment", func(r chi.Router) {
		r.Put("/", c.ChangeAssignment)
	})
}

// @Summary ChangeAssignment
// @Tags assignment
// @Description Remove and/or add segments to user. Deletes happen first. In case of conflict with existing data new values overwrite old ones
// @ID change-assignment
// @Accept json
// @Produce json
// @Param input body dto.UserSegmentChange true "User ID, slugs of segments to add (TTL is optional) and delete. DateTime in RFC3339 format"
// @Success 200
// @Router /assignment [put]
func (c *assignmentController) ChangeAssignment(w http.ResponseWriter, r *http.Request) {
	var ch dto.UserSegmentChange
	if err := json.NewDecoder(r.Body).Decode(&ch); err != nil {
		c.l.Warn(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if err := validator.New().Struct(ch); err != nil {
		c.l.Warn(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}
	for _, v := range ch.SegmentAdd {
		if len(v.UntilDate) > 0 {
			_, err := time.Parse(time.RFC3339, v.UntilDate)
			if err != nil {
				c.l.Warn(err)
				http.Error(w, http.StatusText(400), 400)
				return
			}
		}
	}
	if err := c.s.Assign(context.TODO(), ch.UserID, ch.SegmentAdd, ch.SegmentRemove); err != nil {
		c.l.Warn(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}
}
