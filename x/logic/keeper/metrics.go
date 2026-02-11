package keeper

import (
	"fmt"
	"time"

	"github.com/axone-protocol/prolog/v3/engine"
	"github.com/hashicorp/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"

	"github.com/axone-protocol/axoned/v13/x/logic/interpreter"
	"github.com/axone-protocol/axoned/v13/x/logic/types"
)

var metricsKeys = []string{types.ModuleName, "vm", "predicate"}

const (
	labelPredicate = "predicate"
)

func telemetryPredicateCallCounterHookFn() engine.HookFunc {
	return func(opcode engine.Opcode, operand engine.Term, _ *engine.Env) error {
		if opcode != engine.OpCall {
			return nil
		}

		predicate, ok := stringifyOperand(operand)
		if !ok {
			return nil
		}

		if !interpreter.IsRegistered(predicate) {
			return nil
		}

		telemetry.IncrCounterWithLabels(
			metricsKeys,
			1,
			[]metrics.Label{
				telemetry.NewLabel(labelPredicate, predicate),
			},
		)

		return nil
	}
}

func telemetryPredicateDurationHookFn() engine.HookFunc {
	var predicate string
	var start time.Time
	return func(opcode engine.Opcode, operand engine.Term, _ *engine.Env) error {
		if opcode != engine.OpCall {
			if predicate != "" {
				telemetry.MeasureSince(start, append(metricsKeys, predicate)...)
				predicate = ""
				start = time.Time{}
			}
			return nil
		}

		var ok bool
		if predicate, ok = stringifyOperand(operand); !ok {
			return nil
		}

		start = telemetry.Now()

		return nil
	}
}

// stringifyOperand returns the string representation of the operand if it implements fmt.Stringer.
// It returns an empty string and false if the operand does not have a string representation.
func stringifyOperand(operand engine.Term) (string, bool) {
	if stringer, ok := operand.(fmt.Stringer); ok {
		return stringer.String(), true
	}
	return "", false
}
