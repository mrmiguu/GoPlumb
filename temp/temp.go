package main
import "fmt"
import "time"
import "sync"
func main() {
var wg sync.WaitGroup
here := make(chan bool)
f := func() {
	defer wg.Done()
	time.Sleep(1*time.Second)
	println("here?")
	<-here
	println("here!")
}
wg.Add(1)
go f()
println(`"here"`)
time.Sleep(2*time.Second)
here <- true
wg.Wait()}
func print(args ...interface{}) {
	fmt.Print(args...)
}
func println(args ...interface{}) {
	fmt.Println(args...)
}
