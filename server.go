package main

import (
	"fmt"
	"log"
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

	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		fmt.Println("listner Accept error: ", err)
	// 		continue
	// 	}

	// 	go this.Handler(conn)
	// }

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

		name := strings.Fields(commandToExec)[0]
		args := strings.Fields(commandToExec)[1:]

		// --------------exec command----------------
		// out, err := exec.Command("date").Output()
		// out, err := exec.Command(commandToExec[:comLen-1]).Output()
		// out, err := exec.Command("cat", "/etc/passwd").Output()
		out, err := exec.Command(name, args...).Output()

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
