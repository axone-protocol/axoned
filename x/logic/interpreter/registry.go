package interpreter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/axone-protocol/prolog/v3"
	"github.com/axone-protocol/prolog/v3/engine"
	orderedmap "github.com/wk8/go-ordered-map/v2"

	"github.com/axone-protocol/axoned/v14/x/logic/predicate"
)

// registry is a map from predicate names (in the form of "atom/arity") to predicates functions.
var registry = orderedmap.New[string, any](
	orderedmap.WithInitialData[string, any]([]orderedmap.Pair[string, any]{
		{Key: "call/1", Value: predicate.Call},
		{Key: "catch/3", Value: predicate.Catch},
		{Key: "throw/1", Value: predicate.Throw},
		{Key: "=/2", Value: predicate.Unify},
		{Key: "unify_with_occurs_check/2", Value: predicate.UnifyWithOccursCheck},
		{Key: "subsumes_term/2", Value: predicate.SubsumesTerm},
		{Key: "var/1", Value: predicate.TypeVar},
		{Key: "atom/1", Value: predicate.TypeAtom},
		{Key: "integer/1", Value: predicate.TypeInteger},
		{Key: "float/1", Value: predicate.TypeFloat},
		{Key: "compound/1", Value: predicate.TypeCompound},
		{Key: "acyclic_term/1", Value: predicate.AcyclicTerm},
		{Key: "compare/3", Value: predicate.Compare},
		{Key: "sort/2", Value: predicate.Sort},
		{Key: "keysort/2", Value: predicate.KeySort},
		{Key: "functor/3", Value: predicate.Functor},
		{Key: "arg/3", Value: predicate.Arg},
		{Key: "=../2", Value: predicate.Univ},
		{Key: "copy_term/2", Value: predicate.CopyTerm},
		{Key: "term_variables/2", Value: predicate.TermVariables},
		{Key: "is/2", Value: predicate.Is},
		{Key: "=:=/2", Value: predicate.Equal},
		{Key: "=\\=/2", Value: predicate.NotEqual},
		{Key: "</2", Value: predicate.LessThan},
		{Key: "=</2", Value: predicate.LessThanOrEqual},
		{Key: ">/2", Value: predicate.GreaterThan},
		{Key: ">=/2", Value: predicate.GreaterThanOrEqual},
		{Key: "clause/2", Value: predicate.Clause},
		{Key: "current_predicate/1", Value: predicate.CurrentPredicate},
		{Key: "asserta/1", Value: predicate.Asserta},
		{Key: "assertz/1", Value: predicate.Assertz},
		{Key: "retract/1", Value: predicate.Retract},
		{Key: "abolish/1", Value: predicate.Abolish},
		{Key: "findall/3", Value: predicate.FindAll},
		{Key: "bagof/3", Value: predicate.BagOf},
		{Key: "setof/3", Value: predicate.SetOf},
		{Key: "current_input/1", Value: predicate.CurrentInput},
		{Key: "current_output/1", Value: predicate.CurrentOutput},
		{Key: "set_input/1", Value: predicate.SetInput},
		{Key: "set_output/1", Value: predicate.SetOutput},
		{Key: "open/4", Value: predicate.Open},
		{Key: "close/2", Value: predicate.Close},
		{Key: "flush_output/1", Value: predicate.FlushOutput},
		{Key: "stream_property/2", Value: predicate.StreamProperty},
		{Key: "set_stream_position/2", Value: predicate.SetStreamPosition},
		{Key: "get_char/2", Value: predicate.GetChar},
		{Key: "peek_char/2", Value: predicate.PeekChar},
		{Key: "put_char/2", Value: predicate.PutChar},
		{Key: "get_byte/2", Value: predicate.GetByte},
		{Key: "peek_byte/2", Value: predicate.PeekByte},
		{Key: "put_byte/2", Value: predicate.PutByte},
		{Key: "read_term/3", Value: predicate.ReadTerm3},
		{Key: "write_term/3", Value: predicate.WriteTerm3},
		{Key: "op/3", Value: predicate.Op},
		{Key: "current_op/3", Value: predicate.CurrentOp},
		{Key: "char_conversion/2", Value: predicate.CharConversion},
		{Key: "current_char_conversion/2", Value: predicate.CurrentCharConversion},
		{Key: "\\+/1", Value: predicate.Negate},
		{Key: "repeat/0", Value: predicate.Repeat},
		{Key: "call/2", Value: predicate.Call1},
		{Key: "call/3", Value: predicate.Call2},
		{Key: "call/4", Value: predicate.Call3},
		{Key: "call/5", Value: predicate.Call4},
		{Key: "call/6", Value: predicate.Call5},
		{Key: "call/7", Value: predicate.Call6},
		{Key: "call/8", Value: predicate.Call7},
		{Key: "atom_length/2", Value: predicate.AtomLength},
		{Key: "atom_concat/3", Value: predicate.AtomConcat},
		{Key: "sub_atom/5", Value: predicate.SubAtom},
		{Key: "atom_chars/2", Value: predicate.AtomChars},
		{Key: "atom_codes/2", Value: predicate.AtomCodes},
		{Key: "char_code/2", Value: predicate.CharCode},
		{Key: "number_chars/2", Value: predicate.NumberChars},
		{Key: "number_codes/2", Value: predicate.NumberCodes},
		{Key: "current_prolog_flag/2", Value: predicate.CurrentPrologFlag},
		{Key: "consult/1", Value: predicate.Consult},
		{Key: "phrase/3", Value: predicate.Phrase},
		{Key: "expand_term/2", Value: predicate.ExpandTerm},
		{Key: "append/3", Value: predicate.Append},
		{Key: "length/2", Value: predicate.Length},
		{Key: "between/3", Value: predicate.Between},
		{Key: "succ/2", Value: predicate.Succ},
		{Key: "nth0/3", Value: predicate.Nth0},
		{Key: "nth1/3", Value: predicate.Nth1},
		{Key: "call_nth/2", Value: predicate.CallNth},
		{Key: "./3", Value: predicate.Op3},
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
		{Key: "base64_encoded/3", Value: predicate.Base64Encoded},
		{Key: "base64url/2", Value: predicate.Base64URL},
		{Key: "base64/2", Value: predicate.Base64},
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
