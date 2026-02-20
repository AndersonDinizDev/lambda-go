package services

import (
	"Lambda/internal/models"
	"Lambda/pkg/hashutils"
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"regexp"
	"strings"

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

func (cmd *PdfHanlder) ProcessPDFHandler(ctx context.Context, s3Event events.S3Event) error {

	for _, record := range s3Event.Records {

		bucketName := record.S3.Bucket.Name
		fileName := record.S3.Object.Key

		log.Printf("Processando o arquivo [%s] do bucket [%s]", fileName, bucketName)

		bucketFile, err := cmd.S3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fileName),
		})

		if err != nil {
			log.Printf("Erro ao encontrar o arquivo [%s] no bucket [%s]", fileName, bucketName)
			return err
		}

		defer bucketFile.Body.Close()

		body, err := io.ReadAll(bucketFile.Body)

		if err != nil {
			log.Printf("Erro ao ler o arquivo [%s]", fileName)
			return err
		}

		readingFromMem := bytes.NewReader(body)
		fileSize := int64(len(body))

		f, err := pdf.NewReader(readingFromMem, fileSize)

		if err != nil {
			log.Printf("Erro ao ler o pdf [%s]", fileName)
			return err
		}

		var buf bytes.Buffer

		b, err := f.GetPlainText()

		if err != nil {
			log.Printf("Erro ao obter os textos do pdf [%s]", fileName)
			return err
		}

		buf.ReadFrom(b)
		content := buf.String()

		data, err := extractPdfData(content, fileName)

		if err != nil {
			log.Println(err)
			return err
		}

		statusSave, err := cmd.Dynamo.SaveData(ctx, data)

		if err != nil {
			log.Printf("Erro ao salvar os dados no banco de dados")
			return err
		}

		log.Println(statusSave)

	}

	return nil
}
