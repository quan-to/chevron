package keymagic

import (
	"context"
	"fmt"

	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/chevron/pkg/models"

	"github.com/quan-to/slog"
)

type DatabaseHandler interface {
	AddGPGKey(key models.GPGKey) (string, bool, error)
	FindGPGKeyByEmail(email string, pageStart, pageEnd int) ([]models.GPGKey, error)
	FindGPGKeyByFingerPrint(fingerPrint string, pageStart, pageEnd int) ([]models.GPGKey, error)
	FindGPGKeyByValue(value string, pageStart, pageEnd int) ([]models.GPGKey, error)
	FindGPGKeyByName(name string, pageStart, pageEnd int) ([]models.GPGKey, error)
	FetchGPGKeyByFingerprint(fingerprint string) (*models.GPGKey, error)
}

var pksLog = slog.Scope("PKS")

func dbHandlerFromContext(ctx context.Context) DatabaseHandler {
	dbhI := ctx.Value(tools.CtxDatabaseHandler)
	if dbhI != nil {
		dbh, ok := dbhI.(DatabaseHandler)
		if ok {
			return dbh
		}
	}

	return nil
}

func PKSGetKey(ctx context.Context, fingerPrint string) (string, error) {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pksLog.Tag(requestID)
	log.DebugNote("PKSGetKey(%q)", fingerPrint)
	dbh := dbHandlerFromContext(ctx)
	if dbh == nil {
		return GetSKSKey(fingerPrint)
	}

	v, err := dbh.FetchGPGKeyByFingerprint(fingerPrint)

	if v != nil {
		return v.AsciiArmoredPublicKey, nil
	}

	return "", err
}

func PKSSearchByName(ctx context.Context, name string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	pksLog.DebugNote("PKSSearchByName(%s, %d, %d)", name, pageStart, pageEnd)
	dbh := dbHandlerFromContext(ctx)
	if dbh != nil {
		return dbh.FindGPGKeyByName(name, pageStart, pageEnd)
	}

	return nil, fmt.Errorf("the server does not have database enabled so it cannot serve search")
}

func PKSSearchByFingerPrint(ctx context.Context, fingerPrint string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	pksLog.DebugNote("PKSSearchByFingerPrint(%s, %d, %d)", fingerPrint, pageStart, pageEnd)
	dbh := dbHandlerFromContext(ctx)
	if dbh != nil {
		return dbh.FindGPGKeyByFingerPrint(fingerPrint, pageStart, pageEnd)
	}
	return nil, fmt.Errorf("the server does not have database enabled so it cannot serve search")
}

func PKSSearchByEmail(ctx context.Context, email string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	pksLog.DebugNote("PKSSearchByEmail(%s, %d, %d)", email, pageStart, pageEnd)
	dbh := dbHandlerFromContext(ctx)
	if dbh != nil {
		return dbh.FindGPGKeyByEmail(email, pageStart, pageEnd)
	}
	return nil, fmt.Errorf("the server does not have database enabled so it cannot serve search")
}

func PKSSearch(ctx context.Context, value string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	pksLog.DebugNote("PKSSearch(%s, %d, %d)", value, pageStart, pageEnd)
	dbh := dbHandlerFromContext(ctx)
	if dbh != nil {
		return dbh.FindGPGKeyByValue(value, pageStart, pageEnd)
	}

	return nil, fmt.Errorf("the server does not have database enabled so it cannot serve search")
}

func PKSAdd(ctx context.Context, pubKey string) string {
	requestID := tools.GetRequestIDFromContext(ctx)
	log := pksLog.Tag(requestID)
	log.DebugNote("PKSAdd(---)")
	dbh := dbHandlerFromContext(ctx)
	if dbh != nil {
		key, err := models.AsciiArmored2GPGKey(pubKey)
		if err != nil {
			log.Debug("PKSAdd Error: %s", err)
			return "NOK"
		}

		keys, err := dbh.FindGPGKeyByFingerPrint(key.FullFingerprint, 0, 1)

		if err != nil {
			log.Debug("PKSAdd Error: %s", err)
			return "NOK"
		}

		if len(keys) > 0 {
			log.Info("Tried to add key %s to PKS but already exists.", key.GetShortFingerPrint())
			return "OK"
		}

		log.Info("Adding public key %s to PKS", key.GetShortFingerPrint())
		_, _, err = dbh.AddGPGKey(key)

		if err != nil {
			log.Debug("PKSAdd Error: %s", err)
			return "NOK"
		}

		return "OK"
	}

	res, err := PutSKSKey(pubKey)

	if err != nil {
		log.Debug("PKSAdd Error: %s", err)
	}

	if res {
		return "OK"
	}

	return "NOK"
}
