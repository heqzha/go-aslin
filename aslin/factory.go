package aslin

type Process func(*Context)

type Line struct{
	ctxt *Context
	procs []*Process
}

type Factory struct{
	group []*Line
}

func InitFactory() *Factory{
	return &Factory{
		group:[]*Line{},
	}
}

func (f *Factory)AddLine(c *Context, procs ...*Process){
	line := &Line{
		ctxt:c,
		procs:procs,
	}
	f.group = append(f.group, line)
}
