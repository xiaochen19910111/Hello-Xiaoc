package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type myConfig struct {
	Count int      `json:"count"`
	Name  string   `json:"name"`
	Fruit []string `json:"fruit"`
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	mySession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")}))
	svc := appconfig.New(mySession)
	input := &appconfig.GetConfigurationInput{
		Application:   aws.String("newConfigApp"),
		ClientId:      aws.String(uuid.NewString()),
		Configuration: aws.String("TestConfig1"),
		Environment:   aws.String("Prod"),
	}

	result, err := svc.GetConfiguration(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case appconfig.ErrCodeResourceNotFoundException:
				fmt.Println(appconfig.ErrCodeResourceNotFoundException, aerr.Error())
			case appconfig.ErrCodeInternalServerException:
				fmt.Println(appconfig.ErrCodeInternalServerException, aerr.Error())
			case appconfig.ErrCodeBadRequestException:
				fmt.Println(appconfig.ErrCodeBadRequestException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.GET("/testconfig", func(c echo.Context) error {
		var config myConfig
		json.Unmarshal(result.Content, &config)
		return c.HTML(http.StatusOK, config.Name+" has "+strconv.Itoa(config.Count)+" apples, but he prefers "+config.Fruit[0]+" and "+config.Fruit[1])
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
