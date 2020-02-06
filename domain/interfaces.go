package domain

type UrlGetter interface {
	Get(string) (string, error)
}
