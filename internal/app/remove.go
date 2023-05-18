package app

import (
	"context"
	"database/sql"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	desc "item-service/pkg/gen/go/api/service/v1"
	"log"
)

func (m *ItemServiceServer) Remove(ctx context.Context, req *desc.RemoveRequest) (*desc.RemoveResponse, error) {
	err := m.itemService.Remove(req.Id, req.CampaignId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "404"))
			return nil, status.Error(3, "errors.item.notFound")
		}

		log.Println("remove error: ", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &desc.RemoveResponse{
		Id:         req.Id,
		CampaignId: req.CampaignId,
		Removed:    true,
	}, nil
}
