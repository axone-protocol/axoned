package types

// String implements the Stringer interface.
func (p *Params) String() string {
	return p.Interpreter.String() + "\n" +
		p.Limits.String()
}
