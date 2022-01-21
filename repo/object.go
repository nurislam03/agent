package repo

type ObjectStore interface {
	GetObject(string) ([]byte, error)
}
