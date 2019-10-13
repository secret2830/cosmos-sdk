package ante_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())

	return app, ctx
}

func testSimGasEstimate(t *testing.T, antefn sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx) {
	cc, _ := ctx.CacheContext()
	_, err := antefn(cc, tx, true)

	require.Nil(t, err)
	simGas := cc.GasMeter().GasConsumed()

	_, err = antefn(ctx, tx, false)
	require.Nil(t, err)
	require.True(t, simGas >= ctx.GasMeter().GasConsumed())
}
