package keyBackend

import (
	"fmt"
)

type Void struct {
	folder      string
	prefix      string
	saveEnabled bool
}

func MakeVoidBackend() *Void {
	return &Void{}
}

func (d *Void) Name() string {
	return "Void Backend"
}

func (d *Void) Path() string {
	return "*"
}

func (d *Void) Save(key, data string) error {
	return nil
}

func (d *Void) SaveWithMetadata(key, data, metadata string) error {
	return nil
}

func (d *Void) Read(key string) (data string, metadata string, err error) {
	return "", "", fmt.Errorf("nothing to read")
}

func (d *Void) List() ([]string, error) {
	return make([]string, 0), nil
}
