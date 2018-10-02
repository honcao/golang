package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/honcao/golang/pkg/armhelper"
)

// EnsureResourceGroup ensures the named resouce group exists in the given location.
func (az *AzureClient) EnsureResourceGroup(ctx context.Context, name, location string, managedBy *string) (resourceGroup *armhelper.Group, err error) {
	var tags map[string]*string
	group, err := az.groupsClient.Get(ctx, name)
	if err == nil {
		tags = group.Tags
	}

	response, err := az.groupsClient.CreateOrUpdate(ctx, name, resources.Group{
		Name:      &name,
		Location:  &location,
		ManagedBy: managedBy,
		Tags:      tags,
	})

	response1 := armhelper.Group(response)
	if err != nil {
		return &response1, err
	}

	return &response1, nil
}
