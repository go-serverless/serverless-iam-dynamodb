package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"gopkg.in/go-playground/validator.v9"
	"os"
	"reflect"
	"strings"
	"time"
)

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

type User struct {
	ID            *string `json:"id,omitempty"`
	UserName      *string `json:"user_name,omitempty" validate:"omitempty,min=4,max=20"`
	FirstName     *string `json:"first_name,omitempty"`
	LastName      *string `json:"last_name,omitempty"`
	Age           *int    `json:"age,omitempty"`
	Phone         *string `json:"phone,omitempty"`
	Email         *string `json:"email,omitempty" validate:"omitempty,email"`
	Role          *string `json:"role,omitempty" validate:"omitempty,min=4,max=20"`
	IsActive      *bool   `json:"is_active,omitempty"`
	CreatedAt     *string `json:"created_at,omitempty"`
	ModifiedAt    string  `json:"modified_at,omitempty"`
	DeactivatedAt *string `json:"deactivated_at,omitempty"`
}

func Update(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		tableName = aws.String(os.Getenv("IAM_TABLE_NAME"))
		id        = aws.String(request.PathParameters["id"])
	)

	user := &User{
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

	// Construct update builder
	update := expression.UpdateBuilder{}

	u := reflect.ValueOf(user).Elem()
	t := u.Type()

	for i := 0; i < u.NumField(); i++ {
		f := u.Field(i)
		// check if it is nil
		if !reflect.DeepEqual(f.Interface(), reflect.Zero(f.Type()).Interface()) {
			jsonFieldName := t.Field(i).Name
			// get json field name
			if jsonTag := t.Field(i).Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
				if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
					jsonFieldName = jsonTag[:commaIdx]
				}
			}
			// construct update
			update = update.Set(expression.Name(jsonFieldName), expression.Value(f.Interface()))
		}
	}

	builder := expression.NewBuilder().WithUpdate(update)
	expression, err := builder.Build()

	if err != nil {
		// Error HTTP response
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	// Update a record by id
	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: id,
			},
		},
		ExpressionAttributeNames:  expression.Names(),
		ExpressionAttributeValues: expression.Values(),
		UpdateExpression:          expression.Update(),
		ReturnValues:              aws.String("UPDATED_NEW"),
		TableName:                 tableName,
	}
	_, err = ddbSvc.UpdateItem(input)
	if err != nil {
		fmt.Println("Got error calling UpdateItem:")
		fmt.Println(err.Error())
		// Error HTTP response
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	} else {
		// Success HTTP response
		return events.APIGatewayProxyResponse{
			Body:       request.Body,
			StatusCode: 200,
		}, nil
	}
}

func main() {
	lambda.Start(Update)
}
