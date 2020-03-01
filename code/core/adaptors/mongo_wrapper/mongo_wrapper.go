package mongo_wrapper

import (
	"context"
	"github.com/spacetimi/timi_shared_server/code/config"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
)

/**
 * TODO: Delete mongo_wrapper
 */

/**
 * TODO: Modify all the functions to return err when appropriate
 */

// Constants

type DBSpace int
const (
	SHARED_DB DBSpace = iota
	APP_DB
)

var true_const bool

// Package init

func init() {
	true_const = true
}

// Public API

func Initialize(sharedMongoURL string, sharedDBName string,
				appMongoURL string, appDBName string) {
	_sharedMongoClient = createMongoClient(sharedMongoURL)
	_sharedDBName 	  = sharedDBName
	_appMongoClient    = createMongoClient(appMongoURL)
	_appDBName         = appDBName
}

func GetDataItemByPrimaryKeys(pkValues []interface{}, outDataItemPtr IDataItem) {
	if outDataItemPtr == nil {
		logger.LogError("Data item is nil")
		return
	}

	dataItemDescriptor := outDataItemPtr.GetDescriptor()
	if dataItemDescriptor == nil {
		logger.LogError("Data item descriptor is nil")
		return
	}

	collection := getMongoCollection(dataItemDescriptor.GetDB(),
									 dataItemDescriptor.GetCollectionName())

	primaryKeys := dataItemDescriptor.GetPrimaryKeys()
	if len(primaryKeys) != len(pkValues) {
		logger.LogError("Mismatched number of primary keys and values")
		return
	}

	var filterConditions []bson.E
	for i, pkValue := range pkValues {
		filterConditions = append(filterConditions, bson.E{Key: primaryKeys[i], Value: pkValue})
	}
	filter := bson.D(filterConditions)

	singleResult := collection.FindOne(context.Background(), filter)
	if singleResult.Err() != nil {
		logger.LogError("Error finding object matching primary key: " + singleResult.Err().Error())
		return
	}

	err := singleResult.Decode(outDataItemPtr)
	if err != nil {
		logger.LogError("Could not find and decode object by primary key: " + err.Error())
	}
}

func GetDataItemCollection(outDataItemCollectionPtr IDataItemCollection) {
	GetDataItemCollectionByFilter(nil, nil, outDataItemCollectionPtr)
}

func GetDataItemCollectionByFilter(keys []string, values []interface{},
								   outDataItemCollectionPtr IDataItemCollection) {
	if outDataItemCollectionPtr == nil {
		logger.LogError("Data item collection is nil")
		return
	}

	dataItemDescriptor := outDataItemCollectionPtr.GetDescriptor()
	if dataItemDescriptor == nil {
		logger.LogError("Data item descriptor is nil")
		return
	}

	dataItemFactory := outDataItemCollectionPtr.GetDataItemFactory()
	if dataItemFactory == nil {
		logger.LogError("Data item factory is nil")
		return
	}

	collection := getMongoCollection(dataItemDescriptor.GetDB(),
									 dataItemDescriptor.GetCollectionName())

	if len(keys) != len(values) {
		logger.LogError("Mismatched number of keys and values")
	}

	filter := bson.D{{}}
	if len(keys) != 0 {
		var filterConditions []bson.E
		for i, value := range values {
			filterConditions = append(filterConditions, bson.E{keys[i], value})
		}
		filter = bson.D(filterConditions)
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		logger.LogError("Error reading data items from collection: " + err.Error())
		return
	}

	for cursor.Next(context.Background()) {
		dataItem := dataItemFactory.CreateDataItem()
		err = cursor.Decode(dataItem)
		if err != nil {
			logger.LogError("Error decoding: " + err.Error())
		}
		outDataItemCollectionPtr.AddDataItem(dataItem)
	}
	return
}

