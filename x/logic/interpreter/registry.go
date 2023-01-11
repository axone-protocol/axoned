package interpreter

import (
	goctx "context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ichiban/prolog"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/context"
	"github.com/okp4/okp4d/x/logic/predicate"
)

// wrapContext wraps a context.Context around a predicate(ctx).
func wrapCtx(p any) func(goctx.Context) any {
	return func(ctx goctx.Context) any {
		return p
	}
}

func relax[T any](f func(goctx.Context) T) func(goctx.Context) any {
	return func(ctx goctx.Context) any {
		return f(ctx)
	}
}

// RegistryEntry is the type of registry entry.
type RegistryEntry struct {
	// predicate is the registered predicate(ctx).
	predicate func(goctx.Context) any
	// cost is the cost of the predicate when it is called.
	cost uint64
}

// Registry is a map from predicate names (in the form of "atom/arity") to predicates and their costs.
var Registry = map[string]RegistryEntry{
	"call/1":                    {wrapCtx(engine.Call), 1},
	"catch/3":                   {wrapCtx(engine.Catch), 1},
	"throw/1":                   {wrapCtx(engine.Throw), 1},
	"=/2":                       {wrapCtx(engine.Unify), 1},
	"unify_with_occurs_check/2": {wrapCtx(engine.UnifyWithOccursCheck), 1},
	"subsumes_term/2":           {wrapCtx(engine.SubsumesTerm), 1},
	"var/1":                     {wrapCtx(engine.TypeVar), 1},
	"atom/1":                    {wrapCtx(engine.TypeAtom), 1},
	"integer/1":                 {wrapCtx(engine.TypeInteger), 1},
	"float/1":                   {wrapCtx(engine.TypeFloat), 1},
	"compound/1":                {wrapCtx(engine.TypeCompound), 1},
	"acyclic_term/1":            {wrapCtx(engine.AcyclicTerm), 1},
	"compare/3":                 {wrapCtx(engine.Compare), 1},
	"sort/2":                    {wrapCtx(engine.Sort), 1},
	"keysort/2":                 {wrapCtx(engine.KeySort), 1},
	"functor/3":                 {wrapCtx(engine.Functor), 1},
	"arg/3":                     {wrapCtx(engine.Arg), 1},
	"=../2":                     {wrapCtx(engine.Univ), 1},
	"copy_term/2":               {wrapCtx(engine.CopyTerm), 1},
	"term_variables/2":          {wrapCtx(engine.TermVariables), 1},
	"is/2":                      {wrapCtx(engine.Is), 1},
	"=:=/2":                     {wrapCtx(engine.Equal), 1},
	"=\\=/2":                    {wrapCtx(engine.NotEqual), 1},
	"</2":                       {wrapCtx(engine.LessThan), 1},
	"=</2":                      {wrapCtx(engine.LessThanOrEqual), 1},
	">/2":                       {wrapCtx(engine.GreaterThan), 1},
	">=/2":                      {wrapCtx(engine.GreaterThanOrEqual), 1},
	"clause/2":                  {wrapCtx(engine.Clause), 1},
	"current_predicate/1":       {wrapCtx(engine.CurrentPredicate), 1},
	"asserta/1":                 {wrapCtx(engine.Asserta), 1},
	"assertz/1":                 {wrapCtx(engine.Assertz), 1},
	"retract/1":                 {wrapCtx(engine.Retract), 1},
	"abolish/1":                 {wrapCtx(engine.Abolish), 1},
	"findall/3":                 {wrapCtx(engine.FindAll), 1},
	"bagof/3":                   {wrapCtx(engine.BagOf), 1},
	"setof/3":                   {wrapCtx(engine.SetOf), 1},
	"current_input/1":           {wrapCtx(engine.CurrentInput), 1},
	"current_output/1":          {wrapCtx(engine.CurrentOutput), 1},
	"set_input/1":               {wrapCtx(engine.SetInput), 1},
	"set_output/1":              {wrapCtx(engine.SetOutput), 1},
	"open/4":                    {wrapCtx(engine.Open), 1},
	"close/2":                   {wrapCtx(engine.Close), 1},
	"flush_output/1":            {wrapCtx(engine.FlushOutput), 1},
	"stream_property/2":         {wrapCtx(engine.StreamProperty), 1},
	"set_stream_position/2":     {wrapCtx(engine.SetStreamPosition), 1},
	"get_char/2":                {wrapCtx(engine.GetChar), 1},
	"peek_char/2":               {wrapCtx(engine.PeekChar), 1},
	"put_char/2":                {wrapCtx(engine.PutChar), 1},
	"get_byte/2":                {wrapCtx(engine.GetByte), 1},
	"peek_byte/2":               {wrapCtx(engine.PeekByte), 1},
	"put_byte/2":                {wrapCtx(engine.PutByte), 1},
	"read_term/3":               {wrapCtx(engine.ReadTerm), 1},
	"write_term/3":              {wrapCtx(engine.WriteTerm), 1},
	"op/3":                      {wrapCtx(engine.Op), 1},
	"current_op/3":              {wrapCtx(engine.CurrentOp), 1},
	"char_conversion/2":         {wrapCtx(engine.CharConversion), 1},
	"current_char_conversion/2": {wrapCtx(engine.CurrentCharConversion), 1},
	`\+/1`:                      {wrapCtx(engine.Negate), 1},
	"repeat/0":                  {wrapCtx(engine.Repeat), 1},
	"call/2":                    {wrapCtx(engine.Call1), 1},
	"call/3":                    {wrapCtx(engine.Call2), 1},
	"call/4":                    {wrapCtx(engine.Call3), 1},
	"call/5":                    {wrapCtx(engine.Call4), 1},
	"call/6":                    {wrapCtx(engine.Call5), 1},
	"call/7":                    {wrapCtx(engine.Call6), 1},
	"call/8":                    {wrapCtx(engine.Call7), 1},
	"atom_length/2":             {wrapCtx(engine.AtomLength), 1},
	"atom_concat/3":             {wrapCtx(engine.AtomConcat), 1},
	"sub_atom/5":                {wrapCtx(engine.SubAtom), 1},
	"atom_chars/2":              {wrapCtx(engine.AtomChars), 1},
	"atom_codes/2":              {wrapCtx(engine.AtomCodes), 1},
	"char_code/2":               {wrapCtx(engine.CharCode), 1},
	"number_chars/2":            {wrapCtx(engine.NumberChars), 1},
	"number_codes/2":            {wrapCtx(engine.NumberCodes), 1},
	"set_prolog_flag/2":         {wrapCtx(engine.SetPrologFlag), 1},
	"current_prolog_flag/2":     {wrapCtx(engine.CurrentPrologFlag), 1},
	"halt/1":                    {wrapCtx(engine.Halt), 1},
	"consult/1":                 {wrapCtx(engine.Consult), 1},
	"phrase/3":                  {wrapCtx(engine.Phrase), 1},
	"expand_term/2":             {wrapCtx(engine.ExpandTerm), 1},
	"append/3":                  {wrapCtx(engine.Append), 1},
	"length/2":                  {wrapCtx(engine.Length), 1},
	"between/3":                 {wrapCtx(engine.Between), 1},
	"succ/2":                    {wrapCtx(engine.Succ), 1},
	"nth0/3":                    {wrapCtx(engine.Nth0), 1},
	"nth1/3":                    {wrapCtx(engine.Nth1), 1},
	"call_nth/2":                {wrapCtx(engine.CallNth), 1},
	"chain_id/1":                {relax(predicate.ChainID), 1},
	"block_height/1":            {relax(predicate.BlockHeight), 1},
	"block_time/1":              {relax(predicate.BlockTime), 1},
}

// RegistryNames is the list of the predicate names in the Registry.
var RegistryNames = func() []string {
	names := make([]string, 0, len(Registry))

	for name := range Registry {
		names = append(names, name)
	}
	return names
}()

// Register registers a well-known predicate in the interpreter with support for "instrumentation".
// name is the name of the predicate in the form of "atom/arity".
// inc is the increment function that is called when the predicate is called and which allows to count the cost of
// executing the predicate(ctx).
//
//nolint:lll
func Register(ctx goctx.Context, i *prolog.Interpreter, name string, inc context.IncrementCountByFunc) error {
	if entry, ok := Registry[name]; ok {
		parts := strings.Split(name, "/")
		if len(parts) == 2 {
			atom := engine.NewAtom(parts[0])
			arity, err := strconv.Atoi(parts[1])
			if err != nil {
				return err
			}

			hook := inc.By(entry.cost)
			p := entry.predicate(ctx)

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
