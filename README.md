# Lambda PDF Extractor

Este projeto implementa uma função AWS Lambda em Go para processar arquivos PDF enviados para um bucket S3. A função extrai informações específicas (CPF e valores monetários) e armazena os dados processados em uma tabela DynamoDB. A infraestrutura é gerenciada via Terraform.

## Funcionalidades

- **Trigger Automático:** Execução iniciada por eventos de upload no S3 (`s3:ObjectCreated:*`).
- **Processamento de PDF:** Leitura e extração de texto de arquivos PDF.
- **Extração de Dados:** Identificação de padrões de CPF e valores monetários (R$) via expressões regulares.
- **Persistência:** Armazenamento dos dados extraídos no DynamoDB com ID único gerado via hash (SHA256).

## Arquitetura

O fluxo de dados segue a seguinte arquitetura:

1.  Usuário/Sistema faz upload de um arquivo `.pdf` no bucket S3 configurado.
2.  O S3 envia uma mensagem para o SQS.
3.  A função Lambda mapeia a fila do SQS e inicia a processamento da fila
4.  A função Lambda baixa o arquivo na memória e processa o conteúdo.
5.  Os dados extraídos (CPF, Valor) são salvos na tabela DynamoDB.

## Pré-requisitos

Para executar e implantar este projeto, você precisará das seguintes ferramentas instaladas:

- **Go** (versão 1.26 ou superior)
- **Terraform** (versão 1.0 ou superior)
- **AWS CLI** configurado com credenciais apropriadas

## Estrutura do Projeto

```
.
├── cmd/
│   └── lambda-pdf/    # Código fonte da função Lambda (main.go)
├── infra/             # Código Terraform para provisionamento
├── internal/          # Lógica de negócio e serviços internos
│   ├── models/        # Estruturas de dados
│   └── services/      # Lógica de extração e persistência
├── pkg/               # Pacotes utilitários (hashing)
└── go.mod             # Definição de dependências Go
```

## Configuração e Implantação

### 1. Clonar o Repositório

```bash
git clone <url-do-repositorio>
cd Lambda
```

### 2. Compilar o Binário

A função Lambda utiliza o runtime `provided.al2023` em arquitetura `arm64`. É necessário compilar o binário com as flags corretas e nomeá-lo como `bootstrap`.

```bash
GOOS=linux GOARCH=arm64 go build -o cmd/lambda-pdf/bootstrap cmd/lambda-pdf/main.go
```

> **Nota:** O Terraform espera encontrar o arquivo `bootstrap` no diretório `cmd/lambda-pdf/`.

### 3. Provisionar Infraestrutura

Navegue até o diretório de infraestrutura e execute os comandos do Terraform.

```bash
cd infra

# Inicializar o Terraform (baixa providers)
terraform init

# Visualizar o plano de execução
terraform plan

# Aplicar a infraestrutura
terraform apply
```

> **Atenção:** Verifique o nome do bucket S3 no arquivo `infra/s3.tf`. Buckets S3 devem ter nomes globalmente únicos. Pode ser necessário alterar o valor antes de aplicar.

## Uso

Após a implantação bem-sucedida:

1.  Faça upload de um arquivo PDF contendo um CPF (formato `000.000.000-00`) e um valor monetário (formato `R$ 0,00`) para o bucket S3 criado.
2.  A função Lambda será acionada automaticamente.
3.  Verifique os logs da execução no AWS CloudWatch.
4.  Consulte a tabela DynamoDB (`lambda-go-test`) para ver os dados extraídos.

## Limpeza

Para remover todos os recursos criados e evitar custos adicionais na AWS:

```bash
cd infra
terraform destroy
```
