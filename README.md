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
	fmt.Printf("funcA - line: %d\n", id)

	//Go to next process
	c.Next()
}

func funcB(c *aslin.Context){
	p, existed := c.Get("p")
	if existed{
		id := c.MustGet("id").(int)
		fmt.Printf("funcB - line: %d, p:%d\n", id, p)
		intP := p.(int) + 1
		c.Set("p", intP)
		c.Next()
	}else{
		// Abort current process
		fmt.Println(c.AbortWithError(errors.New("No params")))
	}
}

func funcC(c *aslin.Context){
	//Don't forget call c.Abort() to finish workflow
	defer c.Abort()
	id := c.MustGet("id").(int)
	p := c.MustGet("p")
	fmt.Printf("funcC - line: %d, out:%d\n", id, p)
}

func main(){
	// Create new line
	lIndex1 := aslin.InstFactory.NewLine(funcA, funcB, funcC)
	lIndex2 := aslin.InstFactory.NewLine(funcA, funcB, funcB, funcB, funcC)

	// Set parameters and run
	aslin.InstFactory.Start(lIndex1, aslin.Params{
		"id":1,
	})
	aslin.InstFactory.Start(lIndex2, aslin.Params{
		"id":2,
	})

	// Clear all lines
	defer aslin.InstFactory.Destory()

	for {
		//Wait for all lines stopped
		if aslin.InstFactory.AreAllStop(){
			break
		}
	}
}

```
    Output:
    funcA - line: 1
    funcB - line: 1, p:0
    funcC - line: 1, out:1
    funcA - line: 2
    funcB - line: 2, p:0
    funcB - line: 2, p:1
    funcB - line: 2, p:2
    funcC - line: 2, out:3

## Example of Using Repeat in Workflow
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
	fmt.Printf("funcA - line: %d\n", id)

	//Go to next process
	c.Next()
}

func funcB(c *aslin.Context){
	p, existed := c.Get("p")
	if existed{
		id := c.MustGet("id").(int)
		fmt.Printf("funcB - line: %d, p:%d\n", id, p)
		intP := p.(int) + 1
		c.Set("p", intP)
		c.Next()
	}else{
		// Abort current process
		fmt.Println(c.AbortWithError(errors.New("No params")))
	}
}

func funcD(c *aslin.Context){
	max := c.MustGet("repeat_max").(int)
	repeat, existed := c.Get("repeat")
	if existed{
		if repeat.(int) < max{
			//Repeat workflow at funcB
			defer c.Repeat(1)
			repeat = repeat.(int) + 1
			c.Set("repeat", repeat)
		}else{
			// Reach the max repeat times, abort workflow
			defer c.Abort()
		}
	}else{
		c.Set("repeat", 0)
		defer c.Repeat(1)
	}
	id := c.MustGet("id").(int)
	p := c.MustGet("p")
	fmt.Printf("funcD - line: %d, out:%d\n", id, p)
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
		if aslin.InstFactory.AreAllStop(){
			break
		}
	}

}
```
    Output:
    funcA - line: 1
    funcB - line: 1, p:0
    funcD - line: 1, out:1
    funcB - line: 1, p:1
    funcD - line: 1, out:2
    funcB - line: 1, p:2
    funcD - line: 1, out:3
    funcB - line: 1, p:3
    funcD - line: 1, out:4
    funcB - line: 1, p:4
    funcD - line: 1, out:5
    funcB - line: 1, p:5
    funcD - line: 1, out:6
    funcB - line: 1, p:6
    funcD - line: 1, out:7

## Example of Using Pass in Workflow
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
	fmt.Printf("funcA - line: %d\n", id)

	//Go to next process
	c.Next()
}

func funcE(c *aslin.Context){
	//Don't forget call c.Abort() to finish workflow
	defer c.Abort()
	max := c.MustGet("loop_max").(int)
	for i := 0; i < max; i++{
		p, existed := c.Get("p")
		if existed{
			id := c.MustGet("id").(int)
			fmt.Printf("funcE - line: %d, p:%d\n", id, p)
			intP := p.(int) + 1
			c.Set("p", intP)
			c.Pass()
		}else{
			// Abort current process
			fmt.Println(c.AbortWithError(errors.New("No params")))
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func funcF(c *aslin.Context){
	//Don't call c.Abort() here
	id := c.MustGet("id").(int)
	p := c.MustGet("p")
	fmt.Printf("funcC - line: %d, out:%d\n", id, p)
}

func main(){
	// Create new line
	lIndex1 := aslin.InstFactory.NewLine(funcA, funcE, funcF)

	// Set parameters and run
	aslin.InstFactory.Start(lIndex1, aslin.Params{
		"id":1,
		"loop_max":5,
	})

	// Clear all lines
	defer aslin.InstFactory.Destory()

	for {
		//Wait for all lines stopped
		if aslin.InstFactory.AreAllStop(){
			break
		}
	}
}
```
    Output:
    funcA - line: 1
    funcE - line: 1, p:0
    funcF - line: 1, out:1
    funcE - line: 1, p:1
    funcF - line: 1, out:2
    funcE - line: 1, p:2
    funcF - line: 1, out:3
    funcE - line: 1, p:3
    funcF - line: 1, out:4
    funcE - line: 1, p:4
    funcF - line: 1, out:5
