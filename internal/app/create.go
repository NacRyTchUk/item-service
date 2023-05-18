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

func (m *ItemServiceServer) Create(ctx context.Context, req *desc.CreateRequest) (*md.Item, error) {
	create, err := m.itemService.Create(req.CampaignId, req.Payload.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "404"))
			return nil, status.Error(3, "errors.campaigns.notFound")
		}
		log.Println("create error: ", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &md.Item{
		Id:          create.Id,
		CampaignId:  create.CampaignId,
		Name:        create.Name,
		Description: create.Description,
		Priority:    create.Priority,
		Removed:     create.Removed,
		CreatedAt:   create.CreatedAt.Unix(),
	}, nil
}
