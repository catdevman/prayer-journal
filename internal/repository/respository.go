package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/catdevman/prayer-journal/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type PrayerRepository interface {
	SavePrayer(ctx context.Context, prayer *models.Prayer) error
	GetPrayersByUser(ctx context.Context, userID string, limit int32) ([]models.Prayer, error)
}

type DynamoRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoRepository(client *dynamodb.Client) *DynamoRepository {
	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		// Fallback or panic depending on your style, logging warning for now
		fmt.Println("WARNING: TABLE_NAME env var is not set")
	}
	return &DynamoRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *DynamoRepository) SavePrayer(ctx context.Context, prayer *models.Prayer) error {
	item, err := attributevalue.MarshalMap(prayer)
	if err != nil {
		return fmt.Errorf("failed to marshal prayer: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("failed to put item to dynamodb: %w", err)
	}

	return nil
}

func (r *DynamoRepository) GetPrayersByUser(ctx context.Context, userID string, limit int32) ([]models.Prayer, error) {
	// Query using the UserID (Partition Key)
	// Assuming standard single-table design or simple table where PK=UserId
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("pk = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userID},
		},
		Limit:            aws.Int32(limit),
		ScanIndexForward: aws.Bool(false), // Sort by SK (CreatedAt) descending (newest first)
	}

	out, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query prayers: %w", err)
	}

	var prayers []models.Prayer
	err = attributevalue.UnmarshalListOfMaps(out.Items, &prayers)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal prayers: %w", err)
	}

	return prayers, nil
}
