package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func uploadTest(start, end string) int {
	lead := "peer"
	cmd := "chaincode invoke -C channel47 -n mycc --waitForEvent -c"

	cmdstring := `{"Args":[ "uploadbulktest",`

	cmdstring2 := `"` + start + `"` + `,"` + end + `"]}`

	cmdstring = cmdstring + cmdstring2
	args := []string{}
	for _, each := range (strings.Split(cmd, " ")) {
		args = append(args, each)
	}
	args = append(args, cmdstring)
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
	fmt.Println(bufarray[1])
	if (strings.Contains(bufarray[1], "VALID") == true) {
		return 1
	} else {
		return -1
	}
	return 1
}

func main() {
	/*
	call the function with step and the number of the blocks you need to put

	args[1]: step
	args[2]: Time for sleep
	args[3]: start of the transaction
	args[4]: end of transactions

	*/
	var start uint64
	start, err := strconv.ParseUint(os.Args[3], 10, 64)

	num, err := strconv.ParseUint(os.Args[4], 10, 64)
	SLEEP, err := strconv.Atoi(os.Args[2])

	step, err := strconv.ParseUint(os.Args[1], 10, 64)

	if err != nil {
		return
	}

	//var step uint64

	var end uint64
	var ret int
	var exception bool
	exception = false

	for (true) {
		end = start + step
		out := fmt.Sprintf("uploading... %d -- %d", start, end)
		fmt.Println(out)
		ret = uploadTest(strconv.FormatUint(start, 10), strconv.FormatUint(end, 10))
		if ret == -1 {
			exception = true
			break
		}
		//fmt.Println(out)
		start = end
		if (end+step >= num) {
			break
		}
		time.Sleep(time.Duration(SLEEP)*time.Second)
	}
	if (exception == true) {
		fmt.Println("error in upload the test")
		return
	}
	time.Sleep(time.Duration(SLEEP)*time.Second)
	out := fmt.Sprintf("uploading... %d -- %d", start, num)
	fmt.Println(out)
	ret = uploadTest(strconv.FormatUint(start, 10), strconv.FormatUint(num, 10))
	if ret == -1 {

		fmt.Println("error in upload test")
	}

}
