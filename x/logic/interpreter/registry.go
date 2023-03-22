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

const defaultCost = 1

// Registry is a map from predicate names (in the form of "atom/arity") to predicates and their costs.
var Registry = map[string]RegistryEntry{
	"call/1":                    {engine.Call, defaultCost},
	"catch/3":                   {engine.Catch, defaultCost},
	"throw/1":                   {engine.Throw, defaultCost},
	"=/2":                       {engine.Unify, defaultCost},
	"unify_with_occurs_check/2": {engine.UnifyWithOccursCheck, defaultCost},
	"subsumes_term/2":           {engine.SubsumesTerm, defaultCost},
	"var/1":                     {engine.TypeVar, defaultCost},
	"atom/1":                    {engine.TypeAtom, defaultCost},
	"integer/1":                 {engine.TypeInteger, defaultCost},
	"float/1":                   {engine.TypeFloat, defaultCost},
	"compound/1":                {engine.TypeCompound, defaultCost},
	"acyclic_term/1":            {engine.AcyclicTerm, defaultCost},
	"compare/3":                 {engine.Compare, defaultCost},
	"sort/2":                    {engine.Sort, defaultCost},
	"keysort/2":                 {engine.KeySort, defaultCost},
	"functor/3":                 {engine.Functor, defaultCost},
	"arg/3":                     {engine.Arg, defaultCost},
	"=../2":                     {engine.Univ, defaultCost},
	"copy_term/2":               {engine.CopyTerm, defaultCost},
	"term_variables/2":          {engine.TermVariables, defaultCost},
	"is/2":                      {engine.Is, defaultCost},
	"=:=/2":                     {engine.Equal, defaultCost},
	"=\\=/2":                    {engine.NotEqual, defaultCost},
	"</2":                       {engine.LessThan, defaultCost},
	"=</2":                      {engine.LessThanOrEqual, defaultCost},
	">/2":                       {engine.GreaterThan, defaultCost},
	">=/2":                      {engine.GreaterThanOrEqual, defaultCost},
	"clause/2":                  {engine.Clause, defaultCost},
	"current_predicate/1":       {engine.CurrentPredicate, defaultCost},
	"asserta/1":                 {engine.Asserta, defaultCost},
	"assertz/1":                 {engine.Assertz, defaultCost},
	"retract/1":                 {engine.Retract, defaultCost},
	"abolish/1":                 {engine.Abolish, defaultCost},
	"findall/3":                 {engine.FindAll, defaultCost},
	"bagof/3":                   {engine.BagOf, defaultCost},
	"setof/3":                   {engine.SetOf, defaultCost},
	"current_input/1":           {engine.CurrentInput, defaultCost},
	"current_output/1":          {engine.CurrentOutput, defaultCost},
	"set_input/1":               {engine.SetInput, defaultCost},
	"set_output/1":              {engine.SetOutput, defaultCost},
	"open/4":                    {engine.Open, defaultCost},
	"close/2":                   {engine.Close, defaultCost},
	"flush_output/1":            {engine.FlushOutput, defaultCost},
	"stream_property/2":         {engine.StreamProperty, defaultCost},
	"set_stream_position/2":     {engine.SetStreamPosition, defaultCost},
	"get_char/2":                {engine.GetChar, defaultCost},
	"peek_char/2":               {engine.PeekChar, defaultCost},
	"put_char/2":                {engine.PutChar, defaultCost},
	"get_byte/2":                {engine.GetByte, defaultCost},
	"peek_byte/2":               {engine.PeekByte, defaultCost},
	"put_byte/2":                {engine.PutByte, defaultCost},
	"read_term/3":               {engine.ReadTerm, defaultCost},
	"write_term/3":              {engine.WriteTerm, defaultCost},
	"op/3":                      {engine.Op, defaultCost},
	"current_op/3":              {engine.CurrentOp, defaultCost},
	"char_conversion/2":         {engine.CharConversion, defaultCost},
	"current_char_conversion/2": {engine.CurrentCharConversion, defaultCost},
	`\+/1`:                      {engine.Negate, defaultCost},
	"repeat/0":                  {engine.Repeat, defaultCost},
	"call/2":                    {engine.Call1, defaultCost},
	"call/3":                    {engine.Call2, defaultCost},
	"call/4":                    {engine.Call3, defaultCost},
	"call/5":                    {engine.Call4, defaultCost},
	"call/6":                    {engine.Call5, defaultCost},
	"call/7":                    {engine.Call6, defaultCost},
	"call/8":                    {engine.Call7, defaultCost},
	"atom_length/2":             {engine.AtomLength, defaultCost},
	"atom_concat/3":             {engine.AtomConcat, defaultCost},
	"sub_atom/5":                {engine.SubAtom, defaultCost},
	"atom_chars/2":              {engine.AtomChars, defaultCost},
	"atom_codes/2":              {engine.AtomCodes, defaultCost},
	"char_code/2":               {engine.CharCode, defaultCost},
	"number_chars/2":            {engine.NumberChars, defaultCost},
	"number_codes/2":            {engine.NumberCodes, defaultCost},
	"set_prolog_flag/2":         {engine.SetPrologFlag, defaultCost},
	"current_prolog_flag/2":     {engine.CurrentPrologFlag, defaultCost},
	"halt/1":                    {engine.Halt, defaultCost},
	"consult/1":                 {engine.Consult, defaultCost},
	"phrase/3":                  {engine.Phrase, defaultCost},
	"expand_term/2":             {engine.ExpandTerm, defaultCost},
	"append/3":                  {engine.Append, defaultCost},
	"length/2":                  {engine.Length, defaultCost},
	"between/3":                 {engine.Between, defaultCost},
	"succ/2":                    {engine.Succ, defaultCost},
	"nth0/3":                    {engine.Nth0, defaultCost},
	"nth1/3":                    {engine.Nth1, defaultCost},
	"call_nth/2":                {engine.CallNth, defaultCost},
	"chain_id/1":                {predicate.ChainID, defaultCost},
	"block_height/1":            {predicate.BlockHeight, defaultCost},
	"block_time/1":              {predicate.BlockTime, defaultCost},
	"bank_balances/2":           {predicate.BankBalances, defaultCost},
	"bank_spendable_balances/2": {predicate.BankSpendableBalances, defaultCost},
	"bank_locked_balances/2":    {predicate.BankLockedBalances, defaultCost},
	"did_components/2":          {predicate.DIDComponents, defaultCost},
	"sha_hash/2":                {predicate.SHAHash, defaultCost},
	"hex_bytes/2":               {predicate.HexBytes, defaultCost},
	"bech32_address/2":          {predicate.Bech32Address, defaultCost},
	"source_file/1":             {predicate.SourceFile, defaultCost},
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
