package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"gopkg.in/cheggaaa/pb.v1"
)


func QueryTest(key string, fw *os.File) int {
	lead := "peer"
	cmd := "chaincode query -C channel47 -n mycc  -c "

	cmdstring := `{"Args":[ "query",`

	cmdstring2 := `"` + key + `"]}`

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
	elapse := time.Since(start)
	bufs := fmt.Sprintf("%s", buf)
	if (len(bufs)!=0) {
		elapse_string := strconv.FormatInt(int64(elapse/time.Millisecond), 10)
		result := key + ": " + "OK" + ": " + elapse_string + "\n"
		fw.WriteString(result)
		return 1
	} else {
		elapse_string := strconv.FormatInt(int64(elapse/time.Millisecond), 10)
		result := key + ": " + "FAILED" + ": " + elapse_string + "\n"
		fw.WriteString(result)
		return -1
	}
	return 1
}

func main() {
	/*
	call the function with step and the number of the blocks you need to put

	args[1]: step
	args[2]: Time for sleep
	args[3]: query times
	args[4]: MAX

	*/
	rand.Seed(time.Now().UTC().UnixNano())
	
	fw, err := os.OpenFile("./notesquery.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("errir in create file")
		panic(err)
	}
	defer fw.Close()

	var num uint64
	num, err = strconv.ParseUint(os.Args[3], 10, 64)
	max, err := strconv.ParseUint(os.Args[3], 10, 64)

	SLEEP, err := strconv.Atoi(os.Args[2])

	step, err := strconv.ParseUint(os.Args[1], 10, 64)

	if err != nil {
		return
	}

	//var step uint64

	var ret int
	var i uint64
	bar := pb.StartNew(int(num))
	for i=0;i<num;i++{

		val := rand.Uint64()% max
		ret = QueryTest(strconv.FormatUint(val, 10),fw)
		if ret == -1 {
			fmt.Printf("err in loop %d\n", i)
			break
		}
		if (i%step == 0){
			time.Sleep(time.Duration(SLEEP)*time.Second)
		}
	    bar.Increment()
	}

	bar.FinishPrint("The End!")

}
