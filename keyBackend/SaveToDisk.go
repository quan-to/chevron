package keyBackend

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/slog"
	"io/ioutil"
	"path"
)

var log = slog.Scope("DiskBackend")

type Disk struct {
	folder      string
	prefix      string
	saveEnabled bool
}

func MakeSaveToDiskBackend(folder, prefix string) *Disk {
	saveEnabled := true
	log.Info("Initialized DiskBackend on folder %s with prefix %s", folder, prefix)
	if remote_signer.ReadonlyKeyPath {
		log.Warn("Readonly keypath. Creating temporary storage in disk.")
		tmpFolder, err := ioutil.TempDir("/tmp", "secret")
		if err != nil {
			log.Error("Error creating temporary folder. Disabling save.")
			saveEnabled = false
		} else {
			log.Info("Copying files from %s to %s", folder, tmpFolder)
			err = remote_signer.CopyFiles(folder, tmpFolder)
			if err != nil {
				saveEnabled = false
				log.Error("Cannot copy files from %s to %s: %s", folder, tmpFolder, err)
			} else {
				folder = tmpFolder
			}
		}
	}
	return &Disk{
		folder:      folder,
		prefix:      prefix,
		saveEnabled: saveEnabled,
	}
}

func (d *Disk) Name() string {
	return "Disk Backend"
}

func (d *Disk) Path() string {
	return path.Join(d.folder, d.prefix+"*")
}

func (d *Disk) Save(key, data string) error {
	if !d.saveEnabled {
		log.Warn("Save disabled")
		return nil
	}
	log.Debug("Saving to %s", path.Join(d.folder, d.prefix+key))
	return ioutil.WriteFile(path.Join(d.folder, d.prefix+key), []byte(data), 0600)
}

func (d *Disk) SaveWithMetadata(key, data, metadata string) error {
	log.Debug("Saving to %s", path.Join(d.folder, d.prefix+key))
	err := ioutil.WriteFile(path.Join(d.folder, d.prefix+key), []byte(data), 0600)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(d.folder, "metadata-"+d.prefix+key), []byte(metadata), 0600)
}

func (d *Disk) Read(key string) (data string, metadata string, err error) {
	log.Debug("Reading from %s", path.Join(d.folder, d.prefix+key))
	sdata, err := ioutil.ReadFile(path.Join(d.folder, d.prefix+key))
	if err != nil {
		return "", "", err
	}

	mdata, err := ioutil.ReadFile(path.Join(d.folder, "metadata-"+d.prefix+key))
	if err != nil {
		return string(sdata), "", nil
	}

	return string(sdata), string(mdata), nil
}

func (d *Disk) List() ([]string, error) {
	files, err := ioutil.ReadDir(d.folder)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0)

	for _, file := range files {
		fileName := file.Name()
		if !file.IsDir() && len(fileName) > len(d.prefix) && fileName[:len(d.prefix)] == d.prefix {
			keys = append(keys, fileName[len(d.prefix):])
		}
	}

	return keys, nil
}
