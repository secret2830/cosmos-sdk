package transfer

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	routing "github.com/cosmos/cosmos-sdk/x/ibc/026-routing"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/20-transfer/client/cli"
)

// Name returns the IBC transfer ICS name
func Name() string {
	return SubModuleName
}

// GetTxCmd returns the root tx command for the IBC transfer.
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

var (
	_ routing.ModuleCallbacks = Callbacks{}
)

// Callbacks implements ModuleCallbacks in routing module
type Callbacks struct {
	keeper Keeper
}

// NewCallbacks constructs a Callbacks
func NewCallbacks(k Keeper) Callbacks {
	return Callbacks{
		keeper: k,
	}
}

func (c Callbacks) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	counterparty channeltypes.Counterparty,
	version string,
) error {
	return c.keeper.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, counterparty, version)
}

func (c Callbacks) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	counterparty channeltypes.Counterparty,
	version string,
	counterpartyVersion string,
) error {
	return c.keeper.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, counterparty, version, counterpartyVersion)
}

func (c Callbacks) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	version string,
) error {
	return c.keeper.OnChanOpenAck(ctx, portID, channelID, version)
}

func (c Callbacks) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return c.keeper.OnChanOpenConfirm(ctx, portID, channelID)
}

func (c Callbacks) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return c.keeper.OnChanCloseInit(ctx, portID, channelID)
}

func (c Callbacks) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return c.keeper.OnChanCloseConfirm(ctx, portID, channelID)
}

func (c Callbacks) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) ([]byte, error) {
	return c.keeper.OnRecvPacket(ctx, packet)
}

func (c Callbacks) OnAcknowledgePacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
) error {
	return c.keeper.OnAcknowledgePacket(ctx, packet, acknowledgement)
}

func (c Callbacks) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) error {
	return c.keeper.OnTimeoutPacket(ctx, packet)
}

func (c Callbacks) OnTimeoutPacketClose(
	ctx sdk.Context,
	packet channeltypes.Packet,
) error {
	return c.keeper.OnTimeoutPacketClose(ctx, packet)
}
