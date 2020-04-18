package di_injector

import (
	"errors"
	"fmt"
	"reflect"
)

type DiContainer interface {
	AddToDependencies(dependencies ...interface{}) error
	InjectWithDependencies(object interface{}) error
}

type diContainer struct {
	dependencies []interface{}
}

/* NewDiContainer returns an empty container for your dependencies */
func NewDiContainer() DiContainer {
	return &diContainer{
		dependencies: []interface{}{},
	}
}

/* AddToDependencies lets you add dependencies to your container. To match a dependency to a field on the struct, the
dependency must ether implement the interface specified on the field or be the exact type as the field element */
func (dc *diContainer) AddToDependencies(dependencies ...interface{}) error {
	for _, dependency := range dependencies {

		// Skip nil dependency
		if reflect.TypeOf(dependency) == nil {
			continue
		}
		err := validateDependency(dependency)
		if err != nil {
			return errors.New("Cannot add the dependency " + fmt.Sprint(dependency) + "because " + err.Error())
		}
		dc.dependencies = append(dc.dependencies, dependency)
	}
	return nil
}

/* InjecWithDepedencies receives a pointer to the object you want to inject with dependencies */
func (dc *diContainer) InjectWithDependencies(object interface{}) error {
	err := validateObject(object)
	if err != nil {
		return err
	}
	var result error
	f := func() {
		defer func() {
			if err := recover(); err != nil {
				result = errors.New("Fatal Error at Injection")
			}
		}()
		err := injectObjectWithDependencies(object, dc.dependencies)
		if err != nil {
			result = err
		}
	}
	f()
	return result
}

/* Basic validation over the object that will receive the dependencies */
func validateObject(object interface{}) error {
	/*
		TODO: Validate if fields are exported to deliver a more specific error message
		TODO: Validate if more than one dependency can be injected into one field
	*/
	return nil
}

/* Basic validation over the dependency */
func validateDependency(dependency interface{}) error {
	/*
		TODO: Add validations
	*/
	return nil
}

func injectObjectWithDependencies(object interface{}, dependencies []interface{}) error {
	obj := reflect.ValueOf(object).Elem()
	typ := reflect.TypeOf(object).Elem()
	for i := 0; i < obj.NumField(); i++ {
		f := obj.Field(i)
		t := typ.Field(i)
		if t.Type.Kind() == reflect.Interface && t.Type.Name() == "" {
			return errors.New("Cannot inject into interface{}")
		}
		if t.Tag.Get("inject") == "auto" {
			injectOk := false
			for _, dependency := range dependencies {
				if f.Kind() == reflect.Interface && reflect.TypeOf(dependency).Implements(f.Type()) ||
					reflect.TypeOf(dependency) == f.Type() {

					// Add recursive injection
					if needsInjection(dependency){
						err := injectObjectWithDependencies(dependency, dependencies)
						if err != nil {
							return errors.New("Coudnt inject inner dependency -- " + err.Error())
						}
					}

					value := reflect.ValueOf(dependency)
					f.Set(value)
					injectOk = true
					break
				}
			}
			if !injectOk {
				return errors.New("No dependency injected on field " + t.Name)
			}
		}
	}
	return nil
}

func needsInjection(dependency interface{}) bool {
	if reflect.TypeOf(dependency).Kind() != reflect.Ptr {
		return false
	}
	typ := reflect.TypeOf(dependency).Elem()
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Tag.Get("inject") == "auto" {
			return true
		}
	}
	return false
}