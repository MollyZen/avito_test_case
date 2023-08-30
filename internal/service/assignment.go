package service

import (
	"avito_test_case/internal/datastruct"
	"avito_test_case/internal/dto"
	"avito_test_case/internal/repository"
	"avito_test_case/pkg/logger"
	"context"
	"fmt"
)

type AssignmentUseCase struct {
	userRep repository.UserRepository
	segRep  repository.SegmentRepository
	asRep   repository.AssignmentRepository
	hisRep  repository.HistoryRepository
	l       logger.Logger
}

func NewAssignmentUseCase(userRep repository.UserRepository, segRep repository.SegmentRepository, asRepo repository.AssignmentRepository, hisRepos repository.HistoryRepository, l logger.Logger) *AssignmentUseCase {
	return &AssignmentUseCase{
		userRep: userRep,
		segRep:  segRep,
		asRep:   asRepo,
		hisRep:  hisRepos,
		l:       l,
	}
}

func (uc AssignmentUseCase) Assign(ctx context.Context, userId int64, segToAdd []dto.SegmentToAdd, segToDelete []string) error {
	//adding user if they don't exist
	var err error
	var user datastruct.User
	user, err = uc.userRep.Upsert(ctx, datastruct.User{
		ID: userId,
	})
	if err != nil {
		return err
	}

	//getting ids for segment slugs + checking if they exist
	var segs []datastruct.Segment
	segToAddNames := make([]string, len(segToAdd))
	for i, v := range segToAdd {
		segToAddNames[i] = v.Name
	}
	toAddAndDelete := append(segToAddNames, segToDelete...)
	segs, err = uc.segRep.GetForIds(ctx, toAddAndDelete)
	if err != nil {
		return err
	}
	if len(segs) != len(toAddAndDelete) {
		for i, v := range toAddAndDelete {
			if segs[i].Name != v {
				return fmt.Errorf("segment with slug %s doesn't exist", v)
			}
		}
	}
	segToAddT := segs[0:len(segToAdd)]
	segToDeleteT := segs[len(segToAdd):]

	//deleting assignments
	if len(segToDeleteT) > 0 {
		toDeleteAs := make([]datastruct.Assignment, len(segToDeleteT))
		for i, seg := range segToDeleteT {
			toDeleteAs[i].UserID = user.ID
			toDeleteAs[i].SegmentID = seg.ID
		}
		var deletedAs []datastruct.Assignment
		deletedAs, err = uc.asRep.Delete(ctx, toDeleteAs)
		err = uc.hisRep.CreateAll(ctx, assignmentsToHistory(deletedAs, datastruct.OpDeleted))
		if err != nil {
			return err
		}
	}

	//getting remaining assignments
	var remaining []datastruct.Assignment
	remaining, err = uc.asRep.GetAllForUser(ctx, user.ID)
	if err != nil {
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
				UntilDate: segToAdd[i].UntilDate,
			})
		} else {
			toAdd = append(toAdd, datastruct.Assignment{
				UserID:    user.ID,
				SegmentID: v.ID,
				UntilDate: segToAdd[i].UntilDate,
			})
		}
	}

	//adding assignments without time conflicts
	if len(toAdd) > 0 {
		var created []datastruct.Assignment
		created, err = uc.asRep.Create(ctx, toAdd)
		err = uc.hisRep.CreateAll(ctx, assignmentsToHistory(created, datastruct.OpAdded))
		if err != nil {
			return err
		}
	}

	//updating those with time conflicts
	if len(toUpdate) > 0 {
		var updated []datastruct.Assignment
		updated, err = uc.asRep.Create(ctx, toUpdate)
		err = uc.hisRep.CreateAll(ctx, assignmentsToHistory(updated, datastruct.OpAdded))
		if err != nil {
			return err
		}
	}

	return nil
}

func assignmentsToHistory(as []datastruct.Assignment, code datastruct.OperationCode) []datastruct.History {
	res := make([]datastruct.History, len(as))
	for i, v := range as {
		res[i].UserID = v.UserID
		res[i].SegmentID = v.SegmentID
		res[i].OperationID = int64(code)
	}
	return res
}
