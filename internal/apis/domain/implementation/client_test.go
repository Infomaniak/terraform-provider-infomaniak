package implementation

// import (
// 	"net/http"
// 	"strings"
// 	"terraform-provider-infomaniak/internal/apis/helpers"
// 	"terraform-provider-infomaniak/internal/apis/kaas"

// 	"github.com/jarcoal/httpmock"
// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"
// )

// const (
// 	BaseURI = "https://example.com"
// 	Token   = "ThisIsNotARealToken"

// 	TestEndpointKaases         = `=~^/1/public_clouds/\d+/projects/\d+/kaas\z`
// 	TestEndpointKaas           = `=~^/1/public_clouds/\d+/projects/\d+/kaas/\d+\z`
// 	TestEndpointKaasKubeconfig = `=~^/1/public_clouds/\d+/projects/\d+/kaas/\d+/kube_config\z`
// 	TestEndpointInstancePools  = `=~^/1/public_clouds/\d+/projects/\d+/kaas/\d+/instance_pools\z`
// 	TestEndpointInstancePool   = `=~^/1/public_clouds/\d+/projects/\d+/kaas/\d+/instance_pools/\d+\z`
// )

// func NewSuccessResponse[K any](data K) helpers.NormalizedApiResponse[K] {
// 	return helpers.NormalizedApiResponse[K]{
// 		Result: "success",
// 		Data:   data,
// 	}
// }

// var _ = Describe("KaaS API Client", func() {
// 	Context("General Testing", func() {
// 		client := New(BaseURI, Token, "test")

// 		It("should authenticate properly", func() {
// 			httpmock.ActivateNonDefault(client.resty.Client())
// 			defer httpmock.DeactivateAndReset()

// 			expectedResult := &kaas.Kaas{
// 				Id: 12,
// 			}

// 			httpmock.RegisterResponder("GET", TestEndpointKaas, func(req *http.Request) (resp *http.Response, err error) {
// 				bearer, ok := req.Header["Authorization"]
// 				if !ok || len(bearer) == 0 {
// 					return httpmock.NewBytesResponse(401, []byte("not authorized")), nil
// 				}

// 				if !strings.Contains(bearer[0], Token) {
// 					return httpmock.NewBytesResponse(401, []byte("not authorized")), nil
// 				}

// 				return httpmock.NewJsonResponse(200, NewSuccessResponse(expectedResult))
// 			})

// 			_, err := client.GetKaas(1, 1, 1)
// 			Expect(err).ShouldNot(HaveOccurred())
// 		})

// 		It("should be able to get KaaS", func() {
// 			httpmock.ActivateNonDefault(client.resty.Client())
// 			defer httpmock.DeactivateAndReset()

// 			expectedResult := &kaas.Kaas{
// 				Id: 12,
// 			}

// 			httpmock.RegisterResponder("GET", TestEndpointKaas, httpmock.NewJsonResponderOrPanic(200, NewSuccessResponse(expectedResult)))

// 			kaas, err := client.GetKaas(1, 1, 12)
// 			Expect(err).ShouldNot(HaveOccurred())

// 			Expect(kaas.Id).To(Equal(expectedResult.Id))
// 		})

// 		It("should be able to create KaaS", func() {
// 			httpmock.ActivateNonDefault(client.resty.Client())
// 			defer httpmock.DeactivateAndReset()

// 			expectedResult := 12

// 			httpmock.RegisterResponder("POST", TestEndpointKaases, httpmock.NewJsonResponderOrPanic(200, NewSuccessResponse(expectedResult)))

// 			kaasId, err := client.CreateKaas(&kaas.Kaas{
// 				Project: kaas.KaasProject{
// 					PublicCloudId: 8,
// 					ProjectId:     546,
// 				},
// 			})
// 			Expect(err).ShouldNot(HaveOccurred())

// 			Expect(kaasId).To(Equal(expectedResult))
// 		})

// 		It("should err when the body is of the wrong type", func() {
// 			httpmock.ActivateNonDefault(client.resty.Client())
// 			defer httpmock.DeactivateAndReset()

// 			expectedResult := &kaas.Kaas{
// 				Id: 12,
// 			}

// 			httpmock.RegisterResponder("POST", TestEndpointKaases, httpmock.NewJsonResponderOrPanic(200, NewSuccessResponse(expectedResult)))

// 			_, err := client.CreateKaas(&kaas.Kaas{
// 				Project: kaas.KaasProject{
// 					PublicCloudId: 8,
// 					ProjectId:     546,
// 				},
// 			})
// 			Expect(err).Should(HaveOccurred())
// 		})

// 		It("should be able to get KaaS Instance Pool", func() {
// 			httpmock.ActivateNonDefault(client.resty.Client())
// 			defer httpmock.DeactivateAndReset()

// 			expectedResult := &kaas.InstancePool{
// 				Id: 12,
// 			}

// 			httpmock.RegisterResponder("GET", TestEndpointInstancePool, httpmock.NewJsonResponderOrPanic(200, NewSuccessResponse(expectedResult)))

// 			instancePool, err := client.GetInstancePool(1, 1, 12, 12)
// 			Expect(err).ShouldNot(HaveOccurred())

// 			Expect(instancePool.Id).To(Equal(expectedResult.Id))
// 		})
// 	})
// })
