package main

import (
	"context"
	"log"

	"Lambda/internal/services"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("Erro ao carregar as configurações da AWS: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	dynamoClient := dynamodb.NewFromConfig(cfg)

	dynamoConf := &services.DynamoConfig{
		DynamoClient: dynamoClient,
		TableName:    "lambda-go-teste-1",
	}

	s3Conf := &services.PdfHanlder{
		S3Client: s3Client,
		Dynamo:   dynamoConf,
	}

	lambda.Start(s3Conf.ProcessPDFHandler)
}
