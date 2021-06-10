package utils

import "go.mongodb.org/mongo-driver/bson"

func GetValidator(collection string, validations bson.M) bson.D {
	return bson.D{
		{"collMod", collection},
		{"validator", bson.M{
			"$jsonSchema": bson.M{
				"properties": validations,
			},
		}},
		{"validationLevel", "moderate"},
		{"validationAction", "error"},
	}
}
