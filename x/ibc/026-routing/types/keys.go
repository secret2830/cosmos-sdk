package types

import (
	"fmt"
)

const (
	// SubModuleName defines the IBC routing name
	SubModuleName = "routing"

	// StoreKey is the store key string for IBC routing
	StoreKey = SubModuleName

	// RouterKey is the message route for IBC routing
	RouterKey = SubModuleName

	// QuerierRoute is the querier route for IBC routing
	QuerierRoute = SubModuleName
)

// The following paths are the keys to the store as defined in https://github.com/cosmos/ics/tree/master/spec/ics-003-connection-semantics#store-paths

// CallbackPath defines the path under which callback paths are stored
func CallbackPath(portID string) string {
	return fmt.Sprintf("callbacks/%s", portID)
}

// AuthenticationPath defines the path under which authentication paths are stored
func AuthenticationPath(portID string) string {
	return fmt.Sprintf("authentication/%s", portID)
}

// KeyCallback returns the store key for a set of callbacks of a specified port
func KeyCallback(portID string) []byte {
	return []byte(CallbackPath(portID))
}

// KeyAuthentication returns the store key for an authentication identifier of a specified port
func KeyAuthentication(portID string) []byte {
	return []byte(AuthenticationPath(portID))
}
