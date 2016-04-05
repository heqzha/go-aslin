#Aslin - A Sequence Workflow Framework for Golang

## Example

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
