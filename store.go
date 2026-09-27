package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "gjnetwork"
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
func (s *Store) readStream(key string)(io.ReadCloser, error){
	pathKey := s.PathTransformFunc(key)
	return os.Open(pathKey.FullPath())


}
func (p PathKey) FirstPathName() string{
	paths := strings.Split(p.Pathname, "/")
	if len(paths) == 0{
		return ""
	}
	return paths[0]
}
func (p PathKey) FullPath() string{
	return fmt.Sprintf("%s/%s", p.Pathname, p.FileName)
}
type StoreOpts struct{
	// root is the folder name fo the root containing all the files folders
	root 			  string
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

var DefaultPathTransformFunc = func(key string) PathKey{
	return PathKey{
		Pathname: key,
		FileName: key,
	}
}

func NewStore(opts StoreOpts) *Store{
	if opts.PathTransformFunc == nil{
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	
	if len(opts.root)==0{
		opts.root= defaultRootFolderName	
	}
	return &Store{
		StoreOpts: opts,

	}
}
func (s *Store) Has(key string) bool{
	PathKey := s.PathTransformFunc(key)

	_, err := os.Stat(PathKey.FullPath())
	if err == fs.ErrNotExist{
		return false;
	}
	
	
	return true;
}
func (s *Store) Delete(key string) error{
	pathKey := s.PathTransformFunc(key)

	defer func(){
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()


	//return os.RemoveAll(pathKey.FullPath())

	//fi, err := os.Stat()
	return os.RemoveAll(pathKey.FirstPathName())
}

func (s *Store) writeStream(key string,r io.Reader) error{
	pathKey := s.PathTransformFunc(key)
	
	if err := os.MkdirAll(pathKey.Pathname, os.ModePerm); err != nil{
		return err
	}

	fullPath := pathKey.FullPath()

	f, err := os.Create(fullPath)
	if err != nil{
		return err
	}
	n, err := io.Copy(f,r)
	if err != nil{
		return err
	}
	log.Printf("written (%d) bytes to disk: %s", n, fullPath)
	return nil
}

