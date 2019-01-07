package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
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

func underline(s string, b byte, writer io.Writer) {
	_, err := writer.Write([]byte(s))
	panicIf(err)
	_, err = writer.Write([]byte{'\n'})
	panicIf(err)
	underline := []byte{b}
	for i := 0; i < len(s); i++ {
		_, err := writer.Write(underline)
		panicIf(err)
	}
	_, err = writer.Write([]byte("\n\n"))
	panicIf(err)
}

func h1(s string, writer io.Writer) {
	underline(s, '=', writer)
}

func h2(s string, writer io.Writer) {
	underline(s, '-', writer)
}

func h3(s string, writer io.Writer) {
	writer.Write([]byte("### "))
	writer.Write([]byte(s))
	writer.Write([]byte("\n\n"))
}

func specPath2md(uri string, method string, op *spec.Operation, out io.Writer) {
	writeString := func(s string) {
		out.Write([]byte(s))
	}
	h2(fmt.Sprintf("%s %s", method, uri), out)
	fmt.Fprintf(out, "%s\n\n", op.Summary)
	fmt.Fprintf(out, "%s\n\n", op.Description)

	{
		print := func(label string, things []string) {
			if len(things) > 0 {
				fmt.Fprintf(out, "%s: %s\n", label, strings.Join(op.Schemes, ", "))
			}
		}
		print("Schemes", op.Schemes)
		print("Consumes", op.Schemes)
		print("Produces", op.Schemes)
		print("Tags", op.Schemes)
	}
	fmt.Fprintln(out, "")

	h3("Parameters", out)
	if op.Parameters == nil {
		fmt.Fprintf(out, "*No parameters.*\n\n")
	} else {
		// out.Write([]byte("    "))
		// paramsJSON, err := yaml.MarshalIndent(op.Parameters, "    ", "  ")
		paramscode, err := yaml.Marshal(op.Parameters)
		panicIf(err)
		writeString("```yaml\n")
		out.Write(paramscode)
		writeString("```\n\n")
	}

	h3("Responses", out)
	for status, resp := range op.Responses.StatusCodeResponses {
		fmt.Fprintf(out, "#### %d: %s\n\n", status, http.StatusText(status))
		// writeString("    ")
		// respCode, err := json.MarshalIndent(resp, "    ", "  ")
		respCode, err := yaml.Marshal(resp)
		panicIf(err)
		writeString("```yaml\n")
		out.Write(respCode)
		writeString("```\n\n")
	}
	// &{VendorExtensible:{Extensions:map[]} OperationProps:{Description:This will show all available accounts by default. Consumes:[application/json] Produces:[application/json] Schemes:[https] Tags:[accounts] Summary:Lists accounts filtered by some parameters. ExternalDocs:<nil> ID:listAccounts Deprecated:false Security:[] Parameters:[] Responses:0xc00000d840}}⏎
}

func spec2md(filepath string) {
	specDoc, err := loads.Spec(filepath)
	panicIf(err)
	specDoc, err = specDoc.Expanded()
	panicIf(err)
	s := specDoc.Spec()
	basePath := "/api/v1"
	// searchPath := "/accounts"
	// fullPath := basePath + searchPath

	out, err := os.Create("./docs.md")
	panicIf(err)
	defer out.Close()

	for uri, path := range s.Paths.Paths {
		h1(uri, out)
		ops := map[string]*spec.Operation{
			"GET":     path.Get,
			"POST":    path.Post,
			"DELETE":  path.Delete,
			"PUT":     path.Put,
			"HEAD":    path.Head,
			"OPTIONS": path.Options,
		}
		for method, op := range ops {
			if op != nil {
				specPath2md(basePath+uri, method, op, out)
			}
		}
	}
}

func main() {
	spec2md("../swag/docs/swagger/swagger.json")
	log.Println("wrote docs to docs.md")

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
