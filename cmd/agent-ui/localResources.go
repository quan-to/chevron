package main

import (
	"context"
	"crypto"
	"fmt"
	"github.com/quan-to/chevron/pkg/chevronlib"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/pkg/models"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	MaxFileSize = 2 * 1024 * 1024 // 1MB
)

var (
	pgp              interfaces.PGPManager
	krm              interfaces.KeyRingManager
	executableFolder string
	keysFolder       string
)

var ctx = context.Background()

func init() {
	ex, _ := os.Executable()
	executableFolder = filepath.Dir(ex)
	keysFolder = path.Join(executableFolder, "keys")
}

func Begin() {
	_ = os.Mkdir(keysFolder, os.ModePerm)
	kb := chevronlib.MakeSaveToDiskBackend(nil, keysFolder, "key_")
	krm = chevronlib.MakeKeyRingManager(nil)
	pgp = chevronlib.MakePGPManager(nil, kb, krm)
	pgp.LoadKeys(ctx)
}

func AddPrivateKey(privateKeyData string) (string, error) {
	n, err := pgp.LoadKey(ctx, privateKeyData)
	if err != nil {
		return err.Error(), err
	}

	return fmt.Sprintf("Loaded %d keys", n), nil
}

func UnlockKey(fingerPrint, password string) (string, error) {
	err := pgp.UnlockKey(ctx, fingerPrint, password)
	if err != nil {
		log.Error("Error unlocking key %s: %s", fingerPrint, err)
		return err.Error(), err
	}

	return fmt.Sprintf("Key %s unlocked!", fingerPrint), nil
}

func Sign(fingerPrint, data string) (string, error) {
	key := pgp.GetPrivateKeyInfo(ctx, fingerPrint)

	if key == nil {
		err := fmt.Errorf("key not found")
		return err.Error(), err
	}

	if !key.PrivateKeyIsDecrypted {
		err := fmt.Errorf("key is not decrypted")
		return err.Error(), err
	}

	signature, err := pgp.SignData(ctx, fingerPrint, []byte(data), crypto.SHA512)
	if err != nil {
		return err.Error(), err
	}
	quantoSig := chevronlib.GPG2Quanto(signature, fingerPrint, "SHA512")
	return quantoSig, nil
}

func ListPrivateKeys() ([]models.KeyInfo, error) {
	return pgp.GetLoadedPrivateKeys(ctx), nil
}

func IsFile(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}

	return fi.Mode().IsRegular()
}

func FileSize(name string) int64 {
	fi, err := os.Stat(name)
	if err != nil {
		return 0
	}

	return fi.Size()
}

func GetFileContentType(name string) string {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err = f.Read(buffer)
	if err != nil {
		return "error reading file"
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType
}

func AddKeys(files []string) (bool, []string) {
	errors := make([]string, len(files))
	hasErrors := false

	for i, file := range files {
		if !IsFile(file) {
			errors[i] = fmt.Sprintf("%s is not a regular file", file)
			hasErrors = true
			continue
		}

		fileType := GetFileContentType(file)
		if !strings.Contains(fileType, "text/plain") {
			errors[i] = fmt.Sprintf("invalid file type: %s", fileType)
			hasErrors = true
			continue
		}

		size := FileSize(file)
		if size > MaxFileSize {
			errors[i] = fmt.Sprintf("file size too big: %d", size)
			hasErrors = true
			continue
		}

		data, err := ioutil.ReadFile(file)
		if err != nil {
			errors[i] = err.Error()
			hasErrors = true
			continue
		}

		fingerPrint, err := chevronlib.GetFingerprintFromKey(string(data))
		if err != nil {
			errors[i] = err.Error()
			continue
		}

		log.Info("Saving key %s from %s", fingerPrint, file)
		err = pgp.SaveKey(fingerPrint, string(data), nil)
		if err != nil {
			errors[i] = err.Error()
			hasErrors = true
			continue
		}

		pgp.LoadKeys(ctx)
	}

	return hasErrors, errors
}

func Migrate() {
	storeFolder := path.Join(executableFolder, "store")
	if chevronlib.FolderExists(storeFolder) { // Old key store
		log.Warn("Found \"store\" folder. Migrating keys...")
		err := chevronlib.CopyFiles(storeFolder, keysFolder)
		if err != nil {
			log.Error("Error moving files from store to keys: %s", err)
		} else {
			err = os.RemoveAll(storeFolder)
			if err != nil {
				log.Error("Error removing folder store: %s", err)
			}
		}
	}
}
