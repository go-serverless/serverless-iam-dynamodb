package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type User struct {
	ID            string `json:"id"`
	UserName      string `json:"user_name"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Age           int    `json:"age"`
	Phone         string `json:"phone"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	IsActive      bool   `json:"is_active"`
	CreatedAt     string `json:"created_at"`
	ModifiedAt    string `json:"modified_at"`
	DeactivatedAt string `json:"deactivated_at"`
}

type Response struct {
	Response User `json:"response"`
}

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

func Get(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		tableName = aws.String(os.Getenv("IAM_TABLE_NAME"))
		id        = aws.String(request.PathParameters["id"])
	)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: id,
			},
		},
	})
	if err != nil {
		fmt.Println("Got error calling GetItem:")
		fmt.Println(err.Error())
		// Status Internal Server Error
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}

	// Contruct final response
	user := User{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &user); err != nil {
		fmt.Println("Got error unmarshalling:")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}

	body, _ := json.Marshal(&Response{
		Response: user,
	})

	// Status OK
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Get)
}
