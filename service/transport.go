package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	ErrBadRouting = errors.New("Bad Routing")
)

func MakeHTTPRouter(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoint(s)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/api/puzzle").
		Queries(
			"width", "{width:[0-9]+}",
			"height", "{height:[0-9]+}",
			"words", "{words:[0-9]+}",
		).
		Handler(
			httptransport.NewServer(
				e.GeneratePuzzleEndpoint,
				decodeGeneratePuzzleRequest,
				encodeResponse,
				options...,
			),
		)
	return r

}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	return json.NewEncoder(w).Encode(response)
}
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	// case ErrNotFound:
	// 	return http.StatusNotFound
	case ErrBadRouting:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}

}
func decodeGeneratePuzzleRequest(_ context.Context, r *http.Request) (request interface{}, err error) {

	vars := mux.Vars(r)
	wStr, ok := vars["width"]
	if !ok {
		return nil, ErrBadRouting
	}
	w, e := strconv.Atoi(wStr)
	if e != nil {
		return nil, ErrBadRouting
	}
	hStr, ok := vars["height"]
	if !ok {
		return nil, ErrBadRouting
	}
	h, err := strconv.Atoi(hStr)
	if err != nil {
		return nil, ErrBadRouting
	}
	wordsStr, ok := vars["words"]
	if !ok {
		return nil, ErrBadRouting
	}
	words, err := strconv.Atoi(wordsStr)
	return generatePuzzleRequest{Height: h, Width: w, NumberOfWords: words}, nil
}
