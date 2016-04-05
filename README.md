#Aslin - A Sequence Workflow Framework for Golang

## Linear Workflow Example

```go
package main

import(
    "fmt"
    "errors"

    "github.com/heqzha/go-aslin/aslin"
)

func funcA(c *aslin.Context){
	c.Set("p", 0)
	id := c.MustGet("id").(int)

	c.Next()
}

func funcB(c *aslin.Context){
	p, existed := c.Get("p")
	if existed{
		intP := p.(int) + 1
		c.Set("p", intP)
        c.Next()
	}else{
	    fmt.Println(c.AbortWithError(errors.New("No params")))
	}
}

func funcC(c *aslin.Context){
	defer c.Abort()
	id := c.MustGet("id").(int)
	p := c.MustGet("p")
	fmt.Printf("line: %d, out:%d\n", id, p)
}

func main(){
	// Create new line
	lIndex1 := aslin.InstFactory.NewLine(funcA, funcB, funcC)

	// Set parameters and run
	aslin.InstFactory.Start(lIndex1, aslin.Params{
		"id":1,
	})

	// Clear all lines
	defer aslin.InstFactory.Destory()

	for {
		//Wait for all lines stopped
		if aslin.InstFactory.IsAllStop(){
			break
		}
	}
}

```
## Repeat Workflow Example
```go
package main

import(
    "fmt"
    "errors"

    "github.com/heqzha/go-aslin/aslin"
)

func funcA(c *aslin.Context){
	c.Set("p", 0)
	id := c.MustGet("id").(int)

	c.Next()
}

func funcB(c *aslin.Context){
	p, existed := c.Get("p")
	if existed{
		intP := p.(int) + 1
		c.Set("p", intP)
        c.Next()
	}else{
	    fmt.Println(c.AbortWithError(errors.New("No params")))
	}
}

func funcD(c *aslin.Context){
	max := c.MustGet("repeat_max").(int)
	repeat, existed := c.Get("repeat")
	if existed{
		if repeat.(int) < max{
			defer c.Repeat(1)
			repeat = repeat.(int) + 1
			c.Set("repeat", repeat)
		}else{
			defer c.Abort()
		}
	}else{
		c.Set("repeat", 0)
		defer c.Repeat(1)
	}
	id := c.MustGet("id").(int)
	p := c.MustGet("p")
	fmt.Printf("line: %d, out:%d\n", id, p)
}

func main(){
	// Create new line
	lIndex1 := aslin.InstFactory.NewLine(funcA, funcB, funcD)

	// Set parameters and run
	aslin.InstFactory.Start(lIndex1, aslin.Params{
		"id":1,
		"repeat_max":5,
	})

	// Clear all lines
	defer aslin.InstFactory.Destory()

	for {
		//Wait for all lines stopped
		if aslin.InstFactory.IsAllStop(){
			break
		}
	}
	assert.True(t, true, "True is true")
}
```
