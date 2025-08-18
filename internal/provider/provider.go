package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-firefly/internal/client"
)

// Ensure the implementation satisfies the provider.Provider interface
var _ provider.Provider = &FireflyProvider{}

// FireflyProvider is the provider implementation for Firefly
type FireflyProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// FireflyProviderModel describes the provider data model
type FireflyProviderModel struct {
	AccessKey types.String `tfsdk:"access_key"`
	SecretKey types.String `tfsdk:"secret_key"`
	APIURL    types.String `tfsdk:"api_url"`
}

// New creates a new provider instance
func New() func() provider.Provider {
	return func() provider.Provider {
		return &FireflyProvider{
			version: "dev",
		}
	}
}

// Metadata returns the provider type name
func (p *FireflyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "firefly"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data
func (p *FireflyProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Firefly",
		Attributes: map[string]schema.Attribute{
			"access_key": schema.StringAttribute{
				Description: "The access key for API operations. May also be provided via FIREFLY_ACCESS_KEY environment variable.",
				Required:    true,
				Sensitive:   true,
			},
			"secret_key": schema.StringAttribute{
				Description: "The secret key for API operations. May also be provided via FIREFLY_SECRET_KEY environment variable.",
				Required:    true,
				Sensitive:   true,
			},
			"api_url": schema.StringAttribute{
				Description: "The URL of the Firefly API. May also be provided via FIREFLY_API_URL environment variable.",
				Optional:    true,
			},
		},
	}
}

// Configure prepares the API client using the provider configuration
func (p *FireflyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Firefly client")
	
	// Retrieve provider data from configuration
	var config FireflyProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	// Default values for the API URL if unspecified
	apiURL := "https://api.firefly.ai"
	if !config.APIURL.IsNull() {
		apiURL = config.APIURL.ValueString()
	}
	
	// Check for required configuration
	if config.AccessKey.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_key"),
			"Missing Firefly Access Key",
			"The provider cannot create the Firefly API client without an access_key. "+
				"Please provide a valid access_key or set the FIREFLY_ACCESS_KEY environment variable.",
		)
	}
	
	if config.SecretKey.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret_key"),
			"Missing Firefly Secret Key",
			"The provider cannot create the Firefly API client without a secret_key. "+
				"Please provide a valid secret_key or set the FIREFLY_SECRET_KEY environment variable.",
		)
	}
	
	if resp.Diagnostics.HasError() {
		return
	}
	
	// Create a new client
	c, err := client.NewClient(client.Config{
		AccessKey:  config.AccessKey.ValueString(),
		SecretKey:  config.SecretKey.ValueString(),
		APIURL:     apiURL,
	})
	
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Firefly API client",
			fmt.Sprintf("Unable to create Firefly API client: %s", err),
		)
		return
	}
	
	// Make the client available to resources and data sources
	resp.DataSourceData = c
	resp.ResourceData = c
	
	tflog.Info(ctx, "Configured Firefly client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider
func (p *FireflyProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWorkspacesDataSource,
		NewWorkspaceRunsDataSource,
		NewGuardrailsDataSource,
		NewProjectsDataSource,
		NewProjectDataSource,
		NewVariableSetsDataSource,
		NewVariableSetDataSource,
	}
}

// Resources defines the resources implemented in the provider
func (p *FireflyProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWorkspaceLabelsResource,
		NewGuardrailResource,
		NewProjectResource,
		NewRunnersWorkspaceResource,
		NewVariableSetResource,
	}
}
