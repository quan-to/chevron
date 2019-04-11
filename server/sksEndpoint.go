package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/quan-to/chevron/etc"
	"github.com/quan-to/chevron/keymagic"
	"github.com/quan-to/chevron/models"
	"github.com/quan-to/slog"
	"net/http"
	"strconv"
)

var sksLog = slog.Scope("SKS Endpoint")

type SKSEndpoint struct {
	sm  etc.SMInterface
	gpg etc.PGPInterface
}

func MakeSKSEndpoint(sm etc.SMInterface, gpg etc.PGPInterface) *SKSEndpoint {
	return &SKSEndpoint{
		sm:  sm,
		gpg: gpg,
	}
}

func (sks *SKSEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/getKey", sks.getKey).Methods("GET")
	r.HandleFunc("/searchByName", sks.searchByName).Methods("GET")
	r.HandleFunc("/searchByFingerPrint", sks.searchByFingerPrint).Methods("GET")
	r.HandleFunc("/searchByEmail", sks.searchByEmail).Methods("GET")
	r.HandleFunc("/search", sks.search).Methods("GET")
	r.HandleFunc("/addKey", sks.addKey).Methods("POST")
}

func (sks *SKSEndpoint) getKey(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, geLog)
		}
	}()

	q := r.URL.Query()

	fingerPrint := q.Get("fingerPrint")
	key, _ := sks.gpg.GetPublicKeyAscii(fingerPrint)

	if key == "" {
		NotFound("fingerPrint", fmt.Sprintf("Key with fingerPrint %s was not found", fingerPrint), w, r, sksLog)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(key))
	LogExit(sksLog, r, 200, n)
}

func (sks *SKSEndpoint) searchByName(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, geLog)
		}
	}()

	q := r.URL.Query()
	name := q.Get("name")
	pageStartS := q.Get("pageStart")
	pageEndS := q.Get("pageEnd")

	if name == "" {
		InvalidFieldData("name", "you should provide a name", w, r, sksLog)
		return
	}

	pageStart, err := strconv.ParseInt(pageStartS, 10, 32)
	if err != nil {
		pageStart = models.DefaultPageStart
	}

	pageEnd, err := strconv.ParseInt(pageEndS, 10, 32)
	if err != nil {
		pageEnd = models.DefaultPageEnd
	}

	gpgKeys, err := keymagic.PKSSearchByName(name, int(pageStart), int(pageEnd))

	if err != nil {
		NotFound("name", err.Error(), w, r, sksLog)
		return
	}

	bodyData, err := json.Marshal(gpgKeys)

	if err != nil {
		InternalServerError("There was an internal server error. Please try again", nil, w, r, sksLog)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(sksLog, r, 200, n)
}

func (sks *SKSEndpoint) searchByFingerPrint(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, geLog)
		}
	}()

	q := r.URL.Query()
	fingerPrint := q.Get("fingerPrint")
	pageStartS := q.Get("pageStart")
	pageEndS := q.Get("pageEnd")

	if fingerPrint == "" {
		InvalidFieldData("fingerPrint", "you should provide a fingerPrint", w, r, sksLog)
		return
	}

	pageStart, err := strconv.ParseInt(pageStartS, 10, 32)
	if err != nil {
		pageStart = models.DefaultPageStart
	}

	pageEnd, err := strconv.ParseInt(pageEndS, 10, 32)
	if err != nil {
		pageEnd = models.DefaultPageEnd
	}

	gpgKeys, err := keymagic.PKSSearchByFingerPrint(fingerPrint, int(pageStart), int(pageEnd))

	if err != nil {
		NotFound("fingerPrint", err.Error(), w, r, sksLog)
		return
	}

	bodyData, err := json.Marshal(gpgKeys)

	if err != nil {
		InternalServerError("There was an internal server error. Please try again", nil, w, r, sksLog)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(sksLog, r, 200, n)
}

func (sks *SKSEndpoint) searchByEmail(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, geLog)
		}
	}()

	q := r.URL.Query()
	email := q.Get("email")
	pageStartS := q.Get("pageStart")
	pageEndS := q.Get("pageEnd")

	if email == "" {
		InvalidFieldData("email", "you should provide a email", w, r, sksLog)
		return
	}

	pageStart, err := strconv.ParseInt(pageStartS, 10, 32)
	if err != nil {
		pageStart = models.DefaultPageStart
	}

	pageEnd, err := strconv.ParseInt(pageEndS, 10, 32)
	if err != nil {
		pageEnd = models.DefaultPageEnd
	}

	gpgKeys, err := keymagic.PKSSearchByEmail(email, int(pageStart), int(pageEnd))

	if err != nil {
		NotFound("email", err.Error(), w, r, sksLog)
		return
	}

	bodyData, err := json.Marshal(gpgKeys)

	if err != nil {
		InternalServerError("There was an internal server error. Please try again", nil, w, r, sksLog)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(sksLog, r, 200, n)
}

func (sks *SKSEndpoint) search(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, sksLog)
		}
	}()

	q := r.URL.Query()
	valueData := q.Get("valueData")
	pageStartS := q.Get("pageStart")
	pageEndS := q.Get("pageEnd")

	if valueData == "" {
		InvalidFieldData("email", "you should provide a valueData", w, r, sksLog)
		return
	}

	pageStart, err := strconv.ParseInt(pageStartS, 10, 32)
	if err != nil {
		pageStart = models.DefaultPageStart
	}

	pageEnd, err := strconv.ParseInt(pageEndS, 10, 32)
	if err != nil {
		pageEnd = models.DefaultPageEnd
	}

	gpgKeys, err := keymagic.PKSSearch(valueData, int(pageStart), int(pageEnd))

	if err != nil {
		NotFound("valueData", err.Error(), w, r, sksLog)
		return
	}

	bodyData, err := json.Marshal(gpgKeys)

	if err != nil {
		InternalServerError("There was an internal server error. Please try again", nil, w, r, sksLog)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(sksLog, r, 200, n)
}

func (sks *SKSEndpoint) addKey(w http.ResponseWriter, r *http.Request) {
	InitHTTPTimer(r)

	var data models.SKSAddKey

	if !UnmarshalBodyOrDie(&data, w, r, sksLog) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, sksLog)
		}
	}()

	status := keymagic.PKSAdd(data.PublicKey)

	if status != "OK" {
		InvalidFieldData("PublicKey", "Invalid Public Key specified. Check if its in ASCII Armored Format", w, r, sksLog)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte("OK"))
	LogExit(sksLog, r, 200, n)
}
