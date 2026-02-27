package predicate_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"reflect"
	"strings"
	"testing"
	"time"

	"dario.cat/mergo"
	"github.com/cucumber/godog"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/smartystreets/goconvey/convey/reporting"
	"go.uber.org/mock/gomock"
	"sigs.k8s.io/yaml"

	. "github.com/smartystreets/goconvey/convey"

	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axone-protocol/axoned/v14/x/logic"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/bank"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/composite"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/dual"
	logicembeddedfs "github.com/axone-protocol/axoned/v14/x/logic/fs/embedded"
	logicsyscomet "github.com/axone-protocol/axoned/v14/x/logic/fs/sys/comet"
	logicsysheader "github.com/axone-protocol/axoned/v14/x/logic/fs/sys/header"
	logicvfs "github.com/axone-protocol/axoned/v14/x/logic/fs/vfs"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/wasm"
	"github.com/axone-protocol/axoned/v14/x/logic/keeper"
	logiclib "github.com/axone-protocol/axoned/v14/x/logic/lib"
	logictestutil "github.com/axone-protocol/axoned/v14/x/logic/testutil"
	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

var key = storetypes.NewKVStoreKey(types.StoreKey)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: initializeScenario(t),
		Options: &godog.Options{
			Format:      "pretty",
			Paths:       []string{"features"},
			TestingT:    t,
			Concurrency: 1,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("Test failed")
	}
}

type testCase struct {
	ctx              sdktestutil.TestContext
	accountKeeper    *logictestutil.MockAccountKeeper
	authQueryService *logictestutil.MockAuthQueryService
	bankKeeper       *logictestutil.MockBankKeeper
	wasmKeeper       *logictestutil.MockWasmKeeper
	params           types.Params
	request          types.QueryServiceAskRequest
	got              *types.QueryServiceAskResponse
}

type SmartContractConfiguration struct {
	Message  string `json:"message" yaml:"message"`
	Response string `json:"response" yaml:"response"`
}

type testCaseCtxKey struct{}

func testCaseToContext(ctx context.Context, tc testCase) context.Context {
	return context.WithValue(ctx, testCaseCtxKey{}, &tc)
}

func testCaseFromContext(ctx context.Context) *testCase {
	tc, _ := ctx.Value(testCaseCtxKey{}).(*testCase)

	return tc
}

func givenABlockWithTheFollowingHeader(ctx context.Context, headerConfig *godog.DocString) error {
	tc := testCaseFromContext(ctx)

	header := tc.ctx.Ctx.BlockHeader()

	if err := parseDocStringYaml(headerConfig, &header); err != nil {
		return err
	}

	tc.ctx.Ctx = tc.ctx.Ctx.WithBlockHeader(header)
	tc.ctx.Ctx = tc.ctx.Ctx.WithHeaderInfo(coreheader.Info{
		Height:  header.Height,
		Time:    header.Time,
		ChainID: header.ChainID,
		AppHash: header.AppHash,
	})

	return nil
}

func givenTheModuleConfiguration(ctx context.Context, configuration *godog.DocString) error {
	params := types.Params{}
	if err := json.Unmarshal([]byte(configuration.Content), &params); err != nil {
		return err
	}

	x := testCaseFromContext(ctx).params
	mergedParams := x
	if err := mergo.Merge(&mergedParams, params); err != nil {
		return err
	}

	testCaseFromContext(ctx).params = mergedParams

	return nil
}

func givenTheProgram(ctx context.Context, program *godog.DocString) error {
	testCaseFromContext(ctx).request.Program = program.Content

	return nil
}

func givenASmartContractWithAddress(ctx context.Context, address string, configuration *godog.DocString) error {
	smartContractConfiguration := &SmartContractConfiguration{}
	if err := parseDocStringYaml(configuration, smartContractConfiguration); err != nil {
		return err
	}

	contractAddr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}
	messageWant := map[string]any{}
	if err := json.Unmarshal([]byte(smartContractConfiguration.Message), &messageWant); err != nil {
		return err
	}

	wasmKeeper := testCaseFromContext(ctx).wasmKeeper
	wasmKeeper.EXPECT().
		QuerySmart(gomock.Any(), contractAddr, gomock.Any()).
		DoAndReturn(func(_ context.Context, _ []byte, messageBytes []byte) ([]byte, error) {
			message := map[string]any{}
			if err := json.Unmarshal(messageBytes, &message); err != nil {
				return nil, err
			}

			if reflect.DeepEqual(message, messageWant) {
				return []byte(smartContractConfiguration.Response), nil
			}

			return nil, fmt.Errorf("unexpected message: %v", message)
		}).
		AnyTimes()

	return nil
}

