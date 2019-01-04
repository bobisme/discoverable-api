package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/labstack/echo"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func swaggerMiddleware(filepath string) echo.MiddlewareFunc {
	specDoc, err := loads.Spec("./swagger.json")
	panicIf(err)
	specDoc, err = specDoc.Expanded()
	panicIf(err)
	spec := specDoc.Spec()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			if req.Method != "OPTIONS" {
				return next(c)
			}
			reqURI, err := url.PathUnescape(req.RequestURI)
			if err != nil {
				return err
			}
			searchPath := strings.TrimPrefix(reqURI, spec.BasePath)
			path, ok := spec.Paths.Paths[searchPath]
			if !ok {
				// return c.JSON(http.StatusNotFound, spec.Paths.Paths)
				return c.JSON(http.StatusNotFound, map[string]string{
					"basePath": spec.BasePath,
					"path":     searchPath,
				})
			}
			return c.JSON(200, path)
		}
	}
}

func main() {
	e := echo.New()
	e.Use(swaggerMiddleware("./swagger.json"))
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, e.Routes())
	})
	e.GET("/accounts/:id", ShowAccount)

	e.Logger.Fatal(e.Start(":1323"))
}

// HALLink
// swagger:model HALLink
type HALLink struct {
	HREF      string `json:"href"`
	Title     string `json:"title,omitempty"`
	Name      string `json:"name,omitempty"`
	Templated bool   `json:"templated,omitempty"`
}

// AccountResponseLinks
// swagger:model AccountResponseLinks
type AccountResponseLinks struct {
	Self       HALLink   `json:"self"`
	Collection HALLink   `json:"collection,omitempty"`
	Curies     []HALLink `json:"curies,omitempry"`
}

// AccountResponse
// swagger:response AccountResponse
type AccountResponse struct {
	Links AccountResponseLinks `json:"_links"`
	ID    int                  `json:"id"`
}

type Account struct {
	ID int `json:"id"`
}

// genericError is an error that is used when the required input fails
// validation.
// swagger:response genericError
type genericError struct{}

// validationError is an error that is used when the required input fails
// validation.
// swagger:response validationError
type validationError struct{}

// ListAccountsResponse
// swagger:response ListAccountsResponse
type ListAccountsResponse struct {
	Links    AccountResponseLinks `json:"_links"`
	Embedded struct {
		Items []AccountResponse `json:"items"`
	} `json:"_embedded"`
	Total int `json:"total"`
}

// ListAccounts serves the API for this record store
func ListAccounts(ctx echo.Context) error {
	// swagger:route GET /accounts accounts listAccounts
	//
	// Lists accounts filtered by some parameters.
	//
	// This will show all available accounts by default.
	//
	//     Consumes:
	//       - application/json
	//     Produces:
	//       - application/json
	//     Schemes: https
	//     Responses:
	//       default: genericError
	//       200: ListAccountsResponse
	//       422: validationError
	return ctx.JSON(http.StatusOK, []Account{
		{ID: 123},
		{ID: 234},
		{ID: 345},
	})
}

// ShowAccount serves the API for this record store
func ShowAccount(ctx echo.Context) error {
	// swagger:route GET /accounts/{id} accounts showAccount
	//
	// Get account.
	//
	// This will show the account by the given account ID.
	//
	//     Consumes:
	//     - application/json
	//     Produces:
	//     - application/json
	//     Schemes: https
	//     Responses:
	//       default: genericError
	//       200:
	//         body: AccountResponse
	//       422: validationError
	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	account := &AccountResponse{
		Links: AccountResponseLinks{
			Self:       HALLink{HREF: ctx.Request().RequestURI},
			Collection: HALLink{HREF: "/accounts"},
		},
		ID: aid,
	}
	return ctx.JSON(http.StatusOK, account)
}
