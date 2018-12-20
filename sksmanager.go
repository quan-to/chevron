package remote_signer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetSKSKey(fingerPrint string) string {
	response, err := http.Get(fmt.Sprintf("%s/pks/lookup?op=get&options=mr&search=0x%s", SKSServer, fingerPrint))

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	return string(contents)
}

func PutSKSKey(publicKey string) bool {
	response, err := http.PostForm(SKSServer, url.Values{"keytext": {publicKey}})

	if err != nil {
		panic(err)
	}

	_ = response.Body.Close()
	return response.StatusCode == http.StatusOK
}
