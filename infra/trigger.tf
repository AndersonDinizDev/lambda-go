resource "aws_lambda_event_source_mapping" "allows_sqs_trigger_lambda" {
  event_source_arn = aws_sqs_queue.test_sqs.arn
  function_name = aws_lambda_function.lambda_go.arn
  batch_size = 10
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = aws_s3_bucket.test_bucket.id


  queue {
    queue_arn     = aws_sqs_queue.test_sqs.arn
    events        = ["s3:ObjectCreated:*"]
    filter_prefix = "pdf/"
    filter_suffix = ".pdf"
  }

  depends_on = [aws_sqs_queue_policy.allow_s3_sendMessage]
}