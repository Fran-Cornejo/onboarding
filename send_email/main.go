package main

import (
	"send_email/lambda_send"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(lambda_send.SendEmail)
}
