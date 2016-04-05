package aslin

// Factory contains multiple lines and their contexts
type Factory struct{
	lines []*line
	contexts []*Context
}

// InstFactory is the instance of factory
var InstFactory = new(Factory)

// NewLine creates a new line and returns it's index
func (f *Factory)NewLine(handlers ...HandlerFunc)(int){
	l := new(line)
	l.init(handlers ...)
	f.lines = append(f.lines, l)
	return len(f.lines) - 1
}

// Start initializes context and starts existed line by index
func (f *Factory)Start(i int, p Params){
	l := f.lines[i]
	if l != nil{
		ctxt := new(Context)
		ctxt.init(p, l)
		go ctxt.line.start()
		f.contexts = append(f.contexts, ctxt)
	}
}

// Stop ends a running line by index
func (f *Factory)Stop(i int){
	c := f.contexts[i]
	if c != nil{
		c.Abort()
	}
}

// IsStop returns true, if the line is stopped
func (f *Factory)IsStop(i int) bool{
	c := f.contexts[i]
	if c != nil{
		return c.IsAborted()
	}
	return true
}

// AreAllStop return false, if anyone of lines is not stopped
func (f *Factory)AreAllStop() bool{
	for _, c :=  range f.contexts{
		if c != nil && !c.IsAborted(){
			return false
		}
	}
	return true
}

// Destory ends all lines
func (f *Factory)Destory(){
	for _, c := range f.contexts{
		c.Abort()
	}
}
