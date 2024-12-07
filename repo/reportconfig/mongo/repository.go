package mongo

import (
	"context"

	"github.com/rituraj-junglee/container-critic/models"
	"github.com/rituraj-junglee/container-critic/repo/reportconfig"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repository struct {
	mongoCLient *mongo.Client
	mongodb     string
}

func NewRepository(mongoClient *mongo.Client, mongodb string) reportconfig.Repository {
	return &repository{
		mongoCLient: mongoClient,
		mongodb:     mongodb,
	}
}

func (r *repository) UpdateReportConfig(ctx context.Context, report models.ReportConfig) (err error) {
	filter := bson.M{"_id": report.ReportID}
	opts := options.Update().SetUpsert(true)
	update := bson.M{"$set": report}

	_, err = r.mongoCLient.Database(r.mongodb).
		Collection("reports").
		UpdateOne(ctx, filter, update, opts)
	return
}
func (r *repository) GetReportConfig(ctx context.Context, reportID string) (report models.ReportConfig, err error) {
	filter := bson.M{"_id": reportID}

	err = r.mongoCLient.Database(r.mongodb).
		Collection("reports").
		FindOne(ctx, filter).
		Decode(&report)

	return

}
