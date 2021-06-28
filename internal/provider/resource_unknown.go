package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Temporary for testing (maybe there's a way to make this more test-oriented and keep it?)
func resourceUnknown() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Passthrough but make it unknown until apply-time.",

		CreateContext: resourceUnknownUpsert,
		Read:          schema.Noop,
		UpdateContext: resourceUnknownUpsert,
		Delete:        schema.RemoveFromState,

		Schema: map[string]*schema.Schema{
			"input": {
				Description: "Value to produce at creation.",
				Type:        schema.TypeMap,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"result": {
				Description: "Only known at creation.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"id": {
				Description: "A static value used internally by Terraform, this should not be referenced in configurations.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceUnknownUpsert(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.Set("result", d.Get("input"))
	d.SetId("-")

	return nil
}
