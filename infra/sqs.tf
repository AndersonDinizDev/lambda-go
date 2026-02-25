resource "aws_sqs_queue" "test_sqs" {
  name = "lambda-go-sqs"

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.test_sqs_deadletter.arn
    maxReceiveCount     = 3
  })
}

resource "aws_sqs_queue" "test_sqs_deadletter" {
  name = "lambda-go-sqs-deadletter"
}

resource "aws_sqs_queue_redrive_allow_policy" "test_sqs_redrive_allow_policy" {
  queue_url = aws_sqs_queue.test_sqs_deadletter.id

  redrive_allow_policy = jsonencode({
    redrivePermission = "byQueue",
    sourceQueueArns   = [aws_sqs_queue.test_sqs.arn]
  })
}

resource "aws_sqs_queue_policy" "allow_s3_sendMessage" {
  queue_url = aws_sqs_queue.test_sqs.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "s3.amazonaws.com"
        }
        Action   = "sqs:SendMessage"
        Resource = aws_sqs_queue.test_sqs.arn

        Condition = {
          ArnEquals = {
            "aws:SourceArn" = aws_s3_bucket.test_bucket.arn
          }
        }
      }
    ]
  })
}