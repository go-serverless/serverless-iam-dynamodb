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
	Response []User `json:"response"`
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

func List(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		tableName = aws.String(os.Getenv("IAM_TABLE_NAME"))
	)

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		TableName: tableName,
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		// Status Bad Request
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	// Construct final response
	var users []User
	for _, i := range result.Items {
		user := User{}
		if err := dynamodbattribute.UnmarshalMap(i, &user); err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: 400,
			}, nil
		}
		users = append(users, user)
	}

	body, _ := json.Marshal(&Response{
		Response: users,
	})

	// Status OK
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(List)
}
