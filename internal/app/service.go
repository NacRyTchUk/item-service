package app

import (
	"item-service/internal/service"
	proto_item_service "item-service/pkg/gen/go/api/service/v1"
)

type ItemServiceServer struct {
	proto_item_service.UnimplementedItemServiceServer
	itemService service.ItemService
}

func NewItemServiceServer(
	itemService service.ItemService,
) *ItemServiceServer {
	return &ItemServiceServer{
		itemService: itemService,
	}
}
