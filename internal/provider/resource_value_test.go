package provider

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceScaffolding(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "metadata_value" "foo" {
					inputs = {
						value = 11
					}
					update = false
				}
				resource "metadata_value" "bar" {
					inputs = {
						value = 11
					}
					update = true
				}
				`,
				Check: func(s *terraform.State) error {
					outputs, err := getOutputs("foo")(s)
					if err != nil {
						return err
					}
					if len(outputs) != 0 {
						return fmt.Errorf("wrong number of outputs")
					}
					outputs, err = getOutputs("bar")(s)
					if err != nil {
						return err
					}
					if len(outputs) != 1 {
						return fmt.Errorf("wrong number of outputs")
					}
					if value, ok := outputs["value"]; !ok || value != "11" {
						return fmt.Errorf("outputs value is wrong")
					}
					return nil
				},
			},
		},
	})
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				// resource "metadata_value" "foo" {
				resource "metadata_value" "foo" {
					inputs = {
						value = 12
					}
					update = true
				}
				`,
				Check: func(s *terraform.State) error {
					outputs, err := getOutputs("foo")(s)
					if err != nil {
						return err
					}
					if len(outputs) != 1 {
						return fmt.Errorf("wrong number of outputs")
					}
					if value, ok := outputs["value"]; !ok || value != "12" {
						return fmt.Errorf("outputs value is wrong")
					}
					return nil
				},
			},
			{
				Config: `
				// resource "metadata_value" "foo" {
				resource "metadata_value" "foo" {
					inputs = {
						value = 13
					}
					update = false
				}
				`,
				Check: func(s *terraform.State) error {
					outputs, err := getOutputs("foo")(s)
					if err != nil {
						return err
					}
					if len(outputs) != 1 {
						return fmt.Errorf("wrong number of outputs")
					}
					if value, ok := outputs["value"]; !ok || value != "12" {
						return fmt.Errorf("outputs value is wrong")
					}
					return nil
				},
			},
			{
				Config: `
				// resource "metadata_value" "foo" {
				resource "metadata_value" "foo" {
					inputs = {
						value = 13
					}
					update = true
				}
				`,
				Check: func(s *terraform.State) error {
					outputs, err := getOutputs("foo")(s)
					if err != nil {
						return err
					}
					if len(outputs) != 1 {
						return fmt.Errorf("wrong number of outputs")
					}
					if value, ok := outputs["value"]; !ok || value != "13" {
						return fmt.Errorf("outputs value is wrong")
					}
					return nil
				},
			},
		},
	})
}

func getOutputs(resourceName string) func(s *terraform.State) (map[string]string, error) {
	resourceName = "metadata_value." + resourceName
	return func(s *terraform.State) (map[string]string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return nil, fmt.Errorf("resource not found")
		}
		if rs.Primary.ID == "" {
			return nil, fmt.Errorf("no ID set")
		}

		attrs := rs.Primary.Attributes
		dest := make(map[string]string)
		rawCount, ok := attrs["outputs.%"]
		if !ok {
			return nil, fmt.Errorf("outputs not set")
		}
		count, err := strconv.ParseUint(rawCount, 10, 64)
		if err != nil {
			return nil, err
		}

		for key, value := range attrs {
			if strings.HasPrefix(key, "outputs.") && key != "outputs.%" {
				dest[key[8:]] = value
				if uint64(len(dest)) == count {
					break
				}
			}
		}
		if uint64(len(dest)) != count {
			return nil, fmt.Errorf("unable to identify all the output values")
		}
		return dest, nil
	}
}
