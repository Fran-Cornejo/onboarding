package lambda_create

import (
	"context"
	"create_user/models"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

func CreateUser(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var newContact models.Contacts

	err := json.Unmarshal([]byte(req.Body), &newContact)

	if err != nil {
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
		fmt.Println(err, " | "+err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error al insertar el contacto en DynamoDB",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       fmt.Sprintf("Contacto con ID %s creado exitosamente en la tabla %s", contactID, tableName),
	}, nil
}
