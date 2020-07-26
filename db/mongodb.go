package db

import (
	"context"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mStructTagParser bsoncodec.StructTagParserFunc = func(sf reflect.StructField) (st bsoncodec.StructTags, err error) {
	st, err = bsoncodec.DefaultStructTagParser(sf)

	if _, ok := sf.Tag.Lookup("bson"); ok || err != nil {
		return
	}
	if _, ok := sf.Tag.Lookup("json"); !ok {
		return
	}

	tag := strings.Replace(string(sf.Tag), "json:", "bson:", 1)
	sf.Tag = reflect.StructTag(tag)

	return bsoncodec.DefaultStructTagParser(sf)
}

func NewMongoDBClient(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	codec, err := bsoncodec.NewStructCodec(mStructTagParser)
	if err != nil {
		panic(err)
	}

	registry := bson.NewRegistryBuilder().
		RegisterDefaultEncoder(reflect.Struct, codec).
		RegisterDefaultDecoder(reflect.Struct, codec).Build()
	clientOpts := options.Client().
		SetRegistry(registry).
		ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}

	return client
}
