package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpg bytes"))
	if err := s.writeStream("myspecial picutres", data) ; err != nil{
		t.Error(err)
	}

}
func TestPathTransfromFunc(t *testing.T){
	key := "mombestpicture"
	pathname:= CASPathTransformFunc(key)
	fmt.Println(pathname)
	expectedPathName := "cf5d4/b01c4/d9438/c22c5/6c832/f83bd/3e8c6/304f9"
	if pathname != expectedPathName{
		t.Error(t, "have %s want %s", pathname, expectedPathName)

	}

}