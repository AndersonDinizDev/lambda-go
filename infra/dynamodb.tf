resource "aws_dynamodb_table" "test_table" {
  name         = "lambda-go-test"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }
}