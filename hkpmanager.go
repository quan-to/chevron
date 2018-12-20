package remote_signer

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/quan-to/remote-signer/SLog"
	"github.com/quan-to/remote-signer/models/HKP"
	"net/http"
)

/// HKP Server based on https://tools.ietf.org/html/draft-shaw-openpgp-hkp-00

var hkpLog = SLog.Scope("HKP")

func operationGet(options, searchData string, machineReadable, noModification bool) (error, string) {
	hkpLog.Info("GET(%s, %s, %v, %v)", options, searchData, machineReadable, noModification)
	if searchData[:2] == "0x" {
		var k = PKSGetKey(searchData[2:])
		if k == "" {
			return errors.New("not found"), ""
		}
		return nil, k
	}

	results := PKSSearch(searchData, 0, 1)

	if len(results) > 0 {
		return nil, results[0].AsciiArmoredPublicKey
	}

	return errors.New("not found"), ""
}

func operationIndex(options, searchData string, machineReadable, noModification, showFingerPrint, exactMatch bool) (error, string) {
	hkpLog.Warn("Index(%s, %s, %v, %v, %v, %v) ==> Not Implemented", options, searchData, machineReadable, noModification, showFingerPrint, exactMatch)
	return errors.New("not implemented"), ""
}

func operationVIndex(options, searchData string, machineReadable, noModification, showFingerPrint, exactMatch bool) (error, string) {
	hkpLog.Warn("VIndex(%s, %s, %v, %v, %v, %v) ==> Not Implemented", options, searchData, machineReadable, noModification, showFingerPrint, exactMatch)
	return errors.New("not implemented"), ""
}

func hkpLookup(w http.ResponseWriter, r *http.Request) {
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
		err, result = operationGet(options, search, mr, nm)
	case HKP.OperationIndex:
		err, result = operationIndex(options, search, mr, nm, fingerPrint, exact)
	case HKP.OperationVindex:
		err, result = operationVIndex(options, search, mr, nm, fingerPrint, exact)
	}

	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))

		return
	}

	if result == "" {
		panic("Unknown operation")
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result))
}

func hkpAdd(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func AddHKPEndpoints(r *mux.Router) {
	r.HandleFunc("/pks/lookup", hkpLookup).Methods("GET")
	r.HandleFunc("/pks/add", hkpAdd).Methods("POST")
}
