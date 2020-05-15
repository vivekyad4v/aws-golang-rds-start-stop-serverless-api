package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"rdst/utils"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"
)

// Help function to generate an IAM policy
func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	fmt.Println("Generating policy...")
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	return authResponse
}

func handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := event.AuthorizationToken
	fmt.Println("Validating token...")

	getServerAuthToken, err := utils.GetParameterValue("authorization-token")

	if err != nil {
		log.Error("Unable to fetch server auth token", err)
	}

	getServerAuthToken = "beer " + getServerAuthToken

	switch strings.ToLower(token) {
	case getServerAuthToken:
		return generatePolicy("user", "Allow", event.MethodArn), nil
	case "deny":
		return generatePolicy("user", "Deny", event.MethodArn), nil
	case "unauthorized":
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized") // Return a 401 Unauthorized response
	default:
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Error: Invalid token")
	}
}
