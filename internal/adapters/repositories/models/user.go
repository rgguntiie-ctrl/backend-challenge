package models

import (
	"time"

	"github.com/kanta/backend-challenge/internal/core/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Email     string        `bson:"email" json:"email"`
	Password  string        `bson:"password" json:"-"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

func ToUserModels(u *domain.User) *User {
	var (
		id  bson.ObjectID
		err error
	)

	if u.ID != "" {
		id, err = bson.ObjectIDFromHex(u.ID)
		if err != nil {
			id = bson.NewObjectID()
		}
	} else {
		id = bson.NewObjectID()
	}

	return &User{
		ID:        id,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
}

func ToUserDomain(u *User) *domain.User {
	return &domain.User{
		ID:        u.ID.Hex(),
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
}
