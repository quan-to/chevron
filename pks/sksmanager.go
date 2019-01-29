package pks

import (
	"fmt"
	"github.com/quan-to/remote-signer"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetSKSKey(fingerPrint string) string {
	response, err := http.Get(fmt.Sprintf("%s/pks/lookup?op=get&options=mr&search=0x%s", remote_signer.SKSServer, fingerPrint))

	if err != nil {
		panic(err)
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	_ = response.Body.Close()

	return string(contents)
}

func PutSKSKey(publicKey string) bool {
	response, err := http.PostForm(remote_signer.SKSServer, url.Values{"keytext": {publicKey}})

	if err != nil {
		panic(err)
	}

	_ = response.Body.Close()
	return response.StatusCode == http.StatusOK
}
