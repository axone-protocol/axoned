package interpreter

import (
	"fmt"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"strconv"
	"strings"

	"github.com/ichiban/prolog"
	engine "github.com/ichiban/prolog/engine"

	"github.com/axone-protocol/axoned/v8/x/logic/predicate"
)

// registry is a map from predicate names (in the form of "atom/arity") to predicates functions.
var registry = orderedmap.New[string, any](
	orderedmap.WithInitialData[string, any]([]orderedmap.Pair[string, any]{
		{Key: "call/1", Value: engine.Call},
		{Key: "catch/3", Value: engine.Catch},
		{Key: "throw/1", Value: engine.Throw},
		{Key: "=/2", Value: engine.Unify},
		{Key: "unify_with_occurs_check/2", Value: engine.UnifyWithOccursCheck},
		{Key: "subsumes_term/2", Value: engine.SubsumesTerm},
		{Key: "var/1", Value: engine.TypeVar},
		{Key: "atom/1", Value: engine.TypeAtom},
		{Key: "integer/1", Value: engine.TypeInteger},
		{Key: "float/1", Value: engine.TypeFloat},
		{Key: "compound/1", Value: engine.TypeCompound},
		{Key: "acyclic_term/1", Value: engine.AcyclicTerm},
		{Key: "compare/3", Value: engine.Compare},
		{Key: "sort/2", Value: engine.Sort},
		{Key: "keysort/2", Value: engine.KeySort},
		{Key: "functor/3", Value: engine.Functor},
		{Key: "arg/3", Value: engine.Arg},
		{Key: "=../2", Value: engine.Univ},
		{Key: "copy_term/2", Value: engine.CopyTerm},
		{Key: "term_variables/2", Value: engine.TermVariables},
		{Key: "is/2", Value: engine.Is},
		{Key: "=:=/2", Value: engine.Equal},
		{Key: "=\\=/2", Value: engine.NotEqual},
		{Key: "</2", Value: engine.LessThan},
		{Key: "=</2", Value: engine.LessThanOrEqual},
		{Key: ">/2", Value: engine.GreaterThan},
		{Key: ">=/2", Value: engine.GreaterThanOrEqual},
		{Key: "clause/2", Value: engine.Clause},
		{Key: "current_predicate/1", Value: engine.CurrentPredicate},
		{Key: "asserta/1", Value: engine.Asserta},
		{Key: "assertz/1", Value: engine.Assertz},
		{Key: "retract/1", Value: engine.Retract},
		{Key: "abolish/1", Value: engine.Abolish},
		{Key: "findall/3", Value: engine.FindAll},
		{Key: "bagof/3", Value: engine.BagOf},
		{Key: "setof/3", Value: engine.SetOf},
		{Key: "current_input/1", Value: engine.CurrentInput},
		{Key: "current_output/1", Value: predicate.CurrentOutput},
		{Key: "set_input/1", Value: engine.SetInput},
		{Key: "set_output/1", Value: engine.SetOutput},
		{Key: "open/4", Value: predicate.Open},
		{Key: "open/3", Value: predicate.Open3},
		{Key: "close/2", Value: engine.Close},
		{Key: "flush_output/1", Value: engine.FlushOutput},
		{Key: "stream_property/2", Value: engine.StreamProperty},
		{Key: "set_stream_position/2", Value: engine.SetStreamPosition},
		{Key: "get_char/2", Value: engine.GetChar},
		{Key: "peek_char/2", Value: engine.PeekChar},
		{Key: "put_char/2", Value: engine.PutChar},
		{Key: "get_byte/2", Value: engine.GetByte},
		{Key: "peek_byte/2", Value: engine.PeekByte},
		{Key: "put_byte/2", Value: engine.PutByte},
		{Key: "read_term/3", Value: engine.ReadTerm},
		{Key: "write_term/3", Value: engine.WriteTerm},
		{Key: "op/3", Value: engine.Op},
		{Key: "current_op/3", Value: engine.CurrentOp},
		{Key: "char_conversion/2", Value: engine.CharConversion},
		{Key: "current_char_conversion/2", Value: engine.CurrentCharConversion},
		{Key: "\\+/1", Value: engine.Negate},
		{Key: "repeat/0", Value: engine.Repeat},
		{Key: "call/2", Value: engine.Call1},
		{Key: "call/3", Value: engine.Call2},
		{Key: "call/4", Value: engine.Call3},
		{Key: "call/5", Value: engine.Call4},
		{Key: "call/6", Value: engine.Call5},
		{Key: "call/7", Value: engine.Call6},
		{Key: "call/8", Value: engine.Call7},
		{Key: "atom_length/2", Value: engine.AtomLength},
		{Key: "atom_concat/3", Value: engine.AtomConcat},
		{Key: "sub_atom/5", Value: engine.SubAtom},
		{Key: "atom_chars/2", Value: engine.AtomChars},
		{Key: "atom_codes/2", Value: engine.AtomCodes},
		{Key: "char_code/2", Value: engine.CharCode},
		{Key: "number_chars/2", Value: engine.NumberChars},
		{Key: "number_codes/2", Value: engine.NumberCodes},
		{Key: "set_prolog_flag/2", Value: engine.SetPrologFlag},
		{Key: "current_prolog_flag/2", Value: engine.CurrentPrologFlag},
		{Key: "halt/1", Value: engine.Halt},
		{Key: "consult/1", Value: predicate.Consult},
		{Key: "phrase/3", Value: engine.Phrase},
		{Key: "expand_term/2", Value: engine.ExpandTerm},
		{Key: "append/3", Value: engine.Append},
		{Key: "length/2", Value: engine.Length},
		{Key: "between/3", Value: engine.Between},
		{Key: "succ/2", Value: engine.Succ},
		{Key: "nth0/3", Value: engine.Nth0},
		{Key: "nth1/3", Value: engine.Nth1},
		{Key: "call_nth/2", Value: engine.CallNth},
		{Key: "chain_id/1", Value: predicate.ChainID},
		{Key: "block_height/1", Value: predicate.BlockHeight},
		{Key: "block_time/1", Value: predicate.BlockTime},
		{Key: "bank_balances/2", Value: predicate.BankBalances},
		{Key: "bank_spendable_balances/2", Value: predicate.BankSpendableBalances},
		{Key: "bank_locked_balances/2", Value: predicate.BankLockedBalances},
		{Key: "did_components/2", Value: predicate.DIDComponents},
		{Key: "crypto_data_hash/3", Value: predicate.CryptoDataHash},
		{Key: "hex_bytes/2", Value: predicate.HexBytes},
		{Key: "bech32_address/2", Value: predicate.Bech32Address},
		{Key: "source_file/1", Value: predicate.SourceFile},
		{Key: "json_prolog/2", Value: predicate.JSONProlog},
		{Key: "uri_encoded/3", Value: predicate.URIEncoded},
		{Key: "read_string/3", Value: predicate.ReadString},
		{Key: "eddsa_verify/4", Value: predicate.EDDSAVerify},
		{Key: "ecdsa_verify/4", Value: predicate.ECDSAVerify},
		{Key: "string_bytes/3", Value: predicate.StringBytes},
	}...),
)

// RegistryNames is the list of the predicate names in the Registry.
var RegistryNames = func() []string {
	names := make([]string, 0, registry.Len())

	for name := registry.Oldest(); name != nil; name = name.Next() {
		names = append(names, name.Key)
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
	if p, ok := registry.Get(name); ok {
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
