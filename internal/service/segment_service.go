package service

import (
	"avito_test_case/internal/datastruct"
	"avito_test_case/internal/dto"
	"avito_test_case/internal/misc"
	"avito_test_case/internal/repository"
	"avito_test_case/pkg/logger"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type SegmentService struct {
	db     *pgxpool.Pool
	segRep *repository.PostgresSegmentRepository
	asRep  *repository.PostgresAssignmentRepository
	hisRep *repository.PostgresHistoryRepository
	l      logger.Logger
}

func NewPostgresSegmentService(db *pgxpool.Pool, segRep *repository.PostgresSegmentRepository, asRep *repository.PostgresAssignmentRepository, hisRep *repository.PostgresHistoryRepository, l logger.Logger) *SegmentService {
	return &SegmentService{
		db:     db,
		segRep: segRep,
		asRep:  asRep,
		hisRep: hisRep,
		l:      l,
	}
}

func (s *SegmentService) Create(ctx context.Context, seg dto.Segment) (datastruct.Segment, error) {
	var res datastruct.Segment
	var err error
	var tr pgx.Tx
	tr, err = s.db.BeginTx(context.TODO(), pgx.TxOptions{})
	isactive := false
	var tmp []datastruct.Segment
	if tmp, err = s.segRep.GetAllBySlug(context.TODO(), []string{seg.Slug}); err != nil {
		return datastruct.Segment{}, err
	}
	if len(tmp) > 0 {
		isactive = tmp[0].IsActive
	}

	if res, err = s.segRep.UpsertWithConn(ctx, datastruct.Segment{
		Slug: seg.Slug,
	}, tr.Conn()); err != nil {
		_ = tr.Rollback(context.TODO())
		return datastruct.Segment{}, err
	}

	if seg.Percent > 0. && seg.Percent <= 100. && !isactive {
		var tmp time.Time
		var t pgtype.Timestamptz
		if len(seg.UntilDate) > 0 {
			tmp, _ = time.Parse(time.RFC3339, seg.UntilDate)
			t = pgtype.Timestamptz{Time: tmp}
			t.Valid = true
		} else {
			t = pgtype.Timestamptz{}
		}
		var as []datastruct.Assignment
		//note: returns only user ids
		as, err = s.segRep.AddToPercentOfUsersWithConn(context.TODO(), res.ID, seg.Percent, t, tr.Conn())
		if err != nil {
			_ = tr.Rollback(context.TODO())
			return datastruct.Segment{}, err
		}
		for i, _ := range as {
			as[i].SegmentID = res.ID
			as[i].UntilDate = t
		}
		err = s.hisRep.CreateAllWithConn(context.TODO(), misc.AssignmentsToHistory(as, datastruct.OpAddedRand), tr.Conn())
		if err != nil {
			_ = tr.Rollback(context.TODO())
			return datastruct.Segment{}, err
		}
	}

	_ = tr.Commit(context.TODO())
	return res, nil
}

func (s *SegmentService) Delete(ctx context.Context, seg datastruct.Segment) (datastruct.Segment, error) {
	var tr pgx.Tx
	var err error
	tr, err = s.db.BeginTx(context.TODO(), pgx.TxOptions{})

	var delSeg datastruct.Segment
	if delSeg, err = s.segRep.DeleteBySlugWithConn(ctx, seg.Slug, tr.Conn()); err != nil {
		_ = tr.Rollback(context.TODO())
		return datastruct.Segment{}, err
	}

	var delAssignments []datastruct.Assignment
	if delAssignments, err = s.asRep.DeleteAllForSegWithConn(ctx, delSeg.ID, tr.Conn()); err != nil {
		_ = tr.Rollback(context.TODO())
		return datastruct.Segment{}, err
	}

	if err = s.hisRep.CreateAllWithConn(ctx, misc.AssignmentsToHistory(delAssignments, datastruct.OpSegDeleted), tr.Conn()); err != nil {
		_ = tr.Rollback(context.TODO())
		return datastruct.Segment{}, err
	}

	err = tr.Commit(context.TODO())
	if err != nil {
		return datastruct.Segment{}, err
	}
	return delSeg, nil
}
