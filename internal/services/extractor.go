package services

import (
	"Lambda/internal/models"
	"Lambda/pkg/hashutils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ledongthuc/pdf"
)

var (
	regexCpf   = regexp.MustCompile(`(\d{3}\.?\d{3}\.?\d{3}-?\d{2})`)
	regexValue = regexp.MustCompile(`R\$\s?(\d{1,3}(?:\.\d{3})*,\d{2})`)
)

type PdfHanlder struct {
	S3Client *s3.Client
	Dynamo   *DynamoConfig
}

func extractPdfData(data, fileName string) (models.PdfData, error) {

	var u models.PdfData

	cpf := regexCpf.FindStringSubmatch(data)
	value := regexValue.FindStringSubmatch(data)

	if len(cpf) > 1 {
		cpfBrute := cpf[1]

		cpfClean := strings.ReplaceAll(cpfBrute, ".", "")
		cpfClean = strings.ReplaceAll(cpfClean, "-", "")

		u.Cpf = cpfClean
	} else {
		return u, errors.New("CPF não encontrado no documento")
	}

	if len(value) > 1 {
		valueClean := value[1]

		u.Value = valueClean
	}

	u.Id = hashutils.GenerateSHA256(u.Cpf + "#" + fileName)

	return u, nil

}

func (cmd *PdfHanlder) ProcessPDFHandler(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {

	var wg sync.WaitGroup
	var mu sync.Mutex
	var response events.SQSEventResponse

	for _, SqsRecord := range sqsEvent.Records {

		wg.Add(1)

		go func(record events.SQSMessage) {
			defer wg.Done()

			log.Printf("Iniciando processamento paralelo da msg ID: %s", record.MessageId)

			var s3Event events.S3Event

			err := json.Unmarshal([]byte(record.Body), &s3Event)

			if err != nil {
				log.Printf("Erro ao decodificar JSON do SQS para S3Event: %v", err)

				mu.Lock()

				response.BatchItemFailures = append(response.BatchItemFailures, events.SQSBatchItemFailure{
					ItemIdentifier: record.MessageId,
				})

				mu.Unlock()

				return
			}

			for _, s3Record := range s3Event.Records {
				bucketName := s3Record.S3.Bucket.Name
				fileName := s3Record.S3.Object.Key

				log.Printf("Processando o arquivo [%s] do bucket [%s]", fileName, bucketName)

				bucketFile, err := cmd.S3Client.GetObject(ctx, &s3.GetObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(fileName),
				})

				if err != nil {
					log.Printf("Erro ao encontrar o arquivo [%s] no bucket [%s]. Motivo: [%v]", fileName, bucketName, err)

					mu.Lock()

					response.BatchItemFailures = append(response.BatchItemFailures, events.SQSBatchItemFailure{
						ItemIdentifier: record.MessageId,
					})

					mu.Unlock()

					return
				}

				body, err := io.ReadAll(bucketFile.Body)

				bucketFile.Body.Close()

				if err != nil {
					log.Printf("Erro ao ler o arquivo [%s]. Motivo: [%v]", fileName, err)

					mu.Lock()

					response.BatchItemFailures = append(response.BatchItemFailures, events.SQSBatchItemFailure{
						ItemIdentifier: record.MessageId,
					})

					mu.Unlock()

					return
				}

				readingFromMem := bytes.NewReader(body)
				fileSize := int64(len(body))

				f, err := pdf.NewReader(readingFromMem, fileSize)

				if err != nil {
					log.Printf("Erro ao ler o pdf [%s]. Motivo: [%v]", fileName, err)

					mu.Lock()

					response.BatchItemFailures = append(response.BatchItemFailures, events.SQSBatchItemFailure{
						ItemIdentifier: record.MessageId,
					})

					mu.Unlock()

					return
				}

				var buf bytes.Buffer

				b, err := f.GetPlainText()

				if err != nil {
					log.Printf("Erro ao obter os textos do pdf [%s]. Motivo: [%v]", fileName, err)

					mu.Lock()

					response.BatchItemFailures = append(response.BatchItemFailures, events.SQSBatchItemFailure{
						ItemIdentifier: record.MessageId,
					})

					mu.Unlock()

					return
				}

				buf.ReadFrom(b)
				content := buf.String()

				data, err := extractPdfData(content, fileName)

				if err != nil {
					log.Println(err)

					mu.Lock()

					response.BatchItemFailures = append(response.BatchItemFailures, events.SQSBatchItemFailure{
						ItemIdentifier: record.MessageId,
					})

					mu.Unlock()

					return
				}

				statusSave, err := cmd.Dynamo.SaveData(ctx, data)

				if err != nil {
					log.Printf("Erro ao salvar os dados no banco de dados. Motivo: [%v]", err)

					mu.Lock()

					response.BatchItemFailures = append(response.BatchItemFailures, events.SQSBatchItemFailure{
						ItemIdentifier: record.MessageId,
					})

					mu.Unlock()

					return
				}

				log.Println(statusSave)

			}

		}(SqsRecord)

	}

	wg.Wait()

	log.Println("Todo o lote do SQS foi processado com sucesso!")
	return response, nil
}
