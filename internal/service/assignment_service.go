package service

import (
	"avito_test_case/internal/datastruct"
	"avito_test_case/internal/dto"
	"avito_test_case/internal/errors"
	"avito_test_case/internal/misc"
	"avito_test_case/internal/repository"
	"avito_test_case/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PostgresAssignmentService struct {
	db      *pgxpool.Pool
	userRep *repository.PostgresUserRepository
	segRep  *repository.PostgresSegmentRepository
	asRep   *repository.PostgresAssignmentRepository
	hisRep  *repository.PostgresHistoryRepository
	l       logger.Logger
}

func NewPostgresAssignmentService(db *pgxpool.Pool, userRep *repository.PostgresUserRepository, segRep *repository.PostgresSegmentRepository,
	asRepo *repository.PostgresAssignmentRepository, hisRep *repository.PostgresHistoryRepository,
	l logger.Logger) *PostgresAssignmentService {
	return &PostgresAssignmentService{
		db:      db,
		userRep: userRep,
		segRep:  segRep,
		asRep:   asRepo,
		hisRep:  hisRep,
		l:       l,
	}
}

func (uc PostgresAssignmentService) Assign(ctx context.Context, userID int64, segToAdd []dto.SegmentToAdd, segToDelete []string) error {
	//starting transaction
	var tr pgx.Tx
	var err error
	tr, err = uc.db.BeginTx(context.TODO(), pgx.TxOptions{})

	//adding user if they don't exist
	var user datastruct.User
	user, err = uc.userRep.UpsertWithConn(ctx, datastruct.User{
		ID: userID,
	}, tr.Conn())
	if err != nil {
		_ = tr.Rollback(context.TODO())
		return err
	}

	//getting ids for segment slugs + checking if they exist
	var segs []datastruct.Segment
	segToAddNames := make([]string, len(segToAdd))
	for i, v := range segToAdd {
		segToAddNames[i] = v.Slug
	}
	toAddAndDelete := append(segToAddNames, segToDelete...)
	segs, err = uc.segRep.GetAllBySlugWithConn(ctx, toAddAndDelete, tr.Conn())
	if err != nil {
		_ = tr.Rollback(context.TODO())
		return err
	}
	segNameMap := make(map[string]datastruct.Segment, len(segs))
	for _, v := range segs {
		segNameMap[v.Slug] = v
	}
	segToAddT := make([]datastruct.Segment, len(segToAdd))
	segToDeleteT := make([]datastruct.Segment, len(segToDelete))
	for i, v := range segToAdd {
		mv, ok := segNameMap[v.Slug]
		if !ok {
			_ = tr.Rollback(context.TODO())
			return fmt.Errorf("segment with slug %s doesn't exist", v)
		}
		segToAddT[i] = mv
	}
	for i, v := range segToDelete {
		mv, ok := segNameMap[v]
		if !ok {
			_ = tr.Rollback(context.TODO())
			return fmt.Errorf("segment with slug %s doesn't exist", v)
		}
		segToDeleteT[i] = mv
	}

	//deleting assignments
	if len(segToDeleteT) > 0 {
		toDeleteAs := make([]datastruct.Assignment, len(segToDeleteT))
		for i, seg := range segToDeleteT {
			toDeleteAs[i].UserID = user.ID
			toDeleteAs[i].SegmentID = seg.ID
		}
		var deletedAs []datastruct.Assignment
		deletedAs, err = uc.asRep.DeleteWithConn(ctx, toDeleteAs, tr.Conn())
		err = uc.hisRep.CreateAll(ctx, misc.AssignmentsToHistory(deletedAs, datastruct.OpRemoved))
		if err != nil {
			_ = tr.Rollback(context.TODO())
			return err
		}
	}

	//getting remaining assignments
	var remaining []datastruct.Assignment
	remaining, err = uc.asRep.GetAllForUserWithConn(ctx, user.ID, tr.Conn())
	if err != nil {
		_ = tr.Rollback(context.TODO())
		return err
	}

	//filtering
	var toAdd []datastruct.Assignment
	var toUpdate []datastruct.Assignment
	alreadyExist := make(map[int64]datastruct.Assignment)
	for _, v := range remaining {
		alreadyExist[v.SegmentID] = v
	}
	for i, v := range segToAddT {
		if _, ok := alreadyExist[v.ID]; ok {
			toUpdate = append(toUpdate, datastruct.Assignment{
				UserID:    user.ID,
				SegmentID: v.ID,
			})
			if len(segToAdd[i].UntilDate) != 0 {
				var tmp time.Time
				tmp, err = time.Parse(time.RFC3339, segToAdd[i].UntilDate)
				toUpdate[len(toUpdate)-1].UntilDate =
					pgtype.Timestamptz{Time: tmp}
			}
		} else {
			toAdd = append(toAdd, datastruct.Assignment{
				UserID:    user.ID,
				SegmentID: v.ID,
			})
			if len(segToAdd[i].UntilDate) != 0 {
				var tmp time.Time
				tmp, err = time.Parse(time.RFC3339, segToAdd[i].UntilDate)
				toAdd[len(toAdd)-1].UntilDate =
					pgtype.Timestamptz{Time: tmp}
			}
		}
	}

	//adding assignments without time conflicts
	if len(toAdd) > 0 {
		var created []datastruct.Assignment
		created, err = uc.asRep.CreateWithConn(ctx, toAdd, tr.Conn())
		if err != nil {
			_ = tr.Rollback(context.TODO())
			return err
		}
		err = uc.hisRep.CreateAllWithConn(ctx, misc.AssignmentsToHistory(created, datastruct.OpAdded), tr.Conn())
		if err != nil {
			_ = tr.Rollback(context.TODO())
			return err
		}
	}

	//updating those with time conflicts
	if len(toUpdate) > 0 {
		var updated []datastruct.Assignment
		updated, err = uc.asRep.UpdateWithConn(ctx, toUpdate, tr.Conn())
		if err != nil {
			_ = tr.Rollback(context.TODO())
			return err
		}
		err = uc.hisRep.CreateAllWithConn(ctx, misc.AssignmentsToHistory(updated, datastruct.OpUpdated), tr.Conn())
		if err != nil {
			_ = tr.Rollback(context.TODO())
			return err
		}
	}

	err = tr.Commit(context.TODO())
	if err != nil {
		return err
	}

	return nil
}

