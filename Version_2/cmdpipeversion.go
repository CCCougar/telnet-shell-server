package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/reiver/go-telnet"
)

func main() {
	var handler telnet.Handler = myHandler
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

var myHandler telnet.Handler = internalEchoHandler{}

type internalEchoHandler struct{}

func (handler internalEchoHandler) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {

	var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
	p := buffer[:]
	var commandToExec string
	var comLen int = 0
	cmd := exec.Command("cmd.exe")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Stdout = w
	cmd.Stderr = w
	// stdoutStderr, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// cmd.Run() // Run starts the specified command and waits for it to complete.
	cmd.Start() // Start starts the specified command but does not wait for it to complete.

	for {
		// Print propt
		// if comLen == 0 {
		// 	oi.LongWrite(w, []byte("$ "))
		// }

		_, err := r.Read(p)
		// fmt.Println("You pressed: ", p)

		if err != nil {
			log.Fatal("r.Read error: ", err)
		}

		// Every end of an input is "\x0d\x0a"
		if p[0] != 0x0a {
			comLen++
			commandToExec = commandToExec + string(p[0]) // the format of commandToExec is: "command\x0d"
			continue
		}
		theCommand := commandToExec[:len(commandToExec)-1] // get rid of the last byte "\x0d"
		fmt.Printf("******DEBUG****** >>>>> execting: " + theCommand + "\n")
		io.WriteString(stdin, fmt.Sprintf("%s\n", theCommand))

		commandToExec = ""
		comLen = 0
	}
}

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"io"
// 	"log"
// 	"os"
// 	"os/exec"
// )

// func main() {
// 	cmd := exec.Command("cmd.exe")
// 	stdin, err := cmd.StdinPipe()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	cmd.Start()
// 	inputReader := bufio.NewReader(os.Stdin)
// 	for {
// 		input := ""
// 		input, _ = inputReader.ReadString('\n')
// 		// fmt.Println("****DEBUG****" + input)
// 		io.WriteString(stdin, fmt.Sprintf("%s\n", input))
// 	}
// }
// interactive shell
