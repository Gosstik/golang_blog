package images_service

import (
	"context"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
)

func (h *Handler) ImagesV1Upload(ctx context.Context, req *api_proto.V1ImagesUploadRequest) (*api_proto.V1ImagesUploadResponse, error) {
	return &api_proto.V1ImagesUploadResponse{
		ImageUuid: "image-uuid-mock-12345",
		ImageUrl:  "https://example.com/images/uploaded-image-12345.png",
	}, nil
}
