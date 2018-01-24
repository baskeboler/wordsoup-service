package service

import (
	"context"
	"strings"

	"github.com/baskeboler/wordsoup"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints facade struct with all endpoints
type Endpoints struct {
	GeneratePuzzleEndpoint endpoint.Endpoint
}

// MakeServerEndpoint builds facade struct
func MakeServerEndpoint(s Service) Endpoints {
	return Endpoints{
		GeneratePuzzleEndpoint: MakeGeneratePuzzleEndpoint(s),
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
