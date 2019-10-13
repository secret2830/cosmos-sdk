package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	err "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

var (
	_ TxWithMemo = (*types.StdTx)(nil) // assert StdTx implements TxWithMemo
)

// ValidateBasicDecorator will call tx.ValidateBasic and return any non-nil error.
// If ValidateBasic passes, decorator calls next AnteHandler in chain.
type ValidateBasicDecorator struct{}

func NewValidateBasicDecorator() ValidateBasicDecorator {
	return ValidateBasicDecorator{}
}

func (vbd ValidateBasicDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	beforeGas := ctx.GasMeter().GasConsumed()
	if err := tx.ValidateBasic(); err != nil {
		return ctx, err
	}

	afterGas := ctx.GasMeter().GasConsumed()
	fmt.Printf("VALIDATEBASIC. SIMULATE: %t\n", simulate)
	fmt.Printf("GAS CONSUMED IN DECORATOR: %d\n", afterGas-beforeGas)
	fmt.Printf("TOTAL GAS CONSUMED: %d\n\n", afterGas)
	return next(ctx, tx, simulate)
}

// Tx must have GetMemo() method to use ValidateMemoDecorator
type TxWithMemo interface {
	sdk.Tx
	GetMemo() string
}

// ValidateMemoDecorator will validate memo given the parameters passed in
// If memo is too large decorator returns with error, otherwise call next AnteHandler
// CONTRACT: Tx must implement TxWithMemo interface
type ValidateMemoDecorator struct {
	ak keeper.AccountKeeper
}

func NewValidateMemoDecorator(ak keeper.AccountKeeper) ValidateMemoDecorator {
	return ValidateMemoDecorator{
		ak: ak,
	}
}

func (vmd ValidateMemoDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	beforeGas := ctx.GasMeter().GasConsumed()
	memoTx, ok := tx.(TxWithMemo)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	params := vmd.ak.GetParams(ctx)

	memoLength := len(memoTx.GetMemo())
	if uint64(memoLength) > params.MaxMemoCharacters {
		return ctx, err.Wrapf(err.ErrMemoTooLarge,
			"maximum number of characters is %d but received %d characters",
			params.MaxMemoCharacters, memoLength,
		)
	}

	afterGas := ctx.GasMeter().GasConsumed()
	fmt.Printf("VALIDATEMEMO. SIMULATE: %t\n", simulate)
	fmt.Printf("GAS CONSUMED IN DECORATOR: %d\n", afterGas-beforeGas)
	fmt.Printf("TOTAL GAS CONSUMED: %d\n\n", afterGas)
	return next(ctx, tx, simulate)
}

// ConsumeTxSizeGasDecorator will take in parameters and consume gas proportional to the size of tx
// before calling next AnteHandler
type ConsumeTxSizeGasDecorator struct {
	ak keeper.AccountKeeper
}

func NewConsumeGasForTxSizeDecorator(ak keeper.AccountKeeper) ConsumeTxSizeGasDecorator {
	return ConsumeTxSizeGasDecorator{
		ak: ak,
	}
}

func (cgts ConsumeTxSizeGasDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	beforeGas := ctx.GasMeter().GasConsumed()
	params := cgts.ak.GetParams(ctx)
	ctx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*sdk.Gas(len(ctx.TxBytes())), "txSize")

	afterGas := ctx.GasMeter().GasConsumed()
	fmt.Printf("CONSUMEGASFORTXSIZE. SIMULATE: %t\n", simulate)
	fmt.Printf("GAS CONSUMED IN DECORATOR: %d\n", afterGas-beforeGas)
	fmt.Printf("TOTAL GAS CONSUMED: %d\n\n", afterGas)
	return next(ctx, tx, simulate)
}
