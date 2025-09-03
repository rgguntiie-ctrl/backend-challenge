package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/kanta/backend-challenge/internal/adapters/repositories/models"
	"github.com/kanta/backend-challenge/internal/core/domain"
	"github.com/kanta/backend-challenge/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type userRepository struct {
	client *mongo.Client
	db     string
	col    string
}

func NewUserRepository(client *mongo.Client, db string) ports.UserRepository {
	return &userRepository{
		client: client,
		db:     db,
		col:    "users",
	}
}

func (r *userRepository) Create(user *domain.User) error {
	user.CreatedAt = time.Now()

	m := models.ToUserModels(user)

	_, err := r.client.Database(r.db).Collection(r.col).InsertOne(context.Background(), m)
	return err
}

func (r *userRepository) FindOne(filter map[string]interface{}) (*domain.User, error) {
	bsonFilter := bson.M{}
	for k, v := range filter {
		bsonFilter[k] = v
	}

	var m models.User
	err := r.client.Database(r.db).Collection(r.col).FindOne(context.Background(), bsonFilter).Decode(&m)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}
	return models.ToUserDomain(&m), err
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	var m models.User
	err := r.client.Database(r.db).Collection(r.col).FindOne(context.Background(), bson.M{"_id": id}).Decode(&m)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}
	return models.ToUserDomain(&m), err
}

func (r *userRepository) ListPaginated(skip, limit int) ([]domain.User, error) {
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cur, err := r.client.Database(r.db).Collection(r.col).Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var users []domain.User
	for cur.Next(context.Background()) {
		var m models.User
		if err := cur.Decode(&m); err == nil {
			users = append(users, *models.ToUserDomain(&m))
		}
	}
	return users, cur.Err()
}

func (r *userRepository) Update(user *domain.User) error {
	update := bson.M{
		"$set": bson.M{
			"name":  user.Name,
			"email": user.Email,
		},
	}
	_, err := r.client.Database(r.db).Collection(r.col).UpdateOne(context.Background(), bson.M{"_id": user.ID}, update)
	return err
}

func (r *userRepository) Delete(id string) error {
	_, err := r.client.Database(r.db).Collection(r.col).DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *userRepository) Count() (int64, error) {
	return r.client.Database(r.db).Collection(r.col).CountDocuments(context.Background(), bson.M{})
}
