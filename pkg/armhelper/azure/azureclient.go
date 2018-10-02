package azure

import (
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/honcao/golang/pkg/armhelper"
)

// AzureClient implements the `ACSEngineClient` interface.
// This client is backed by real Azure clients talking to an ARM endpoint.
type AzureClient struct {
	groupsClient   resources.GroupsClient
	environment    azure.Environment
	subscriptionID string
}

// NewAzureClientWithClientSecret returns an AzureClient via client_id and client_secret
func NewAzureClientWithClientSecret(env azure.Environment, subscriptionID, tenantID, clientID, clientSecret string) (armhelper.ACSEngineClient, error) {
	oauthConfig, tenantID, err := getOAuthConfig(env, subscriptionID, tenantID)
	if err != nil {
		return nil, err
	}

	armSpt, err := adal.NewServicePrincipalToken(*oauthConfig, clientID, clientSecret, env.ServiceManagementEndpoint)
	if err != nil {
		return nil, err
	}
	graphSpt, err := adal.NewServicePrincipalToken(*oauthConfig, clientID, clientSecret, env.GraphEndpoint)
	if err != nil {
		return nil, err
	}
	graphSpt.Refresh()

	return getClient(env, subscriptionID, tenantID, armSpt, graphSpt), nil
}

func getClient(env azure.Environment, subscriptionID, tenantID string, armSpt *adal.ServicePrincipalToken, graphSpt *adal.ServicePrincipalToken) *AzureClient {
	c := &AzureClient{
		environment:    env,
		subscriptionID: subscriptionID,
		groupsClient:   resources.NewGroupsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID),
	}

	authorizer := autorest.NewBearerAuthorizer(armSpt)
	c.groupsClient.Authorizer = authorizer

	return c
}

func getOAuthConfig(env azure.Environment, subscriptionID string, tenantID string) (*adal.OAuthConfig, string, error) {

	oauthConfig, err := adal.NewOAuthConfig(env.ActiveDirectoryEndpoint, tenantID)
	if err != nil {
		return nil, "", err
	}

	return oauthConfig, tenantID, nil
}
