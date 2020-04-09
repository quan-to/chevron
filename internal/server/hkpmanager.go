package server

import (
	"context"
	"errors"
	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/internal/models/HKP"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/quan-to/slog"
)

/// HKP Server based on https://tools.ietf.org/html/draft-shaw-openpgp-hkp-00

func operationGet(ctx context.Context, options, searchData string, machineReadable, noModification bool) (error, string) {
	if searchData[:2] == "0x" {
		k, _ := keymagic.PKSGetKey(ctx, searchData[2:])
		if k == "" {
			return errors.New("not found"), ""
		}
		return nil, k
	}

	results, err := keymagic.PKSSearch(searchData, 0, 1)

	if err != nil {
		return err, ""
	}

	if len(results) > 0 {
		return nil, results[0].AsciiArmoredPublicKey
	}

	return errors.New("not found"), ""
}

func operationIndex(options, searchData string, machineReadable, noModification, showFingerPrint, exactMatch bool) (error, string) {
	//hkpLog.Warn("Index(%s, %s, %v, %v, %v, %v) ==> Not Implemented", options, searchData, machineReadable, noModification, showFingerPrint, exactMatch)
	return errors.New("not implemented"), ""
}

func operationVIndex(options, searchData string, machineReadable, noModification, showFingerPrint, exactMatch bool) (error, string) {
	//hkpLog.Warn("VIndex(%s, %s, %v, %v, %v, %v) ==> Not Implemented", options, searchData, machineReadable, noModification, showFingerPrint, exactMatch)
	return errors.New("not implemented"), ""
}

func hkpLookup(log slog.Instance, w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log = wrapLogWithRequestID(log.SubScope("HKP"), r)

	InitHTTPTimer(log, r)
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
		err, result = operationGet(ctx, options, search, mr, nm)
	case HKP.OperationIndex:
		log.WithFields(map[string]interface{}{
			"options":     options,
			"search":      search,
			"mr":          mr,
			"nm":          nm,
			"fingerPrint": fingerPrint,
			"exact":       exact,
		}).Await("Running operation Index")
		err, result = operationIndex(options, search, mr, nm, fingerPrint, exact)
	case HKP.OperationVindex:
		log.WithFields(map[string]interface{}{
			"options":     options,
			"search":      search,
			"mr":          mr,
			"nm":          nm,
			"fingerPrint": fingerPrint,
			"exact":       exact,
		}).Await("Running operation Vindex")
		err, result = operationVIndex(options, search, mr, nm, fingerPrint, exact)
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
	_, _ = w.Write([]byte(result))
	LogExit(log, r, http.StatusOK, len(result))
}

func hkpAdd(log slog.Instance, w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log = wrapLogWithRequestID(log.SubScope("HKP"), r)

	InitHTTPTimer(log, r)
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
	_, _ = w.Write([]byte(result))
	LogExit(log, r, http.StatusOK, len(result))
}

// AddHKPEndpoints attach the HKP /lookup and /add endpoints to the specified router with the specified log wrapped into the calls
func AddHKPEndpoints(log slog.Instance, r *mux.Router) {
	r.HandleFunc("/lookup", wrapWithLog(log, hkpLookup)).Methods("GET")
	r.HandleFunc("/add", wrapWithLog(log, hkpAdd)).Methods("POST")
}
