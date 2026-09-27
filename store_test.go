package main

import (
	"bytes"
	"io"
	"testing"
)

func TestPathTransfromFunc(t *testing.T) {
	key := "momsbestpicture"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "cf5d4b01c4d9438c22c56c832f83bd3e8c6304f9"
	expectedPathName := "cf5d4/b01c4/d9438/c22c5/6c832/f83bd/3e8c6/304f9"
	if pathKey.Pathname != expectedPathName {
		t.Error(t, "have %s want %s", pathKey.Pathname, expectedPathName)

	}
	if pathKey.FileName != expectedOriginalKey {
		t.Error(t, "have %s want %s", pathKey.Pathname, expectedOriginalKey)

	}

}
func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momspecials"
	data := []byte("some jpg bytes")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err:= s.Read(key)
	if err != nil{
		t.Error(err)
	}

	b, err := io.ReadAll(r)

	if string(b) != string(data){
		t.Errorf("want %s have %s", data, b)
	}
	s.Delete(key)
	
}


func TestStoreDeleteKey(t *testing.T){
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momspecials"
	data := []byte("some jpg bytes")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	if err := s.Delete(key); err != nil{
		t.Error(err)
	}

}