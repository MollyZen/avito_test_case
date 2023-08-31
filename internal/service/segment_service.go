package service

import (
	"avito_test_case/internal/datastruct"
	"avito_test_case/internal/misc"
	"avito_test_case/internal/repository"
	"avito_test_case/pkg/logger"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (s *SegmentService) Create(ctx context.Context, seg datastruct.Segment) (datastruct.Segment, error) {
	var res datastruct.Segment
	var err error
	if res, err = s.segRep.Create(ctx, seg); err != nil {
		return datastruct.Segment{}, err
	}
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
