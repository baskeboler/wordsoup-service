package service

import (
	"context"
	"net/http"
	"strings"

	"github.com/baskeboler/auth"
	"github.com/baskeboler/wordsoup"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Endpoints facade struct with all endpoints
type Endpoints struct {
	GeneratePuzzleEndpoint endpoint.Endpoint
	LoginEndpoint          endpoint.Endpoint
}

// MakeServerEndpoint builds facade struct
func MakeServerEndpoint(s Service) Endpoints {
	return Endpoints{
		GeneratePuzzleEndpoint: MakeGeneratePuzzleEndpoint(s),
		LoginEndpoint:          MakeLoginEndpoint(s),
	}
}

// MakeGeneratePuzzleEndpoint builds the GeneratePuzzle endpoint
func MakeGeneratePuzzleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(generatePuzzleRequest)

		ws, e := s.GeneratePuzzle(ctx, req.Width, req.Height, req.NumberOfWords)
		if e != nil {
			return nil, ErrBadRouting
		}

		return generatePuzzleResponse{
			Puzzle: ws,
			Err:    e,
			Rows:   strings.Split(ws.String(), "\n"),
		}, nil
	}
}
func MakeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loginRequest)

		k, e := s.Login(ctx, req.Name, req.Password)
		return loginResponse{Key: k, Err: e}, e
	}
}

// httptransport.RequestFunc
func AuthEndpoint() httptransport.RequestFunc {
	m := auth.New()
	return func(ctx context.Context, request *http.Request) context.Context {
		ctx2 := httptransport.PopulateRequestContext(ctx, request)
		key := ctx2.Value(httptransport.ContextKeyRequestAuthorization)
		err := m.ValidateKey(key.(string))
		if err != nil {
			ctx2 := context.WithValue(ctx2, "AuthFailed", true)

			return ctx2
		}
		return ctx2
	}
}

type generatePuzzleRequest struct {
	Width, Height, NumberOfWords int
}

type generatePuzzleResponse struct {
	Puzzle *wordsoup.WordSoup `json:"puzzle,omitempty"`
	Rows   []string           `json:"rows,omitempty"`
	Err    error              `json:"err,omitempty"`
}

func (r generatePuzzleResponse) error() error {
	return r.Err
}

type loginRequest struct {
	Name     string
	Password string
	Err      error
}

type loginResponse struct {
	Key *auth.LoginKey `json:"key,omitempty"`
	Err error          `json:"err,omitempty"`
}

func (r loginResponse) error() error {
	return r.Err
}
