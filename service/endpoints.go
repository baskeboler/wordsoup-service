package service

import (
	"context"
	"strings"

	"github.com/baskeboler/wordsoup"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GeneratePuzzleEndpoint endpoint.Endpoint
}

func MakeServerEndpoint(s Service) Endpoints {
	return Endpoints{
		GeneratePuzzleEndpoint: MakeGeneratePuzzleEndpoint(s),
	}
}

func MakeGeneratePuzzleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// var ws interface{}
		req := request.(generatePuzzleRequest)
		// if !ok {
		// 	return nil, ErrBadRouting
		// }

		ws, e := s.GeneratePuzzle(ctx, req.Width, req.Height, req.NumberOfWords)
		// ws2, ok := ws.(*wordsoup.WordSoup)
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
