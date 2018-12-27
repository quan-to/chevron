package main

import (
	"crypto"
	"github.com/quan-to/remote-signer"
	"github.com/quan-to/remote-signer/SLog"
	"log"
)

var slog = SLog.Scope("RemoteSigner")

func main() {
	//r := mux.NewRouter()
	//remote_signer.AddHKPEndpoints(r)
	//
	//srv := &http.Server{
	//	Addr:    "0.0.0.0:8090",
	//	Handler: r, // Pass our instance of gorilla/mux in.
	//}
	//
	//if err := srv.ListenAndServe(); err != nil {
	//	log.Println(err)
	//}

	k := remote_signer.MakePGPManager()
	k.LoadKeys()

	err := k.UnlockKey("0016A9CA870AFA59", "I think you will never guess")

	if err != nil {
		slog.Error(err)
	}
	for _, v := range k.GetLoadedPrivateKeys() {
		slog.Info("		Key: %s - Private Key Decrypted: %t", v.Identifier, v.PrivateKeyIsDecrypted)
	}

	data := []byte("huebr for the win!")
	data2 := []byte("huebr for the win")

	signature, err := k.SignData("0016A9CA870AFA59", data, crypto.SHA512)

	if err != nil {
		slog.Error(err)
		return
	}

	slog.Info("Signature: \n%s", signature)
	valid, err := k.VerifySignature(data2, signature)
	if err != nil {
		slog.Error(err)
	}

	if valid {
		log.Println("Signature is valid")
	}
}