func givenTheQuery(ctx context.Context, query *godog.DocString) error {
	testCaseFromContext(ctx).request.Query = query.Content

	return nil
}

func givenTheAccountHasTheFollowingBalances(ctx context.Context, address string, balancesTable *godog.Table) error {
	addr, coins, err := addressAndCoinsFromBalancesTable(address, balancesTable)
	if err != nil {
		return err
	}

	bankKeeper := testCaseFromContext(ctx).bankKeeper
	bankKeeper.EXPECT().
		GetAllBalances(gomock.Any(), addr).
		Return(coins).
		AnyTimes()

	return nil
}

func givenTheAccountHasTheFollowingSpendableBalances(ctx context.Context, address string, balancesTable *godog.Table) error {
	addr, coins, err := addressAndCoinsFromBalancesTable(address, balancesTable)
	if err != nil {
		return err
	}

	bankKeeper := testCaseFromContext(ctx).bankKeeper
	bankKeeper.EXPECT().
		SpendableCoins(gomock.Any(), addr).
		Return(coins).
		AnyTimes()

	return nil
}

func whenTheQueryIsRun(ctx context.Context) error {
	tc := testCaseFromContext(ctx)

	tc.wasmKeeper.EXPECT().
		QuerySmart(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, contractAddr []byte, _ []byte) ([]byte, error) {
			return nil, fmt.Errorf("not existing contract: %s", contractAddr)
		}).
		AnyTimes()

	queryClient, err := newQueryClient(ctx)
	if err != nil {
		return err
	}

	got, err := queryClient.Ask(context.Background(), &tc.request)
	if err != nil {
		return err
	}

	tc.got = got

	return nil
}

func whenTheQueryIsRunLimitedToNSolutions(ctx context.Context, n int) error {
	request := testCaseFromContext(ctx).request
	request.Limit = uint64(n) //nolint:gosec // disable G115

	testCaseFromContext(ctx).request = request

	return whenTheQueryIsRun(ctx)
}

func theAnswerWeGetIs(ctx context.Context, want *godog.DocString) error {
	got := testCaseFromContext(ctx).got
	wantAnswer := &types.QueryServiceAskResponse{}
	if err := parseDocStringYaml(want, wantAnswer); err != nil {
		return err
	}

	return assert(got, ShouldResemble, wantAnswer)
}

func assert(actual any, assertion Assertion, expected ...any) error {
	msg := assertion(actual, expected...)
	if msg == "" {
		return nil
	}

	failureView := reporting.FailureView{}
	if err := json.Unmarshal([]byte(msg), &failureView); err != nil {
		return err
	}
	sb := strings.Builder{}
	sb.WriteString("assertion failed\n\u001B[0m")
	fmt.Fprintf(&sb, "Actual:\n%s\n", failureView.Actual)
	fmt.Fprintf(&sb, "Expected:\n%s\n", failureView.Expected)
	sb.WriteString("Diff:\n")

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(failureView.Actual, failureView.Expected, false)
	sb.WriteString(dmp.DiffPrettyText(diffs))

	return errors.New(sb.String())
}

func initializeScenario(t *testing.T) func(ctx *godog.ScenarioContext) {
	sdk.GetConfig().SetBech32PrefixForAccount("axone", "axonepub")

	return func(ctx *godog.ScenarioContext) {
		ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
			testCtx := sdktestutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
			ctrl := gomock.NewController(t)
			accountKeeper := logictestutil.NewMockAccountKeeper(ctrl)
			bankKeeper := logictestutil.NewMockBankKeeper(ctrl)
			wasmKeeper := logictestutil.NewMockWasmKeeper(ctrl)

			header := testCtx.Ctx.BlockHeader()
			header.ChainID = "axone-testchain-1"
			header.Height = 42
			header.Time = time.Date(2024, 4, 10, 10, 44, 27, 0, time.UTC)
			testCtx.Ctx = testCtx.Ctx.WithBlockHeader(header)
			testCtx.Ctx = testCtx.Ctx.WithHeaderInfo(coreheader.Info{
				Height:  header.Height,
				Time:    header.Time,
				ChainID: header.ChainID,
				AppHash: header.AppHash,
			})

			tc := testCase{
				ctx:           testCtx,
				accountKeeper: accountKeeper,
				bankKeeper:    bankKeeper,
				wasmKeeper:    wasmKeeper,
				params:        logicKeeperParams(),
			}

			return testCaseToContext(ctx, tc), nil
		})

		ctx.Given(`the module configuration:`, givenTheModuleConfiguration)
		ctx.Given(`a block with the following header:`, givenABlockWithTheFollowingHeader)
		ctx.Given(`the CosmWasm smart contract "([^"]+)" and the behavior:`, givenASmartContractWithAddress)
		ctx.Given(`the account "([^"]+)" has the following balances:`, givenTheAccountHasTheFollowingBalances)
		ctx.Given(`the account "([^"]+)" has the following spendable balances:`, givenTheAccountHasTheFollowingSpendableBalances)
		ctx.Given(`the query:`, givenTheQuery)
		ctx.Given(`the program:`, givenTheProgram)
		ctx.When(`^the query is run$`, whenTheQueryIsRun)
		ctx.When(`^the query is run \(limited to (\d+) solutions\)$`, whenTheQueryIsRunLimitedToNSolutions)
		ctx.Then(`the answer we get is:`, theAnswerWeGetIs)
	}
}

