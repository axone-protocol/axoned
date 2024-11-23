package interpreter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/axone-protocol/prolog"
	"github.com/axone-protocol/prolog/engine"
	orderedmap "github.com/wk8/go-ordered-map/v2"

	"github.com/axone-protocol/axoned/v10/x/logic/predicate"
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
		{Key: "read_term/3", Value: predicate.ReadTerm3},
		{Key: "write_term/3", Value: predicate.WriteTerm3},
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
		{Key: "current_prolog_flag/2", Value: engine.CurrentPrologFlag},
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
		{Key: "./3", Value: engine.Op3},
		{Key: "block_header/1", Value: predicate.BlockHeader},
		{Key: "chain_id/1", Value: predicate.ChainID}, //nolint:staticcheck // Deprecated but still exposed for compatibility.
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
		{Key: "term_to_atom/2", Value: predicate.TermToAtom},
		{Key: "atomic_list_concat/2", Value: predicate.AtomicListConcat2},
		{Key: "atomic_list_concat/3", Value: predicate.AtomicListConcat3},
		{Key: "json_read/2", Value: predicate.JSONRead},
		{Key: "json_write/2", Value: predicate.JSONWrite},
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

// IsRegistered returns true if the predicate with the given name is registered in the interpreter.
// Registered predicates are built-in predicates that are available in the interpreter.
func IsRegistered(name string) bool {
	_, ok := registry.Get(name)
	return ok
}

// Register registers a well-known predicate in the interpreter.
//
//nolint:lll
func Register(i *prolog.Interpreter, name string) error {
	if p, ok := registry.Get(name); ok {
		parts := strings.Split(name, "/")
		if len(parts) == 2 {
			atom := engine.NewAtom(parts[0])
			arity, err := strconv.Atoi(parts[1])
			if err != nil {
				return err
			}

			switch arity {
			case 0:
				i.Register0(atom, p.(func(*engine.VM, engine.Cont, *engine.Env) *engine.Promise))
			case 1:
				i.Register1(atom, p.(func(*engine.VM, engine.Term, engine.Cont, *engine.Env) *engine.Promise))
			case 2:
				i.Register2(atom, p.(func(*engine.VM, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise))
			case 3:
				i.Register3(atom, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise))
			case 4:
				i.Register4(atom, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise))
			case 5:
				i.Register5(atom, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise))
			case 6:
				i.Register6(atom, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise))
			case 7:
				i.Register7(atom, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise))
			case 8:
				i.Register8(atom, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise))
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
