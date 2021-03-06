package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceValue() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Value resource for the metadata provider.",

		CreateContext: resourceValueCreate,
		Read:          schema.Noop,
		UpdateContext: resourceValueUpdate,
		Delete:        schema.RemoveFromState,

		Schema: map[string]*schema.Schema{
			"update": {
				Description: "Whether to update the metadata to the current inputs value.",
				Type:        schema.TypeBool,
				Optional:    true,
				// Technically, this is a given value by the config. This just lets us clear it from
				// state.
				Computed: true,
			},

			"inputs": {
				Description: "The metadata to store.",
				Type:        schema.TypeMap,
				Optional:    true,
				// Technically, this is a given value by the config. This just lets us clear it from
				// state.
				Computed: true,
			},

			"outputs": {
				Description: "The stored metadata.",
				Type:        schema.TypeMap,
				Computed:    true,
			},
		},

		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			// TODO: handle apply-time inputs.
			if d.Get("update") == true {
				d.SetNew("outputs", d.Get("inputs"))
			} else if d.Get("outputs") == nil {
				d.SetNew("outputs", map[string]string{})
			}
			d.Clear("update")
			d.Clear("inputs")
			return nil
		},
	}
}

func resourceValueCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Not used to uniquely identify the resource, so the value doesn't matter.
	d.SetId("0")

	if !d.HasChange("outputs") {
		if d.Get("update") == true {
			d.Set("outputs", d.Get("inputs"))
		} else {
			d.Set("outputs", map[string]string{})
		}
	}
	d.Set("update", nil)
	d.Set("inputs", nil)

	return nil
}

func resourceValueUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange("outputs") {
		if d.Get("update") == true {
			d.Set("outputs", d.Get("inputs"))
		} else {
			prior_outputs, ok := d.GetOk("outputs")
			if !ok {
				d.Set("outputs", map[string]string{})
			} else {
				d.Set("outputs", prior_outputs)
			}
		}
	}
	d.Set("update", nil)
	d.Set("inputs", nil)

	return nil
}
