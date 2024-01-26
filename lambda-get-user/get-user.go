package lambdas

import (
	"log"
	"onboarding/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetUser() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("AWS_REGION")},
	)

	svc := dynamodb.New(sess)
	contactUid := "asd"
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
		log.Fatalf("Error buscando el campo %s", err)
	}

	var contactUser models.Response

	err = dynamodbattribute.UnmarshalMap(result.Item, &contactUser)

}
