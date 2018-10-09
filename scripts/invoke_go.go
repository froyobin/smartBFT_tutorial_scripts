package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/cheggaaa/pb"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var writelock *sync.RWMutex

func invokejs(line string, f *os.File) {
	lead := "peer"
	cmd := "chaincode invoke -C channel47 -n mycc --waitForEvent -c"

	cmdstring := `{"Args":[ "uploaddomain",`

	cmdstring2 := `"` + line + `"` + `,"111111"]}`
	//	fmt.Println(cmdstring+cmdstring2)
	cmdstring = cmdstring + cmdstring2
	args := []string{}
	for _, each := range (strings.Split(cmd, " ")) {
		args = append(args, each)
	}
	args = append(args, cmdstring)
	start := time.Now()
	process := exec.Command(lead, args...)
	stdin, err := process.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stdin.Close()
	buf := new(bytes.Buffer) // THIS STORES THE NODEJS OUTPUT
	process.Stdout = nil
	process.Stderr = buf

	if err = process.Start(); err != nil {
		fmt.Println("An error occured: ", err)
	}

	process.Wait()
	bufs := fmt.Sprintf("%s", buf)
	bufarray := strings.Split(bufs, "\n")
	retval := "FAILED"
	fmt.Println(bufarray[1])
	if (strings.Contains(bufarray[1], "VALID") == true) {
		retval = "OK"
	}
	elapse := time.Since(start)
	//fmt.Println(elapse)
	elapse_string := strconv.FormatInt(int64(elapse/time.Millisecond), 10)
	result := line + ": " + retval + ": " + elapse_string + "\n"
	writelock.Lock()
	f.WriteString(result)
	writelock.Unlock()
}

func addJob(f *os.File, jobs chan<- string) {

	r := bufio.NewReader(f)
	line, err := r.ReadString('\n') // line defined once
	line = strings.TrimSuffix(line, "\n")
	jobs <- line
	for err != io.EOF {
		//invokejs(line, fw)
		line, err = r.ReadString('\n') // line defined once
		line = strings.TrimSuffix(line, "\n")
		if (len(line)) == 0 {
			continue
		}
		jobs <- line
		//bar.Increment()
		// fmt.Print(line)              // or any stuff

	}
	close(jobs)
}

func doJob(
	jobs <-chan string, dones chan<- struct{},
	fw *os.File, i int, bar *pb.ProgressBar) {


	for element := range jobs {

		invokejs(element,fw)
		bar.Increment()

	}
	dones <- struct{}{}

}

func main() {

	var worker = runtime.NumCPU()
	working := worker
	jobs := make(chan string, worker)
	dones := make(chan struct{}, worker)

	done := false
	writelock = new(sync.RWMutex)

	fw, err := os.OpenFile("./notes.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("errir in create file")
		panic(err)
	}
	defer fw.Close()
	out, err := exec.Command("wc", "-l", "./url.dat").Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	num, err := strconv.Atoi(strings.Split(string(out), " ")[0])
	bar := pb.StartNew(num)
	fmt.Println(num)
	f, err := os.Open("url.dat")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	go addJob(f, jobs)


	for i := 0; i < worker; i++ {
		go doJob(jobs, dones, fw, i, bar)
	}
	for {
		<-dones
		working -= 1
		if working <= 0 {
			done = true
		}
		if done == true {
			break
		}

	}

	bar.FinishPrint("The End!")
	return

}
