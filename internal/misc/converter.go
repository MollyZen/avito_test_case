package misc

import "avito_test_case/internal/datastruct"

func AssignmentsToHistory(as []datastruct.Assignment, code datastruct.OperationCode) []datastruct.History {
	res := make([]datastruct.History, len(as))
	for i, v := range as {
		res[i].UserID = v.UserID
		res[i].SegmentID = v.SegmentID
		res[i].OperationID = int64(code)
	}
	return res
}
