// Code generated for package migrations by go-bindata DO NOT EDIT. (@generated)
// sources:
// migrations/000001_create_users_table.down.sql
// migrations/000001_create_users_table.up.sql
// migrations/000002_create_gpgkey_table.down.sql
// migrations/000002_create_gpgkey_table.up.sql
// migrations/000003_create_gpgkeyuid_table.down.sql
// migrations/000003_create_gpgkeyuid_table.up.sql
package migrations

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var __000001_create_users_tableDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x4d\xce\x48\xcc\x4b\x4f\x2d\x4e\x2d\x51\x28\x4a\x4c\x4e\x2d\xaa\x48\xc9\xb1\x4a\x2e\x4a\x4d\x2c\x49\x8d\x2f\x2d\x4e\x2d\x2a\x8e\x2f\x49\x4c\xca\x49\xe5\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\x48\xce\x48\x2d\x2b\xca\xcf\x03\x4b\x5b\x73\x01\x02\x00\x00\xff\xff\xec\x78\xf3\x92\x41\x00\x00\x00")

func _000001_create_users_tableDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000001_create_users_tableDownSql,
		"000001_create_users_table.down.sql",
	)
}

func _000001_create_users_tableDownSql() (*asset, error) {
	bytes, err := _000001_create_users_tableDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000001_create_users_table.down.sql", size: 65, mode: os.FileMode(436), modTime: time.Unix(1607539736, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000001_create_users_tableUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x90\xc1\x4a\x03\x31\x10\x86\xef\x79\x8a\x39\x76\xc1\xbe\x80\x3d\xad\x36\x42\x31\x6e\x65\xd9\x05\x7b\x0a\xd3\x64\xda\x0d\x64\xb3\xcb\x24\x69\xeb\xdb\x8b\x15\x75\xad\xe2\xa1\xb7\x39\x7c\xdf\xc0\xff\xcd\xe7\xa6\xc3\xb0\xa7\x48\x09\x18\x0d\xf1\xc9\xfa\x5b\xc3\x84\x89\x74\x8e\xc4\x51\x27\xdc\x7a\x12\xf7\xb5\x2c\x1b\x09\x4d\x79\xa7\x24\x98\x8e\x0e\x3c\x84\x33\x00\x33\x01\x00\xf0\x7e\x6a\x67\x21\x67\x67\xa1\x5a\x37\x10\xb2\xf7\xf0\x5c\xaf\x9e\xca\x7a\x03\x8f\x72\x73\xf3\x8d\xed\x5c\xd8\x13\x8f\xec\x42\x82\x03\xb2\xe9\x90\xcf\x4a\xd5\x2a\x35\xc1\x46\x8c\xf1\x38\xb0\x85\xed\x6b\x22\xfc\x8b\xd8\x65\xef\x75\xc0\x9e\xfe\x7b\xf3\x31\xc6\x6a\x4c\x90\x5c\x4f\x31\x61\x3f\x7e\x71\xb0\x94\x0f\x65\xab\x1a\x08\xc3\x71\x56\x4c\xac\x3c\xda\x2b\x2c\x4b\x9e\x7e\x5b\xad\x52\xa2\x58\x88\xcf\x86\xab\x6a\x29\x5f\x7e\x34\x9c\x16\xd1\xce\x9e\x60\x5d\x5d\x34\xbe\xa4\x8a\x85\x78\x0b\x00\x00\xff\xff\x01\x59\x93\x2e\xba\x01\x00\x00")

func _000001_create_users_tableUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000001_create_users_tableUpSql,
		"000001_create_users_table.up.sql",
	)
}

func _000001_create_users_tableUpSql() (*asset, error) {
	bytes, err := _000001_create_users_tableUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000001_create_users_table.up.sql", size: 442, mode: os.FileMode(436), modTime: time.Unix(1607539736, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000002_create_gpgkey_tableDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x4d\xce\x48\xcc\x4b\x4f\x2d\x4e\x2d\x51\x28\x4a\x4c\x4e\x2d\xaa\x48\xc9\xb1\x4a\x2e\x4a\x4d\x2c\x49\x8d\x4f\x2f\x48\xcf\x4e\xad\x8c\x2f\x49\x4c\xca\x49\xe5\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\x48\xce\x48\x2d\x2b\xca\xcf\x83\xca\x5b\x73\x01\x02\x00\x00\xff\xff\x6d\xa5\xec\x36\x44\x00\x00\x00")

func _000002_create_gpgkey_tableDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000002_create_gpgkey_tableDownSql,
		"000002_create_gpgkey_table.down.sql",
	)
}

