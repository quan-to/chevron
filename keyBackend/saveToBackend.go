package keyBackend

type Backend interface {
	Save(key, data string) error
	SaveWithMetadata(key, data, metadata string) error
	Read(key string) (data string, metadata string, err error)
	List() ([]string, error)
	Name() string
	Path() string
}
