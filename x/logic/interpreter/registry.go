package interpreter

import (
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ichiban/prolog"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/predicate"
)

// RegistryEntry is the type of registry entry.
type RegistryEntry struct {
	// predicate is the registered predicate(ctx).
	predicate any
	// cost is the cost of the predicate when it is called.
	cost uint64
}

// Registry is a map from predicate names (in the form of "atom/arity") to predicates and their costs.
var Registry = map[string]RegistryEntry{
	"call/1":                    {engine.Call, 1},
	"catch/3":                   {engine.Catch, 1},
	"throw/1":                   {engine.Throw, 1},
	"=/2":                       {engine.Unify, 1},
	"unify_with_occurs_check/2": {engine.UnifyWithOccursCheck, 1},
	"subsumes_term/2":           {engine.SubsumesTerm, 1},
	"var/1":                     {engine.TypeVar, 1},
	"atom/1":                    {engine.TypeAtom, 1},
	"integer/1":                 {engine.TypeInteger, 1},
	"float/1":                   {engine.TypeFloat, 1},
	"compound/1":                {engine.TypeCompound, 1},
	"acyclic_term/1":            {engine.AcyclicTerm, 1},
	"compare/3":                 {engine.Compare, 1},
	"sort/2":                    {engine.Sort, 1},
	"keysort/2":                 {engine.KeySort, 1},
	"functor/3":                 {engine.Functor, 1},
	"arg/3":                     {engine.Arg, 1},
	"=../2":                     {engine.Univ, 1},
	"copy_term/2":               {engine.CopyTerm, 1},
	"term_variables/2":          {engine.TermVariables, 1},
	"is/2":                      {engine.Is, 1},
	"=:=/2":                     {engine.Equal, 1},
	"=\\=/2":                    {engine.NotEqual, 1},
	"</2":                       {engine.LessThan, 1},
	"=</2":                      {engine.LessThanOrEqual, 1},
	">/2":                       {engine.GreaterThan, 1},
	">=/2":                      {engine.GreaterThanOrEqual, 1},
	"clause/2":                  {engine.Clause, 1},
	"current_predicate/1":       {engine.CurrentPredicate, 1},
	"asserta/1":                 {engine.Asserta, 1},
	"assertz/1":                 {engine.Assertz, 1},
	"retract/1":                 {engine.Retract, 1},
	"abolish/1":                 {engine.Abolish, 1},
	"findall/3":                 {engine.FindAll, 1},
	"bagof/3":                   {engine.BagOf, 1},
	"setof/3":                   {engine.SetOf, 1},
	"current_input/1":           {engine.CurrentInput, 1},
	"current_output/1":          {engine.CurrentOutput, 1},
	"set_input/1":               {engine.SetInput, 1},
	"set_output/1":              {engine.SetOutput, 1},
	"open/4":                    {engine.Open, 1},
	"close/2":                   {engine.Close, 1},
	"flush_output/1":            {engine.FlushOutput, 1},
	"stream_property/2":         {engine.StreamProperty, 1},
	"set_stream_position/2":     {engine.SetStreamPosition, 1},
	"get_char/2":                {engine.GetChar, 1},
	"peek_char/2":               {engine.PeekChar, 1},
	"put_char/2":                {engine.PutChar, 1},
	"get_byte/2":                {engine.GetByte, 1},
	"peek_byte/2":               {engine.PeekByte, 1},
	"put_byte/2":                {engine.PutByte, 1},
	"read_term/3":               {engine.ReadTerm, 1},
	"write_term/3":              {engine.WriteTerm, 1},
	"op/3":                      {engine.Op, 1},
	"current_op/3":              {engine.CurrentOp, 1},
	"char_conversion/2":         {engine.CharConversion, 1},
	"current_char_conversion/2": {engine.CurrentCharConversion, 1},
	`\+/1`:                      {engine.Negate, 1},
	"repeat/0":                  {engine.Repeat, 1},
	"call/2":                    {engine.Call1, 1},
	"call/3":                    {engine.Call2, 1},
	"call/4":                    {engine.Call3, 1},
	"call/5":                    {engine.Call4, 1},
	"call/6":                    {engine.Call5, 1},
	"call/7":                    {engine.Call6, 1},
	"call/8":                    {engine.Call7, 1},
	"atom_length/2":             {engine.AtomLength, 1},
	"atom_concat/3":             {engine.AtomConcat, 1},
	"sub_atom/5":                {engine.SubAtom, 1},
	"atom_chars/2":              {engine.AtomChars, 1},
	"atom_codes/2":              {engine.AtomCodes, 1},
	"char_code/2":               {engine.CharCode, 1},
	"number_chars/2":            {engine.NumberChars, 1},
	"number_codes/2":            {engine.NumberCodes, 1},
	"set_prolog_flag/2":         {engine.SetPrologFlag, 1},
	"current_prolog_flag/2":     {engine.CurrentPrologFlag, 1},
	"halt/1":                    {engine.Halt, 1},
	"consult/1":                 {engine.Consult, 1},
	"phrase/3":                  {engine.Phrase, 1},
	"expand_term/2":             {engine.ExpandTerm, 1},
	"append/3":                  {engine.Append, 1},
	"length/2":                  {engine.Length, 1},
	"between/3":                 {engine.Between, 1},
	"succ/2":                    {engine.Succ, 1},
	"nth0/3":                    {engine.Nth0, 1},
	"nth1/3":                    {engine.Nth1, 1},
	"call_nth/2":                {engine.CallNth, 1},
	"chain_id/1":                {predicate.ChainID, 1},
	"block_height/1":            {predicate.BlockHeight, 1},
	"block_time/1":              {predicate.BlockTime, 1},
	"bank_balances/2":           {predicate.BankBalances, 1},
	"bank_spendable_balances/2": {predicate.BankSpendableBalances, 1},
	"bank_locked_balances/2":    {predicate.BankLockedBalances, 1},
	"did_components/2":          {predicate.DIDComponents, 1},
	"crypto_hash/2":             {predicate.CryptoHash, 1},
}

