package app

import (
	"context"
	"database/sql"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	md "item-service/pkg/gen/go/api/model/v1"
	desc "item-service/pkg/gen/go/api/service/v1"
	"log"
)

func (m *ItemServiceServer) Update(ctx context.Context, req *desc.UpdateRequest) (*md.Item, error) {

	update, err := m.itemService.Update(req.Id, req.CampaignId, req.Payload.Name, req.Payload.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "404"))
			return nil, status.Error(3, "errors.item.notFound")
		}
		log.Println("update error: ", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &md.Item{
		Id:          update.Id,
		CampaignId:  update.CampaignId,
		Name:        update.Name,
		Description: update.Description,
		Priority:    update.Priority,
		Removed:     update.Removed,
		CreatedAt:   update.CreatedAt.Unix(),
	}, nil
}
