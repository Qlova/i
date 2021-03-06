package Concept

import "reflect"
import . "github.com/qlova/script"
import "github.com/qlova/script/compiler"

import (
	//"fmt"
	//"os"
	//"bytes"
	//"runtime/debug"
)

func Create(c *compiler.Compiler, Name string, concept Concept, args func(*compiler.Compiler)) Func {

	if _, ok := Functions[Name]; !ok {
		f := c.Func(func() {
			c.GainScope()
			if args != nil {
				args(c)
			}
			
			c.CompileCache(concept.Cache)
			c.LoseScope()
		}, Name)
		Functions[Name] = f
	}

	return Functions[Name]
}

func CreateAndCall(c *compiler.Compiler, Name string, concept Concept) compiler.Type {
	
	if len(concept.Arguments) == 0 {
		c.Expecting("(")
		c.Expecting(")")
		
		f := Create(c, Name, concept, nil)
		
		return f.Run()

	}
	
	//TODO change name depending on types.
	c.Expecting("(")
	
	var Arguments []Type
	var Reflections []reflect.Type
	for c.Peek() != ")" {
		
		var expression = c.ScanExpression()
		Arguments = append(Arguments, expression)
		Reflections = append(Reflections, reflect.TypeOf(expression))
		
		//TODO do fancy stuff like lists, arrays and function subtypes. 
		//Convert them into Go equivilents.
		
		if c.Peek() != ")" {
			c.Expecting(",")
		}
	}
	
	c.Expecting(")")
	
	f := Create(c, Name, concept, func(c *compiler.Compiler) {
		for i := range Arguments {
			c.SetVariable(concept.Arguments[i], Arguments[i].Value().Arg(concept.Arguments[i]))
		}
	})
	
	if f.HasReturnValue() {
		return f.Call(Arguments...) //TODO deal with returns
	}

	f.Run(Arguments...)
	return nil
}
