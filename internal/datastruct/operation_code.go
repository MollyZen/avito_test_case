package datastruct

type OperationCode int64

const (
	OpAdded OperationCode = iota
	OpRemoved
	OpExpired
	OpUpdated
	OpSegDeleted
	OpAddedRand
)

func OperationCodeToName(code OperationCode) string {
	var res string
	switch code {
	case OpAdded:
		res = "Добавление"
	case OpRemoved:
		res = "Удаление"
	case OpExpired:
		res = "Удаление по истечении времени"
	case OpUpdated:
		res = "Обновление значения"
	case OpSegDeleted:
		res = "Удален вместе с сегментом"
	}
	return res
}
