package lambda_read

import (
	"context"
	"encoding/json"
	"fmt"
	"get_user/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	svc := dynamodb.New(sess)
	contactUid, ok := req.PathParameters["id"]
	if !ok {
		fmt.Println("Error: ID no proporcionado en la solicitud")
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Error obteniendo el item de DynamoDB: ID no proporcionado",
		}, nil
	}

	tableName := "ContactsFC"
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(contactUid),
			},
		},
	})

	if err != nil {
		fmt.Println("Error al obtener el item de DynamoDB:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error al obtener el item de DynamoDB: %v", err),
		}, nil
	}

	if result.Item == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       fmt.Sprintf("No se encontró el elemento con el ID: %v", contactUid),
		}, nil
	}

	var contact models.Response
	err = dynamodbattribute.UnmarshalMap(result.Item, &contact)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       fmt.Sprintf("Error: %v", err),
		}, nil
	}

	return responseBody(200, contact), nil
}

func responseBody(statusCode int, body any) events.APIGatewayProxyResponse {
	contact := body
	json, err := json.Marshal(contact)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error: %v", err),
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(json),
	}
}

func errorResponse(statusCode int, err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       fmt.Sprintf("Error: %v", err),
	}
}
