package mongo_user_repository

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ThePositree/billing_manager/internal/adapter/repository/user/mongo/dto"
	"github.com/ThePositree/billing_manager/internal/model/user"
	model_user "github.com/ThePositree/billing_manager/internal/model/user"
	"github.com/ThePositree/billing_manager/internal/usecase"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ usecase.UserRepository = &userRepository{}

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

type userRepository struct {
	coll   *mongo.Collection
	client *mongo.Client
	mutex  sync.RWMutex
	cache  map[string]user.User
}

var ErrNoData = errors.New("no data")

func (u *userRepository) GetNoDataError() error {
	return ErrNoData
}

func (u *userRepository) Create(ctx context.Context, user model_user.User) (model_user.User, error) {
	userDto := dto.NewUserDTOFromModel(user)

	_, err := u.coll.InsertOne(ctx, userDto)
	if err != nil {
		return model_user.User{}, fmt.Errorf("mongo insert one: %w", err)
	}

	u.mutex.Lock()
	u.cache[user.Id] = user
	u.mutex.Unlock()

	return user, nil
}

func (u *userRepository) GetByTelegramUN(ctx context.Context, telegramUN string) (model_user.User, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	for _, user := range u.cache {
		if user.TelegramUN == telegramUN {
			return user, nil
		}
	}
	return model_user.User{}, ErrNoData
}

func (u *userRepository) Delete(ctx context.Context, id string) (model_user.User, error) {
	result := u.coll.FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: id}})
	err := result.Err()
	if errors.Is(mongo.ErrNoDocuments, err) {
		return model_user.User{}, ErrNoData
	}
	if err != nil {
		return model_user.User{}, fmt.Errorf("mongo find one and delete: %w", err)
	}

	var userDTO dto.User

	if err := result.Decode(&userDTO); err != nil {
		return model_user.User{}, fmt.Errorf("result decode: %w", err)
	}

	user, err := userDTO.ToModel()
	if err != nil {
		return model_user.User{}, fmt.Errorf("dto to model: %w", err)
	}
	u.mutex.Lock()
	delete(u.cache, user.Id)
	u.mutex.Unlock()

	return user, nil
}

func (u *userRepository) Get(ctx context.Context, id string) (model_user.User, error) {
	u.mutex.RLock()
	user, ok := u.cache[id]
	u.mutex.RUnlock()
	if ok {
		return user, nil
	}
	result := u.coll.FindOne(ctx, bson.D{{Key: "_id", Value: id}})

	err := result.Err()
	if errors.Is(mongo.ErrNoDocuments, err) {
		return model_user.User{}, ErrNoData
	}
	if err != nil {
		return model_user.User{}, fmt.Errorf("mongo find one: %w", err)
	}

	var userDTO dto.User

	if err := result.Decode(&userDTO); err != nil {
		return model_user.User{}, fmt.Errorf("result decode: %w", err)
	}

	user, err = userDTO.ToModel()
	if err != nil {
		return model_user.User{}, fmt.Errorf("dto to model: %w", err)
	}
	u.mutex.Lock()
	u.cache[user.Id] = user
	u.mutex.Unlock()

	return user, nil
}

func (u *userRepository) GetAll(ctx context.Context) ([]model_user.User, error) {
	if len(u.cache) != 0 {
		var result []model_user.User
		for _, user := range u.cache {
			result = append(result, user)
		}
		return result, nil
	}
	cursor, err := u.coll.Find(ctx, bson.D{})
	if err != nil {
		return []model_user.User{}, fmt.Errorf("mongo find: %w", err)
	}
	defer cursor.Close(ctx)

	var users []model_user.User

	for {
		if cursor.TryNext(context.TODO()) {
			var result dto.User
			if err := cursor.Decode(&result); err != nil {
				return []model_user.User{}, fmt.Errorf("result decode: %w", err)
			}
			user, err := result.ToModel()
			if err != nil {
				return []model_user.User{}, fmt.Errorf("dto to model: %w", err)
			}
			users = append(users, user)
			continue
		}
		if err := cursor.Err(); err != nil {
			return []model_user.User{}, fmt.Errorf("cursor error: %w", err)
		}
		if cursor.ID() == 0 {
			break
		}
	}

	return users, nil
}

func New(ctx context.Context, logger zerolog.Logger, client *mongo.Client, cfg Config) (*userRepository, error) {
	err := cfg.Validate()
	if err != nil {
		return &userRepository{}, fmt.Errorf("config validate: %w", err)
	}

	userRepo := &userRepository{
		client: client,
	}

	if err = userRepo.client.Ping(ctx, nil); err != nil {
		return &userRepository{}, fmt.Errorf("mongo ping: %w", err)
	}

	coll := userRepo.client.Database(cfg.Database).Collection(cfg.Collection)

	userRepo.coll = coll

	cache := map[string]user.User{}

	users, err := userRepo.GetAll(ctx)
	if err != nil {
		return &userRepository{}, fmt.Errorf("get all users: %w", err)
	}

	for _, user := range users {
		cache[user.Id] = user
	}

	userRepo.cache = cache

	return userRepo, nil
}
