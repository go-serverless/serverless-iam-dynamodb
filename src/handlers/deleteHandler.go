package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

var svc *dynamodb.DynamoDB

func init() {
	region := os.Getenv("AWS_REGION")
	// Initialize a session
	if session, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to initialize a session to AWS: %s", err.Error()))
	} else {
		// Create DynamoDB client
		svc = dynamodb.New(session)
	}
}

func Delete(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		tableName = aws.String(os.Getenv("IAM_TABLE_NAME"))
		id        = aws.String(request.PathParameters["id"])
	)

	// Delete a record by id
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: id,
			},
		},
		TableName: tableName,
	}
	_, err := svc.DeleteItem(input)
	if err != nil {
		fmt.Println("Got error calling DeleteItem:")
		fmt.Println(err.Error())

		// Error HTTP response
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	} else {
		// Success HTTP response
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
		}, nil
	}
}

func main() {
	lambda.Start(Delete)
}
