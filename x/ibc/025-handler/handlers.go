package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/025-handler/types"
)

// BindPort binds the calling module to the given port
func BindPort(portKeeper types.PortKeeper, portID string) sdk.CapabilityKey {
	return portKeeper.BindPort(portID)
}
