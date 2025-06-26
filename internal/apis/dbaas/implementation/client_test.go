package implementation

import (
	"net/http"
	"strings"
	"terraform-provider-infomaniak/internal/apis/dbaas"
	"terraform-provider-infomaniak/internal/apis/helpers"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	BaseURI = "https://example.com"
	Token   = "ThisIsNotARealToken"

	TestEndpointDBaaSes = `=~^/1/public_clouds/\d+/projects/\d+/dbaas\z`
	TestEndpointDBaaS   = `=~^/1/public_clouds/\d+/projects/\d+/dbaas/\d+\z`
)

func NewSuccessResponse[K any](data K) helpers.NormalizedApiResponse[K] {
	return helpers.NormalizedApiResponse[K]{
		Result: "success",
		Data:   data,
	}
}

var _ = Describe("DBaaS API Client", func() {
	Context("General Testing", func() {
		client := New(BaseURI, Token, "test")

		It("should authenticate properly", func() {
			httpmock.ActivateNonDefault(client.resty.Client())
			defer httpmock.DeactivateAndReset()

			expectedResult := &dbaas.DBaaS{
				Id: 12,
			}

			httpmock.RegisterResponder("GET", TestEndpointDBaaS, func(req *http.Request) (resp *http.Response, err error) {
				bearer, ok := req.Header["Authorization"]
				if !ok || len(bearer) == 0 {
					return httpmock.NewBytesResponse(401, []byte("not authorized")), nil
				}

				if !strings.Contains(bearer[0], Token) {
					return httpmock.NewBytesResponse(401, []byte("not authorized")), nil
				}

				return httpmock.NewJsonResponse(200, NewSuccessResponse(expectedResult))
			})

			_, err := client.GetDBaaS(1, 1, 1)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
