package handler

import (
	tid "github.com/taeho-io/go-taeho/id"
	"github.com/taeho-io/id"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateHandlerFunc func(ctx context.Context, req *id.CreateRequest) (*id.CreateResponse, error)

func Create(tid tid.ID) CreateHandlerFunc {
	return func(ctx context.Context, req *id.CreateRequest) (*id.CreateResponse, error) {
		newID, err := tid.Generate()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &id.CreateResponse{
			Id: newID,
		}, nil
	}
}
