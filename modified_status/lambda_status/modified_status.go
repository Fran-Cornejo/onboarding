package lambda_status

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type SNSMessage struct {
	ID string `json:"id"`
}

func StatusProccess(ctx context.Context, snsEvent events.SNSEvent) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return err
	}

	dynamoDBClient := dynamodb.New(sess)

	for _, record := range snsEvent.Records {
		fmt.Println("SNS record:", record)

		/* messageText := record.SNS.Message
		messageSubject := record.SNS.Subject

		fmt.Println(" SNS completo :", record.SNS)
		fmt.Println("Message Text:", messageText)
		fmt.Println("Message subject:", messageSubject) */

		inputString := record.SNS.Subject
		inputStringID := strings.TrimPrefix(inputString, "Nuevo Registro de : ")
		fmt.Println("inputStringID: ", inputStringID)

		snsMessage := SNSMessage{
			ID: inputStringID,
		}
		fmt.Println("snsMessage:", snsMessage)

		// Actualizar el estado a "PROCESSED" en DynamoDB
		_, err = dynamoDBClient.UpdateItem(&dynamodb.UpdateItemInput{
			TableName: aws.String("ContactsFC"),
			Key: map[string]*dynamodb.AttributeValue{
				"id": {
					S: aws.String(snsMessage.ID),
				},
			},
			UpdateExpression: aws.String("SET #s = :status"),
			ExpressionAttributeNames: map[string]*string{
				"#s": aws.String("status"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":status": {
					S: aws.String("PROCESSED"),
				},
			},
		})
		fmt.Println("Mensaje SNS Procesado: ", snsMessage)
		if err != nil {
			fmt.Println("Error actualizando el status:", err)
			return err
		}
	}

	return nil
}
