package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type PostCache interface {
	SetFirstPage(ctx context.Context, postId uint, post []domain.Post) error    // 设置个人用户的第一页帖子缓存
	SetPubFirstPage(ctx context.Context, postId uint, post []domain.Post) error // 设置公开用户的第一页帖子缓存
	DelFirstPage(ctx context.Context, postId uint) error                        // 删除用户的第一页帖子缓存
	GetDetail(ctx context.Context, postId uint) (domain.Post, error)            // 根据ID获取一个帖子详情缓存
	SetDetail(ctx context.Context, post domain.Post) error                      // 设置一个帖子详情缓存
	GetPubDetail(ctx context.Context, postId uint) (domain.Post, error)         // 根据ID获取一个已发布的帖子详情缓存
	SetPubDetail(ctx context.Context, post domain.Post) error                   // 设置一个已发布的帖子详情缓存
	DelPubFirstPage(ctx context.Context, postId uint) error
}

type postCache struct {
	cmd redis.Cmdable
	l   *zap.Logger
}

func NewPostCache(cmd redis.Cmdable, l *zap.Logger) PostCache {
	return &postCache{
		cmd: cmd,
		l:   l,
	}
}

// SetFirstPage 设置个人第一页帖子摘要
func (p *postCache) SetFirstPage(ctx context.Context, postId uint, post []domain.Post) error {
	for i := 0; i < len(post); i++ {
		post[i].Content = post[i].Abstract()
	}
	val, err := json.Marshal(post)
	if err != nil {
		p.l.Error("序列化失败", zap.Error(err))
		return err
	}
	key := fmt.Sprintf("post:first:%d", postId)
	return p.cmd.Set(ctx, key, val, time.Minute*15).Err()
}

// GetPubFirstPage 获取第一页公开帖子摘要
func (p *postCache) GetPubFirstPage(ctx context.Context, postId uint) ([]domain.Post, error) {
	var dp []domain.Post
	key := fmt.Sprintf("post:pub:first:%d", postId)
	val, err := p.cmd.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		p.l.Warn("缓存未命中", zap.String("key", key))
		return nil, nil
	} else if err != nil {
		p.l.Warn("缓存获取失败", zap.Error(err), zap.String("key", key))
		return nil, err
	}
	if er := json.Unmarshal(val, &dp); er != nil {
		p.l.Error("反序列化失败", zap.Error(er), zap.String("key", key))
		return nil, er
	}
	return dp, nil
}

// SetPubFirstPage 设置第一页公开帖子摘要
func (p *postCache) SetPubFirstPage(ctx context.Context, postId uint, post []domain.Post) error {
	for i := 0; i < len(post); i++ {
		post[i].Content = post[i].Abstract()
	}
	val, err := json.Marshal(post)
	if err != nil {
		p.l.Error("序列化失败", zap.Error(err))
		return err
	}
	key := fmt.Sprintf("post:pub:first:%d", postId)
	return p.cmd.Set(ctx, key, val, time.Minute*15).Err()
}

// DelFirstPage 删除第一页帖子摘要
func (p *postCache) DelFirstPage(ctx context.Context, postId uint) error {
	key1 := fmt.Sprintf("post:first:%d", postId)
	key2 := fmt.Sprintf("post:pub:first:%d", postId)
	err := p.cmd.Del(ctx, key1).Err()
	if err != nil {
		p.l.Error("delete cache failed", zap.Error(err))
		return err
	}
	er := p.cmd.Del(ctx, key2).Err()
	if er != nil {
		p.l.Error("delete cache failed", zap.Error(er))
		return er
	}
	return nil
}

// DelPubFirstPage DelPunFirstPage 删除第一页帖子摘要
func (p *postCache) DelPubFirstPage(ctx context.Context, postId uint) error {
	key := fmt.Sprintf("post:pub:first:%d", postId)
	return p.cmd.Del(ctx, key).Err()
}

// GetDetail 获取帖子详情缓存
func (p *postCache) GetDetail(ctx context.Context, postId uint) (domain.Post, error) {
	var dp domain.Post
	key := fmt.Sprintf("post:detail:%d", postId)
	val, err := p.cmd.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		p.l.Warn("缓存未命中", zap.String("key", key))
		return dp, nil
	} else if err != nil {
		p.l.Warn("缓存获取失败", zap.Error(err), zap.String("key", key))
		return dp, err
	}
	if er := json.Unmarshal(val, &dp); er != nil {
		p.l.Error("反序列化失败", zap.Error(er), zap.String("key", key))
		return dp, er
	}
	return dp, nil
}

// SetDetail 设置帖子详情缓存
func (p *postCache) SetDetail(ctx context.Context, post domain.Post) error {
	val, err := json.Marshal(post)
	if err != nil {
		p.l.Error("序列化失败", zap.Error(err))
		return err
	}
	key := fmt.Sprintf("post:detail:%d", post.ID)
	return p.cmd.Set(ctx, key, val, time.Minute*15).Err()
}

// GetPubDetail 获取公开帖子详情缓存
func (p *postCache) GetPubDetail(ctx context.Context, postId uint) (domain.Post, error) {
	var dp domain.Post
	key := fmt.Sprintf("post:pub:detail:%d", postId)
	val, err := p.cmd.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		p.l.Warn("缓存未命中", zap.String("key", key))
		return dp, nil
	} else if err != nil {
		p.l.Warn("缓存获取失败", zap.Error(err), zap.String("key", key))
		return dp, err
	}
	if er := json.Unmarshal(val, &dp); er != nil {
		p.l.Error("反序列化失败", zap.Error(er), zap.String("key", key))
		return dp, er
	}
	return dp, nil
}

// SetPubDetail 设置公开帖子详情缓存
func (p *postCache) SetPubDetail(ctx context.Context, post domain.Post) error {
	key := fmt.Sprintf("post:pub:detail:%d", post.ID)
	val, err := json.Marshal(post)
	if err != nil {
		p.l.Error("序列化失败", zap.Error(err))
		return err
	}
	return p.cmd.Set(ctx, key, val, time.Minute*15).Err()
}
