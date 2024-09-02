package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// TestAccSdkProvider_billing_project is a series of acc tests asserting how the SDK provider handles billing_project arguments
// It is SDK specific because the HCL used provisions SDK-implemented resources
// It is a counterpart to TestAccFwProvider_billing_project
func TestAccSdkProvider_billing_project(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"config takes precedence over environment variables":                                                           testAccSdkProvider_billing_project_configPrecedenceOverEnvironmentVariables,
		"when billing_project is unset in the config, environment variables are used in a given order":                 testAccSdkProvider_billing_project_precedenceOrderEnvironmentVariables, // GOOGLE_BILLING_PROJECT
		"when billing_project is set to an empty string in the config the value isn't ignored and results in an error": testAccSdkProvider_billing_project_emptyStringValidation,
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

func testAccSdkProvider_billing_project_configPrecedenceOverEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	billingProject := "my-billing-project-id"

	// ensure all possible billing_project env vars set; show they aren't used instead
	t.Setenv("GOOGLE_BILLING_PROJECT", billingProject)

	providerBillingProject := "foobar"

	context := map[string]interface{}{
		"billing_project": providerBillingProject,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Apply-time error; bad value in config is used over of good values in ENVs
				Config: testAccSdkProvider_billing_projectInProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "billing_project", providerBillingProject),
				)},
		},
	})
}

func testAccSdkProvider_billing_project_precedenceOrderEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API
	/*
		These are all the ENVs for billing_project
		GOOGLE_BILLING_PROJECT
	*/

	GOOGLE_BILLING_PROJECT := "GOOGLE_BILLING_PROJECT"

	context := map[string]interface{}{}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// GOOGLE_BILLING_PROJECT is used if set
				PreConfig: func() {
					t.Setenv("GOOGLE_BILLING_PROJECT", GOOGLE_BILLING_PROJECT) //used
				},
				Config: testAccSdkProvider_billing_projectInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "billing_project", GOOGLE_BILLING_PROJECT),
				),
			},
		},
	})
}

func testAccSdkProvider_billing_project_emptyStringValidation(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	billingProject := "my-billing-project-id"

	// ensure all billing_project env vars set
	t.Setenv("GOOGLE_BILLING_PROJECT", billingProject)

	context := map[string]interface{}{
		"billing_project": "", // empty string used
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSdkProvider_billing_projectInProviderBlock(context),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected a non-empty string"),
			},
		},
	})
}

// testAccSdkProvider_billing_projectInProviderBlock allows setting the billing_project argument in a provider block.
// This function uses data.google_provider_config_sdk because it is implemented with the SDKv2
func testAccSdkProvider_billing_projectInProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	billing_project = "%{billing_project}"
}

data "google_provider_config_sdk" "default" {}

output "billing_project" {
  value = data.google_provider_config_sdk.default.billing_project
  sensitive = true
}
`, context)
}

// testAccSdkProvider_billing_projectInEnvsOnly allows testing when the billing_project argument
// is only supplied via ENVs
func testAccSdkProvider_billing_projectInEnvsOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_provider_config_sdk" "default" {}

output "billing_project" {
  value = data.google_provider_config_sdk.default.billing_project
  sensitive = true
}
`, context)
}