func ApplyDataItem(dataItemPtr IDataItem) {
	if dataItemPtr == nil {
		logger.LogError("Data item is nil")
		return
	}

	if dataItemPtr.Dirty() {
		insertOnDuplicateUpdateDataItem(dataItemPtr)
	}
}

func ApplyDataItemCollection(dataItemCollectionPtr IDataItemCollection) {
	if dataItemCollectionPtr == nil {
		logger.LogError("Data item collection is nil")
		return
	}

	// If NOT production, test to make sure that if the collection is not marked dirty, then it truly isn't dirty
	// Will catch potential bugs
	if config.GetEnvironmentConfiguration().AppEnvironment != config.PRODUCTION {
		if !dataItemCollectionPtr.Dirty() {
			for _, dataItem := range dataItemCollectionPtr.GetDataItems() {
				if dataItem.Dirty() {
					logger.LogError("Potential bug: data item collection is not marked dirty," +
						                          "but contains dirty data items")
				}
			}
		}
	}

	if dataItemCollectionPtr.Dirty() {
		for _, dataItem := range dataItemCollectionPtr.GetDataItems() {
			ApplyDataItem(dataItem)
		}
		dataItemCollectionPtr.SetDirty(false)
	}
}

// Private methods

var _sharedMongoClient *mongo.Client
var _sharedDBName string
var _appMongoClient    *mongo.Client
var _appDBName         string

func createMongoClient(mongoURL string) *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		logger.LogFatal("Unable to create mongo client: " + err.Error())
	}
	return client
}

func getMongoCollection(dbType DBSpace, collectionName string) *mongo.Collection {
	var database *mongo.Database
	switch dbType {
	case SHARED_DB: database = _sharedMongoClient.Database(_sharedDBName);
	case APP_DB: database = _appMongoClient.Database(_appDBName);
	}
	if database == nil {
		logger.LogError("Unknown db type")
		return nil
	}

	collection := database.Collection(collectionName)
	if collection == nil {
		logger.LogError("Unknown collection: " + collectionName)
		return nil
	}

	return collection
}

func insertOnDuplicateUpdateDataItem(dataItemPtr IDataItem) {

	if dataItemPtr == nil {
		logger.LogError("Data item is nil")
		return
	}

	dataItemDescriptor := dataItemPtr.GetDescriptor()
	if dataItemDescriptor == nil {
		logger.LogError("Data item descriptor is nil")
		return
	}

	collection := getMongoCollection(dataItemDescriptor.GetDB(),
									 dataItemDescriptor.GetCollectionName())
	keys := dataItemDescriptor.GetPrimaryKeys()
	bsonMRepresentation := marshalStructPtrToBson(dataItemPtr)

	var filterConditions []bson.E
	for _, key := range keys {
		value, ok := bsonMRepresentation[key]
		if !ok {
			logger.LogError("Not found value for key: " + key)
		}
		filterConditions = append(filterConditions, bson.E{key, value})
	}
	filter := bson.D(filterConditions)

	_, err := collection.UpdateOne(context.Background(), filter,
								   bson.D{
										{"$set", bsonMRepresentation},
								   },
								   &options.UpdateOptions{Upsert:&true_const})

	if err != nil {
		logger.LogError("Error in insert/update: " + err.Error())
	}

	dataItemPtr.SetDirty(false)
}

func marshalStructPtrToBson(s interface{}) bson.M {
	bsonMRepresentation := make(map[string]interface{})

	p := reflect.ValueOf(s)
	if p.Kind() != reflect.Ptr {
		logger.LogError("Not a struct ptr")
		return bsonMRepresentation
	}

	v := reflect.Indirect(p)
	if v.Kind() != reflect.Struct {
		logger.LogError("Not a struct")
		return bsonMRepresentation
	}

	numFields := v.NumField()
	for i := 0; i < numFields; i++ {
		value := v.Field(i)
		if value.CanInterface() {
			bsonMRepresentation[strings.ToLower(v.Type().Field(i).Name)] = value.Interface()
		}
	}

	return bsonMRepresentation
}


