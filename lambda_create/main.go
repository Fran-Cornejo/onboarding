package main

import (
	"lambda_create/lambdas"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	lambda.Start(lambdas.CreateUser)
}
