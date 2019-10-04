package keymagic

import (
	"context"
	"fmt"

	remote_signer "github.com/quan-to/chevron"
	"github.com/quan-to/chevron/database"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/slog"
)

var pksLog = slog.Scope("PKS")

func PKSGetKey(ctx context.Context, fingerPrint string) (string, error) {
	requestID := remote_signer.GetRequestIDFromContext(ctx)
	log := pksLog.Tag(requestID)
	log.DebugNote("PKSGetKey(%s)", fingerPrint)
	if !remote_signer.EnableRethinkSKS {
		return GetSKSKey(fingerPrint)
	}

	conn := database.GetConnection()
	v, err := models.GetGPGKeyByFingerPrint(conn, fingerPrint)

	if v != nil {
		return v.AsciiArmoredPublicKey, nil
	}

	return "", err
}

func PKSSearchByName(name string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	pksLog.DebugNote("PKSSearchByName(%s, %d, %d)", name, pageStart, pageEnd)
	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()
		return models.SearchGPGKeyByName(conn, name, pageStart, pageEnd)
	}
	return nil, fmt.Errorf("the server does not have RethinkDB enabled so it cannot serve search")
}

func PKSSearchByFingerPrint(fingerPrint string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	pksLog.DebugNote("PKSSearchByFingerPrint(%s, %d, %d)", fingerPrint, pageStart, pageEnd)
	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()
		return models.SearchGPGKeyByFingerPrint(conn, fingerPrint, pageStart, pageEnd)
	}
	return nil, fmt.Errorf("the server does not have RethinkDB enabled so it cannot serve search")
}

func PKSSearchByEmail(email string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	pksLog.DebugNote("PKSSearchByEmail(%s, %d, %d)", email, pageStart, pageEnd)
	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()
		return models.SearchGPGKeyByEmail(conn, email, pageStart, pageEnd)
	}
	return nil, fmt.Errorf("the server does not have RethinkDB enabled so it cannot serve search")
}

func PKSSearch(value string, pageStart, pageEnd int) ([]models.GPGKey, error) {
	pksLog.DebugNote("PKSSearch(%s, %d, %d)", value, pageStart, pageEnd)
	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()
		return models.SearchGPGKeyByValue(conn, value, pageStart, pageEnd)
	}
	return nil, fmt.Errorf("the server does not have RethinkDB enabled so it cannot serve search")
}

func PKSAdd(ctx context.Context, pubKey string) string {
	requestID := remote_signer.GetRequestIDFromContext(ctx)
	log := pksLog.Tag(requestID)
	log.DebugNote("PKSAdd(---)")
	if remote_signer.EnableRethinkSKS {
		conn := database.GetConnection()
		key, err := models.AsciiArmored2GPGKey(pubKey)
		if err != nil {
			log.Debug("PKSAdd Error: %s", err)
			return "NOK"
		}

		keys, err := models.SearchGPGKeyByFingerPrint(conn, key.FullFingerPrint, 0, 1)

		if err != nil {
			log.Debug("PKSAdd Error: %s", err)
			return "NOK"
		}

		if len(keys) > 0 {
			log.Info("Tried to add key %s to PKS but already exists.", key.GetShortFingerPrint())
			return "OK"
		}

		log.Info("Adding public key %s to PKS", key.GetShortFingerPrint())
		_, _, err = models.AddGPGKey(conn, key)

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
