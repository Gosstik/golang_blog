package posts_service

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
)

// NOTE: for hardcoded responses, will be removed in task 2

var mockUsers = []*api_proto.User{
	{
		Uuid:      "user-uuid-1",
		Nickname:  "john_doe",
		Name:      "John",
		Surname:   "Doe",
		AvatarUrl: "https://example.com/avatars/john.png",
	},
	{
		Uuid:      "user-uuid-2",
		Nickname:  "jane_smith",
		Name:      "Jane",
		Surname:   "Smith",
		AvatarUrl: "https://example.com/avatars/jane.png",
	},
	{
		Uuid:      "user-uuid-3",
		Nickname:  "bob_wilson",
		Name:      "Bob",
		Surname:   "Wilson",
		AvatarUrl: "https://example.com/avatars/bob.png",
	},
}

var mockPosts = []*api_proto.BlogPost{
	{
		Uuid:             "post-uuid-1",
		Author:           mockUsers[0],
		CreatedAt:        timestamppb.New(time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC)),
		UpdatedAt:        nil,
		LikesCount:       5,
		ContentText:      "Hello everyone! This is my first blog post about Go programming.",
		ContentImageUrls: []string{"https://example.com/images/go-gopher.png"},
		LikedByMe:        false,
	},
	{
		Uuid:             "post-uuid-2",
		Author:           mockUsers[1],
		CreatedAt:        timestamppb.New(time.Date(2026, 4, 19, 15, 30, 0, 0, time.UTC)),
		UpdatedAt:        timestamppb.New(time.Date(2026, 4, 19, 16, 0, 0, 0, time.UTC)),
		LikesCount:       12,
		ContentText:      "Just finished reading a great book on microservices architecture!",
		ContentImageUrls: nil,
		LikedByMe:        true,
	},
	{
		Uuid:        "post-uuid-3",
		Author:      mockUsers[2],
		CreatedAt:   timestamppb.New(time.Date(2026, 4, 18, 8, 0, 0, 0, time.UTC)),
		UpdatedAt:   nil,
		LikesCount:  3,
		ContentText: "Check out this beautiful sunset I captured!",
		ContentImageUrls: []string{
			"https://example.com/images/sunset1.png",
			"https://example.com/images/sunset2.png",
		},
		LikedByMe: false,
	},
}

func findUserByUuid(uuid string) (*api_proto.User, error) {
	for _, u := range mockUsers {
		if u.Uuid == uuid {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user with uuid %q not found", uuid)
}
