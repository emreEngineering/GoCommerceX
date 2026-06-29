package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"GoCommerceX/cart-service/internal/domain"
	"GoCommerceX/cart-service/internal/ports"

	"github.com/redis/go-redis/v9"
)

type RedisCartRepository struct {
	client    *redis.Client
	keyPrefix string
}

func NewRedisCartRepository(client *redis.Client, keyPrefix string) *RedisCartRepository {
	return &RedisCartRepository{client: client, keyPrefix: keyPrefix}
}

func (r *RedisCartRepository) key(userID string) string {
	return fmt.Sprintf("%s%s", r.keyPrefix, userID)
}

func (r *RedisCartRepository) Save(ctx context.Context, cart domain.Cart) error {
	cart.UpdatedAt = time.Now()
	if cart.CreatedAt.IsZero() {
		cart.CreatedAt = cart.UpdatedAt
	}

	payload, err := json.Marshal(cart)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key(cart.UserID), payload, 0).Err()
}

func (r *RedisCartRepository) FindByUserID(ctx context.Context, userID string) (domain.Cart, error) {
	value, err := r.client.Get(ctx, r.key(userID)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return domain.Cart{}, ports.ErrCartNotFound
		}
		return domain.Cart{}, err
	}

	var cart domain.Cart
	if err := json.Unmarshal(value, &cart); err != nil {
		return domain.Cart{}, err
	}
	return cart, nil
}

func (r *RedisCartRepository) Update(ctx context.Context, cart domain.Cart) error {
	if cart.ID == "" {
		return ports.ErrCartNotFound
	}
	return r.Save(ctx, cart)
}

func (r *RedisCartRepository) Delete(ctx context.Context, userID string) error {
	result := r.client.Del(ctx, r.key(userID))
	if result.Err() != nil {
		return result.Err()
	}
	if result.Val() == 0 {
		return ports.ErrCartNotFound
	}
	return nil
}
