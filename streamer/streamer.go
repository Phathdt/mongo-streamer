package streamer

import (
	"context"
	"encoding/json"
	sctx "github.com/phathdt/service-context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mongo-streamer/plugins/mongoc"
	"mongo-streamer/shared/common"
)

type streamer struct{}

func New() *streamer {
	return &streamer{}
}

func (s *streamer) Run(sc sctx.ServiceContext) error {
	ctx := context.Background()
	defer ctx.Done()

	comp := sc.MustGet(common.KeyMongo).(mongoc.MongoComp)

	database := comp.GetClient().Database(comp.GetDbName())
	collection := database.Collection(comp.GetCollectionName())

	//resumeTokenDb := comp.GetClient().Database(fmt.Sprintf("%s-resume-tokens", comp.GetDbName()))
	//resumeTokenCollection := resumeTokenDb.Collection(comp.GetCollectionName())
	opts := options.ChangeStream().SetFullDocument("updateLookup").SetFullDocumentBeforeChange("whenAvailable")
	stream, err := collection.Watch(ctx, mongo.Pipeline{}, opts)
	if err != nil {
		return err
	}

	logger := sc.Logger("")

	for stream.Next(ctx) {
		bsonRaw := stream.Current

		// Unmarshal BSON to JSON
		type changeset struct {
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
		var bsonMap changeset
		if err = bson.Unmarshal(bsonRaw, &bsonMap); err != nil {
			return err
		}

		// Marshal JSON
		jsonData, err := json.MarshalIndent(bsonMap, "", "  ")
		if err != nil {
			return err
		}

		logger.Info(string(jsonData))
	}

	return nil
}
