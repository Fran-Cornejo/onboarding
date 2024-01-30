package main

import (
	"modified_status/lambda_status"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	lambda.Start(lambda_status.StatusProccess)
}
