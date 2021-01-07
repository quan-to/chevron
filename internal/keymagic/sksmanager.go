package keymagic

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/quan-to/chevron/internal/config"
)

func GetSKSKey(fingerPrint string) (string, error) {
	response, err := http.Get(fmt.Sprintf("%s/pks/lookup?op=get&options=mr&search=0x%s", config.SKSServer, fingerPrint))

	if err != nil {
		return "", err
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	_ = response.Body.Close()

	return string(contents), nil
}

func PutSKSKey(publicKey string) (bool, error) {
	response, err := http.PostForm(config.SKSServer, url.Values{"keytext": {publicKey}})

	if err != nil {
		return false, err
	}

	_ = response.Body.Close()
	return response.StatusCode == http.StatusOK, nil
}