func _000002_create_gpgkey_tableDownSql() (*asset, error) {
	bytes, err := _000002_create_gpgkey_tableDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000002_create_gpgkey_table.down.sql", size: 68, mode: os.FileMode(436), modTime: time.Unix(1607539736, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000002_create_gpgkey_tableUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x92\x31\x6f\xc2\x30\x10\x85\x77\xff\x8a\x1b\x89\x04\x43\x97\x0e\x65\x4a\x93\x43\x42\x75\x43\x15\x82\x54\x26\xcb\xc4\x47\xb0\x30\x21\x32\x0e\x85\x7f\x5f\xa5\x14\x82\x53\xa1\x0c\xf5\x96\xe4\x7d\x79\xcf\xef\x6e\x34\xca\x37\xb2\x2c\xe8\x40\x0e\xac\xcc\xc9\x9e\x94\x79\xc9\x2d\x49\x47\xa2\xa8\x8a\x2d\x9d\x85\x93\x2b\x43\x2c\x4a\x31\xcc\x10\xb2\xf0\x95\x23\xe4\x1b\x3a\xda\x7d\xd9\x28\xc4\x96\xce\x6c\xc0\x00\x00\x7e\x9f\x84\x56\xe0\x9f\xba\xbe\xbe\x4a\x66\x19\x24\x0b\xce\xe1\x23\x9d\xbe\x87\xe9\x12\xde\x70\x39\xf4\xe0\x75\x6d\x8c\x58\xeb\xb2\x20\x5b\x59\x5d\x3a\x38\x4a\x9b\x6f\xa4\xbd\x83\x3b\x40\xab\x7d\x7a\x6e\x3e\xf4\x01\x5b\x3a\xaf\xb4\x3b\xb4\xf1\x1a\x17\x3f\x9e\x0f\x54\xd2\xd2\x4d\xe2\xdf\x27\xc5\x09\xa6\x98\x44\x38\xef\x76\x02\x83\xb6\x8e\x00\x66\x09\xc4\xc8\x31\x43\x88\xc2\x79\x14\xc6\x38\x64\xbe\x45\xbd\x32\x3a\xff\xe1\x2e\xc7\xd1\xc9\x75\x52\x58\x7d\x6c\xa6\x72\xd3\x5c\x24\x9e\xe6\x32\x38\x25\xe4\x35\xad\xd3\x3b\x3a\x38\xb9\xab\xda\xe6\x63\x9c\x84\x0b\x9e\x41\xb9\xff\x1a\x04\xbe\x45\x5d\xa9\xff\xe0\x8a\x0c\x3d\xc2\x17\x9c\xb3\x60\xcc\xae\x6b\x34\x4d\x62\xfc\xec\x56\xf6\x67\xf6\x42\xab\x53\x53\xdd\xc3\x6a\xbb\x40\x30\xee\x31\xb8\xdf\x95\xfe\xbf\xdf\xab\x83\x31\xfb\x0e\x00\x00\xff\xff\x5e\xe4\x7f\x02\x2c\x03\x00\x00")

func _000002_create_gpgkey_tableUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000002_create_gpgkey_tableUpSql,
		"000002_create_gpgkey_table.up.sql",
	)
}

func _000002_create_gpgkey_tableUpSql() (*asset, error) {
	bytes, err := _000002_create_gpgkey_tableUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000002_create_gpgkey_table.up.sql", size: 812, mode: os.FileMode(436), modTime: time.Unix(1607621856, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000003_create_gpgkeyuid_tableDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x4d\xce\x48\xcc\x4b\x4f\x2d\x4e\x2d\x51\x28\x4a\x4c\x4e\x2d\xaa\x48\xc9\xb1\x4a\x2e\x4a\x4d\x2c\x49\x8d\x4f\x2f\x48\xcf\x4e\xad\x2c\xcd\x4c\x89\x2f\x49\x4c\xca\x49\xe5\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\x48\xce\x48\x2d\x2b\xca\xcf\x03\x29\x89\xcf\x4e\xad\x8c\x2f\xcd\x4c\xb1\xe6\x02\x04\x00\x00\xff\xff\xa0\x1d\xf9\xc4\x4c\x00\x00\x00")

func _000003_create_gpgkeyuid_tableDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000003_create_gpgkeyuid_tableDownSql,
		"000003_create_gpgkeyuid_table.down.sql",
	)
}

