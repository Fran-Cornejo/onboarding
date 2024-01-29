package main

import (
	"get_user/lambda_read"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	lambda.Start(lambda_read.GetUser)
}
