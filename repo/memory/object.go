package memory

import (
	"github.com/nurislam03/agent/repo"
	"os"
)

type objectStore struct {}

//NewObjectStore ...
func NewObjectStore() repo.ObjectStore{
	return &objectStore{}
}

func (o *objectStore) GetObject(s string) ([]byte, error) {
	return os.ReadFile("./files/golang.png")
}