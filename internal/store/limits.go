package store

const (
	defaultListLimit = int32(100)
	maxListLimit     = int32(500)
)

func boundedListLimit(limit int) int32 {
	if limit <= 0 || limit > int(maxListLimit) {
		return defaultListLimit
	}
	return int32(limit)
}
