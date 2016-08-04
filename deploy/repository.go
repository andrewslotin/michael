package deploy

type Repository interface {
	All(key string) []Deploy
}
