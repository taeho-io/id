package handler

import (
	tid "github.com/taeho-io/go-taeho/id"
	"github.com/taeho-io/id"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NewHandlerFunc func(ctx context.Context, req *id.NewRequest) (*id.NewResponse, error)

func New(tid tid.ID) NewHandlerFunc {
	return func(ctx context.Context, req *id.NewRequest) (*id.NewResponse, error) {
		newID, err := tid.Generate()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &id.NewResponse{
			Id: newID,
		}, nil
	}
}
