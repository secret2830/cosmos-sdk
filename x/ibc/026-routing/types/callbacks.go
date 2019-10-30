package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
)

// ModuleCallbacks defines the callback interfaces which must be
// implemented in modules which need to interact with IBC
type ModuleCallbacks interface {
	OnChanOpenInit(
		ctx sdk.Context,
		order channeltypes.Order,
		connectionHops []string,
		portID,
		channelID string,
		counterparty channeltypes.Counterparty,
		version string,
	) error

	OnChanOpenTry(
		ctx sdk.Context,
		order channeltypes.Order,
		connectionHops []string,
		portID,
		channelID string,
		counterparty channeltypes.Counterparty,
		version string,
		counterpartyVersion string,
	) error

	OnChanOpenAck(
		ctx sdk.Context,
		portID,
		channelID string,
		version string,
	) error

	OnChanOpenConfirm(
		ctx sdk.Context,
		portID,
		channelID string,
	) error

	OnChanCloseInit(
		ctx sdk.Context,
		portID,
		channelID string,
	) error

	OnChanCloseConfirm(
		ctx sdk.Context,
		portID,
		channelID string,
	) error

	OnRecvPacket(
		ctx sdk.Context,
		packet channeltypes.Packet,
	) ([]byte, error)

	OnAcknowledgePacket(
		ctx sdk.Context,
		packet channeltypes.Packet,
		acknowledgement []byte,
	) error

	OnTimeoutPacket(
		ctx sdk.Context,
		packet channeltypes.Packet,
	) error

	OnTimeoutPacketClose(
		ctx sdk.Context,
		packet channeltypes.Packet,
	) error
}
