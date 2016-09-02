package deploy

import "time"

type Repository interface {
	All(key string) []Deploy
	Since(key string, startTime time.Time) []Deploy
}
