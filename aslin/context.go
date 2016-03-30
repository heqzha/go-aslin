package aslin

import (
	"math"
)

const abortIndex int8 = math.MaxInt8 / 2

type Params map[string]interface{}

type Context struct {
	params   Params
	line AsLine
	index    int8

	Errors   errors
	Accepted []string
}

/************************************/
/********** CONTEXT CREATION ********/
/************************************/

func (c *Context) reset() {
	c.params = make(Params)
	c.line = nil
	c.index = -1
	c.Errors = c.Errors[0:0]
	c.Accepted = nil
}

// Copy returns a copy of the current context that can be safely used outside the request's scope.
// This have to be used then the context has to be passed to a goroutine.
func (c *Context) Copy() *Context {
	var cp Context = *c
	cp.index = abortIndex
	cp.line = nil
	return &cp
}

// LineName returns the main line's name. For example if the line is "handleGetUsers()", this
// function will return "main.handleGetUsers"
func (c *Context) LineName() string {
	return nameOfFunction(c.line.Last())
}

/************************************/
/*********** FLOW CONTROL ***********/
/************************************/

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in github.
func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// IsAborted returns true if the currect context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized. If the
// authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() {
	c.index = abortIndex
}

// AbortWithStatus calls `Abort()` and writes the headers with the specified status code.
// For example, a failed attempt to authentificate a request could use: context.AbortWithStatus(401).
func (c *Context) AbortWithStatus(code int) {
	c.Status(code)
	c.Abort()
}

// AbortWithError calls `AbortWithStatus()` and `Error()` internally. This method stops the chain, writes the status code and
// pushes the specified error to `c.Errors`.
// See Context.Error() for more details.
func (c *Context) AbortWithError(code int, err error) *Error {
	c.AbortWithStatus(code)
	return c.Error(err)
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
