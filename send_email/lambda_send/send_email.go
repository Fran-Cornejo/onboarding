package lambda_send

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sns"
)

type dynamoDBImage struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Status    string `json:"status"`
}

func SendEmail(ctx context.Context, event events.DynamoDBEvent) error {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return err
	}

	svc := sns.New(sess)
	topicArn := "arn:aws:sns:us-east-1:620097380428:ContactsTopicFC"

	for _, record := range event.Records {
		if record.EventName == "INSERT" {
			dynamoDBImg := dynamoDBImage{}
			err := unmarshalStreamImage(record.Change.NewImage, &dynamoDBImg)
			if err != nil {
				return err
			}

			params := &sns.PublishInput{
				Subject: aws.String(fmt.Sprintf("Nuevo Registro de : %s", dynamoDBImg.ID)),
				Message: aws.String(fmt.Sprintf("Nuevo usuario registrado! \nID: %s\nFirstName: %s\nLastName: %s",
					dynamoDBImg.ID, dynamoDBImg.FirstName, dynamoDBImg.LastName)),
				TopicArn: aws.String(topicArn),
			}

			_, err = svc.Publish(params)
			if err != nil {
				fmt.Println("Error enviando un mensaje:", err)
				return err
			}
		}
	}

	return nil
}

func unmarshalStreamImage(attribute map[string]events.DynamoDBAttributeValue, v interface{}) error {
	av := make(map[string]*dynamodb.AttributeValue)

	for key, value := range attribute {
		av[key] = &dynamodb.AttributeValue{S: aws.String(value.String())}
	}

	return dynamodbattribute.UnmarshalMap(av, v)
}
