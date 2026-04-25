package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// LikesRepository uses Redis Sets to store likes.
// Key: "post:{postUUID}:likes", members: user UUIDs.
// SCARD gives total count, SISMEMBER checks if user liked.
type LikesRepository struct {
	client *redis.Client
}

func NewLikesRepository(client *redis.Client) *LikesRepository {
	return &LikesRepository{client: client}
}

func likesKey(postUUID string) string {
	return fmt.Sprintf("post:%s:likes", postUUID)
}

func (r *LikesRepository) Like(ctx context.Context, postUUID, userUUID string) (int64, error) {
	if err := r.client.SAdd(ctx, likesKey(postUUID), userUUID).Err(); err != nil {
		return 0, err
	}
	return r.client.SCard(ctx, likesKey(postUUID)).Result()
}

func (r *LikesRepository) Unlike(ctx context.Context, postUUID, userUUID string) (int64, error) {
	if err := r.client.SRem(ctx, likesKey(postUUID), userUUID).Err(); err != nil {
		return 0, err
	}
	return r.client.SCard(ctx, likesKey(postUUID)).Result()
}

func (r *LikesRepository) GetLikesCount(ctx context.Context, postUUID string) (int64, error) {
	return r.client.SCard(ctx, likesKey(postUUID)).Result()
}

func (r *LikesRepository) IsLikedByUser(ctx context.Context, postUUID, userUUID string) (bool, error) {
	return r.client.SIsMember(ctx, likesKey(postUUID), userUUID).Result()
}

func (r *LikesRepository) GetLikedUserUUIDs(ctx context.Context, postUUID string, limit int, cursor int64) ([]string, int64, error) {
	keys, nextCursor, err := r.client.SScan(ctx, likesKey(postUUID), uint64(cursor), "", int64(limit)).Result()
	if err != nil {
		return nil, 0, err
	}
	return keys, int64(nextCursor), nil
}

func (r *LikesRepository) GetLikesCountBatch(ctx context.Context, postUUIDs []string) (map[string]int64, error) {
	pipe := r.client.Pipeline()
	cmds := make(map[string]*redis.IntCmd, len(postUUIDs))

	for _, postUUID := range postUUIDs {
		cmds[postUUID] = pipe.SCard(ctx, likesKey(postUUID))
	}

	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return nil, err
	}

	result := make(map[string]int64, len(postUUIDs))
	for postUUID, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		result[postUUID] = val
	}
	return result, nil
}

func (r *LikesRepository) IsLikedByUserBatch(ctx context.Context, postUUIDs []string, userUUID string) (map[string]bool, error) {
	pipe := r.client.Pipeline()
	cmds := make(map[string]*redis.BoolCmd, len(postUUIDs))

	for _, postUUID := range postUUIDs {
		cmds[postUUID] = pipe.SIsMember(ctx, likesKey(postUUID), userUUID)
	}

	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return nil, err
	}

	result := make(map[string]bool, len(postUUIDs))
	for postUUID, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		result[postUUID] = val
	}
	return result, nil
}
