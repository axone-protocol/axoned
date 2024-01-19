package logic

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/cosmos/cosmos-sdk/version"

	"github.com/okp4/okp4d/x/logic/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface for the logic module.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.GrpcQueryServiceDesc().ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Ask",
					Use:       "ask [query]",
					Short:     "Executes a logic query and returns the solution(s) found.",
					Long: `Executes the [query] and return the solution(s) found.

Optionally, a program can be transmitted, which will be compiled before the query is processed.

Since the query is without any side-effect, the query is not executed in the context of a transaction and no fee
is charged for this, but the execution is constrained by the current limits configured in the module (that you can
query).`,
					Example: fmt.Sprintf(`$ %s query %s ask "chain_id(X)." # returns the chain-id`,
						version.AppName,
						types.ModuleName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "query"}},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"program": {
							Shorthand: "p",
							Usage:     "The program to compile before the query.",
						},
					},
				},
			},
		},
	}
}
