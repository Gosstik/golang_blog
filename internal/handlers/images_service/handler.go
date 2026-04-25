package images_service

import (
	api_proto "github.com/Gosstik/golang_blog/api/proto"
)

type Handler struct {
	api_proto.UnimplementedImagesServiceServer
}

func NewHandler() *Handler {
	return &Handler{}
}
