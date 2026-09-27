package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "gjnetwork"

//CAS PATH TRANSFORM FUNCTION
func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	blocksize	:= 5
	slicelen := len(hashStr)/blocksize
	paths := make([]string, slicelen)
	for i :=0;i< slicelen;i++{
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from: to]
		
	}
	return PathKey{
		Pathname: strings.Join(paths, "/"),
		FileName: hashStr,

	}

}

type PathTransformFunc func(string) PathKey

type PathKey struct{
	Pathname string
	FileName string


}
//CLEAR
func (s *Store) Clear() error {
	
	return os.RemoveAll(s.Root)
}
//FIRST PATH NAME
func (p PathKey) FirstPathName() string{
	paths := strings.Split(p.Pathname, "/")
	if len(paths) == 0{
		return ""
	}
	return paths[0]
}
//FULL PATH
func (p PathKey) FullPath() string{
	return fmt.Sprintf("%s/%s", p.Pathname, p.FileName)
}
type StoreOpts struct{
	// Root is the folder name fo the Root containing all the files folders
	Root 			  string
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}
//DEFAULT PATH TRANSFORM FUNCTION
var DefaultPathTransformFunc = func(key string) PathKey{
	return PathKey{
		Pathname: key,
		FileName: key,
	}
}
//NEW STORE
func NewStore(opts StoreOpts) *Store{
	if opts.PathTransformFunc == nil{
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	
	if len(opts.Root)==0{
		opts.Root= defaultRootFolderName	
	}
	return &Store{
		StoreOpts: opts,

	}
}
//HAS
func (s *Store) Has(key string) bool{
	PathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, PathKey.FullPath())

	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)	
}
//DELETE 
func (s *Store) Delete(key string) error{
	pathKey := s.PathTransformFunc(key)

	defer func(){
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()

	firstPathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FirstPathName())

	return os.RemoveAll(firstPathNameWithRoot)

}

func (s *Store) Write(key string, r io.Reader)error{
	return s.writeStream(key, r)
}

//READ
func (s *Store) Read(key string) (io.Reader, error){
	f, err := s.readStream(key)
	if err != nil{
		return nil, err
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	
	return buf,  nil
}

//READ STREAM
func (s *Store) readStream(key string)(io.ReadCloser, error){
	pathKey := s.PathTransformFunc(key)
	pathKeyWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	return os.Open(pathKeyWithRoot)
	
}
//WRITE STREAM
func (s *Store) writeStream(key string,r io.Reader) error{
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.Pathname)

	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil{
		return err
	}
	
	fullPathWithRoot :=  fmt.Sprintf("%s/%s", s.Root,pathKey.FullPath())

	f, err := os.Create(fullPathWithRoot)
	if err != nil{
		return err
	}
	n, err := io.Copy(f,r)
	if err != nil{
		return err
	}
	defer f.Close()
	log.Printf("written (%d) bytes to disk: %s", n, fullPathWithRoot)
	return nil
}


