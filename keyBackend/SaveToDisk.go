package saveTo

import (
    "github.com/quan-to/remote-signer/SLog"
    "io/ioutil"
    "path"
)

var slog = SLog.Scope("DiskBackend")

type Disk struct {
    folder string
    prefix string
}

func MakeSaveToDiskBackend(folder, prefix string) *Disk {
    slog.Info("Initialized DiskBackend on folder %s with prefix %s", folder, prefix)
    return &Disk{
        folder: folder,
        prefix: prefix,
    }
}

func (d *Disk) Save(key, data string) error {
    return ioutil.WriteFile(path.Join(d.folder, key), []byte(data), 0600)
}

func (d *Disk) Read(key string) (string ,error) {
    data, err := ioutil.ReadFile(path.Join(d.folder, key))
    if err != nil {
        return "", err
    }

    return string(data), nil
}

func (d *Disk) List(keyPrefix string) ([]string, error) {
    files, err := ioutil.ReadDir(d.folder)
    if err != nil {
        return nil, err
    }

    keys := make([]string, 0)

    for _, file := range files {
        fileName := file.Name()
        if !file.IsDir() && len(fileName) > len(d.prefix) && fileName[:len(d.prefix)] == d.prefix {
            keys = append(keys, fileName)
        }
    }

    return keys, nil
}