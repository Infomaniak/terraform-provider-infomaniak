package registry

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var resources []func() resource.Resource
var datasources []func() datasource.DataSource

func RegisterResource(F func() resource.Resource) {
	resources = append(resources, F)
}

func RegisterDataSource(F func() datasource.DataSource) {
	datasources = append(datasources, F)
}

func GetResources() []func() resource.Resource {
	return resources
}

func GetDataSources() []func() datasource.DataSource {
	return datasources
}
