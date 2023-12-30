package streamer

import (
	"context"
	"encoding/json"
	"fmt"
	sctx "github.com/phathdt/service-context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mongo-streamer/plugins/mongoc"
	"mongo-streamer/plugins/watermillapp"
	"mongo-streamer/shared/common"
	"os"
)

type streamer struct{}

func New() *streamer {
	return &streamer{}
}

type ResumeToken struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Token string             `json:"token" bson:"token"`
}

type ChangeDocument struct {
	ID struct {
		Data string `json:"_data" bson:"_data"`
	} `json:"_id" bson:"_id"`
	DocumentKey struct {
		ID string `json:"_id" bson:"_id"`
	} `json:"document_key" bson:"documentKey"`
	FullDocument             map[string]interface{} `json:"full_document" bson:"fullDocument"`
	FullDocumentBeforeChange map[string]interface{} `json:"full_document_before_change" bson:"fullDocumentBeforeChange"`
	OperationType            string                 `json:"operation_type" bson:"operationType"`
}

func (s *streamer) Run(sc sctx.ServiceContext) error {
	ctx := context.Background()
	defer ctx.Done()

	streamName := os.Getenv("STREAM_NAME")
	fullDocumentBeforeChange := os.Getenv("FullDocumentBeforeChange")

	comp := sc.MustGet(common.KeyMongo).(mongoc.MongoComp)
	publisher := sc.MustGet(common.KeyNatsPub).(watermillapp.Publisher)

	database := comp.GetClient().Database(comp.GetDbName())
	collection := database.Collection(comp.GetCollectionName())

	resumeTokenDb := comp.GetClient().Database(fmt.Sprintf("%s-resume-tokens", comp.GetDbName()))
	resumeTokenCollection := resumeTokenDb.Collection(comp.GetCollectionName())

	// Check if the resumeToken is already stored
	resumeToken := getStoredResumeToken(resumeTokenCollection)

	// Create a Change Stream with resumeAfter
	opts := options.ChangeStream().
		SetFullDocument("updateLookup")

	if fullDocumentBeforeChange == "true" {
		opts = opts.SetFullDocumentBeforeChange("whenAvailable")
	}

	if resumeToken != nil {
		opts = opts.SetResumeAfter(bson.D{{"_data", resumeToken}})
	}

	stream, err := collection.Watch(ctx, mongo.Pipeline{}, opts)
	if err != nil {
		return err
	}

	for stream.Next(ctx) {
		bsonRaw := stream.Current

		var bsonMap ChangeDocument
		if err = bson.Unmarshal(bsonRaw, &bsonMap); err != nil {
			return err
		}

		// Marshal JSON
		jsonData, err := json.MarshalIndent(bsonMap, "", "  ")
		if err != nil {
			return err
		}

		if err = publisher.PublishRaw(streamName, jsonData); err != nil {
			return err
		}

		if err = storeResumeToken(bsonMap.ID.Data, resumeTokenCollection); err != nil {
			return err
		}
	}

	return nil
}

func storeResumeToken(resumeToken string, collection *mongo.Collection) error {
	doc := ResumeToken{Token: resumeToken}
	if _, err := collection.InsertOne(context.Background(), doc); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Function to get the stored resumeToken
func getStoredResumeToken(collection *mongo.Collection) *string {
	opts := options.FindOne().SetSort(bson.D{{"_id", -1}})
	var result ResumeToken

	err := collection.FindOne(context.Background(), bson.M{}, opts).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil
	} else if err != nil {
		return nil
	}

	return &result.Token
}
