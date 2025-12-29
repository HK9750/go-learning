package main

import (
	"errors"
	"fmt"
	"os"
)

type ValidationError struct {
	field string
	message string
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("Error in field %s with message %s",v.field,v.message)
}

type NotFoundError struct {
	resource string
	id int
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("Resource %s of id %v not found",n.resource,n.id)
}

func getUser(id int) (string,error) {
	if id < 0 {
		return "", NotFoundError{resource: "user", id: id}
	}
	if id == 0 {
		return "", ValidationError{field: "id", message: "id cannot be zero"}
	}
	return "user", nil
}

func main() {
	_,err := getUser(-1)

	if err == nil {
		return
	}

	var vErr ValidationError;
	var nErr NotFoundError;

	if errors.As(err,&vErr) {
		fmt.Println("Error:",vErr.Error())
		os.Exit(2)
	}

	if errors.As(err,&nErr) {
		fmt.Println("Error:",nErr.Error())
		os.Exit(3)
	}

	fmt.Println("System Error",err)
	os.Exit(1)
}