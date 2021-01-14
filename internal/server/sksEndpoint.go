package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/quan-to/chevron/internal/agent"

	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/pkg/interfaces"
	"github.com/quan-to/chevron/pkg/models"

	"github.com/gorilla/mux"
	"github.com/quan-to/slog"
)

type SKSEndpoint struct {
	sm  interfaces.SecretsManager
	gpg interfaces.PGPManager
	log slog.Instance
	dbh DatabaseHandler
}
type DatabaseHandler agent.DatabaseHandler

// MakeSKSEndpoint creates a handler for SKS Server Endpoint
func MakeSKSEndpoint(log slog.Instance, sm interfaces.SecretsManager, gpg interfaces.PGPManager, dbHandler DatabaseHandler) *SKSEndpoint {
	if log == nil {
		log = slog.Scope("SKS")
	} else {
		log = log.SubScope("SKS")
	}

	return &SKSEndpoint{
		sm:  sm,
		gpg: gpg,
		log: log,
		dbh: dbHandler,
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

// Get GPG Key godoc
// @id pks-get-key
// @tags Public Key Server, Key Store
// @Summary Fetches a GPG Public Key
// @Produce plain
// @param fingerPrint query string true "Fingerprint of the key you want to fetch"
// @Success 200 {string} result "gpg public key"
// @Failure default {object} QuantoError.ErrorObject
// @Router /sks/getKey [get]
func (sks *SKSEndpoint) getKey(w http.ResponseWriter, r *http.Request) {
	log := wrapLogWithRequestID(sks.log, r)
	InitHTTPTimer(log, r)
	ctx := wrapContextWithRequestID(r)
	ctx = wrapContextWithDatabaseHandler(sks.dbh, ctx)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	q := r.URL.Query()

	fingerPrint := q.Get("fingerPrint")
	key, _ := sks.gpg.GetPublicKeyASCII(ctx, fingerPrint)

	if key == "" {
		NotFound("fingerPrint", fmt.Sprintf("Key with fingerPrint %s was not found", fingerPrint), w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(key))
	LogExit(log, r, 200, n)
}

// Search GPG Key by Name godoc
// @id pks-search-by-name
// @tags Public Key Server, Key Store
// @Summary Searches for GPG Keys by its identifier name
// @Produce json
// @param name query string true "Name of the Key to be fetched"
// @param pageStart query int false "Pagination Start Index (default: 0)"
// @param pageEnd query int false "Pagination End Index (default: 100)"
// @Success 200 {object} []models.GPGKey
// @Failure default {object} QuantoError.ErrorObject
// @Router /sks/searchByName [get]
func (sks *SKSEndpoint) searchByName(w http.ResponseWriter, r *http.Request) {
	log := wrapLogWithRequestID(sks.log, r)
	InitHTTPTimer(log, r)
	ctx := wrapContextWithRequestID(r)
	ctx = wrapContextWithDatabaseHandler(sks.dbh, ctx)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	q := r.URL.Query()
	name := q.Get("name")
	pageStartS := q.Get("pageStart")
	pageEndS := q.Get("pageEnd")

	if name == "" {
		InvalidFieldData("name", "you should provide a name", w, r, log)
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

	gpgKeys, err := keymagic.PKSSearchByName(ctx, name, int(pageStart), int(pageEnd))

	if err != nil {
		NotFound("name", err.Error(), w, r, log)
		return
	}

	bodyData, err := json.Marshal(gpgKeys)

	if err != nil {
		InternalServerError("There was an internal server error. Please try again", nil, w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(log, r, 200, n)
}

// Search GPG Key by Fingerprint godoc
// @id pks-search-by-fingerprint
// @tags Public Key Server, Key Store
// @Summary Searches for GPG Keys by its fingerprint
// @Produce json
// @param fingerPrint query string true "Fingerprint to be fetched"
// @param pageStart query int false "Pagination Start Index (default: 0)"
// @param pageEnd query int false "Pagination End Index (default: 100)"
// @Success 200 {object} []models.GPGKey
// @Failure default {object} QuantoError.ErrorObject
// @Router /sks/searchByFingerPrint [get]
func (sks *SKSEndpoint) searchByFingerPrint(w http.ResponseWriter, r *http.Request) {
	log := wrapLogWithRequestID(sks.log, r)
	InitHTTPTimer(log, r)
	ctx := wrapContextWithRequestID(r)
	ctx = wrapContextWithDatabaseHandler(sks.dbh, ctx)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	q := r.URL.Query()
	fingerPrint := q.Get("fingerPrint")
	pageStartS := q.Get("pageStart")
	pageEndS := q.Get("pageEnd")

	if fingerPrint == "" {
		InvalidFieldData("fingerPrint", "you should provide a fingerPrint", w, r, log)
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

	gpgKeys, err := keymagic.PKSSearchByFingerPrint(ctx, fingerPrint, int(pageStart), int(pageEnd))

	if err != nil {
		NotFound("fingerPrint", err.Error(), w, r, log)
		return
	}

	bodyData, err := json.Marshal(gpgKeys)

	if err != nil {
		InternalServerError("There was an internal server error. Please try again", nil, w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(log, r, 200, n)
}

// Search GPG Key by Email godoc
// @id pks-search-by-email
// @tags Public Key Server, Key Store
// @Summary Searches for GPG Keys by its email
// @Produce json
// @param email query string true "Email of the Key to be fetched"
// @param pageStart query int false "Pagination Start Index (default: 0)"
// @param pageEnd query int false "Pagination End Index (default: 100)"
// @Success 200 {object} []models.GPGKey
// @Failure default {object} QuantoError.ErrorObject
// @Router /sks/searchByEmail [get]
func (sks *SKSEndpoint) searchByEmail(w http.ResponseWriter, r *http.Request) {
	log := wrapLogWithRequestID(sks.log, r)
	InitHTTPTimer(log, r)
	ctx := wrapContextWithRequestID(r)
	ctx = wrapContextWithDatabaseHandler(sks.dbh, ctx)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	q := r.URL.Query()
	email := q.Get("email")
	pageStartS := q.Get("pageStart")
	pageEndS := q.Get("pageEnd")

	if email == "" {
		InvalidFieldData("email", "you should provide a email", w, r, log)
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

	gpgKeys, err := keymagic.PKSSearchByEmail(ctx, email, int(pageStart), int(pageEnd))

	if err != nil {
		NotFound("email", err.Error(), w, r, log)
		return
	}

	bodyData, err := json.Marshal(gpgKeys)

	if err != nil {
		InternalServerError("There was an internal server error. Please try again", nil, w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(log, r, 200, n)
}

// Search GPG Key by Value godoc
// @id pks-search-by-value
// @tags Public Key Server, Key Store
// @Summary Searches for GPG Keys by any field
// @Produce json
// @param valueData query string true "Value of the Key to be fetched"
// @param pageStart query int false "Pagination Start Index (default: 0)"
// @param pageEnd query int false "Pagination End Index (default: 100)"
// @Success 200 {object} []models.GPGKey
// @Failure default {object} QuantoError.ErrorObject
// @Router /sks/search [get]
func (sks *SKSEndpoint) search(w http.ResponseWriter, r *http.Request) {
	log := wrapLogWithRequestID(sks.log, r)
	InitHTTPTimer(log, r)
	ctx := wrapContextWithRequestID(r)
	ctx = wrapContextWithDatabaseHandler(sks.dbh, ctx)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	q := r.URL.Query()
	valueData := q.Get("valueData")
	pageStartS := q.Get("pageStart")
	pageEndS := q.Get("pageEnd")

	if valueData == "" {
		InvalidFieldData("email", "you should provide a valueData", w, r, log)
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

	gpgKeys, err := keymagic.PKSSearch(ctx, valueData, int(pageStart), int(pageEnd))

	if err != nil {
		NotFound("valueData", err.Error(), w, r, log)
		return
	}

	bodyData, err := json.Marshal(gpgKeys)

	if err != nil {
		InternalServerError("There was an internal server error. Please try again", nil, w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(log, r, 200, n)
}

// Add Public Key godoc
// @id pks-add-public-key
// @tags Public Key Server, Key Store
// @Summary Adds a GPG Public Key
// @Accept json
// @Produce plain
// @Param message body models.SKSAddKey true "GPG Public Key in an Armored format"
// @Success 200 {string} result "OK"
// @Failure default {object} QuantoError.ErrorObject
// @Router /sks/addKey [post]
func (sks *SKSEndpoint) addKey(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(sks.log, r)
	InitHTTPTimer(log, r)
	ctx = wrapContextWithDatabaseHandler(sks.dbh, ctx)

	var data models.SKSAddKey

	if !UnmarshalBodyOrDie(&data, w, r, log) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	status := keymagic.PKSAdd(ctx, data.PublicKey)

	if status != "OK" {
		InvalidFieldData("PublicKey", "Invalid Public Key specified. Check if its in ASCII Armored Format", w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte("OK"))
	LogExit(log, r, 200, n)
}
