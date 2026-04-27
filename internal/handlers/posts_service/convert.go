package posts_service

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
	"github.com/Gosstik/golang_blog/internal/entities"
)

func userToProto(u *entities.User) *api_proto.User {
	return &api_proto.User{
		Uuid:      u.UUID,
		Nickname:  u.Nickname,
		Name:      u.Name,
		Surname:   u.Surname,
		AvatarUrl: u.AvatarUrl,
	}
}

func postToProto(p *entities.BlogPost, likesCount int64, likedByMe bool) *api_proto.BlogPost {
	pb := &api_proto.BlogPost{
		Uuid:             p.UUID,
		Author:           userToProto(&p.Author),
		CreatedAt:        timestamppb.New(p.CreatedAt),
		LikesCount:       likesCount,
		ContentText:      p.ContentText,
		ContentImageUrls: []string(p.ContentImageUrls),
		LikedByMe:        likedByMe,
	}
	if p.UpdatedAt != nil {
		pb.UpdatedAt = timestamppb.New(*p.UpdatedAt)
	}
	return pb
}
