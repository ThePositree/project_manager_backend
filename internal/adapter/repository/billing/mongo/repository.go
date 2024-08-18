package mongo_billing_repository

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ThePositree/billing_manager/internal/adapter/repository/billing/mongo/dto"
	model_billing "github.com/ThePositree/billing_manager/internal/model/billing"
	"github.com/ThePositree/billing_manager/internal/usecase"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrNoData = errors.New("no data")

type Config struct {
	Database   string
	Collection string
}

func (cfg Config) Validate() error {
	if cfg.Database == "" {
		return fmt.Errorf("database name cannot be empty")
	}
	if cfg.Collection == "" {
		return fmt.Errorf("collection name cannot be empty")
	}
	return nil
}

var _ usecase.BillingRepository = &billingRepository{}

type billingRepository struct {
	coll   *mongo.Collection
	client *mongo.Client
	mutex  sync.RWMutex
	cache  map[string]model_billing.Billing
}

func (u *billingRepository) GetNoDataError() error {
	return ErrNoData
}

func (u *billingRepository) GetByUserId(ctx context.Context, userId string) ([]model_billing.Billing, error) {
	var result []model_billing.Billing
	u.mutex.RLock()
	for _, billing := range u.cache {
		if billing.UserId == userId {
			result = append(result, billing)
		}
	}
	u.mutex.RUnlock()
	return result, nil
}

func (u *billingRepository) Update(ctx context.Context, billing model_billing.Billing) (model_billing.Billing, error) {
	billingDto := dto.NewBillingDTOFromModel(billing)

	result := u.coll.FindOneAndReplace(ctx, bson.D{{Key: "_id", Value: billing.Id}}, billingDto)
	if err := result.Err(); err != nil {
		return model_billing.Billing{}, fmt.Errorf("mongo find one and replace: %w", err)
	}

	u.mutex.Lock()
	u.cache[billing.Id] = billing
	u.mutex.Unlock()

	return billing, nil
}

func (u *billingRepository) Create(ctx context.Context, billing model_billing.Billing) (model_billing.Billing, error) {
	billingDto := dto.NewBillingDTOFromModel(billing)

	_, err := u.coll.InsertOne(ctx, billingDto)
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("mongo insert one: %w", err)
	}

	u.mutex.Lock()
	u.cache[billing.Id] = billing
	u.mutex.Unlock()

	return billing, nil
}

func (u *billingRepository) Delete(ctx context.Context, id string) (model_billing.Billing, error) {
	result := u.coll.FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: id}})
	if err := result.Err(); err != nil {
		return model_billing.Billing{}, fmt.Errorf("mongo find one and delete: %w", err)
	}

	var billingDTO dto.Billing

	if err := result.Decode(&billingDTO); err != nil {
		return model_billing.Billing{}, fmt.Errorf("result decode: %w", err)
	}

	billing, err := billingDTO.ToModel()
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("dto to model: %w", err)
	}
	u.mutex.Lock()
	delete(u.cache, billing.Id)
	u.mutex.Unlock()

	return billing, nil
}

func (u *billingRepository) Get(ctx context.Context, id string) (model_billing.Billing, error) {
	u.mutex.RLock()
	billing, ok := u.cache[id]
	u.mutex.RUnlock()
	if ok {
		return billing, nil
	}
	result := u.coll.FindOne(ctx, bson.D{{Key: "_id", Value: id}})

	err := result.Err()
	if errors.Is(mongo.ErrNoDocuments, err) {
		return model_billing.Billing{}, ErrNoData
	}
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("mongo find one: %w", err)
	}

	var billingDTO dto.Billing

	if err := result.Decode(&billingDTO); err != nil {
		return model_billing.Billing{}, fmt.Errorf("result decode: %w", err)
	}

	billing, err = billingDTO.ToModel()
	if err != nil {
		return model_billing.Billing{}, fmt.Errorf("dto to model: %w", err)
	}
	u.mutex.Lock()
	u.cache[billing.Id] = billing
	u.mutex.Unlock()

	return billing, nil
}

func (u *billingRepository) GetAll(ctx context.Context) ([]model_billing.Billing, error) {
	if len(u.cache) != 0 {
		var result []model_billing.Billing
		for _, billing := range u.cache {
			result = append(result, billing)
		}
		return result, nil
	}
	cursor, err := u.coll.Find(ctx, bson.D{})
	if err != nil {
		return []model_billing.Billing{}, fmt.Errorf("mongo find: %w", err)
	}
	defer cursor.Close(ctx)

	var billings []model_billing.Billing

	for {
		if cursor.TryNext(context.TODO()) {
			var result dto.Billing
			if err := cursor.Decode(&result); err != nil {
				return []model_billing.Billing{}, fmt.Errorf("result decode: %w", err)
			}
			billing, err := result.ToModel()
			if err != nil {
				return []model_billing.Billing{}, fmt.Errorf("dto to model: %w", err)
			}
			billings = append(billings, billing)
			continue
		}
		if err := cursor.Err(); err != nil {
			return []model_billing.Billing{}, fmt.Errorf("cursor error: %w", err)
		}
		if cursor.ID() == 0 {
			break
		}
	}

	return billings, nil
}

func New(ctx context.Context, logger zerolog.Logger, client *mongo.Client, cfg Config) (*billingRepository, error) {
	err := cfg.Validate()
	if err != nil {
		return &billingRepository{}, fmt.Errorf("config validate: %w", err)
	}

	billingRepo := &billingRepository{
		client: client,
	}

	if err = billingRepo.client.Ping(ctx, nil); err != nil {
		return &billingRepository{}, fmt.Errorf("mongo ping: %w", err)
	}

	coll := billingRepo.client.Database(cfg.Database).Collection(cfg.Collection)

	billingRepo.coll = coll

	cache := map[string]model_billing.Billing{}

	billings, err := billingRepo.GetAll(ctx)
	if err != nil {
		return &billingRepository{}, fmt.Errorf("get all users: %w", err)
	}

	for _, billing := range billings {
		cache[billing.Id] = billing
	}

	billingRepo.cache = cache

	return billingRepo, nil
}
