package app

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	md "item-service/pkg/gen/go/api/model/v1"
	desc "item-service/pkg/gen/go/api/service/v1"
	"log"
)

func (m *ItemServiceServer) List(ctx context.Context, req *desc.ListRequest) (*desc.ListResponse, error) {
	list, err := m.itemService.List()
	if err != nil {
		log.Println("list error: ", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	items := make([]*md.Item, len(list))
	for i, v := range list {
		items[i] = &md.Item{
			Id:          v.Id,
			CampaignId:  v.CampaignId,
			Name:        v.Name,
			Description: v.Description,
			Priority:    v.Priority,
			Removed:     v.Removed,
			CreatedAt:   v.CreatedAt.Unix(),
		}
	}

	return &desc.ListResponse{Items: items}, nil
}
