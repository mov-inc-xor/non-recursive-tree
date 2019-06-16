package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

//TODO: replace / for os.PathSeparator

type stack struct {
	arr []string
	sz  int
}

func newStack() *stack {
	return &stack{
		arr: make([]string, 0),
		sz:  0,
	}
}

func (s *stack) push(str string) {
	s.arr = append(s.arr, str)
	s.sz++
}

func (s *stack) pop() (str string) {
	if s.sz == 0 {
		str = ""
		return
	}
	s.sz--
	str = s.arr[s.sz]
	s.arr = s.arr[:s.sz]
	return
}

func (s *stack) size() int {
	return s.sz
}

func dirTree(out interface{}, path string, printFiles bool) error {
	st := newStack()
	st.push(path)
	pathArr := make([]string, 0)
	fileSize := make(map[string]int64)
	for st.size() != 0 {
		p := st.pop()
		pathArr = append(pathArr, p)
		dir, err := os.Open(p)
		if err != nil {
			return err
		}
		if stat, _ := dir.Stat(); !stat.IsDir() {
			fileSize[p] = stat.Size()
			continue
		}
		infos, err := dir.Readdir(-1)
		if err != nil {
			return err
		}
		for i := range infos {
			if !printFiles && !infos[i].IsDir() {
				continue
			}
			st.push(p + string(os.PathSeparator) + infos[i].Name())
		}
	}
	sort.Strings(pathArr)

	lastDirs := make(map[string]string)

	for i := range pathArr {
		dirs := strings.Split(pathArr[i], "/")
		dirs = dirs[:len(dirs)-1]
		outerDir := strings.Join(dirs, "/")
		lastDirs[outerDir] = pathArr[i]
	}

	for i := range pathArr {
		dirs := strings.Split(pathArr[i], "/")
		if len(dirs) < 2 {
			continue
		}
		innerDir := dirs[0] + "/" + dirs[1]
		outerDir := dirs[0]
		prefix := ""
		for j := range dirs[:len(dirs)-1] {
			if lastDirs[outerDir] == innerDir {
				if j == len(dirs)-2 {
					prefix += "└───"
				} else {
					prefix += "\t"
				}
			} else {
				if j == len(dirs)-2 {
					prefix += "├───"
				} else {
					prefix += "│\t"
				}
			}
			if j+2 >= len(dirs) {
				break
			}
			innerDir += "/" + dirs[j+2]
			outerDir += "/" + dirs[j+1]
		}

		if size, ok := fileSize[pathArr[i]]; ok {
			if size == 0 {
				fmt.Fprintln(out.(io.Writer), prefix+dirs[len(dirs)-1]+" (empty)")
			} else {
				fmt.Fprintln(out.(io.Writer), prefix+dirs[len(dirs)-1]+" ("+strconv.FormatInt(size, 10)+"b"+")")
			}
			continue
		}
		fmt.Fprintln(out.(io.Writer), prefix+dirs[len(dirs)-1])
	}

	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
