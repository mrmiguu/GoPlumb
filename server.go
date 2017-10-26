package main

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"sync"

	"github.com/mrmiguu/sock"
)

var (
	importSync = regexp.MustCompile(`\bsync\.\b`)
	importTime = regexp.MustCompile(`\btime\.\b`)

	head = `package main
`

	tail = `
func print(args ...interface{}) {
	fmt.Print(args...)
}
func println(args ...interface{}) {
	fmt.Println(args...)
}
`
)

func init() {
	err := os.Mkdir("temp", os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func run(R chan<- string, E chan<- error) {
	cmd := exec.Command("go", "run", "temp/temp(+imports).go")
	sout, _ := cmd.StdoutPipe()
	serr, _ := cmd.StderrPipe()
	scanout := bufio.NewScanner(sout)
	scanerr := bufio.NewScanner(serr)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for scanout.Scan() {
			o := scanout.Text()
			println(o)
			R <- o
		}
	}()
	go func() {
		defer wg.Done()
		for scanerr.Scan() {
			e := scanerr.Text()
			println(e)
			E <- errors.New(e)
		}
	}()
	cmd.Start()
	wg.Wait()
	cmd.Wait()
}

func addImports() {
	cmd := exec.Command("goimports", "temp/temp.go")
	src, err := os.Create("temp/temp(+imports).go")
	if err != nil {
		println(err.Error())
		return
	}
	b, err := cmd.Output()
	if err != nil {
		println(err.Error())
		return
	}
	src.Write(b)
	src.Close()
}

func main() {
	F := sock.Rstring()
	R := sock.Wstring()
	E := sock.Werror()

	for f := range F {
		if len(f) == 0 {
			continue
		}

		src, _ := os.Create("temp/temp.go")
		src.WriteString(head)
		src.WriteString(`func main() {
`)
		src.WriteString(f)
		src.WriteString(`}`)
		src.WriteString(tail)
		src.Close()

		addImports()

		run(R, E)
	}
}
