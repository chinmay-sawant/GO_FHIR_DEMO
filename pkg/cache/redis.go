// Package cache provides caching utilities using Redis.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// RedisCache defines the contract for cache operations
type RedisCache interface {
	GetPatient(id string) (*fhir.Patient, error)
	SetPatient(ctx context.Context, id string, patient *fhir.Patient, expiration time.Duration) error
	DeletePatient(ctx context.Context, id string) error
	Ping(ctx context.Context) error
	Close() error
}


// RedisCacheImpl implements RedisCache using Redis
type RedisCacheImpl struct {
	client *redis.Client
}

// Config holds Redis configuration
type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(config Config) RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})

	return &RedisCacheImpl{
		client: client,
	}
}

// GetPatient retrieves a patient from Redis cache
func (r *RedisCacheImpl) GetPatient(id string) (*fhir.Patient, error) {
	key := "patient:" + id
	result, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get patient from cache: %w", err)
	}

	var patient fhir.Patient
	if err := json.Unmarshal([]byte(result), &patient); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached patient: %w", err)
	}

	return &patient, nil
}

// SetPatient stores a patient in Redis cache
func (r *RedisCacheImpl) SetPatient(ctx context.Context, id string, patient *fhir.Patient, expiration time.Duration) error {
	key := "patient:" + id
	data, err := json.Marshal(patient)
	if err != nil {
		return fmt.Errorf("failed to marshal patient for cache: %w", err)
	}

	if err := r.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set patient in cache: %w", err)
	}

	return nil
}

// DeletePatient removes a patient from Redis cache
func (r *RedisCacheImpl) DeletePatient(ctx context.Context, id string) error {
	key := "patient:" + id
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete patient from cache: %w", err)
	}
	return nil
}

// Ping tests the Redis connection
func (r *RedisCacheImpl) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Close closes the Redis connection
func (r *RedisCacheImpl) Close() error {
	return r.client.Close()
}
