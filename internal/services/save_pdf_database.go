package services

import (
	"Lambda/internal/models"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type TableTest struct {
	DynamoDBClient *dynamodb.Client
	TableName      string
}

func (tableTest *TableTest) SaveData(ctx context.Context, pdf models.PdfData) error {

	item, err := attributevalue.MarshalMap(pdf)

	if err != nil {
		return err
	}

	_, err = tableTest.DynamoDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableTest.TableName),
		Item:      item,
	})

	if err != nil {
		return err
	}

	return nil
}
