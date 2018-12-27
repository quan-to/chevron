package main

import (
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

	key, err := k.GeneratePGPKey("TEST", "12345", 16384)

	if err != nil {
		slog.Error(err)
	}

	log.Println(key)
}
