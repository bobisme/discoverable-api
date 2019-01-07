package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/bobisme/discoverable-api/swag/docs"
	"github.com/go-openapi/loads"
	"github.com/labstack/echo"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/xeipuuv/gojsonschema"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath
func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/accounts/:id", ShowAccount)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/paramexamples", ParamExample())

	e.Logger.Fatal(e.Start(":1323"))
}

type Account struct {
	Id int `json:"id"`
}

// ShowAccount godoc
// @Summary Show an account
// @Description get string by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {object} main.Account
// @Failure 400 {object} main.HTTPError
// @Failure 404 {object} main.HTTPError
// @Failure 500 {object} main.HTTPError
// @Router /accounts/{id} [get]
func ShowAccount(ctx echo.Context) error {
	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	account := &Account{Id: aid}
	return ctx.JSON(http.StatusOK, account)
}

func NewError(ctx echo.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

type ParamExampleRequest struct {
	// A is the first parameter I have listed
	A string `json:"a,omitempty" binding:"required"`
	B string `json:"b,omitempty"`
	C int    `json:"c" minimum:"10"`
}

// ParamExample godoc
// @Summary Show an account
// @Description get string by ID
// @ID get-string-by-int-poop
// @Accept  json
// @Produce  json
// @Param data body main.ParamExampleRequest true "body"
// @Success 200 {object} main.Account
// @Router /paramexamples [get]
func ParamExample() echo.HandlerFunc {
	specDoc, err := loads.Spec("./docs/swagger/swagger.json")
	panicIf(err)
	specDoc, err = specDoc.Expanded()
	panicIf(err)
	specc := specDoc.Spec()
	return func(c echo.Context) error {
		path, ok := specc.Paths.Paths[c.Path()]
		if !ok {
			return c.JSON(http.StatusNotFound, specc.Paths.Paths)
		}
		pees := path.Get.Parameters
		paramsSchema, err := pees[0].JSONLookup("schema")
		panicIf(err)
		paramJson, err := json.Marshal(paramsSchema)
		panicIf(err)
		fmt.Println(string(paramJson))
		schemaLoader := gojsonschema.NewBytesLoader(paramJson)
		fmt.Println(schemaLoader.JsonSource())
		schema, err := gojsonschema.NewSchema(schemaLoader)
		panicIf(err)
		derp := new(ParamExampleRequest)
		if err := c.Bind(derp); err != nil {
			return err
		}
		result, err := schema.Validate(gojsonschema.NewGoLoader(derp))
		panicIf(err)
		if result.Valid() {
			return c.String(http.StatusOK, "The document is valid")
		}

		c.Response().Header().Set(
			echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		c.Response().Write([]byte(
			"The document is not valid. see errors :\n"))
		for _, err := range result.Errors() {
			// Err implements the ResultError interface
			fmt.Fprintf(c.Response(), "- %s\n", err)
		}
		return nil
	}
}

type HALLink struct {
	HREF      string `json:"href"`
	Title     string `json:"title"`
	Templated bool   `json:"templated"`
}

type ParamsResource struct {
	Links map[string]HALLink `json:"_links"`
	A     string             `json:"a"`
	B     string             `json:"b"`
	C     int                `json:"c"`
}

type ParamsPage struct {
	Links    map[string]HALLink `json:"_links"`
	Embedded []ParamsResource   `json:"_embedded"`
	Total    int                `json:"total"`
}

// ParamsExample godoc
// @Summary Show an account
// @Description get string by ID
// @ID get-string-by-int-poop
// @Accept  json
// @Produce  json
// @Param data body main.ParamExampleRequest true "body"
// @Success 200 {object} []main.Account
// @Router /paramexamples [get]
func ParamsExample() echo.HandlerFunc {
	specDoc, err := loads.Spec("./docs/swagger/swagger.json")
	panicIf(err)
	specDoc, err = specDoc.Expanded()
	panicIf(err)
	specc := specDoc.Spec()
	return func(c echo.Context) error {
		path, ok := specc.Paths.Paths[c.Path()]
		if !ok {
			return c.JSON(http.StatusNotFound, specc.Paths.Paths)
		}
		pees := path.Get.Parameters
		paramsSchema, err := pees[0].JSONLookup("schema")
		panicIf(err)
		paramJson, err := json.Marshal(paramsSchema)
		panicIf(err)
		fmt.Println(string(paramJson))
		schemaLoader := gojsonschema.NewBytesLoader(paramJson)
		fmt.Println(schemaLoader.JsonSource())
		schema, err := gojsonschema.NewSchema(schemaLoader)
		panicIf(err)
		derp := new(ParamExampleRequest)
		if err := c.Bind(derp); err != nil {
			return err
		}
		result, err := schema.Validate(gojsonschema.NewGoLoader(derp))
		panicIf(err)
		if result.Valid() {
			return c.String(http.StatusOK, "The document is valid")
		}

		c.Response().Header().Set(
			echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		c.Response().Write([]byte(
			"The document is not valid. see errors :\n"))
		for _, err := range result.Errors() {
			// Err implements the ResultError interface
			fmt.Fprintf(c.Response(), "- %s\n", err)
		}
		return nil
	}
}

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