func newQueryClient(ctx context.Context) (types.QueryServiceClient, error) {
	tc := testCaseFromContext(ctx)

	encCfg := moduletestutil.MakeTestEncodingConfig(logic.AppModuleBasic{})
	logicKeeper := keeper.NewKeeper(
		encCfg.Codec,
		encCfg.InterfaceRegistry,
		key,
		key,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		tc.accountKeeper,
		tc.authQueryService,
		tc.bankKeeper,
		func(ctx context.Context) (fs.FS, error) {
			legacyFS := composite.NewFS()
			legacyFS.Mount(wasm.Scheme, wasm.NewFS(ctx, tc.wasmKeeper))

			pathFS := logicvfs.New()
			mounts := []struct {
				path string
				fs   fs.FS
			}{
				{"/v1/lib", logicembeddedfs.NewFS(logiclib.Files)},
				{"/v1/sys/header", logicsysheader.NewFS(ctx)},
				{"/v1/sys/comet", logicsyscomet.NewFS(ctx)},
				{"/v1/bank", bank.NewFS(ctx)},
			}
			for _, m := range mounts {
				if err := pathFS.Mount(m.path, m.fs); err != nil {
					return nil, fmt.Errorf("failed to mount %s: %w", m.path, err)
				}
			}

			return dual.NewFS(pathFS, legacyFS), nil
		})

	if err := logicKeeper.SetParams(tc.ctx.Ctx, tc.params); err != nil {
		return nil, err
	}

	queryHelper := baseapp.NewQueryServerTestHelper(tc.ctx.Ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServiceServer(queryHelper, logicKeeper)
	queryClient := types.NewQueryServiceClient(queryHelper)
	return queryClient, nil
}

func logicKeeperParams() types.Params {
	params := types.DefaultParams()
	limits := params.Limits
	limits.MaxResultCount = 10
	params.Limits = limits
	return params
}

func parseDocStringYaml(docString *godog.DocString, v any) error {
	const yamlMediaType = "yaml"
	if strings.TrimSpace(docString.MediaType) != yamlMediaType {
		return fmt.Errorf("unsupported media type: %s. Want %s", docString.MediaType, yamlMediaType)
	}

	if err := yaml.UnmarshalStrict([]byte(docString.Content), v); err != nil {
		return err
	}

	return nil
}

func addressAndCoinsFromBalancesTable(address string, balancesTable *godog.Table) (sdk.AccAddress, sdk.Coins, error) {
	if len(balancesTable.Rows) < 1 {
		return nil, nil, fmt.Errorf("balances table must have at least a header row")
	}

	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid address %s: %w", address, err)
	}

	coins := sdk.NewCoins()
	for i := 1; i < len(balancesTable.Rows); i++ {
		row := balancesTable.Rows[i]
		if len(row.Cells) != 2 {
			return nil, nil, fmt.Errorf("each balance row must have exactly 2 cells (denom, amount)")
		}

		denom := row.Cells[0].Value
		amountStr := row.Cells[1].Value

		if denom == "" || amountStr == "" {
			return nil, nil, fmt.Errorf("row %d: denom and amount must both be set", i)
		}

		amount, ok := math.NewIntFromString(amountStr)
		if !ok {
			return nil, nil, fmt.Errorf("row %d: invalid amount %s", i, amountStr)
		}

		if err := sdk.ValidateDenom(denom); err != nil {
			return nil, nil, fmt.Errorf("row %d: invalid denom %s: %w", i, denom, err)
		}

		if amount.IsNegative() {
			return nil, nil, fmt.Errorf("row %d: amount must be non-negative, got %s", i, amountStr)
		}

		coins = coins.Add(sdk.NewCoin(denom, amount))
	}

	return addr, coins.Sort(), nil
}
