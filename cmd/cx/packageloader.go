package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/blake2b"
)

type File struct {
	FileName   string
	Length     uint32
	Content    []byte
	Blake2Hash string
}

type Package struct {
	PackageName string
	Files       []string
}

type PackageList struct {
	Packages []string
}

func main() {
	dirList := []string{}
	err := filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirList = append(dirList, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	pl := PackageList{}
	var importList []string

	for _, path := range dirList {
		fileList, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}

		p := Package{}

		for i := 1; i < len(fileList); i++ {
			if fileList[i].IsDir() {
				continue
			}
			samePackage, packageName, err := comparePackages(path+"/"+fileList[i].Name(), path+"/"+fileList[i-1].Name())
			if err != nil {
				log.Fatal(err)
			}
			if !samePackage {
				log.Print("Files in directory " + path + " are not all in the same package.\nSource of the error: " + fileList[i].Name())
				log.Fatal("ErrMismatchedPackageFiles")
			}
			newImports, err := getImports(path + "/" + fileList[i].Name())
			if err != nil {
				log.Fatal(err)
			}
			importList = append(importList, newImports...)
			f := File{
				FileName: fileList[i].Name(),
				Length:   uint32(fileList[i].Size()),
			}
			h := blake2b.Sum512(f.Content)
			f.Blake2Hash = string(h[:])
			f.Content, err = ioutil.ReadFile(path + "/" + fileList[i].Name())
			if err != nil {
				log.Fatal(err)
			}
			p.PackageName = packageName
			p.hashFile(f)
		}
		pl.hashPackage(p)
	}
}

func (p Package) hashFile(f File) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	err := enc.Encode(f)
	if err != nil {
		return err
	}
	h := blake2b.Sum512(buffer.Bytes())
	p.Files = append(p.Files, string(h[:]))
	return nil
}

func (pl PackageList) hashPackage(p Package) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	err := enc.Encode(p)
	if err != nil {
		return err
	}
	h := blake2b.Sum512(buffer.Bytes())
	pl.Packages = append(pl.Packages, string(h[:]))
	return nil
}

func comparePackages(filepath1 string, filepath2 string) (bool, string, error) {
	file1, err := os.Open(filepath1)
	if err != nil {
		return false, "", err
	}
	defer file1.Close()
	file2, err := os.Open(filepath2)
	if err != nil {
		return false, "", err
	}
	defer file2.Close()
	scanner1 := bufio.NewScanner(file1)
	scanner2 := bufio.NewScanner(file2)
	scanner1.Split(bufio.ScanWords)
	scanner2.Split(bufio.ScanWords)
	for scanner1.Text() != "package" {
		scanner1.Scan()
	}
	for scanner2.Text() != "package" {
		scanner2.Scan()
	}
	scanner1.Scan()
	scanner2.Scan()
	return scanner1.Text() == scanner2.Text(), scanner1.Text(), nil
}

func getImports(filepath string) (imports []string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return []string{}, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		if scanner.Text() == "import" {
			scanner.Scan()
			imports = append(imports, scanner.Text())
		}
	}
	return imports, nil
}
