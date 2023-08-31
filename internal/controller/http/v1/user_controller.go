package v1

import (
	"avito_test_case/internal/dto"
	my_errors "avito_test_case/internal/errors"
	"avito_test_case/internal/service"
	"avito_test_case/pkg/logger"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"time"
)

type userController struct {
	s *service.PostgresAssignmentService
	l logger.Logger
}

func NewUserController(handler *chi.Mux, s *service.PostgresAssignmentService, l logger.Logger) {
	c := &userController{s, l}
	_ = handler.Route("/user", func(r chi.Router) {
		r.Get("/", c.getAssignments)
		r.Get("/history", c.getHistory)
	})
}

// @Summary GetAssignments
// @Tags user
// @Description Returns current assignments of user
// @ID get-assignments
// @Accept json
// @Produce json
// @Param input body dto.User true "User ID"
// @Success 200 {object} dto.UserSegmentGet
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500
// @Router /user [get]
func (c *userController) getAssignments(w http.ResponseWriter, r *http.Request) {
	var u dto.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, http.StatusText(500), 500)
		c.l.Warn(err)
		return
	}
	if err := validator.New().Struct(u); err != nil {
		_ = json.NewEncoder(w).Encode(dto.NewErrorResponse(err))
		http.Error(w, http.StatusText(400), 400)
		return
	}
	var err error
	var res dto.UserSegmentGet
	if res, err = c.s.GetUserAssignments(context.TODO(), u.UserID); err != nil {
		var noSuchUserError *my_errors.NoSuchUserError
		if errors.As(err, &noSuchUserError) {
			_ = json.NewEncoder(w).Encode(dto.NewErrorResponse(err))
			http.Error(w, http.StatusText(400), 400)
		} else {
			http.Error(w, http.StatusText(500), 500)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, http.StatusText(500), 500)
		c.l.Warn(err)
		return
	}
}

// @Summary GetHistory
// @Tags user
// @Description Returns all assignment action history for this user
// @ID get-history
// @Accept json
// @Produce json
// @Param input body dto.UserHistoryGet true "User ID with year and month to get history of"
// @Success 200 {object} dto.UserHistory
// @Failure 500
// @Router /user/history [get]
func (c *userController) getHistory(w http.ResponseWriter, r *http.Request) {
	var u dto.UserHistoryGet
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, http.StatusText(500), 500)
		c.l.Warn(err)
		return
	}
	if err := validator.New().Struct(u); err != nil {
		_ = json.NewEncoder(w).Encode(dto.NewErrorResponse(err))
		http.Error(w, http.StatusText(400), 400)
		return
	}
	var err error
	var res dto.UserHistory
	start := time.Time{}.AddDate(int(u.Year)-1, int(u.Month)-1, 0)
	end := start.AddDate(0, 1, 0)
	if res, err = c.s.GetUserHistory(context.TODO(), u.UserID, start, end); err != nil {
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	csvW := csv.NewWriter(w)
	csvW.Comma = ';'
	for _, v := range res.Records {
		if err := csvW.Write([]string{strconv.FormatInt(v.UserID, 10), v.Segment, v.Operation, v.Timestamp.Format(time.RFC3339)}); err != nil {
			http.Error(w, http.StatusText(500), 500)
			c.l.Warn(err)
			return
		}
	}
	csvW.Flush()
}
