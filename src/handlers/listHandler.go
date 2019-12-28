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
	ID            string  `json:"id" validate:"required"`
	UserName      string  `json:"user_name" validate:"required,min=4,max=20"`
	FirstName     *string `json:"first_name,omitempty"`
	LastName      *string `json:"last_name,omitempty"`
	Age           *int    `json:"age,omitempty"`
	Phone         *string `json:"phone,omitempty"`
	Email         string  `json:"email" validate:"required,email"`
	Role          string  `json:"role" validate:"required,min=4,max=20"`
	IsActive      bool    `json:"is_active" validate:"required"`
	CreatedAt     string  `json:"created_at" validate:"required`
	ModifiedAt    string  `json:"modified_at" validate:"required`
	DeactivatedAt *string `json:"deactivated_at,omitempty"`
}

type ListUsersResponse struct {
	Users []User `json: "users"`
}

var ddbSvc *dynamodb.DynamoDB

func init() {
	region := os.Getenv("AWS_REGION")
	// Initialize a session
	if session, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to initialize a session to AWS: %s", err.Error()))
	} else {
		// Create DynamoDB client
		ddbSvc = dynamodb.New(session)
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
	result, err := ddbSvc.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	// Construct users from response
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

	// Success HTTP response
	body, _ := json.Marshal(&ListUsersResponse{
		Users: users,
	})
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(List)
}
