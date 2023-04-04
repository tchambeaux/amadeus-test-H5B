package controller_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"test_H5B/controller"
	"test_H5B/output"
)

type response struct {
	Code int
	Body []byte
}

type request struct {
	Method      string
	ParamsName  []string
	ParamsValue []string
	Handler     func(ec echo.Context) error
}

type tester struct {
	Input  request
	Output response
}

var validController = &controller.Controller{
	MemCache: map[string]*output.Search{
		"test": {
			WordFound:       true,
			NumOccurrences:  42,
			LineOccurrences: []int{42},
		},
	},
}

var testSearchGet = map[string]tester{
	"ok with result": {
		Input: request{
			Method:      echo.GET,
			ParamsName:  []string{"searchWord"},
			ParamsValue: []string{"test"},
			Handler:     validController.SearchGet,
		},
		Output: response{
			Code: http.StatusOK,
			Body: []byte("{\"wordFound\":true,\"numOccurrences\":42,\"lineOccurrences\":[42]}\n"),
		},
	},
	"ok without result": {
		Input: request{
			Method:      echo.GET,
			ParamsName:  []string{"searchWord"},
			ParamsValue: []string{"noresult"},
			Handler:     validController.SearchGet,
		},
		Output: response{
			Code: http.StatusOK,
			Body: []byte("{\"wordFound\":false,\"numOccurrences\":0,\"lineOccurrences\":[]}\n"),
		},
	},
}

func doRequest(i request) (*httptest.ResponseRecorder, error) {
	e := echo.New()
	req, err := http.NewRequest(i.Method, "/", nil)
	if err != nil {
		return nil, err
	}

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames(i.ParamsName...)
	c.SetParamValues(i.ParamsValue...)

	err = i.Handler(c)
	he, ok := err.(*echo.HTTPError)
	if ok {
		rec.Result().StatusCode = he.Code
	}

	return rec, nil
}

func execute(t *testing.T, v tester) {
	resp, err := doRequest(v.Input)
	if err != nil {
		t.Fatalf("Test error occured")
	}
	assert.Equal(t, v.Output.Code, resp.Result().StatusCode)
	data, _ := io.ReadAll(resp.Result().Body)
	assert.Equal(t, v.Output.Body, data)
}

func TestSearchGet(t *testing.T) {
	for k, v := range testSearchGet {
		t.Run(k, func(t *testing.T) {
			execute(t, v)
		})
	}
}

func TestNewController(t *testing.T) {
	out := controller.New(validController.MemCache)
	assert.Equal(t, validController, out)
}
