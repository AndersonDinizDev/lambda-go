data "archive_file" "lambda-go-binary" {
  type        = "zip"
  source_file = "${path.module}/../cmd/lambda-pdf/bootstrap"
  output_path = "${path.module}/../cmd/lambda-pdf/bootstrap.zip"
}

resource "aws_lambda_function" "lambda_go" {
  function_name = "lambda-go-test"

  role = aws_iam_role.lambda_role.arn

  filename = data.archive_file.lambda-go-binary.output_path

  runtime = "provided.al2023"
  handler = "bootstrap"

  source_code_hash = data.archive_file.lambda-go-binary.output_base64sha256

  memory_size = 128
  timeout     = 10

  architectures = ["arm64"]

  environment {
    variables = {
      DYNAMODB_TABLE = aws_dynamodb_table.test_table.name
    }
  }
}