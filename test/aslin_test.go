package test

import (
	"fmt"
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"

	"github.com/heqzha/go-aslin/aslin"
)

func funcA(c *aslin.Context){
	c.Set("p", 0)
	id := c.MustGet("id").(int)
	fmt.Printf("line: %d\n", id)

	c.Next()
}

func funcB(c *aslin.Context){
	p, existed := c.Get("p")
	if existed{
		id := c.MustGet("id").(int)
		fmt.Printf("line: %d, p:%d\n", id, p)
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

func TestAslin(t *testing.T){
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
		if aslin.InstFactory.IsAllStop(){
			break
		}
	}
	assert.True(t, true, "True is true")
}
