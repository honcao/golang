package converter

import (
	azsresources "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/honcao/golang/pkg/armhelper"
)

// ConvertGroupType converts resource group type to common resource group type
func ConvertGroupType(azsg *azsresources.Group) *armhelper.Group {
	return &armhelper.Group{
		Name:      azsg.Name,
		Location:  azsg.Location,
		ManagedBy: azsg.ManagedBy,
		Tags:      azsg.Tags,
	}
}
