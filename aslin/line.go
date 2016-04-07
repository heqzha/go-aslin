package aslin

import(
	"reflect"
)

type HandlerFunc func(*Context)

type node struct{
	handler HandlerFunc
	ch chan *Context
}

func (n *node)init(f HandlerFunc){
	n.handler = f
	n.ch = make(chan *Context)
}

func (n *node)copy()(*node){
	var cp = *n
	cp.ch = make(chan *Context)
	return &cp
}

func (n *node)destory(){
	close(n.ch)

	n.handler = nil
	n.ch = nil
}

func (n *node)in(c *Context){
	go func(){
		n.ch <- c
	}()
}

type line struct{
	nodes []*node
}

func (l *line)getHandler(index int)(f HandlerFunc){
	if index < len(l.nodes){
		return l.nodes[index].handler
	}
	return nil
}

func (l *line)copy()(*line){
	var cp = new(line)
	for _, n := range l.nodes{
		cp.nodes = append(cp.nodes, n.copy())
	}
	return cp
}

func (l *line)add(f HandlerFunc){
	var n = new(node)
	n.init(f)
	l.nodes = append(l.nodes, n)
}

// next return the next node of i,
// if the next one reaches the end of line,
// return nil and true, otherwise return
// node[i+1] and false
func (l *line)next(i int)(*node, bool){
	if i + 1 < len(l.nodes){
		return l.nodes[i + 1], false
	}
	//reach the end of line
	return nil, true
}

// do actives handler i and send context to it
func (l *line)do(i int, c *Context) bool{
	if i < len(l.nodes){
		l.nodes[i].in(c)
		return true
	}
	return false
}

func (l *line)size() int{
	return len(l.nodes)
}

func (l *line)init(handlers ...HandlerFunc){
	l.nodes = make([]*node, 0)
	if len(handlers) >= int(abortIndex){
		panic("too many handlers")
	}

	for _, h := range handlers{
		l.add(h)
	}
}

func (l *line)stop(){
	for _, n := range l.nodes{
		n.destory()
	}
	l.nodes = nil
}

func (l *line)start(){
	cases := make([]reflect.SelectCase, len(l.nodes))
	for i, n := range l.nodes{
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(n.ch)}
	}

	for {
		chosen, value, ok := reflect.Select(cases)
		// ok will be true if the channel has not been closed.
		if ok{
			n := l.nodes[chosen]
			c := value.Interface().(*Context)
			c.index = int8(chosen)
			go (n.handler)(c)
		}else{
			// end process if any of channels in line has been closed.
			break
		}
	}
}
