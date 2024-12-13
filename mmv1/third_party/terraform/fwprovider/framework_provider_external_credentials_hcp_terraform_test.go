package fwprovider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// TestAccFwProvider_external_credentials_hcp_terraform is a series of acc tests asserting how the PF provider handles external_credentials_hcp_terraform arguments
// It is PF specific because the HCL used provisions PF-implemented resources
// It is a counterpart to TestAccSdkProvider_external_credentials_hcp_terraform
func TestAccFwProvider_external_credentials_hcp_terraform(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		"external_credentials_hcp_terraform can be set in config": testAccFwProvider_external_credentials_hcp_terraform_configSet,

		// Schema-level validation
		"external_credentials_hcp_terraform conflicts with other primary credentials fields":                 testAccFwProvider_external_credentials_hcp_terraform_conflicts,
		"external_credentials_hcp_terraform's nested fields are required and cannot be set as empty strings": testAccFwProvider_external_credentials_hcp_terraform_requiredValues,

		// Usage
		// No usage test implemented - this needs to be done with a simple PF-implemented resource/data source
		// and all we have currently is Firebase data sources
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccFwProvider_external_credentials_hcp_terraform_configSet(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{
		"audience":              "foobar-audience",
		"service_account_email": "foobar-service_account_email",
		"identity_token":        "foobar-identity_token",
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFwProvider_external_credentials_hcp_terraformInProviderBlock(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "external_credentials_hcp_terraform.audience", context["audience"].(string)),
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "external_credentials_hcp_terraform.service_account_email", context["service_account_email"].(string)),
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "external_credentials_hcp_terraform.identity_token", context["identity_token"].(string)),
				),
			},
		},
	})
}

func testAccFwProvider_external_credentials_hcp_terraform_conflicts(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	contextWithAccessToken := map[string]interface{}{
		"fields": `
external_credentials_hcp_terraform {
  audience = "foo"
  service_account_email = "foo"
  identity_token = "foo"
}
access_token = "foo"
`,
	}

	contextWithCredentials := map[string]interface{}{
		"fields": `
external_credentials_hcp_terraform {
  audience = "foo"
  service_account_email = "foo"
  identity_token = "foo"
}
credentials = "foo"
`,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// You cannot set both external_credentials_hcp_terraform and access_token
				Config:      testAccFwProvider_external_credentials_hcp_terraform_setFields(contextWithAccessToken),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
			{
				// You cannot set both external_credentials_hcp_terraform and credentials
				Config:      testAccFwProvider_external_credentials_hcp_terraform_setFields(contextWithCredentials),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
		},
	})
}

func testAccFwProvider_external_credentials_hcp_terraform_requiredValues(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	contextEmptyAudience := map[string]interface{}{
		"audience":              "",
		"service_account_email": "foobar-service_account_email",
		"identity_token":        "foobar-identity_token",
	}

	contextWithAudienceUnset := map[string]interface{}{
		"fields": `
external_credentials_hcp_terraform {
  service_account_email = "foo"
  identity_token = "foo"
}
`,
	}

	contextEmptyEmail := map[string]interface{}{
		"audience":              "foobar-audience",
		"service_account_email": "",
		"identity_token":        "foobar-identity_token",
	}

	contextWithEmailUnset := map[string]interface{}{
		"fields": `
external_credentials_hcp_terraform {
  audience = "foo"
  identity_token = "foo"
}
`,
	}

	contextEmptyIdentityToken := map[string]interface{}{
		"audience":              "foobar-audience",
		"service_account_email": "foobar-service_account_email",
		"identity_token":        "",
	}

	contextWithIdentityTokenUnset := map[string]interface{}{
		"fields": `
external_credentials_hcp_terraform {
  audience = "foo"
  service_account_email = "foo"
}
`,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// external_credentials_hcp_terraform.audience cannot be an empty string
				Config:      testAccFwProvider_external_credentials_hcp_terraformInProviderBlock(contextEmptyAudience),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("audience was set to ``"),
			},
			{
				// external_credentials_hcp_terraform.audience is Required
				Config:      testAccFwProvider_external_credentials_hcp_terraform_setFields(contextWithAudienceUnset),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("The argument \"audience\" is required"),
			},
			{
				// external_credentials_hcp_terraform.service_account_email cannot be an empty string
				Config:      testAccFwProvider_external_credentials_hcp_terraformInProviderBlock(contextEmptyEmail),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("service_account_email was set to ``"),
			},
			{
				// external_credentials_hcp_terraform.service_account_email is Required
				Config:      testAccFwProvider_external_credentials_hcp_terraform_setFields(contextWithEmailUnset),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("The argument \"service_account_email\" is required"),
			},
			{
				// external_credentials_hcp_terraform.identity_token cannot be an empty string
				Config:      testAccFwProvider_external_credentials_hcp_terraformInProviderBlock(contextEmptyIdentityToken),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("identity_token was set to ``"),
			},
			{
				// external_credentials_hcp_terraform.identity_token is Required
				Config:      testAccFwProvider_external_credentials_hcp_terraform_setFields(contextWithIdentityTokenUnset),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("The argument \"identity_token\" is required"),
			},
		},
	})
}

// testAccFwProvider_external_credentials_hcp_terraformInProviderBlock allows setting the external_credentials_hcp_terraform argument in a provider block.
// This function uses data.google_provider_config_plugin_framework because it is implemented with the PF
func testAccFwProvider_external_credentials_hcp_terraformInProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  external_credentials_hcp_terraform {
    audience = "%{audience}"
    service_account_email = "%{service_account_email}"
    identity_token = "%{identity_token}"
  }
}
data "google_provider_config_plugin_framework" "default" {}
`, context)
}

// testAccFwProvider_external_credentials_hcp_terraform_setFields allows setting multiple fields in the provider
// block to test conflict validation in the provider schema
func testAccFwProvider_external_credentials_hcp_terraform_setFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
%{fields}
}
data "google_provider_config_plugin_framework" "default" {}
`, context)
}