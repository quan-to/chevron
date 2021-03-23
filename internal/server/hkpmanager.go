package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/pkg/models/HKP"

	"github.com/gorilla/mux"
	"github.com/quan-to/slog"
)

/// HKP Server based on https://tools.ietf.org/html/draft-shaw-openpgp-hkp-00

// HKP Lookup godoc
// @id hkp-lookup
// @tags SKS
// @Summary GPG SKS Keyserver lookup
// @Accept plain
// @Produce plain
// @param op query string true "HKP Operation. Valid values: get, index, vindex"
// @param options query string true "HKP Operation options. Valid values: mr, nm"
// @param search query string true "HKP Search Value"
// @Success 200 {string} result "result of the query"
// @Failure default {object} QuantoError.ErrorObject
// @Router /pks/lookup [get]
func operationGet(ctx context.Context, options, searchData string, machineReadable, noModification bool) (string, error) {
	if searchData[:2] == "0x" {
		k, _ := keymagic.PKSGetKey(ctx, searchData[2:])
		if k == "" {
			return "", errors.New("not found")
		}
		return k, nil
	}

	results, err := keymagic.PKSSearch(ctx, searchData, 0, 1)

	if err != nil {
		return "", nil
	}

	if len(results) > 0 {
		return results[0].AsciiArmoredPublicKey, nil
	}

	return "", errors.New("not found")
}

func operationIndex(options, searchData string, machineReadable, noModification, showFingerPrint, exactMatch bool) (string, error) {
	//hkpLog.Warn("Index(%s, %s, %v, %v, %v, %v) ==> Not Implemented", options, searchData, machineReadable, noModification, showFingerPrint, exactMatch)
	return "", errors.New("not implemented")
}

func operationVIndex(options, searchData string, machineReadable, noModification, showFingerPrint, exactMatch bool) (string, error) {
	//hkpLog.Warn("VIndex(%s, %s, %v, %v, %v, %v) ==> Not Implemented", options, searchData, machineReadable, noModification, showFingerPrint, exactMatch)
	return "", errors.New("not implemented")
}

func hkpLookup(log slog.Instance, w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log = wrapLogWithRequestID(log.SubScope("HKP"), r)

	q := r.URL.Query()
	op := q.Get("op")
	options := q.Get("options")
	mr := q.Get("mr") == "true" || q.Get("mr") == "1"
	nm := q.Get("nm") == "true" || q.Get("nm") == "1"
	fingerPrint := q.Get("fingerprint") == "on"
	exact := q.Get("exact") != ""
	search := q.Get("search")

	result := ""
	var err error

	switch op {
	case HKP.OperationGet:
		log.WithFields(map[string]interface{}{
			"options": options,
			"search":  search,
			"mr":      mr,
			"nm":      nm,
		}).Await("Running operation GET")
		result, err = operationGet(ctx, options, search, mr, nm)
	case HKP.OperationIndex:
		log.WithFields(map[string]interface{}{
			"options":     options,
			"search":      search,
			"mr":          mr,
			"nm":          nm,
			"fingerPrint": fingerPrint,
			"exact":       exact,
		}).Await("Running operation Index")
		result, err = operationIndex(options, search, mr, nm, fingerPrint, exact)
	case HKP.OperationVindex:
		log.WithFields(map[string]interface{}{
			"options":     options,
			"search":      search,
			"mr":          mr,
			"nm":          nm,
			"fingerPrint": fingerPrint,
			"exact":       exact,
		}).Await("Running operation Vindex")
		result, err = operationVIndex(options, search, mr, nm, fingerPrint, exact)
	}

	log.Done("Finished operation")

	if err != nil {
		if err.Error() == "not found" {
			CatchAllRouter(w, r, log)
			return
		}

		if err.Error() == "not implemented" {
			NotImplemented(w, r, log)
			return
		}

		InternalServerError("Internal Server Error", err.Error(), w, r, log)
		return
	}

	if result == "" {
		panic("Unknown operation")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

// HKP Add key godoc
// @id hkp-add
// @tags SKS
// @Summary GPG SKS Keyserver add public key
// @Accept plain
// @Produce plain
// @param publickey body string true "GPG Public Key"
// @Success 200 {string} result "OK"
// @Failure default {object} QuantoError.ErrorObject
// @Router /pks/add [post]
func hkpAdd(log slog.Instance, w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log = wrapLogWithRequestID(log.SubScope("HKP"), r)

	log.Await("Parsing Form Fields")
	err := r.ParseForm()
	log.Done("Parsed")
	if err != nil {
		InvalidFieldData("request", "Invalid request to add key. Please check if its HKP standard", w, r, log)
		return
	}

	key := r.Form.Get("keytext")
	log.Await("Adding key")
	result := keymagic.PKSAdd(ctx, key)
	log.Done("Key add result: %s", result)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

// AddHKPEndpoints attach the HKP /lookup and /add endpoints to the specified router with the specified log wrapped into the calls
func AddHKPEndpoints(log slog.Instance, dbHandler DatabaseHandler, r *mux.Router) {
	lookup := wrapRequestContextWithDatabaseHandler(dbHandler, hkpLookup)
	add := wrapRequestContextWithDatabaseHandler(dbHandler, hkpAdd)
	r.HandleFunc("/lookup", wrapWithLog(log, lookup)).Methods("GET")
	r.HandleFunc("/add", wrapWithLog(log, add)).Methods("POST")
}
