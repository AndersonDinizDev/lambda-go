package services

import (
	"Lambda/internal/models"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoConfig struct {
	DynamoClient *dynamodb.Client
	TableName    string
}

func (cmd *DynamoConfig) SaveData(ctx context.Context, pdf models.PdfData) (bool, error) {

	item, err := attributevalue.MarshalMap(pdf)

	if err != nil {
		return false, err
	}

	_, err = cmd.DynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(cmd.TableName),
		Item:      item,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
