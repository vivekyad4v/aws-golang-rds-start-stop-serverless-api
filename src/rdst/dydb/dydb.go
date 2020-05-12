package dydb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Item struct for DB actions
type Item struct {
	UUID         string
	DbIdentifier string
	Status       string
	CreatedAt    string
	Error        string
}

// PutItem in dynamodb table
func PutItem(i Item) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess)

	ItemAVMap, err := dynamodbattribute.MarshalMap(i)

	if err != nil {
		fmt.Printf("Marshalling: ERROR: %v\n", err.Error())
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String("happay-dev-rdst-tbl"),
		Item:      ItemAVMap,
	}

	_, err = db.PutItem(params)
	if err != nil {
		fmt.Printf("Put item: ERROR: %v\n", err.Error())
	}

	fmt.Println("Put item: Success")
}
