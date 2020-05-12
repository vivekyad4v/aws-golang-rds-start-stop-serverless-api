package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"rdst/dydb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/rs/xid"
)

var err error
var errstrings []string

// BodyRequest requested json file
type BodyRequest struct {
	Type   string   `json:"type"`
	Values []string `json:"values"`
}

// GetCurrentTime get the current time
func GetCurrentTime() (t string) {
	const (
		layoutISO = "2006-01-02 15:04:05.000000"
	)
	currentTime := time.Now()
	n := currentTime.Format(layoutISO)
	fmt.Println(n)
	return n
}

// GetUUID - Use as primary key
func GetUUID() (u string) {
	id := xid.New()
	return id.String()
}

// ActionDBInstance stops the DB instances
func ActionDBInstance(i string, t string) (e error) {
	svc := rds.New(session.New())
	GT := GetCurrentTime()
	GU := GetUUID()
	var s string

	switch t {
	case "stop":
		input := &rds.StopDBInstanceInput{
			DBInstanceIdentifier: &i,
		}
		result, err := svc.StopDBInstance(input)
		fmt.Println("Stopping DB instance - ", i)
		if err == nil {
			s = "stopped"
			fmt.Println(result, "\n [Stopped]: Putting Item..", t)
			InputItem := dydb.Item{
				UUID:         GU,
				DbIdentifier: i,
				Status:       s,
				CreatedAt:    GT,
				Error:        err.Error(),
			}
			dydb.PutItem(InputItem)
			fmt.Println("[Stopped]: Finished putting item...")
		} else {
			fmt.Println(err)
		}
		return err

	case "start":
		input := &rds.StartDBInstanceInput{
			DBInstanceIdentifier: &i,
		}
		result, err := svc.StartDBInstance(input)
		fmt.Println("Starting DB instance - ", i)
		if err == nil {
			s = "started"
			fmt.Println(result, "\n [Started]: Putting Item..", t)
			InputItem := dydb.Item{
				UUID:         GU,
				DbIdentifier: i,
				Status:       s,
				CreatedAt:    GT,
				Error:        err.Error(),
			}
			dydb.PutItem(InputItem)
			fmt.Println("[Started]: Finished putting item...")
		} else {
			fmt.Println(err)
		}
		return err

	default:
		fmt.Println("[Error]: Instance Request Type doesn't match - ", t)
		return err
	}
}

// Handler API Gateway request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var bodyRequest BodyRequest

	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	for _, el := range bodyRequest.Values {
		err := ActionDBInstance(el, bodyRequest.Type)
		if err != nil {
			errstrings = append(errstrings, err.Error())
		}
	}

	if errstrings != nil {
		errstring := strings.Join(errstrings, " ")
		return events.APIGatewayProxyResponse{Body: errstring, StatusCode: 408}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(bodyRequest.Type), StatusCode: 200}, nil
}
