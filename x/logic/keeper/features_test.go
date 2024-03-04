package keeper_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/golang/mock/gomock"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/smartystreets/goconvey/convey/reporting"
	"gopkg.in/yaml.v3"

	. "github.com/smartystreets/goconvey/convey"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/okp4/okp4d/v7/x/logic"
	"github.com/okp4/okp4d/v7/x/logic/keeper"
	logictestutil "github.com/okp4/okp4d/v7/x/logic/testutil"
	"github.com/okp4/okp4d/v7/x/logic/types"
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
	t       testing.TB
	ctx     sdktestutil.TestContext
	request types.QueryServiceAskRequest
	got     *types.Answer
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

func givenTheProgram(ctx context.Context, program *godog.DocString) error {
	testCaseFromContext(ctx).request.Program = program.Content

	return nil
}

func givenTheQuery(ctx context.Context, query *godog.DocString) error {
	testCaseFromContext(ctx).request.Query = query.Content

	return nil
}

func whenTheQueryIsRun(ctx context.Context) error {
	queryClient, err := newQueryClient(ctx)
	if err != nil {
		return err
	}

	tc := testCaseFromContext(ctx)
	got, err := queryClient.Ask(context.Background(), &tc.request)
	if err != nil {
		return err
	}

	tc.got = got.Answer

	return nil
}

func theAnswerWeGetIs(ctx context.Context, want *godog.DocString) error {
	got := testCaseFromContext(ctx).got
	wantAnswer := &types.Answer{}
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
	return func(ctx *godog.ScenarioContext) {
		ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
			testCtx := sdktestutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

			tc := testCase{
				t:   t,
				ctx: testCtx,
			}

			return testCaseToContext(ctx, tc), nil
		})

		ctx.Given(`a block with the following header:`, givenABlockWithTheFollowingHeader)
		ctx.Given(`the query:`, givenTheQuery)
		ctx.Given(`the program:`, givenTheProgram)
		ctx.When(`the query is run`, whenTheQueryIsRun)
		ctx.Then(`the answer we get is:`, theAnswerWeGetIs)
	}
}

func newQueryClient(ctx context.Context) (types.QueryServiceClient, error) {
	tc := testCaseFromContext(ctx)

	ctrl := gomock.NewController(tc.t)
	accountKeeper := logictestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := logictestutil.NewMockBankKeeper(ctrl)
	fsProvider := logictestutil.NewMockFS(ctrl)
	encCfg := moduletestutil.MakeTestEncodingConfig(logic.AppModuleBasic{})

	logicKeeper := keeper.NewKeeper(
		encCfg.Codec,
		key,
		key,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		accountKeeper,
		bankKeeper,
		func(_ context.Context) fs.FS {
			return fsProvider
		},
	)
	if err := logicKeeper.SetParams(tc.ctx.Ctx, types.DefaultParams()); err != nil {
		return nil, err
	}

	queryHelper := baseapp.NewQueryServerTestHelper(tc.ctx.Ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServiceServer(queryHelper, logicKeeper)
	queryClient := types.NewQueryServiceClient(queryHelper)
	return queryClient, nil
}

func atoi64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}
