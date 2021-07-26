package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"
)

func main() {

	var handler telnet.Handler = EchoHandler

	server := &telnet.Server{
		Addr:    ":5555",
		Handler: handler,
	}

	err := server.ListenAndServe()
	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}

}

var EchoHandler telnet.Handler = internalEchoHandler{}

type internalEchoHandler struct{}

func (handler internalEchoHandler) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {

	var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
	p := buffer[:]
	var commandToExec string
	var comLen int = 0

	for {
		if comLen == 0 {
			oi.LongWrite(w, []byte("$ "))
		}
		_, err := r.Read(p)
		if err != nil {
			log.Fatal("r.Read error: ", err)
		}
		if p[0] != 0x0a {
			comLen++
			// fmt.Println(p[0])
			commandToExec = commandToExec + string(p[0])
			continue
		}
		fmt.Printf("******DEBUG****** >>>>>execting: " + commandToExec + "\n")

		arg := strings.Fields(commandToExec)

		name := arg[0]
		args := arg[1:]

		switch name {

		case "cd":
			// 'cd' to home dir with empty path not yet supported.
			if len(args) < 1 {
				// errors.New("path required")
				oi.LongWrite(w, []byte("path required"))
				commandToExec = ""
				comLen = 0
				continue
			}
			err := os.Chdir(args[0])
			if err != nil {
				fmt.Println(err)
			}
			commandToExec = ""
			comLen = 0
			continue
		case "dir":
			files, _ := ioutil.ReadDir("./")
			for _, f := range files {
				oi.LongWrite(w, []byte(f.Name()+"\n"))
			}
			commandToExec = ""
			comLen = 0
			continue
		case "exit":
			os.Exit(0)
			break
		case "execute":
			if len(args) < 1 {
				// errors.New("path required")
				oi.LongWrite(w, []byte("what file to execute?"))
				commandToExec = ""
				comLen = 0
				continue
			}
			// There maybe some problems
			output, err := exec.Command(args[0]).Output()
			if err != nil {
				fmt.Println(err)
			}
			oi.LongWrite(w, []byte(output))
			commandToExec = ""
			comLen = 0
			continue
		}

		out, err := exec.Command(name, args...).Output()
		// out, err := exec.Command("/bin/bash", "-c", commandToExec).Output()

		// name := "echo"
		// args := []string{"hello", "world"}
		// out, err := exec.Command(name, args...).Output()

		if err != nil {
			// log.Fatal(err)
			// oi.LongWrite(w, []byte("Your input is not a valid command"))
			fmt.Println("exec.Command error: ", err)
			commandToExec = ""
			comLen = 0
			continue
		}
		// ------------------------------------------

		if len(commandToExec) > 0 {
			// oi.LongWrite(w, p[:n])
			oi.LongWrite(w, []byte(out))
		}
		commandToExec = ""
		comLen = 0
		if nil != err {
			break
		}
	}
}
