package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrAlredyAdded        = errors.New("product already added")
	ErrNoCartProductFound = errors.New("no cart product found")
)

type RedisRepository struct {
	rds     *redis.Client
	cartTTL time.Duration
}

func NewRedisRepository(rds *redis.Client, cartTTL time.Duration) *RedisRepository {
	return &RedisRepository{rds: rds, cartTTL: cartTTL}
}

func (r RedisRepository) AddToCart(ctx context.Context, UUID, productID string, stock int) error {
	const op = "repository.AddToCart"

	isAdded, err := r.IsAdded(ctx, UUID, productID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if isAdded {
		return ErrAlredyAdded
	}

	err = r.rds.HSet(ctx, UUID, productID, stock).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	err = r.rds.Expire(ctx, UUID, r.cartTTL).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r RedisRepository) DeleteFromCart(ctx context.Context, UUID, productID string) error {
	const op = "repository.DeleteFromCart"
	deletedNum, err := r.rds.HDel(ctx, UUID, productID).Result()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if deletedNum == 0 {
		return ErrNoCartProductFound
	}
	return nil
}

func (r RedisRepository) UpdateStock(ctx context.Context, UUID, productID string, stock int) error {
	const op = "repository.UpdateStock"
	err := r.rds.HSet(ctx, UUID, productID, stock).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r RedisRepository) IsAdded(ctx context.Context, UUID, productID string) (bool, error) {
	const op = "repository.IsAdded"
	num, err := r.rds.HGet(ctx, UUID, productID).Result()
	if err != nil && err != redis.Nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return num != "", nil
}

func (r RedisRepository) GetCart(ctx context.Context, UUID string) (map[string]string, error) {
	const op = "repository.GetCart"
	cart, err := r.rds.HGetAll(ctx, UUID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return make(map[string]string), nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return cart, nil
}
