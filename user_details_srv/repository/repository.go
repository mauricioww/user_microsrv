package repository

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/mauricioww/user_microsrv/helpers"
	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDetailsRepository interface {
	SetUserDetails(ctx context.Context, info entities.UserDetails) (bool, error)
	GetUserDetails(ctx context.Context, user_id int) (entities.UserDetails, error)
}

type userDetailsRepository struct {
	db     *mongo.Database
	logger log.Logger
}

func NewUserDetailsRepository(mongo_db *mongo.Database, l log.Logger) UserDetailsRepository {
	return &userDetailsRepository{
		db:     mongo_db,
		logger: log.With(l, "repository", "mysql"),
	}
}

func (detailsRepo userDetailsRepository) SetUserDetails(ctx context.Context, details entities.UserDetails) (bool, error) {
	collection := detailsRepo.db.Collection("information")
	var err error

	if helpers.NoExists(collection, ctx, details.UserId) {
		_, err = collection.InsertOne(ctx, helpers.BuildInsertBson(details))
	} else {
		_, err = collection.UpdateByID(ctx, details.UserId, helpers.BuildUpdateBson(details))
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (detailsRepo userDetailsRepository) GetUserDetails(ctx context.Context, user_id int) (entities.UserDetails, error) {
	collection := detailsRepo.db.Collection("information")
	var results entities.UserDetails

	if helpers.NoExists(collection, ctx, user_id) {
		err := errors.New("User does not exists")
		return results, err
	} else if err := collection.FindOne(ctx, bson.D{{"_id", user_id}}).Decode(&results); err != nil {
		return results, err
	}

	return results, nil
}
