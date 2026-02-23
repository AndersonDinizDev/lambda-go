package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testScene struct {
	data          string
	fileName      string
	error         bool
	expectedCpf   string
	expectedValue string
	testName      string
}

func TestExtractPdfData(t *testing.T) {

	scenarioTest := []testScene{

		{"Fatura do mês. CPF: 123.456.789-00. Total a pagar R$ 1.500,50 na data de hoje.", "archive.pdf", false, "12345678900", "1.500,50", "Extrair CPF/Valor corretamente"},
		{"Fatura do mês. 123.456.789-00. Total a pagar 1.500,50 na data de hoje.", "archive.pdf", false, "12345678900", "", "Retornar o valor vazio"},
		{"Fatura do mês. 123.456.789-00. Total a pagar R$ 1.500,50 na data de hoje.", "archive.pdf", false, "12345678900", "1.500,50", "Encontrar CPF sem a chave (cpf)"},
		{"Fatura do mês. 12345678900. Total a pagar R$ 1.500,50 na data de hoje.", "archive.pdf", false, "12345678900", "1.500,50", "Encontrar o CPF desformatado"},
		{"Fatura do mês. 123.456.789. Total a pagar R$ 1.500,50 na data de hoje.", "archive.pdf", true, "", "", "Retornar erro do CPF não encontrado"},
	}

	for _, scene := range scenarioTest {

		t.Run(scene.testName, func(t *testing.T) {

			receivedExtractorData, err := extractPdfData(scene.data, scene.fileName)

			assert.False(t, (err != nil) != scene.error, "Não deveria retornar erro neste cenário")
			assert.Equal(t, scene.expectedCpf, receivedExtractorData.Cpf, "O CPF devem ser iguais")
			assert.Equal(t, scene.expectedValue, receivedExtractorData.Value, "O valores devem ser iguais")
		})

	}
}
