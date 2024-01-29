package main

import (
	"create_user/lambda_create"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	lambda.Start(lambda_create.CreateUser)
}
