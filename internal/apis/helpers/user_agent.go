package helpers

import "fmt"

const (
	userAgentFormat = "terraform-provider-infomaniak/%s (resty; +https://github.com/Infomaniak/terraform-provider-infomaniak)"
)

func GetUserAgent(version string) string {
	return fmt.Sprintf(userAgentFormat, version)
}
