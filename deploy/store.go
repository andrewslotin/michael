package deploy

type Store interface {
	Get(key string) (d Deploy, ok bool)
	Set(key string, d Deploy)
}
