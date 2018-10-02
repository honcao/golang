package armhelper

import (
	"context"
)

// ACSEngineClient is the interface used to talk to an Azure environment.
// This interface exposes just the subset of Azure APIs and clients needed for
// ACS-Engine.
type ACSEngineClient interface {

	// EnsureResourceGroup ensures the specified resource group exists in the specified location
	EnsureResourceGroup(ctx context.Context, resourceGroup, location string, managedBy *string) (*Group, error)
}
