package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	handler "github.com/cosmos/cosmos-sdk/x/ibc/025-handler"
	"github.com/cosmos/cosmos-sdk/x/ibc/026-routing/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

// Keeper defines the IBC routing keeper
type Keeper struct {
	storeKey  sdk.StoreKey
	cdc       *codec.Codec
	codespace sdk.CodespaceType
	prefix    []byte // prefix bytes for accessing the store

	clientKeeper     types.ClientKeeper
	connectionKeeper types.ConnectionKeeper
	channelKeeper    types.ChannelKeeper
	portKeeper       types.PortKeeper
}

// NewKeeper creates a new IBC connection Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType,
	clientKeeper types.ClientKeeper, connectionKeeper types.ConnectionKeeper,
	channelKeeper types.ChannelKeeper, portKeeper types.PortKeeper) Keeper {
	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		codespace:        sdk.CodespaceType(fmt.Sprintf("%s/%s", codespace, types.DefaultCodespace)), // "ibc/routing",
		prefix:           []byte(types.SubModuleName + "/"),                                          // "routing/"
		clientKeeper:     clientKeeper,
		connectionKeeper: connectionKeeper,
		channelKeeper:    channelKeeper,
		portKeeper:       portKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s/%s", ibctypes.ModuleName, types.SubModuleName))
}

// GetCallback returns callbacks for the specified port
func (k Keeper) GetCallbacks(ctx sdk.Context, portID string) (types.ModuleCallbacks, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.prefix)
	bz := store.Get(types.KeyCallback(portID))
	if bz == nil {
		return nil, false
	}

	var callbacks types.ModuleCallbacks
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &callbacks)
	return callbacks, true
}

// SetCallback sets callbacks to the store
func (k Keeper) SetCallbacks(ctx sdk.Context, portID string, callbacks types.ModuleCallbacks) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.prefix)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(callbacks)
	store.Set(types.KeyCallback(portID), bz)
}

// GetAuthentication returns the authentication identifier for the specified port
func (k Keeper) GetAuthentication(ctx sdk.Context, portID string) (string, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.prefix)
	bz := store.Get(types.KeyAuthentication(portID))
	if bz == nil {
		return "", false
	}

	var authID string
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &authID)
	return authID, true
}

// SetAuthentication sets the authentication identifier for the specified port
func (k Keeper) SetAuthentication(ctx sdk.Context, portID string, authID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.prefix)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(authID)
	store.Set(types.KeyAuthentication(portID), bz)
}

// BindPort binds the calling module to a given port and register callbacks
func (k Keeper) BindPort(ctx sdk.Context, portID string, callbacks types.ModuleCallbacks) {
	callbacks, found := k.GetCallbacks(ctx, portID)
	if !found {
		panic("the port has been bound")
	}

	_ = handler.BindPort(k.portKeeper, portID)

	k.SetCallbacks(ctx, portID, callbacks)

	// TODO: generate capability key and store
}
