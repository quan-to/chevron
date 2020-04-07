package bootstrap

import (
	"github.com/quan-to/chevron/internal/models"
	"github.com/quan-to/chevron/internal/tools"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func AddSubkeysToGPGKey(conn *r.Session) {
	l := log.SubScope("SubKey")
	l.Await("Running")
	keys, err := models.FetchKeysWithoutSubKeys(conn)
	if err != err {
		l.Fatal(err)
	}

	l.Note("Got %d keys to fill subkeys", len(keys))

	for _, k := range keys {
		fps, err := tools.GetFingerPrintsFromKey(k.AsciiArmoredPublicKey)
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

	l.Done("Done")
}
