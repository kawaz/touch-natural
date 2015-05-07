package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"os"
	"time"
)

func main() {
	rootPath := os.Args[1]
	_, err := touchNatural(rootPath)
	if err != nil {
		fmt.Println(err)
	}
}

func touchNatural(rootPath string) (maxModTime time.Time, err error) {
	maxModTime = time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	fi, err := os.Lstat(rootPath)
	if err != nil {
		return maxModTime, err
	}
	if fi.IsDir() {
		fis, err := ioutil.ReadDir(rootPath)
		if err != nil {
			return maxModTime, err
		}
		if(len(fis) == 0) {
			return fi.ModTime(), nil
		}
		for _, fi := range fis {
			childPath := path.Join(rootPath, fi.Name())
			mtime, err := touchNatural(childPath)
			if err != nil {
				return maxModTime, err
			}
			if mtime.After(maxModTime) {
				maxModTime = mtime
			}
		}
		if(!fi.ModTime().Equal(maxModTime)) {
			fmt.Println(fi.ModTime(), "->", maxModTime, rootPath)
			os.Chtimes(rootPath, time.Now(), maxModTime)
		}
	} else {
		maxModTime = fi.ModTime()
	}
	return maxModTime, nil
}
