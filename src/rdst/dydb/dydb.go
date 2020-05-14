package dydb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"
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
func PutItem(i Item) (Error error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess)

	ItemAVMap, err := dynamodbattribute.MarshalMap(i)
	if err != nil {
		log.Error("Marshalling: ERROR: ", err)
		return err
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String("happay-dev-rdst-tbl"),
		Item:      ItemAVMap,
	}
	_, err = db.PutItem(params)
	if err != nil {
		log.Error("Unable to put item: ERROR: ", err)
		return err
	}

	log.Info("Put item: Success")
	return err
}