func (uc PostgresAssignmentService) GetUserAssignments(ctx context.Context, userID int64) (dto.UserSegmentGet, error) {
	var err error
	if _, err = uc.userRep.Get(context.TODO(), userID); err != nil {
		return dto.UserSegmentGet{}, &errors.NoSuchUserError{UserID: userID}
	}

	var as []datastruct.Assignment
	if as, err = uc.asRep.GetAllForUser(context.TODO(), userID); err != nil {
		return dto.UserSegmentGet{}, err
	}

	segIds := make([]int64, len(as))
	for i, v := range as {
		segIds[i] = v.SegmentID
	}
	var segs []datastruct.Segment
	if segs, err = uc.segRep.GetAllById(context.TODO(), segIds); err != nil {
		return dto.UserSegmentGet{}, err
	}

	segSlugMap := make(map[int64]string, len(segs))
	for _, v := range segs {
		segSlugMap[v.ID] = v.Slug
	}

	var res dto.UserSegmentGet
	res.UserID = userID
	segRes := make([]dto.SegmentToAdd, len(as))
	for i, v := range as {
		segRes[i].Slug = segSlugMap[v.SegmentID]
		if v.UntilDate.Valid {
			segRes[i].UntilDate = v.UntilDate.Time.Format(time.RFC3339)
		}
	}
	res.SegmentAdded = segRes

	return res, nil
}

func (uc PostgresAssignmentService) GetUserHistory(ctx context.Context, userID int64, start, end time.Time) (dto.UserHistory, error) {
	var err error
	if _, err = uc.userRep.Get(context.TODO(), userID); err != nil {
		return dto.UserHistory{}, err
	}

	var his []datastruct.History
	if his, err = uc.hisRep.GetAllForUserPeriod(context.TODO(), userID, start, end); err != nil {
		return dto.UserHistory{}, err
	}

	segIds := make([]int64, len(his))
	for i, v := range his {
		segIds[i] = v.SegmentID
	}
	var segs []datastruct.Segment
	if segs, err = uc.segRep.GetAllById(context.TODO(), segIds); err != nil {
		return dto.UserHistory{}, err
	}

	segSlugMap := make(map[int64]string, len(segs))
	for _, v := range segs {
		segSlugMap[v.ID] = v.Slug
	}

	var res dto.UserHistory
	hisRes := make([]dto.UserHistoryRecord, len(his))
	for i, v := range his {
		hisRes[i].UserID = userID
		hisRes[i].Segment = segSlugMap[v.SegmentID]
		hisRes[i].Operation = datastruct.OperationCodeToName(datastruct.OperationCode(v.OperationID))
		hisRes[i].Timestamp = v.TimeStamp
	}
	res.Records = hisRes

	return res, nil
}
