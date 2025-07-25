package mongodb

import (
	"auth/internal/entity"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportMongo struct {
	db *mongo.Database
}

func NewReportMongo(db *mongo.Database) *ReportMongo {
	return &ReportMongo{db: db}
}

func (r *ReportMongo) CreateReport(report *entity.Report) error {
	collection := r.db.Collection("reports")

	// creation date
	report.Created_at = time.Now()

	_, err := collection.InsertOne(context.Background(), report)
	if err != nil {
		return fmt.Errorf("failed to create report: %w", err)
	}

	return nil
}

func (r *ReportMongo) GetUserReports(userID uuid.UUID) ([]*entity.Report, error) {
	collection := r.db.Collection("reports")

	filter := map[string]interface{}{"user_id": userID}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find reports: %w", err)
	}
	defer cursor.Close(context.Background())

	var reports []*entity.Report
	for cursor.Next(context.Background()) {
		var report entity.Report
		if err := cursor.Decode(&report); err != nil {
			return nil, fmt.Errorf("failed to decode report: %w", err)
		}
		reports = append(reports, &report)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return reports, nil
}

func (r *ReportMongo) SetAnonimousIdReport(clientGeneratedID string, userID uuid.UUID) error {
	collection := r.db.Collection("reports")

	filter := bson.M{
		"client_generated_id": clientGeneratedID,
		"user_id":             nil, // чтобы не тронуть уже привязанные
	}

	update := bson.M{
		"$set": bson.M{
			"user_id": userID,
		},
	}

	_, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to set anonymous ID: %w", err)
	}

	return nil
}

func (r *ReportMongo) GetUserIdAndPriceByReportId(reportID string) (uuid.UUID, float64, error) {
	collection := r.db.Collection("reports")

	filter := bson.M{"report_id": reportID}
	var report entity.Report
	err := collection.FindOne(context.Background(), filter).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return uuid.Nil, 0, fmt.Errorf("no report found with ID %s", reportID)
		}
		return uuid.Nil, 0, fmt.Errorf("failed to find report: %w", err)
	}

	userID, err := uuid.Parse(report.User_id)
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("invalid user ID in report: %w", err)
	}

	return userID, report.Price, nil
}

func (r *ReportMongo) PurchaseReport(reportID string) error {
	collection := r.db.Collection("reports")

	filter := bson.M{
		"report_id":    reportID,
		"is_purchased": false, // <--- защита от повторной покупки
	}
	update := bson.M{"$set": bson.M{"is_purchased": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to purchase report: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no report found with ID %s", reportID)
	}

	return nil
}
