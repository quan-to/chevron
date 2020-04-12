package keybackend

import (
	"fmt"
	"github.com/quan-to/chevron/pkg/interfaces"
)

// voidBackend is a Key Backend that does nothing
type voidBackend struct{}

// MakeVoidBackend creates a new KeyBackend that does nothing
func MakeVoidBackend() interfaces.Backend {
	return &voidBackend{}
}

// Name returns the name of the KeyBackend
func (d *voidBackend) Name() string {
	return "voidBackend Backend"
}

// Path returns the path of the current KeyBackend
func (d *voidBackend) Path() string {
	return "*"
}

// Save saves a key to the backend
func (d *voidBackend) Save(key, data string) error {
	return nil
}

// SaveWithMetadata saves a key to backend storing some metadata with it
func (d *voidBackend) SaveWithMetadata(key, data, metadata string) error {
	return nil
}

// Delete a key from the backend
func (d *voidBackend) Delete(key string) error {
	return nil
}

// Read reads a key from the backend
func (d *voidBackend) Read(key string) (data string, metadata string, err error) {
	return "", "", fmt.Errorf("nothing to read")
}

// List lists the stored keys
func (d *voidBackend) List() ([]string, error) {
	return make([]string, 0), nil
}
