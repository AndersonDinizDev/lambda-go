resource "aws_lambda_permission" "allows_s3_trigger_lambda" {
  statement_id = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_go.function_name
  principal     = "s3.amazonaws.com"

  source_arn = aws_s3_bucket.test_bucket.arn
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = aws_s3_bucket.test_bucket.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.lambda_go.arn

    events = ["s3:ObjectCreated:*"]

    filter_prefix = "pdf/"
    filter_suffix = ".pdf"
  }
  depends_on = [aws_lambda_permission.allows_s3_trigger_lambda]
}