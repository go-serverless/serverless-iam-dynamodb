package main

import (
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	jwt "github.com/dgrijalva/jwt-go"
	"os"
)

func Auth(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := request.AuthorizationToken
	tokenSlice := strings.Split(token, " ")
	var tokenString string
	if len(tokenSlice) > 1 {
		tokenString = tokenSlice[len(tokenSlice)-1]
	}

	k := os.Getenv("JWT_SECRET_KEY")

	var jwtToken, _ = jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		// TODO: validate the expected algo
		return []byte(k), nil
	})

	if jwtToken != nil && !jwtToken.Valid {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
	}

	return generatePolicy("user", "Allow", request.MethodArn), nil
}

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
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

func main() {
	lambda.Start(Auth)
}
