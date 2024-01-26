package lambdas

import (
	"context"
	"encoding/json"
	"fmt"
	"lambda_create/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

func CreateUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var newContact models.Contacts
	fmt.Println(req.Body)

	err := json.Unmarshal([]byte(req.Body), &newContact)
	if err != nil {
		fmt.Println(req.Body)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error al decodificar el body de la solicitud",
		}, nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error al crear la sesion de DynamoDB",
		}, nil
	}

	// Crear un cliente de DynamoDB
	svc := dynamodb.New(sess)
	fmt.Println(&svc)
	tableName := "ContactsFC"

	//Creo un UUID para el contacto y genero el valor en la tabla
	contactID := uuid.NewString()
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(contactID),
			},
			"firstName": {
				S: aws.String(newContact.FirstName),
			},
			"lastName": {
				S: aws.String(newContact.LastName),
			},
			"status": {
				S: aws.String("CREATED"),
			},
		},
	}

	//Inserto el contacto en la tabla
	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Aca esta rompiendo svc.PutItem(input)")
		fmt.Println(err, " | "+err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error al insertar el contacto en DynamoDB",
		}, nil
	}

	responseBody := fmt.Sprintf("Contacto con ID %s creado exitosamente en la tabla %s", contactID, tableName)
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       responseBody,
	}, nil
}
