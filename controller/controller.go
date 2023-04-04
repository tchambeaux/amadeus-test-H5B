package controller

import (
	"net/http"
	"test_H5B/output"

	"github.com/labstack/echo/v4"
)

func New(cache map[string]*output.Search) *Controller {
	return &Controller{
		MemCache: cache,
	}
}

// search controller method
func (control *Controller) SearchGet(c echo.Context) error {
	var searchWord string
	out := &output.Search{
		LineOccurrences: []int{},
	}
	if err := echo.PathParamsBinder(c).String("searchWord", &searchWord).BindError(); err != nil {
		// Info: On fails should be 500 => this error should never occur
		return c.JSON(http.StatusInternalServerError, output.Error{Error: err.Error()})
	}
	if _, found := control.MemCache[searchWord]; found {
		out = control.MemCache[searchWord]
	}
	return c.JSON(http.StatusOK, out)
}
