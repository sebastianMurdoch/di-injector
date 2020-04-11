package di_injector

import (
	"errors"
	"fmt"
	"reflect"
)

type diContainer struct {
	dependencies []interface{}
}

func NewDiContainer() *diContainer {
	return &diContainer{
		dependencies: []interface{}{},
	}
}

func (dc *diContainer) AddToDependencies(dependency interface{}) error{
	err := validateDependency(dependency)
	if err != nil {
		return errors.New("Cannot add the dependency " + fmt.Sprint(dependency) + "because " + err.Error())
	}
	dc.dependencies = append(dc.dependencies, dependency)
	return nil
}

func (dc *diContainer) InjecWithDepedencies(object interface{}) error{
	err := validateObject()
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
		obj := reflect.ValueOf(object).Elem()
		typ := reflect.TypeOf(object).Elem()
		for i := 0; i < obj.NumField(); i++ {
			f := obj.Field(i)
			t := typ.Field(i)
			if t.Type.Kind() == reflect.Interface && t.Type.Name() == "" {
				result = errors.New("Cannot inject into interface{}")
				return
			}
			if t.Tag.Get("inject")=="auto"{
				injectOk := false
				for _, dependency := range dc.dependencies {
					if f.Kind() == reflect.Interface && reflect.TypeOf(dependency).Implements(f.Type()) ||
						reflect.TypeOf(dependency) == f.Type(){
						value := reflect.ValueOf(dependency)
						f.Set(value)
						injectOk = true
						break
					}
				}
				if !injectOk {
					result = errors.New("No dependency injected on field " + t.Name)
					return
				}
			}

		}
	}
	f()
	return result
}

func validateObject() error {
	/*
	TODO: Validate if fields are exported to deliver a more specific error message
	TODO: Validate if more than one dependency can be injected into one field
	*/
	return nil
}

/* Basic validation over the dependency */
func validateDependency(dependency interface{}) error  {
	dType := reflect.TypeOf(dependency)
	if dType == nil{
		return errors.New("dependency type cannot be nil")
	}
	return nil
}
