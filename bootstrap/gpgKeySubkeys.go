package bootstrap

import (
	"github.com/quan-to/chevron"
	"github.com/quan-to/chevron/models"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func AddSubkeysToGPGKey(conn *r.Session) {
	l := log.SubScope("SubKey")
	l.Info("Running")
	keys, err := models.FetchKeysWithoutSubKeys(conn)
	if err != err {
		l.Fatal(err)
	}

	l.Note("Got %d keys to fill subkeys", len(keys))

	for _, k := range keys {
		fps, err := remote_signer.GetFingerPrintsFromKey(k.AsciiArmoredPublicKey)
		if err != nil {
			l.Error("Error getting fingerprints from key %s: %s", k.FullFingerPrint, err)
			_ = k.Delete(conn)
			continue
		}
		l.Note("Base: %s Keys: %v", k.GetShortFingerPrint(), fps)
		k.Subkeys = fps
		err = k.Save(conn)
		if err != nil {
			l.Error("Error saving key: %s", err)
		}
	}

	l.Info("Done")
}
