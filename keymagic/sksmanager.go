package keymagic

import (
	"fmt"
	"github.com/quan-to/chevron"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetSKSKey(fingerPrint string) (string, error) {
	response, err := http.Get(fmt.Sprintf("%s/pks/lookup?op=get&options=mr&search=0x%s", remote_signer.SKSServer, fingerPrint))

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
	response, err := http.PostForm(remote_signer.SKSServer, url.Values{"keytext": {publicKey}})

	if err != nil {
		return false, err
	}

	_ = response.Body.Close()
	return response.StatusCode == http.StatusOK, nil
}
