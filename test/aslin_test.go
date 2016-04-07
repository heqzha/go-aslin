package test

import (
	"fmt"
	"time"
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"

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

func funcE(c *aslin.Context){
	//Don't forget call c.Abort() to finish workflow
	defer c.Abort()
	max := c.MustGet("loop_max").(int)
	for i := 0; i < max; i++{
		p, existed := c.Get("p")
		if existed{
			id := c.MustGet("id").(int)
			p := p.(int) + 1
			fmt.Printf("funcE - line: %d, p:%d\n", id, p)
			c.Set("p", p)
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
	fmt.Printf("funcF - line: %d, out:%d\n", id, p)
}

func funcG(c *aslin.Context){
	id := c.MustGet("id").(int)
	p := c.MustGet("p").(int)
	for{
		p *= 2
		fmt.Printf("funcG - line: %d, out:%d\n", id, p)
		time.Sleep(500 * time.Millisecond)
	}
}

// func TestAslin(t *testing.T){
//	fmt.Println("=====================================================================")
//	// Create new line
//	lIndex1 := aslin.InstFactory.NewLine(funcA, funcB, funcC)
//	lIndex2 := aslin.InstFactory.NewLine(funcA, funcB, funcB, funcB, funcC)

//	// Set parameters and run
//	aslin.InstFactory.Start(lIndex1, aslin.Params{
//		"id":1,
//	})
//	aslin.InstFactory.Start(lIndex2, aslin.Params{
//		"id":2,
//	})

//	// Clear all lines
//	defer aslin.InstFactory.Destory()

//	for {
//		//Wait for all lines stopped
//		if aslin.InstFactory.AreAllStop(){
//			break
//		}
//	}
//	assert.True(t, true, "True is true")
// }

// func TestAslinRepeat(t *testing.T){
//	fmt.Println("=====================================================================")
//	// Create new line
//	lIndex1 := aslin.InstFactory.NewLine(funcA, funcB, funcD)

//	// Set parameters and run
//	aslin.InstFactory.Start(lIndex1, aslin.Params{
//		"id":1,
//		"repeat_max":5,
//	})

//	// Clear all lines
//	defer aslin.InstFactory.Destory()

//	for {
//		//Wait for all lines stopped
//		if aslin.InstFactory.AreAllStop(){
//			break
//		}
//	}
//	assert.True(t, true, "True is true")
// }

func TestAslinPass(t *testing.T){
	fmt.Println("=====================================================================")
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
	assert.True(t, true, "True is true")
}

// func TestAslinPass2(t *testing.T){
//	fmt.Println("=====================================================================")
//	// Create new line
//	lIndex1 := aslin.InstFactory.NewLine(funcA, funcE, funcG)

//	// Set parameters and run
//	aslin.InstFactory.Start(lIndex1, aslin.Params{
//		"id":1,
//		"loop_max":5,
//	})

//	// Clear all lines
//	defer aslin.InstFactory.Destory()

//	for {
//		//Wait for all lines stopped
//		if aslin.InstFactory.AreAllStop(){
//			break
//		}
//	}
//	assert.True(t, true, "True is true")
// }