// RegistryNames is the list of the predicate names in the Registry.
var RegistryNames = func() []string {
	names := make([]string, 0, len(Registry))

	for name := range Registry {
		names = append(names, name)
	}
	return names
}()

// Register registers a well-known predicate in the interpreter with support for consumption measurement.
// name is the name of the predicate in the form of "atom/arity".
// meter is the gas meter object that is called when the predicate is called and which allows to count the cost of
// executing the predicate(ctx).
//
//nolint:lll
func Register(i *prolog.Interpreter, name string, meter sdk.GasMeter) error {
	if entry, ok := Registry[name]; ok {
		parts := strings.Split(name, "/")
		if len(parts) == 2 {
			atom := engine.NewAtom(parts[0])
			arity, err := strconv.Atoi(parts[1])
			if err != nil {
				return err
			}

			hook := func() sdk.Gas {
				meter.ConsumeGas(entry.cost, fmt.Sprintf("predicate %s", name))

				return meter.GasRemaining()
			}
			p := entry.predicate

			switch arity {
			case 0:
				i.Register0(atom, Instrument0(hook, p.(func(*engine.VM, engine.Cont, *engine.Env) *engine.Promise)))
			case 1:
				i.Register1(atom, Instrument1(hook, p.(func(*engine.VM, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
			case 2:
				i.Register2(atom, Instrument2(hook, p.(func(*engine.VM, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
			case 3:
				i.Register3(atom, Instrument3(hook, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
			case 4:
				i.Register4(atom, Instrument4(hook, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
			case 5:
				i.Register5(atom, Instrument5(hook, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
			case 6:
				i.Register6(atom, Instrument6(hook, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
			case 7:
				i.Register7(atom, Instrument7(hook, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
			case 8:
				i.Register8(atom, Instrument8(hook, p.(func(*engine.VM, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Term, engine.Cont, *engine.Env) *engine.Promise)))
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
