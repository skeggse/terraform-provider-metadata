package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

func New() func() *schema.Provider {
	return func() *schema.Provider {
		return &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				"metadata_value": resourceValue(),
			},
		}
	}
}
