package keeper_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"dario.cat/mergo"
	"github.com/cucumber/godog"
	"github.com/golang/mock/gomock"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/smartystreets/goconvey/convey/reporting"
	"gopkg.in/yaml.v3"

	. "github.com/smartystreets/goconvey/convey"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axone-protocol/axoned/v8/x/logic"
	"github.com/axone-protocol/axoned/v8/x/logic/fs/composite"
	"github.com/axone-protocol/axoned/v8/x/logic/fs/wasm"
	"github.com/axone-protocol/axoned/v8/x/logic/keeper"
	logictestutil "github.com/axone-protocol/axoned/v8/x/logic/testutil"
	"github.com/axone-protocol/axoned/v8/x/logic/types"
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
	ctx           sdktestutil.TestContext
	accountKeeper *logictestutil.MockAccountKeeper
	bankKeeper    *logictestutil.MockBankKeeper
	wasmKeeper    *logictestutil.MockWasmKeeper
	params        types.Params
	request       types.QueryServiceAskRequest
	got           *types.QueryServiceAskResponse
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

func givenABlockWithTheFollowingHeader(ctx context.Context, table *godog.Table) error {
	tc := testCaseFromContext(ctx)

	header := tc.ctx.Ctx.BlockHeader()
	for _, row := range table.Rows {
		switch row.Cells[0].Value {
		case "Height":
			height, err := atoi64(row.Cells[1].Value)
			if err != nil {
				return err
			}
			header.Height = height
		case "Time":
			sec, err := atoi64(row.Cells[1].Value)
			if err != nil {
				return err
			}
			header.Time = time.Unix(sec, 0)
		default:
			return fmt.Errorf("unknown field: %s", row.Cells[0].Value)
		}
	}
	tc.ctx.Ctx = tc.ctx.Ctx.WithBlockHeader(header)

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
	if strings.TrimSpace(configuration.MediaType) != "yaml" {
		return fmt.Errorf("unsupported media type: %s. Want %s", configuration.MediaType, "yaml")
	}
	if err := yaml.Unmarshal([]byte(configuration.Content), &smartContractConfiguration); err != nil {
		return err
	}

	contractAddr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}
	messageWant := map[string]interface{}{}
	if err := json.Unmarshal([]byte(smartContractConfiguration.Message), &messageWant); err != nil {
		return err
	}

	wasmKeeper := testCaseFromContext(ctx).wasmKeeper
	wasmKeeper.EXPECT().
		QuerySmart(gomock.Any(), contractAddr, gomock.Any()).
		DoAndReturn(func(_ context.Context, _ []byte, messageBytes []byte) ([]byte, error) {
			message := map[string]interface{}{}
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

	limit := sdkmath.NewUint(uint64(n))
	request.Limit = &limit

	testCaseFromContext(ctx).request = request

	return whenTheQueryIsRun(ctx)
}

func theAnswerWeGetIs(ctx context.Context, want *godog.DocString) error {
	got := testCaseFromContext(ctx).got
	wantAnswer := &types.QueryServiceAskResponse{}
	if err := yaml.Unmarshal([]byte(want.Content), &wantAnswer); err != nil {
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
	sb.WriteString(fmt.Sprintf("Actual:\n%s\n", failureView.Actual))
	sb.WriteString(fmt.Sprintf("Expected:\n%s\n", failureView.Expected))
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
		key,
		key,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		tc.accountKeeper,
		tc.bankKeeper,
		func(ctx context.Context) fs.FS {
			vfs := composite.NewFS()
			vfs.Mount(wasm.Scheme, wasm.NewFS(ctx, tc.wasmKeeper))

			return vfs
		},
	)

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
	maxResultCount := sdkmath.NewUint(10)
	limits.MaxResultCount = &maxResultCount
	params.Limits = limits
	return params
}

func atoi64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}
