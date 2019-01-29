package saveTo


type Backend interface {
    Save(key, data string) error
    Read(key string) (string ,error)
    List(keyPrefix string) ([]string, error)
}