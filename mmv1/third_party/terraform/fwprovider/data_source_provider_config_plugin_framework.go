package fwprovider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	"github.com/hashicorp/terraform-provider-google/google/fwresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Ensure the data source satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &GoogleProviderConfigPluginFrameworkDataSource{}
	_ datasource.DataSourceWithConfigure = &GoogleProviderConfigPluginFrameworkDataSource{}
	_ fwresource.LocationDescriber       = &GoogleProviderConfigPluginFrameworkModel{}
)

func NewGoogleProviderConfigPluginFrameworkDataSource() datasource.DataSource {
	return &GoogleProviderConfigPluginFrameworkDataSource{}
}

type GoogleProviderConfigPluginFrameworkDataSource struct {
	providerConfig *transport_tpg.Config
}

type GoogleProviderConfigPluginFrameworkModel struct {
	Credentials                        string   `tfsdk:"credentials"`
	AccessToken                        string   `tfsdk:"access_token"`
	ImpersonateServiceAccount          string   `tfsdk:"impersonate_service_account"`
	ImpersonateServiceAccountDelegates []string `tfsdk:"impersonate_service_account_delegates"`
	Project                            string   `tfsdk:"project"`
	BillingProject                     string   `tfsdk:"billing_project"`
	Region                             string   `tfsdk:"region"`
	Zone                               string   `tfsdk:"zone"`
	Scopes                             []string `tfsdk:"scopes"`
	//	omit Batching
	UserProjectOverride                       bool              `tfsdk:"user_project_override"`
	RequestReason                             string            `tfsdk:"request_reason"`
	RequestTimeout                            time.Duration     `tfsdk:"request_timeout"`
	UniverseDomain                            string            `tfsdk:"universe_domain"`
	DefaultLabels                             map[string]string `tfsdk:"default_labels"`
	AddTerraformAttributionLabel              bool              `tfsdk:"add_terraform_attribution_label"`
	TerraformAttributionLabelAdditionStrategy string            `tfsdk:"terraform_attribution_label_addition_strategy"`
}

func (m *GoogleProviderConfigPluginFrameworkModel) GetLocationDescription(providerConfig *transport_tpg.Config) fwresource.LocationDescription {
	return fwresource.LocationDescription{
		RegionSchemaField: types.StringValue("region"),
		ZoneSchemaField:   types.StringValue("zone"),
		ProviderRegion:    types.StringValue(providerConfig.Region),
		ProviderZone:      types.StringValue(providerConfig.Zone),
	}
}

func (d *GoogleProviderConfigPluginFrameworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_config_plugin_framework"
}

func (d *GoogleProviderConfigPluginFrameworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		Description:         "Use this data source to access the configuration of the Google Cloud provider. This data source is implemented with the SDK.",
		MarkdownDescription: "Use this data source to access the configuration of the Google Cloud provider. This data source is implemented with the SDK.",
		Attributes: map[string]schema.Attribute{
			// Start of user inputs
			"access_token": schema.StringAttribute{
				Description:         "The access_token argument used to configure the provider",
				MarkdownDescription: "The access_token argument used to configure the provider",
				Computed:            true,
				Sensitive:           true,
			},
			"credentials": schema.StringAttribute{
				Description:         "The credentials argument used to configure the provider",
				MarkdownDescription: "The credentials argument used to configure the provider",
				Computed:            true,
				Sensitive:           true,
			},
			"impersonate_service_account": schema.StringAttribute{
				Description:         "The impersonate_service_account argument used to configure the provider",
				MarkdownDescription: "The impersonate_service_account argument used to configure the provider.",
				Computed:            true,
			},
			"impersonate_service_account_delegates": schema.ListAttribute{
				ElementType:         types.StringType,
				Description:         "The impersonate_service_account_delegates argument used to configure the provider",
				MarkdownDescription: "The impersonate_service_account_delegates argument used to configure the provider.",
				Computed:            true,
			},
			"project": schema.StringAttribute{
				Description:         "The project argument used to configure the provider",
				MarkdownDescription: "The project argument used to configure the provider.",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				Description:         "The region argument used to configure the provider.",
				MarkdownDescription: "The region argument used to configure the provider.",
				Computed:            true,
			},
			"billing_project": schema.StringAttribute{
				Description:         "The billing_project argument used to configure the provider.",
				MarkdownDescription: "The billing_project argument used to configure the provider.",
				Computed:            true,
			},
			"zone": schema.StringAttribute{
				Description:         "The zone argument used to configure the provider.",
				MarkdownDescription: "The zone argument used to configure the provider.",
				Computed:            true,
			},
			"universe_domain": schema.StringAttribute{
				Description:         "The universe_domain argument used to configure the provider.",
				MarkdownDescription: "The universe_domain argument used to configure the provider.",
				Computed:            true,
			},
			"scopes": schema.ListAttribute{
				ElementType:         types.StringType,
				Description:         "The scopes argument used to configure the provider.",
				MarkdownDescription: "The scopes argument used to configure the provider.",
				Computed:            true,
			},
			"user_project_override": schema.BoolAttribute{
				Description:         "The user_project_override argument used to configure the provider.",
				MarkdownDescription: "The user_project_override argument used to configure the provider.",
				Computed:            true,
			},
			"request_reason": schema.StringAttribute{
				Description:         "The request_reason argument used to configure the provider.",
				MarkdownDescription: "The request_reason argument used to configure the provider.",
				Computed:            true,
			},
			"request_timeout": schema.StringAttribute{
				Description:         "The request_timeout argument used to configure the provider.",
				MarkdownDescription: "The request_timeout argument used to configure the provider.",
				Computed:            true,
			},
			"default_labels": schema.MapAttribute{
				ElementType:         types.StringType,
				Description:         "The default_labels argument used to configure the provider.",
				MarkdownDescription: "The default_labels argument used to configure the provider.",
				Computed:            true,
			},
			"add_terraform_attribution_label": schema.BoolAttribute{
				Description:         "The add_terraform_attribution_label argument used to configure the provider.",
				MarkdownDescription: "The add_terraform_attribution_label argument used to configure the provider.",
				Computed:            true,
			},
			"terraform_attribution_label_addition_strategy": schema.StringAttribute{
				Description:         "The terraform_attribution_label_addition_strategy argument used to configure the provider.",
				MarkdownDescription: "The terraform_attribution_label_addition_strategy argument used to configure the provider.",
				Computed:            true,
			},
			// End of user inputs

			// Note - this data source excludes the default and custom endpoints for individual services
		},
	}
}

func (d *GoogleProviderConfigPluginFrameworkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *transport_tpg.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	// Required for accessing project, region, zone and tokenSource
	d.providerConfig = p
}

func (d *GoogleProviderConfigPluginFrameworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GoogleProviderConfigPluginFrameworkModel
	var metaData *fwmodels.ProviderMetaModel

	// Read Provider meta into the meta model
	resp.Diagnostics.Append(req.ProviderMeta.Get(ctx, &metaData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Copy all values from the provider config into this data source
	data.Credentials = d.providerConfig.Credentials
	data.AccessToken = d.providerConfig.AccessToken
	data.ImpersonateServiceAccount = d.providerConfig.ImpersonateServiceAccount
	data.ImpersonateServiceAccountDelegates = d.providerConfig.ImpersonateServiceAccountDelegates
	data.Project = d.providerConfig.Project
	data.Region = d.providerConfig.Region
	data.BillingProject = d.providerConfig.BillingProject
	data.Zone = d.providerConfig.Zone
	data.UniverseDomain = d.providerConfig.UniverseDomain
	data.Scopes = d.providerConfig.Scopes
	data.UserProjectOverride = d.providerConfig.UserProjectOverride
	data.RequestReason = d.providerConfig.RequestReason
	data.RequestTimeout = time.Duration(d.providerConfig.RequestTimeout)
	data.DefaultLabels = d.providerConfig.DefaultLabels
	data.AddTerraformAttributionLabel = d.providerConfig.AddTerraformAttributionLabel
	data.TerraformAttributionLabelAdditionStrategy = d.providerConfig.TerraformAttributionLabelAdditionStrategy

	// Warn users against using this data source
	resp.Diagnostics.Append(diag.NewWarningDiagnostic(
		"Data source google_provider_config_plugin_framework should not be used",
		"Data source google_provider_config_plugin_framework is intended to be used only in acceptance tests for the provider. Instead, please use the google_client_config data source to access provider configuration details, or open a GitHub issue requesting new features in that datasource. Please go to: https://github.com/hashicorp/terraform-provider-google/issues/new/choose",
	))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
