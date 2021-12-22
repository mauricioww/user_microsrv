package helpers

import (
	"context"

	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func BuildUpdateBson(d entities.UserDetails) bson.D {
	return bson.D{
		{"$set", extracFields(d)},
	}
}

func BuildInsertBson(d entities.UserDetails) bson.M {
	b := extracFields(d)
	b["_id"] = d.UserId
	return b
}

func NoExists(coll *mongo.Collection, ctx context.Context, id int) bool {
	var results bson.M
	err := coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&results)
	return err == mongo.ErrNoDocuments
}

func extracFields(d entities.UserDetails) bson.M {
	return bson.M{
		"country":       d.Country,
		"city":          d.City,
		"mobile_number": d.MobileNumber,
		"married":       d.Married,
		"height":        d.Height,
		"weigth":        d.Weigth,
	}
}
