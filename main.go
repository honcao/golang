package main

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest/azure"

	azarmhelper "github.com/honcao/golang/pkg/armhelper/azure"
	azsarmhelper "github.com/honcao/golang/pkg/armhelper/azurestack"
)

func main() {
	/*
		fn := decorator.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r.Method = "PUT"
			return r, nil
		})

		rep, _ := fn.Prepare(&http.Request{})
		fmt.Println(rep.Method)

		p := decorator.CreatePreparer(decorator.SetMethod("PUT"))
		r, _ := p.Prepare(&http.Request{})
		fmt.Println(r.Method)
	*/
	env, err := azure.EnvironmentFromName("AzurePublicCloud")
	if err != nil {
		panic(err)
	}

	armclient, err := azarmhelper.NewAzureClientWithClientSecret(env, "9ee2ec52-83c0-405e-a009-6636ead37acd", "72f988bf-86f1-41af-91ab-2d7cd011db47", "85115f84-ef7b-4ddb-b44d-b3a9d3b1990d", "y9liZF65vOyPgpjqJLUnOnjRRH7i4rCA+EMhPAM4dac=")

	rg, _ := armclient.EnsureResourceGroup(context.Background(), "honcaorg1", "eastus", nil)
	fmt.Println(*rg.Name)

	azsarmclient, err := azsarmhelper.NewAzureClientWithClientSecret(env, "9ee2ec52-83c0-405e-a009-6636ead37acd", "72f988bf-86f1-41af-91ab-2d7cd011db47", "85115f84-ef7b-4ddb-b44d-b3a9d3b1990d", "y9liZF65vOyPgpjqJLUnOnjRRH7i4rCA+EMhPAM4dac=")

	rg2, err := azsarmclient.EnsureResourceGroup(context.Background(), "honcaorg2", "westus", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*rg2.Name)

	/*
		os.Setenv("AZURE_TENANT_ID", "72f988bf-86f1-41af-91ab-2d7cd011db47")
		os.Setenv("AZURE_CLIENT_ID", "85115f84-ef7b-4ddb-b44d-b3a9d3b1990d")
		os.Setenv("AZURE_CLIENT_SECRET", "y9liZF65vOyPgpjqJLUnOnjRRH7i4rCA+EMhPAM4dac=")
		subID := "9ee2ec52-83c0-405e-a009-6636ead37acd"

		groupClient := resources.NewGroupsClient(subID)

		authorizer, err := auth.NewAuthorizerFromEnvironment()
		if err == nil {
			groupClient.Authorizer = authorizer
		}

		groupClient.CreateOrUpdate(
			context.Background(),
			"honcaogoclientrg1",
			resources.Group{
				Location: to.StringPtr("eastus"),
			},
		)
	*/

	/*
		// create a VirtualNetworks client
		vnetClient := network.NewVirtualNetworksClient(subID)

		// create an authorizer from env vars or Azure Managed Service Idenity
		authorizer, err := auth.NewAuthorizerFromEnvironment()
		if err == nil {
			vnetClient.Authorizer = authorizer
		}

		// call the VirtualNetworks CreateOrUpdate API
		vnetClient.CreateOrUpdate(context.Background(),
			"<resourceGroupName>",
			"<vnetName>",
			network.VirtualNetwork{
				Location: to.StringPtr("<azureRegion>"),
				VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
					AddressSpace: &network.AddressSpace{
						AddressPrefixes: &[]string{"10.0.0.0/8"},
					},
					Subnets: &[]network.Subnet{
						{
							Name: to.StringPtr("<subnet1Name>"),
							SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
								AddressPrefix: to.StringPtr("10.0.0.0/16"),
							},
						},
						{
							Name: to.StringPtr("<subnet2Name>"),
							SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
								AddressPrefix: to.StringPtr("10.1.0.0/16"),
							},
						},
					},
				},
			})
	*/
}
