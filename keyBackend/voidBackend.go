package keyBackend

import (
	"fmt"
)

// Void is a Key Backend that does nothing
type Void struct{}

// MakeVoidBackend creates a new KeyBackend that does nothing
func MakeVoidBackend() *Void {
	return &Void{}
}

// Name returns the name of the KeyBackend
func (d *Void) Name() string {
	return "Void Backend"
}

// Path returns the path of the current KeyBackend
func (d *Void) Path() string {
	return "*"
}

// Save saves a key to the backend
func (d *Void) Save(key, data string) error {
	return nil
}

// SaveWithMetadata saves a key to backend storing some metadata with it
func (d *Void) SaveWithMetadata(key, data, metadata string) error {
	return nil
}

// Delete a key from the backend
func (d *Void) Delete(key string) error {
	return nil
}

// Read reads a key from the backend
func (d *Void) Read(key string) (data string, metadata string, err error) {
	return "", "", fmt.Errorf("nothing to read")
}

// List lists the stored keys
func (d *Void) List() ([]string, error) {
	return make([]string, 0), nil
}
