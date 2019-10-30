package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/20-transfer/types"
)

// nolint: unused
func (k Keeper) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	counterparty channeltypes.Counterparty,
	version string,
) error {
	if order != channeltypes.UNORDERED {
		return types.ErrInvalidChannelOrder(k.codespace, order.String())
	}

	if counterparty.PortID != types.BoundPortID {
		return types.ErrInvalidPort(k.codespace, portID)
	}

	if strings.TrimSpace(version) != "" {
		return types.ErrInvalidVersion(k.codespace, fmt.Sprintf("invalid version: %s", version))
	}

	// NOTE: as the escrow address is generated from both the port and channel IDs
	// there's no need to store it on a map.
	return nil
}

// nolint: unused
func (k Keeper) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	counterparty channeltypes.Counterparty,
	version string,
	counterpartyVersion string,
) error {
	if order != channeltypes.UNORDERED {
		return types.ErrInvalidChannelOrder(k.codespace, order.String())
	}

	if counterparty.PortID != types.BoundPortID {
		return types.ErrInvalidPort(k.codespace, portID)
	}

	if strings.TrimSpace(version) != "" {
		return types.ErrInvalidVersion(k.codespace, fmt.Sprintf("invalid version: %s", version))
	}

	if strings.TrimSpace(counterpartyVersion) != "" {
		return types.ErrInvalidVersion(k.codespace, fmt.Sprintf("invalid counterparty version: %s", version))
	}

	// NOTE: as the escrow address is generated from both the port and channel IDs
	// there's no need to store it on a map.
	return nil
}

// nolint: unused
func (k Keeper) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	version string,
) error {
	if strings.TrimSpace(version) != "" {
		return types.ErrInvalidVersion(k.codespace, fmt.Sprintf("invalid version: %s", version))
	}

	return nil
}

// nolint: unused
func (k Keeper) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// no-op
	return nil
}

// nolint: unused
func (k Keeper) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// no-op
	return nil
}

// nolint: unused
func (k Keeper) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// no-op
	return nil
}

// onRecvPacket is called when an FTTransfer packet is received
// nolint: unused
func (k Keeper) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) ([]byte, error) {
	var data types.PacketData

	err := data.UnmarshalJSON(packet.Data())
	if err != nil {
		return nil, types.ErrInvalidPacketData(k.codespace)
	}

	err = k.ReceiveTransfer(
		ctx, packet.SourcePort(), packet.SourceChannel(),
		packet.DestPort(), packet.DestChannel(), data,
	)
	if err != nil {
		return nil, err
	}

	return []byte{0x0}, nil
}

// nolint: unused
func (k Keeper) OnAcknowledgePacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
) error {
	// no-op
	return nil
}

// nolint: unused
func (k Keeper) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) error {
	var data types.PacketData

	err := data.UnmarshalJSON(packet.Data())
	if err != nil {
		return types.ErrInvalidPacketData(k.codespace)
	}

	// check the denom prefix
	prefix := types.GetDenomPrefix(packet.SourcePort(), packet.SourcePort())
	coins := make(sdk.Coins, len(data.Amount))
	for i, coin := range data.Amount {
		coin := coin
		if !strings.HasPrefix(coin.Denom, prefix) {
			return sdk.ErrInvalidCoins(fmt.Sprintf("%s doesn't contain the prefix '%s'", coin.Denom, prefix))
		}
		coins[i] = sdk.NewCoin(coin.Denom[len(prefix):], coin.Amount)
	}

	if data.Source {
		escrowAddress := types.GetEscrowAddress(packet.DestPort(), packet.DestChannel())
		return k.bankKeeper.SendCoins(ctx, escrowAddress, data.Sender, coins)
	}

	// mint from supply
	err = k.supplyKeeper.MintCoins(ctx, types.GetModuleAccountName(), data.Amount)
	if err != nil {
		return err
	}

	return k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.GetModuleAccountName(), data.Sender, data.Amount)
}

// nolint: unused
func (k Keeper) OnTimeoutPacketClose(_ sdk.Context, _ channeltypes.Packet) error {
	panic("can't happen, only unordered channels allowed")
}
