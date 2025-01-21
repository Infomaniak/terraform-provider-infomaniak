package endpoints

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndpointCompilation(t *testing.T) {
	var endpoint = NewEndpoint(http.MethodGet, "/api/{user_id}/{org.id}/as")

	c, err := endpoint.Compile(nil, "test1", "test2")
	if err != nil {
		t.Errorf("this error shouldn't happen : err")
	}

	assert.Equal(t, "/api/test1/test2/as", c.URL, "URL should match when compiling")
}

func TestEndpointCompilationParamMismatch(t *testing.T) {
	var endpoint = NewEndpoint(http.MethodGet, "/api/{user_id}/{org.id}/as")

	_, err := endpoint.Compile(nil, "test1")
	if err == nil {
		t.Errorf("should have errored when putting only one param instead of two")
	}
	_, err = endpoint.Compile(nil, "test2", "test1", "test3")
	if err == nil {
		t.Errorf("should have errored when putting three params instead of two")
	}
}

func TestEndpointCompilationWithParam(t *testing.T) {
	var endpoint = NewEndpoint(http.MethodGet, "/api/{user_id}/{org.id}/as")

	c, err := endpoint.Compile(map[string]any{
		"test": "tttt",
	}, "test1", "test2")
	if err != nil {
		t.Errorf("this error shouldn't happen : err")
	}

	assert.Equal(t, "/api/test1/test2/as?test=tttt", c.URL, "URL should match when compiling")
}
