package main

import (
	"encoding/json"
	"strings"
	"time"

	"rdst/dydb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

var (
	getCurrentTime  string
	getUUID         string
	errstrings      []string
	err             error
	listInstanceIDs []string
	listClusterIDs  []string
	setEngine       string
)

// BodyRequest requested json file
type BodyRequest struct {
	Type   string   `json:"type"`
	Values []string `json:"values"`
}

// GetCurrentTime get the current time
func GetCurrentTime() (getCurrentTime string) {
	const (
		layoutISO = "2006-01-02 15:04:05.000000"
	)
	getCurrentTime = time.Now().Format(layoutISO)
	return getCurrentTime
}

// GetUUID - Use as primary key
func GetUUID() (getUUID string) {
	getUUID = xid.New().String()
	return getUUID
}

// ActionDBInstance stops the DB instances
func ActionDBInstance(instanceID string, actionType string) (Error error) {
	rdssvc := rds.New(session.New())
	log.SetFormatter(&log.JSONFormatter{})

	switch actionType {

	case "stop", "start":
		getCurrentTime = GetCurrentTime()
		getUUID = GetUUID()

		setEngine = "postgres"
		input := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: &instanceID,
		}
		result, err := rdssvc.DescribeDBInstances(input)
		if err != nil {
			log.WithFields(log.Fields{
				"input": result,
			}).Info("Engine type doesn't exists, set type - aurora")
			setEngine = "aurora"
		}

		if actionType == "stop" && setEngine == "postgres" {
			input := &rds.StopDBInstanceInput{
				DBInstanceIdentifier: &instanceID,
			}
			_, err = rdssvc.StopDBInstance(input)
			if err != nil {
				log.Error("unable to stop instance Error - ", instanceID, err)
				actionType = "stop-error"
			}
		}

		if actionType == "start" && setEngine == "postgres" {
			input := &rds.StartDBInstanceInput{
				DBInstanceIdentifier: &instanceID,
			}
			_, err = rdssvc.StartDBInstance(input)
			if err != nil {
				log.Error("unable to start instance Error - ", instanceID, err)
				actionType = "start-error"
			}
		}

		if actionType == "stop" && setEngine == "aurora" {
			input := &rds.StopDBClusterInput{
				DBClusterIdentifier: &instanceID,
			}
			_, err = rdssvc.StopDBCluster(input)
			if err != nil {
				log.Error("unable to stop cluster Error - ", instanceID, err)
				actionType = "stop-error"
			}
		}

		if actionType == "start" && setEngine == "aurora" {
			input := &rds.StartDBClusterInput{
				DBClusterIdentifier: &instanceID,
			}
			_, err = rdssvc.StartDBCluster(input)
			if err != nil {
				log.Error("unable to start cluster Error - ", instanceID, err)
				actionType = "start-error"
			}
		}

		setError := "NIL"
		if err != nil {
			setError = err.Error()
		}

		inputItem := dydb.Item{
			UUID:         getUUID,
			DbIdentifier: instanceID,
			Status:       actionType,
			CreatedAt:    getCurrentTime,
			Error:        setError,
		}

		err = dydb.PutItem(inputItem)
		if err != nil {
			log.Error("unable to put item Error - ", instanceID, err)
			log.WithFields(log.Fields{
				"input": inputItem,
			}).Error("unable to put item Error!")
		}

	case "stopall", "startall":
		if actionType == "stopall" {
			// POSTGRES ACTION
			setEngine = "postgres"
			input := &rds.DescribeDBInstancesInput{
				Filters: []*rds.Filter{
					{
						Name: aws.String("engine"),
						Values: []*string{
							aws.String(setEngine),
						},
					},
				},
			}
			result, err := rdssvc.DescribeDBInstances(input)
			if err != nil {
				log.Error("unable to fetch all instances - ", instanceID, err)
				actionType = "stopall-error"
			}

			for _, i := range result.DBInstances {
				listInstanceIDs = append(listInstanceIDs, *i.DBInstanceIdentifier)
			}

			log.Info("Instance list to action - ", listInstanceIDs)
			for _, i := range listInstanceIDs {
				getCurrentTime = GetCurrentTime()
				getUUID = GetUUID()
				input := &rds.StopDBInstanceInput{
					DBInstanceIdentifier: &i,
				}
				_, err = rdssvc.StopDBInstance(input)
				if err != nil {
					log.Error("unable to stopall instance Error - ", i, err)
					actionType = "stopall-error"
				}
				inputItem := dydb.Item{
					UUID:         getUUID,
					DbIdentifier: i,
					Status:       actionType,
					CreatedAt:    getCurrentTime,
					Error:        err.Error(),
				}
				err = dydb.PutItem(inputItem)
				if err != nil {
					log.Error("unable to put item Error - ", i, err)
					log.WithFields(log.Fields{
						"input": inputItem,
					}).Error("unable to put item Error!")
				}
			}

			// AURORA ACTION
			setEngine = "aurora"
			inputa := &rds.DescribeDBClustersInput{
				Filters: []*rds.Filter{
					{
						Name: aws.String("engine"),
						Values: []*string{
							aws.String(setEngine),
						},
					},
				},
			}
			resulta, err := rdssvc.DescribeDBClusters(inputa)
			if err != nil {
				log.Error("unable to fetch all instances - ", instanceID, err)
				actionType = "stopall-error"
			}
			for _, i := range resulta.DBClusters {
				listClusterIDs = append(listInstanceIDs, *i.DBClusterIdentifier)
			}

			log.Info("Clusters list to action - ", listClusterIDs)
			for _, i := range listClusterIDs {
				getCurrentTime = GetCurrentTime()
				getUUID = GetUUID()
				input := &rds.StopDBClusterInput{
					DBClusterIdentifier: &i,
				}
				_, err = rdssvc.StopDBCluster(input)
				if err != nil {
					log.Error("unable to stopall cluster Error - ", i, err)
					actionType = "stopall-error"
				}
				inputItem := dydb.Item{
					UUID:         getUUID,
					DbIdentifier: i,
					Status:       actionType,
					CreatedAt:    getCurrentTime,
					Error:        err.Error(),
				}
				err = dydb.PutItem(inputItem)
				if err != nil {
					log.Error("unable to put item Error - ", i, err)
					log.WithFields(log.Fields{
						"input": inputItem,
					}).Error("unable to put item Error!")
				}
			}
		}

		if actionType == "startall" {
			setEngine = "postgres"
			input := &rds.DescribeDBInstancesInput{
				Filters: []*rds.Filter{
					{
						Name: aws.String("engine"),
						Values: []*string{
							aws.String(setEngine),
						},
					},
				},
			}
			result, err := rdssvc.DescribeDBInstances(input)
			if err != nil {
				log.Error("unable to fetch all instances - ", instanceID, err)
				actionType = "startall-error"
			}

			for _, i := range result.DBInstances {
				listInstanceIDs = append(listInstanceIDs, *i.DBInstanceIdentifier)
			}

			log.Info("Instance list to action - ", listInstanceIDs)
			for _, i := range listInstanceIDs {
				getCurrentTime = GetCurrentTime()
				getUUID = GetUUID()
				input := &rds.StartDBInstanceInput{
					DBInstanceIdentifier: &i,
				}
				_, err = rdssvc.StartDBInstance(input)
				if err != nil {
					log.Error("unable to startall instance Error - ", i, err)
					actionType = "startall-error"
				}
				inputItem := dydb.Item{
					UUID:         getUUID,
					DbIdentifier: i,
					Status:       actionType,
					CreatedAt:    getCurrentTime,
					Error:        err.Error(),
				}
				err = dydb.PutItem(inputItem)
				if err != nil {
					log.Error("unable to put item Error - ", instanceID, err)
					log.WithFields(log.Fields{
						"input": inputItem,
					}).Error("unable to put item Error!")
				}
			}

			// AURORA ACTION
			setEngine = "aurora"
			inputa := &rds.DescribeDBClustersInput{
				Filters: []*rds.Filter{
					{
						Name: aws.String("engine"),
						Values: []*string{
							aws.String(setEngine),
						},
					},
				},
			}
			resulta, err := rdssvc.DescribeDBClusters(inputa)
			if err != nil {
				log.Error("unable to fetch all instances - ", instanceID, err)
				actionType = "Startall-error"
			}
			for _, i := range resulta.DBClusters {
				listClusterIDs = append(listInstanceIDs, *i.DBClusterIdentifier)
			}

			log.Info("Clusters list to action - ", listClusterIDs)
			for _, i := range listClusterIDs {
				getCurrentTime = GetCurrentTime()
				getUUID = GetUUID()
				input := &rds.StartDBClusterInput{
					DBClusterIdentifier: &i,
				}
				_, err = rdssvc.StartDBCluster(input)
				if err != nil {
					log.Error("unable to Startall cluster Error - ", i, err)
					actionType = "Startall-error"
				}
				inputItem := dydb.Item{
					UUID:         getUUID,
					DbIdentifier: i,
					Status:       actionType,
					CreatedAt:    getCurrentTime,
					Error:        err.Error(),
				}
				err = dydb.PutItem(inputItem)
				if err != nil {
					log.Error("unable to put item Error - ", i, err)
					log.WithFields(log.Fields{
						"input": inputItem,
					}).Error("unable to put item Error!")
				}
			}
		}

	default:
		log.Info("[Info]: actionType doesn't match - ", actionType)
	}

	return err
}

// Handler API Gateway request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var bodyRequest BodyRequest
	log.SetFormatter(&log.JSONFormatter{})

	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		log.Error("Unable to unmarshal JSON", err)
		log.WithFields(log.Fields{
			"input": []byte(request.Body),
		}).Error("JSON unmarshal error!")

		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 405}, nil
	}

	for _, instanceID := range bodyRequest.Values {
		err := ActionDBInstance(instanceID, bodyRequest.Type)
		if err != nil {
			errstrings = append(errstrings, err.Error())
		}
	}

	if errstrings != nil {
		errstring := strings.Join(errstrings, " ")
		return events.APIGatewayProxyResponse{Body: errstring, StatusCode: 406}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(bodyRequest.Type), StatusCode: 200}, nil
}
