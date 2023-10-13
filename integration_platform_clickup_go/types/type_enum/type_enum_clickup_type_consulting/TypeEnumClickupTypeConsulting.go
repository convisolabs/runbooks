package TypeEnumClickupTypeConsulting

const (
	EPIC  int = 0
	STORE int = 1
	TASK  int = 2
)

func ToString(value int) string {
	switch value {
	case EPIC:
		return "EPIC"
	case STORE:
		return "STORE"
	case TASK:
		return "TASK"
	}
	return "unknown"
}
