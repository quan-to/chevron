package kubernetes

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getWithToken(url, token string) (string, error) {
	// skipcq: GSC-G402
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // Sad :(
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	bodyText, err := ioutil.ReadAll(res.Body)

	return string(bodyText), err
}
