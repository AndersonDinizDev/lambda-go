package models

type PdfData struct {
	Id    string `json:"id" dynamodbav:"id"`
	Cpf   string `json:"cpf" dynamodbav:"cpf"`
	Value string `json:"value" dynamodbav:"value"`
}
