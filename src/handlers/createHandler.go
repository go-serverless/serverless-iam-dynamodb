package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gopkg.in/go-playground/validator.v9"

	"github.com/satori/go.uuid"

	utils "../utils"
)

type User struct {
	ID            string  `json:"id" validate:"required"`
	UserName      string  `json:"user_name" validate:"required,min=4,max=20"`
	FirstName     *string `json:"first_name,omitempty"`
	LastName      *string `json:"last_name,omitempty"`
	Age           *int    `json:"age,omitempty"`
	Phone         *string `json:"phone,omitempty"`
	Password      string  `json:"password" validate:"required,min=4,max=50"`
	Email         string  `json:"email" validate:"required,email"`
	Role          string  `json:"role" validate:"required,min=4,max=20"`
	IsActive      bool    `json:"is_active" validate:"required"`
	CreatedAt     string  `json:"created_at,omitempty"`
	ModifiedAt    string  `json:"modified_at,omitempty"`
	DeactivatedAt *string `json:"deactivated_at,omitempty"`
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

func Create(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		id        = uuid.Must(uuid.NewV4()).String()
		tableName = aws.String(os.Getenv("IAM_TABLE_NAME"))
	)

	// Initialize user
	user := &User{
		ID:         id,
		IsActive:   true,
		Role:       "user",
		CreatedAt:  time.Now().String(),
		ModifiedAt: time.Now().String(),
	}

	// Parse request body
	json.Unmarshal([]byte(request.Body), user)

	// Validate user struct
	var validate *validator.Validate
	validate = validator.New()
	err := validate.Struct(user)
	if err != nil {
		// Error HTTP response
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	// Encrypt password
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		fmt.Println("Got error calling HashPassword:")
		fmt.Println(err.Error())
		// Error HTTP response
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	// Write to DynamoDB
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		fmt.Println("Got error calling MarshalMap:")
		fmt.Println(err.Error())
		// Error HTTP response
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: tableName,
	}
	if _, err := ddbSvc.PutItem(params); err != nil {
		// Error HTTP response
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	} else {
		body, _ := json.Marshal(user)
		// Success HTTP response
		return events.APIGatewayProxyResponse{
			Body:       string(body),
			StatusCode: 200,
		}, nil
	}
}

func main() {
	lambda.Start(Create)
}