func _000003_create_gpgkeyuid_tableDownSql() (*asset, error) {
	bytes, err := _000003_create_gpgkeyuid_tableDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000003_create_gpgkeyuid_table.down.sql", size: 76, mode: os.FileMode(436), modTime: time.Unix(1607539736, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000003_create_gpgkeyuid_tableUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x91\xdf\x6a\xc2\x30\x14\xc6\xef\xf3\x14\xe7\xd2\x82\xbe\xc0\xbc\xca\xda\x23\xc8\xb2\x3a\x6a\x85\x79\x15\xb2\xe4\xa0\xc1\x34\x96\x98\x38\x7d\xfb\xa1\xa3\x73\xac\x1d\x0e\x96\xbb\x24\xbf\xef\x3b\x7f\xbe\xc9\x44\x6f\x95\xdf\xd0\x81\x22\x04\xa5\x29\x9c\x8c\x7b\xd0\x81\x54\x24\xb9\x69\x37\x3b\x3a\x27\x6b\x64\x54\x6f\x8e\x58\x5e\x21\xaf\x11\x6a\xfe\x28\x10\xf4\x96\x8e\x61\xef\x2f\x90\xdc\xd1\x59\x26\x6b\xd8\x88\x01\x00\x7c\x7b\x91\xd6\xc0\xd7\x49\xa9\xbb\x95\x8b\x1a\x7c\x72\x0e\x5e\xaa\xf9\x33\xaf\xd6\xf0\x84\xeb\x71\x4f\xeb\x55\x43\x9d\xf6\xa8\x82\xde\xaa\xd0\x87\xa8\x51\xd6\xdd\x83\x0c\x1d\x74\xb0\x6d\xb4\x7b\xff\x3b\xd4\xaa\x40\x3e\x0e\xb5\x5a\xae\x84\x80\x0a\x67\x58\x61\x99\xe3\xf2\xe7\xe8\x30\xea\x6c\xac\xc9\x60\x51\x42\x81\x02\x6b\x84\x9c\x2f\x73\x5e\xe0\x98\xf5\x4a\x7d\xee\xd7\x48\x15\x01\xa2\x6d\xe8\x10\x55\xd3\xde\x4a\x15\x38\xe3\x2b\x51\x83\xdf\xbf\x8f\xb2\x7e\xa3\xa9\x35\xff\x50\x1b\x72\x34\xa0\x5e\x09\xc1\xb2\x29\xeb\x32\x9e\x97\x05\xbe\x0e\x65\x7c\x4d\x45\x5a\x73\xba\x0c\x3a\xf0\x7f\x5b\x46\x07\x67\xd3\xfb\xa6\xd7\x14\xff\xec\x7a\xa5\xb3\x29\xfb\x08\x00\x00\xff\xff\x70\x54\x0d\x96\xbc\x02\x00\x00")

func _000003_create_gpgkeyuid_tableUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000003_create_gpgkeyuid_tableUpSql,
		"000003_create_gpgkeyuid_table.up.sql",
	)
}

func _000003_create_gpgkeyuid_tableUpSql() (*asset, error) {
	bytes, err := _000003_create_gpgkeyuid_tableUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000003_create_gpgkeyuid_table.up.sql", size: 700, mode: os.FileMode(436), modTime: time.Unix(1607539736, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"000001_create_users_table.down.sql":     _000001_create_users_tableDownSql,
	"000001_create_users_table.up.sql":       _000001_create_users_tableUpSql,
	"000002_create_gpgkey_table.down.sql":    _000002_create_gpgkey_tableDownSql,
	"000002_create_gpgkey_table.up.sql":      _000002_create_gpgkey_tableUpSql,
	"000003_create_gpgkeyuid_table.down.sql": _000003_create_gpgkeyuid_tableDownSql,
	"000003_create_gpgkeyuid_table.up.sql":   _000003_create_gpgkeyuid_tableUpSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"000001_create_users_table.down.sql":     &bintree{_000001_create_users_tableDownSql, map[string]*bintree{}},
	"000001_create_users_table.up.sql":       &bintree{_000001_create_users_tableUpSql, map[string]*bintree{}},
	"000002_create_gpgkey_table.down.sql":    &bintree{_000002_create_gpgkey_tableDownSql, map[string]*bintree{}},
	"000002_create_gpgkey_table.up.sql":      &bintree{_000002_create_gpgkey_tableUpSql, map[string]*bintree{}},
	"000003_create_gpgkeyuid_table.down.sql": &bintree{_000003_create_gpgkeyuid_tableDownSql, map[string]*bintree{}},
	"000003_create_gpgkeyuid_table.up.sql":   &bintree{_000003_create_gpgkeyuid_tableUpSql, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
