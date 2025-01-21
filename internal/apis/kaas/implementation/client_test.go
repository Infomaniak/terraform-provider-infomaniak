package implementation

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha"
	"github.com/vitorsalgado/mocha/expect"
	"github.com/vitorsalgado/mocha/reply"
)

func Test_ClientMockServer(t *testing.T) {
	m := mocha.New(t)
	m.Start()
	scope := m.AddMocks(
		mocha.Request().
			Method(http.MethodGet).
			URL(expect.URLPath("/1/public_clouds/test1/kaas/test2/")).
			Reply(reply.Created().BodyString("hello world")),
	)

	client := New(m.URL())

	compiledRoute, err := GetKaas.Compile(nil, "test1", "test2")
	if err != nil {
		t.Fatalf("got error when compiling route : %v", err)
	}

	_, err = client.Do(compiledRoute, nil)
	if err != nil {
		t.Fatalf("got error when sending request : %v", err)
	}

	assert.True(t, scope.Called())
}
