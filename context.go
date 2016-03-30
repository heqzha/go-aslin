package aslin

type Params map[string]interface{}

type Context struct {
	params   Params
	line AsLine
	index    int8

	Errors   errorMsgs
	Accepted []string
}

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// Set is used to store a new key/value pair exclusivelly for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value interface{}) {
	if c.params == nil {
		c.params = make(Params)
	}
	c.params[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	if c.params != nil {
		value, exists = c.params[key]

	}
	return
}

// Returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value

	}
	panic("Key \"" + key + "\" does not exist")
}
