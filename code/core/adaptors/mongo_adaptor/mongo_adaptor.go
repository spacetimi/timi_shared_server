package mongo_adaptor

import (
	"context"
	"errors"
	"fmt"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"github.com/spacetimi/timi_shared_server/utils/reflection_utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/** Package init **/
func init() {
    kTrueConst = true
    kAtomicIncrementReturnDocumentOption = options.After
}

type Config struct {
	SharedMongoURL string
	SharedDBName string
	AppMongoURL string
	AppDBName string
}

func Initialize(cfg Config) {
    _sharedMongoClient = createMongoClient(cfg.SharedMongoURL)
    _sharedDBName 	  = cfg.SharedDBName
    _appMongoClient    = createMongoClient(cfg.AppMongoURL)
    _appDBName         = cfg.AppDBName
}

func GetDataItemByPrimaryKeys(dbSpace DBSpace,
                              collectionName string,
                              primaryKeys []string, primaryKeyValues []interface{},
                              outDataItemPtr interface{},
                              ctx context.Context) error {

    if outDataItemPtr == nil {
		return errors.New("data item pointer is null")
	}

    collection, err := getMongoCollection(dbSpace, collectionName)
    if err != nil {
    	return errors.New("error getting collection: " + err.Error())
	}

	if len(primaryKeys) != len(primaryKeyValues) {
		return errors.New(fmt.Sprintf("mismatched number of primary keys(%d) and values(%d)", len(primaryKeys), len(primaryKeyValues)))
	}

    var filterConditions []bson.E
	for i, primaryKeyValue := range primaryKeyValues {
		filterConditions = append(filterConditions, bson.E{Key: primaryKeys[i], Value: primaryKeyValue})
	}
	filter := bson.D(filterConditions)

	singleResult := collection.FindOne(ctx, filter)
	if singleResult.Err() != nil {
		return errors.New("error finding object matching primary key: " + singleResult.Err().Error())
	}

	err = singleResult.Decode(outDataItemPtr)
	if err != nil {
		return errors.New("error decoding object: " + err.Error())
	}

	return nil
}

func GetDataItemsByFilter(dbSpace DBSpace,
						  collectionName string,
						  keys []string,
						  values []interface{},
						  dataItemFactory func() interface{},
						  ctx context.Context) ([]interface{}, error) {

	collection, err := getMongoCollection(dbSpace, collectionName)
	if err != nil {
		return nil, errors.New("error finding collection: " + err.Error())
	}

	if len(keys) != len(values) {
		return nil, errors.New(fmt.Sprintf("mismatched number of primary keys(%d) and values(%d)", len(keys), len(values)))
	}

	filter := bson.D{{}}
	if len(keys) != 0 {
		var filterConditions []bson.E
		for i, value := range values {
			filterConditions = append(filterConditions, bson.E {
				Key: keys[i],
				Value: value,
			})
		}
		filter = bson.D(filterConditions)
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, errors.New("error reading data items from collection: " + err.Error())
	}

	var dataItems []interface{}
	for cursor.Next(context.Background()) {
		dataItem := dataItemFactory()
		err = cursor.Decode(dataItem)
		if err != nil {
		    return nil, errors.New("error decoding data item: " + err.Error())
		}
		dataItems = append(dataItems, dataItem)
	}

	return dataItems, nil
}

func WriteDataItemByPrimaryKeys(dbSpace DBSpace,
								collectionName string,
								primaryKeys []string,
								dataItemPtr interface{},
								ctx context.Context) error {

	if dataItemPtr == nil {
		return errors.New("data item pointer is null")
	}

	err := insertOnDuplicateUpdateDataItem(dbSpace, collectionName, primaryKeys, dataItemPtr, ctx)
	if err != nil {
		return errors.New("error during insert/update of data item: " + err.Error())
	}
	return nil
}

func AtomicIncrement(dbSpace DBSpace,
					 collectionName string,
					 documentPrimaryKey string,
					 documentPrimaryKeyValue interface{},
					 fieldName string,
					 delta interface{},
					 ctx context.Context) (interface{}, error) {

	collection, err := getMongoCollection(dbSpace, collectionName)
	if err != nil {
		return nil, errors.New("error finding collection: " + err.Error())
	}

	filter := bson.M { documentPrimaryKey: documentPrimaryKeyValue }
	update := bson.M { "$inc": bson.M { fieldName: delta } }

	result := collection.FindOneAndUpdate(ctx, filter, update,
										  &options.FindOneAndUpdateOptions{Upsert:&kTrueConst,
																 		   ReturnDocument:&kAtomicIncrementReturnDocumentOption})

	if result.Err() != nil {
		return nil, errors.New("error performing FindOneAndUpdate: " + result.Err().Error())
	}

	returnedDocument := bson.M{}
	err = result.Decode(&returnedDocument)
	if err != nil {
		return nil, errors.New("error decoding returned document: " + err.Error())
	}

	newFieldValue, ok := returnedDocument[fieldName]
	if !ok {
		return nil, errors.New("failed to find field name in updated document")
	}

	return newFieldValue, nil
}

// TODO: Use contexts correctly. Don't just use context.Background


/*********** Private **********************************************************/

var _sharedMongoClient *mongo.Client
var _sharedDBName string
var _appMongoClient    *mongo.Client
var _appDBName         string

var kTrueConst bool
var kAtomicIncrementReturnDocumentOption options.ReturnDocument

func createMongoClient(mongoURL string) *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		logger.LogFatal("Unable to create mongo client: " + err.Error())
	}
	return client
}

func getMongoCollection(dbSpace DBSpace, collectionName string) (*mongo.Collection, error) {
	var database *mongo.Database
	switch dbSpace {
	case SHARED_DB: database = _sharedMongoClient.Database(_sharedDBName)
	case APP_DB: database = _appMongoClient.Database(_appDBName)
	}
	if database == nil {
		return nil, errors.New("invalid db space")
	}

	collection := database.Collection(collectionName)
	if collection == nil {
		return nil, errors.New("unknown collection")
	}

	return collection, nil
}

func insertOnDuplicateUpdateDataItem(dbSpace DBSpace,
									 collectionName string,
									 primaryKeys []string,
									 dataItemPtr interface{},
									 ctx context.Context) error {
	collection, err := getMongoCollection(dbSpace, collectionName)
	if err != nil {
	    return errors.New("error finding collection: " + err.Error())
	}

	bsonMRepresentation, err := reflection_utils.MarshalStructPtrToBson(dataItemPtr)
	if err != nil {
		return errors.New("error serializing data item: " + err.Error())
	}

    var filterConditions []bson.E
	for _, key := range primaryKeys {
		value, ok := bsonMRepresentation[key]
		if !ok {
			return errors.New("missing value for primary key: " + key)
		}
		filterConditions = append(filterConditions, bson.E{
			Key: key,
			Value: value,
		})
	}
	filter := bson.D(filterConditions)

    _, err = collection.UpdateOne(ctx, filter,
								  bson.D{
									  {"$set", bsonMRepresentation},
								  },
								  &options.UpdateOptions{Upsert:&kTrueConst})

    if err != nil {
    	return errors.New("error updating data item: " + err.Error())
	}

    return nil
}


