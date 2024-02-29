package interpreter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ichiban/prolog"
	"github.com/ichiban/prolog/engine"

	"github.com/okp4/okp4d/x/logic/predicate"
)

// registry is a map from predicate names (in the form of "atom/arity") to predicates functions.
var registry = map[string]any{
	"call/1":                    engine.Call,
	"catch/3":                   engine.Catch,
	"throw/1":                   engine.Throw,
	"=/2":                       engine.Unify,
	"unify_with_occurs_check/2": engine.UnifyWithOccursCheck,
	"subsumes_term/2":           engine.SubsumesTerm,
	"var/1":                     engine.TypeVar,
	"atom/1":                    engine.TypeAtom,
	"integer/1":                 engine.TypeInteger,
	"float/1":                   engine.TypeFloat,
	"compound/1":                engine.TypeCompound,
	"acyclic_term/1":            engine.AcyclicTerm,
	"compare/3":                 engine.Compare,
	"sort/2":                    engine.Sort,
	"keysort/2":                 engine.KeySort,
	"functor/3":                 engine.Functor,
	"arg/3":                     engine.Arg,
	"=../2":                     engine.Univ,
	"copy_term/2":               engine.CopyTerm,
	"term_variables/2":          engine.TermVariables,
	"is/2":                      engine.Is,
	"=:=/2":                     engine.Equal,
	"=\\=/2":                    engine.NotEqual,
	"</2":                       engine.LessThan,
	"=</2":                      engine.LessThanOrEqual,
	">/2":                       engine.GreaterThan,
	">=/2":                      engine.GreaterThanOrEqual,
	"clause/2":                  engine.Clause,
	"current_predicate/1":       engine.CurrentPredicate,
	"asserta/1":                 engine.Asserta,
	"assertz/1":                 engine.Assertz,
	"retract/1":                 engine.Retract,
	"abolish/1":                 engine.Abolish,
	"findall/3":                 engine.FindAll,
	"bagof/3":                   engine.BagOf,
	"setof/3":                   engine.SetOf,
	"current_input/1":           engine.CurrentInput,
	"current_output/1":          engine.CurrentOutput,
	"set_input/1":               engine.SetInput,
	"set_output/1":              engine.SetOutput,
	"open/4":                    predicate.Open,
	"close/2":                   engine.Close,
	"flush_output/1":            engine.FlushOutput,
	"stream_property/2":         engine.StreamProperty,
	"set_stream_position/2":     engine.SetStreamPosition,
	"get_char/2":                engine.GetChar,
	"peek_char/2":               engine.PeekChar,
	"put_char/2":                engine.PutChar,
	"get_byte/2":                engine.GetByte,
	"peek_byte/2":               engine.PeekByte,
	"put_byte/2":                engine.PutByte,
	"read_term/3":               engine.ReadTerm,
	"write_term/3":              engine.WriteTerm,
	"op/3":                      engine.Op,
	"current_op/3":              engine.CurrentOp,
	"char_conversion/2":         engine.CharConversion,
	"current_char_conversion/2": engine.CurrentCharConversion,
	`\+/1`:                      engine.Negate,
	"repeat/0":                  engine.Repeat,
	"call/2":                    engine.Call1,
	"call/3":                    engine.Call2,
	"call/4":                    engine.Call3,
	"call/5":                    engine.Call4,
	"call/6":                    engine.Call5,
	"call/7":                    engine.Call6,
	"call/8":                    engine.Call7,
	"atom_length/2":             engine.AtomLength,
	"atom_concat/3":             engine.AtomConcat,
	"sub_atom/5":                engine.SubAtom,
	"atom_chars/2":              engine.AtomChars,
	"atom_codes/2":              engine.AtomCodes,
	"char_code/2":               engine.CharCode,
	"number_chars/2":            engine.NumberChars,
	"number_codes/2":            engine.NumberCodes,
	"set_prolog_flag/2":         engine.SetPrologFlag,
	"current_prolog_flag/2":     engine.CurrentPrologFlag,
	"halt/1":                    engine.Halt,
	"consult/1":                 engine.Consult,
	"phrase/3":                  engine.Phrase,
	"expand_term/2":             engine.ExpandTerm,
	"append/3":                  engine.Append,
	"length/2":                  engine.Length,
	"between/3":                 engine.Between,
	"succ/2":                    engine.Succ,
	"nth0/3":                    engine.Nth0,
	"nth1/3":                    engine.Nth1,
	"call_nth/2":                engine.CallNth,
	"chain_id/1":                predicate.ChainID,
	"block_height/1":            predicate.BlockHeight,
	"block_time/1":              predicate.BlockTime,
	"bank_balances/2":           predicate.BankBalances,
	"bank_spendable_balances/2": predicate.BankSpendableBalances,
	"bank_locked_balances/2":    predicate.BankLockedBalances,
	"did_components/2":          predicate.DIDComponents,
	"crypto_data_hash/3":        predicate.CryptoDataHash,
	"hex_bytes/2":               predicate.HexBytes,
	"bech32_address/2":          predicate.Bech32Address,
	"source_file/1":             predicate.SourceFile,
	"json_prolog/2":             predicate.JSONProlog,
	"uri_encoded/3":             predicate.URIEncoded,
	"read_string/3":             predicate.ReadString,
	"eddsa_verify/4":            predicate.EDDSAVerify,
	"ecdsa_verify/4":            predicate.ECDSAVerify,
	"string_bytes/3":            predicate.StringBytes,
}

// RegistryNames is the list of the predicate names in the Registry.
var RegistryNames = func() []string {
	names := make([]string, 0, len(registry))

	for name := range registry {
		names = append(names, name)
	}
	return names
}()

type Hook = func(functor string) func(env *engine.Env) error

// Register registers a well-known predicate in the interpreter with support for consumption measurement.
// name is the name of the predicate in the form of "atom/arity".
// cost is the cost of executing the predicate.
// meter is the gas meter object that is called when the predicate is called and which allows to count the cost of
// executing the predicate(ctx).
//
//nolint:lll
func Register(i *prolog.Interpreter, name string, hook Hook) error {
	if p, ok := registry[name]; ok {
		parts := strings.Split(name, "/")
		if len(parts) == 2 {
			atom := engine.NewAtom(parts[0])
			arity, err := strconv.Atoi(parts[1])
			if err != nil {
				return err
			}

			invariant := hook(name)

			switch arity {
			case 0:
				i.Register0(atom, Instrument0(invariant, p.(func(*engine.VM, engine.Cont, *engine.Env) *engine.Promise)))
			case 1:
				i.Register1(atom, Instrument1(invariant, p.(func(*engine.VM, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
			case 2:
				i.Register2(atom, Instrument2(invariant, p.(func(*engine.VM, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
			case 3:
				i.Register3(atom, Instrument3(invariant, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
			case 4:
				i.Register4(atom, Instrument4(invariant, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
			case 5:
				i.Register5(atom, Instrument5(invariant, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
			case 6:
				i.Register6(atom, Instrument6(invariant, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
			case 7:
				i.Register7(atom, Instrument7(invariant, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
			case 8:
				i.Register8(atom, Instrument8(invariant, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
			default:
				panic(fmt.Sprintf("unsupported arity: %s", name))
			}
		} else {
			panic(fmt.Sprintf("invalid name: %s", name))
		}

		return nil
	}

	return fmt.Errorf("unknown predicate %s", name)
}
