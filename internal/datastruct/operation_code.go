package datastruct

type OperationCode int64

const (
	OpAdded OperationCode = iota
	OpDeleted
	OpExpired
	OpUpdated
)
