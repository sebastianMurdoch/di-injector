package di_injector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
Tests if the injection executes without error in a normal scenario.
*/
func TestDiContainer_InjectWithDependencies(t *testing.T) {
	assert := assert.New(t)
	str := "sample string"
	c := NewDiContainer()
	a := A{Dependency0:str}
	r := RunnerImpl{}
	err := c.AddToDependencies(str)
	if err != nil {
		t.Fail()
		return
	}
	err = c.AddToDependencies(&r)
	if err != nil {
		t.Fail()
		return
	}
	err = c.InjectWithDependencies(&a)
	if err != nil {
		t.Fail()
		return
	}
	assert.Equal(a.Dependency0, str)
	assert.Equal(a.Dependency1, str)
	assert.Equal(a.Dependency2, &r)
}

/*
Tests that the panic is cached
*/
func TestDiContainer_InjectWithDependencies_Panic(t *testing.T) {
	assert := assert.New(t)
	str := "sample string"
	c := NewDiContainer()
	a := A{Dependency0:str}
	err := c.InjectWithDependencies(a)
	assert.Equal(err.Error(), "Fatal Error at Injection")
}

/*
Tests the error result when no injections occurred
*/
func TestDiContainer_InjectWithDependencies_NoInjections(t *testing.T) {
	assert := assert.New(t)
	str := "sample string"
	c := NewDiContainer()
	a := A{Dependency0:str}
	err := c.InjectWithDependencies(&a)
	assert.Equal(err.Error(), "No dependency injected on field Dependency1")
}

/*
Tests that nil dependencies cannot be injected
*/
func TestNewDiContainer_NilDependency(t *testing.T) {
	assert := assert.New(t)
	c := NewDiContainer()
	err := c.AddToDependencies(nil)
	assert.NotNil(err)
}

/*
Checks that any of the fields to be injected couldn't be interface{}. This kind of fields implement every dependency so
this is not allowed
*/
func TestNewDiContainer_NoInterface(t *testing.T) {
	assert := assert.New(t)
	c := NewDiContainer()
	b := B{Dependency0:""}
	err := c.InjectWithDependencies(&b)
	assert.Equal(err.Error(), "Cannot inject into interface{}")
}

/*
Tests recursive injection
*/
func TestDiContainer_InjectWithDependencies_recursion(t *testing.T) {
	assert := assert.New(t)
	str := "sample string"
	c := NewDiContainer()
	a := A{Dependency0:str}
	r := C{}
	err := c.AddToDependencies(str)
	if err != nil {
		t.Fail()
		return
	}
	err = c.AddToDependencies(&r)
	if err != nil {
		t.Fail()
		return
	}
	err = c.InjectWithDependencies(&a)
	if err != nil {
		t.Fail()
		return
	}
	assert.Equal(a.Dependency0, str)
	assert.Equal(a.Dependency1, str)
	assert.Equal(a.Dependency2, &r)
	assert.Equal(a.Dependency2.Run(), str)
}
/*
Tests recursive injection
*/
func TestDiContainer_InjectWithDependencies_recursionError(t *testing.T) {
	assert := assert.New(t)
	str := "sample string"
	c := NewDiContainer()
	a := A{Dependency0:str}
	r := D{}
	err := c.AddToDependencies(str)
	if err != nil {
		t.Fail()
		return
	}
	err = c.AddToDependencies(&r)
	if err != nil {
		t.Fail()
		return
	}
	err = c.InjectWithDependencies(&a)
	assert.NotNil(err)
}




