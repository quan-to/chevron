package keyBackend

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/slog"
	"io/ioutil"
	"os"
	"path"
)

type Disk struct {
	folder      string
	prefix      string
	saveEnabled bool
	log         slog.Instance
}

// MakeSaveToDiskBackend creates an instance of DiskBackend that stores keys in files
func MakeSaveToDiskBackend(log slog.Instance, folder, prefix string) *Disk {
	if log == nil {
		log = slog.Scope("Disk")
	} else {
		log = log.SubScope("Disk")
	}

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
		log:         log,
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
		d.log.Warn("Save disabled")
		return nil
	}

	d.log.DebugAwait("Saving to %s", path.Join(d.folder, d.prefix+key))

	err := ioutil.WriteFile(path.Join(d.folder, d.prefix+key), []byte(data), 0600)

	if err != nil {
		d.log.ErrorDone("Error saving to %s: %s", path.Join(d.folder, d.prefix+key), err)
	}

	return err
}

func (d *Disk) SaveWithMetadata(key, data, metadata string) error {
	d.log.DebugAwait("Saving to %s", path.Join(d.folder, d.prefix+key))
	err := ioutil.WriteFile(path.Join(d.folder, d.prefix+key), []byte(data), 0600)
	if err != nil {
		d.log.ErrorDone("Error saving to %s: %s", path.Join(d.folder, d.prefix+key), err)
		return err
	}

	err = ioutil.WriteFile(path.Join(d.folder, "metadata-"+d.prefix+key), []byte(metadata), 0600)

	if err != nil {
		d.log.ErrorDone("Error saving to %s: %s", path.Join(d.folder, d.prefix+key), err)
	}

	return err
}

// Delete deletes a file key and metadata from the disk
func (d *Disk) Delete(key string) error {
	d.log.DebugAwait("Deleting %s", path.Join(d.folder, d.prefix+key))
	_, err := ioutil.ReadFile(path.Join(d.folder, d.prefix+key))
	if err != nil {
		d.log.ErrorDone("Error reading to %s: %s, file not exist to delete", path.Join(d.folder, d.prefix+key), err)
		return err
	}

	err = os.Remove(path.Join(d.folder, d.prefix+key))
	if err != nil {
		d.log.ErrorDone("Error deleting from %s: %s", path.Join(d.folder, d.prefix+key), err)
	}

	_ = os.Remove(path.Join(d.folder, "metadata-"+d.prefix+key))

	return err
}

func (d *Disk) Read(key string) (data string, metadata string, err error) {
	d.log.DebugAwait("Reading from %s", path.Join(d.folder, d.prefix+key))
	sdata, err := ioutil.ReadFile(path.Join(d.folder, d.prefix+key))
	if err != nil {
		d.log.ErrorDone("Error reading to %s: %s", path.Join(d.folder, d.prefix+key), err)
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
