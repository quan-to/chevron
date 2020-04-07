package interfaces

// Backend is a interface for storing / reading keys
type Backend interface {
	// Save saves a key to the backend
	Save(key, data string) error
	// SaveWithMetadata saves a key to backend storing some metadata with it
	SaveWithMetadata(key, data, metadata string) error
	// Delete delete a key from backend
	Delete(key string) error
	// Read reads a key from the backend
	Read(key string) (data string, metadata string, err error)
	// List lists the stored keys
	List() ([]string, error)
	// Name returns the name of the KeyBackend
	Name() string
	// Path returns the path of the current KeyBackend
	Path() string
}
