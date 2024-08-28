package fwprovider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccFwProvider_credentials(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"credentials can be configured as a path to a credentials JSON file":                                       testAccFwProvider_credentials_validJsonFilePath,
		"configuring credentials as a path to a non-existent file results in an error":                             testAccFwProvider_credentials_badJsonFilepathCausesError,
		"config takes precedence over environment variables":                                                       testAccFwProvider_credentials_configPrecedenceOverEnvironmentVariables,
		"when credentials is unset in the config, environment variables are used in a given order":                 testAccFwProvider_credentials_precedenceOrderEnvironmentVariables, // GOOGLE_CREDENTIALS, GOOGLE_CLOUD_KEYFILE_JSON, GCLOUD_KEYFILE_JSON, GOOGLE_APPLICATION_CREDENTIALS
		"when credentials is set to an empty string in the config the value isn't ignored and results in an error": testAccFwProvider_credentials_emptyStringValidation,
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

func testAccFwProvider_credentials_validJsonFilePath(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	// unset all credentials env vars
	for _, v := range envvar.CredsEnvVars {
		t.Setenv(v, "")
	}

	credentials := transport_tpg.TestFakeCredentialsPath

	context := map[string]interface{}{
		"credentials":   credentials,
		"resource_name": "tf-test-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:             testAccFwProvider_credentialsInProviderBlock(context),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFwProvider_credentials_badJsonFilepathCausesError(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	// unset all credentials env vars
	for _, v := range envvar.CredsEnvVars {
		t.Setenv(v, "")
	}

	pathToMissingFile := "./this/path/does/not/exist.json" // Doesn't exist

	context := map[string]interface{}{
		"credentials":   pathToMissingFile,
		"resource_name": "tf-test-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// credentials is a path to a json, but if that file doesn't exist so there's an error
				Config:      testAccFwProvider_credentialsInProviderBlock(context),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("JSON credentials are not valid: invalid character '.' looking for beginning of value"),
			},
		},
	})
}

func testAccFwProvider_credentials_configPrecedenceOverEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	credentials := envvar.GetTestCredsFromEnv()

	// ensure all possible credentials env vars set; show they aren't used
	for _, v := range envvar.CredsEnvVars {
		t.Setenv(v, credentials)
	}

	pathToMissingFile := "./this/path/does/not/exist.json" // Doesn't exist

	context := map[string]interface{}{
		"credentials":   pathToMissingFile,
		"resource_name": "tf-test-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccFwProvider_credentialsInProviderBlock(context),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("JSON credentials are not valid: invalid character '.' looking for beginning of value"),
			},
		},
	})
}

func testAccFwProvider_credentials_precedenceOrderEnvironmentVariables(t *testing.T) {
	/*
		These are all the ENVs for credentials, and they are in order of precedence.
		GOOGLE_CREDENTIALS
		GOOGLE_CLOUD_KEYFILE_JSON
		GCLOUD_KEYFILE_JSON
		GOOGLE_APPLICATION_CREDENTIALS
		GOOGLE_USE_DEFAULT_CREDENTIALS
	*/

	goodCredentials := envvar.GetTestCredsFromEnv()
	badCreds := acctest.GenerateFakeCredentialsJson("test")

	context := map[string]interface{}{
		"resource_name": "tf-test-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Error as all ENVs set to 'bad' creds
				PreConfig: func() {
					for _, v := range envvar.CredsEnvVars {
						t.Setenv(v, badCreds)
					}
				},
				Config:      testAccFwProvider_credentialsInEnvsOnly(context),
				ExpectError: regexp.MustCompile("private key should be a PEM or plain PKCS1 or PKCS8"),
			},
			{
				// GOOGLE_CREDENTIALS is used 1st if set
				PreConfig: func() {
					// good
					t.Setenv("GOOGLE_CREDENTIALS", goodCredentials) //used
					// bad
					t.Setenv("GOOGLE_CLOUD_KEYFILE_JSON", badCreds)
					t.Setenv("GCLOUD_KEYFILE_JSON", badCreds)
					t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", badCreds)
				},
				Config: testAccFwProvider_credentialsInEnvsOnly(context),
			},
			{
				// GOOGLE_CLOUD_KEYFILE_JSON is used 2nd
				PreConfig: func() {
					// unset
					t.Setenv("GOOGLE_CREDENTIALS", "")
					// good
					t.Setenv("GOOGLE_CLOUD_KEYFILE_JSON", goodCredentials) //used
					// bad
					t.Setenv("GCLOUD_KEYFILE_JSON", badCreds)
					t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", badCreds)

				},
				Config: testAccFwProvider_credentialsInEnvsOnly(context),
			},
			{
				// GOOGLE_CLOUD_KEYFILE_JSON is used 3rd
				PreConfig: func() {
					// unset
					t.Setenv("GOOGLE_CREDENTIALS", "")
					t.Setenv("GOOGLE_CLOUD_KEYFILE_JSON", "")
					// good
					t.Setenv("GCLOUD_KEYFILE_JSON", goodCredentials) //used
					// bad
					t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", badCreds)
				},
				Config: testAccFwProvider_credentialsInEnvsOnly(context),
			},
			{
				// GOOGLE_APPLICATION_CREDENTIALS is used 4th
				PreConfig: func() {
					// unset
					t.Setenv("GOOGLE_CREDENTIALS", "")
					t.Setenv("GOOGLE_CLOUD_KEYFILE_JSON", "")
					t.Setenv("GCLOUD_KEYFILE_JSON", "")
					// good
					t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", goodCredentials) //used
				},
				Config: testAccFwProvider_credentialsInEnvsOnly(context),
			},
		},
	})
}

func testAccFwProvider_credentials_emptyStringValidation(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	credentials := envvar.GetTestCredsFromEnv()

	// ensure all credentials env vars set
	for _, v := range envvar.CredsEnvVars {
		t.Setenv(v, credentials)
	}

	context := map[string]interface{}{
		"credentials":   "", // empty string used
		"resource_name": "tf-test-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccFwProvider_credentialsInProviderBlock(context),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected a non-empty string"),
			},
		},
	})
}

func testAccFwProvider_credentialsInProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	credentials = "%{credentials}"
}

// We provision something in the test so that VCR works; need an API interaction
resource "google_service_account" "default" {
  account_id   = "%{resource_name}"
  display_name = "Testing, provisioned by testAccFwProvider_credentialsInProviderBlock"
}
`, context)
}

func testAccFwProvider_credentialsInEnvsOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
// We provision something in the test so that VCR works; need an API interaction
resource "google_service_account" "default" {
  account_id   = "%{resource_name}"
  display_name = "Testing, provisioned by testAccFwProvider_credentialsInProviderBlock"
}
`, context)
}
