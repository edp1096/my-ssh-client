package main // import "ssh-client"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	gotty "github.com/mattn/go-tty"
	"golang.org/x/crypto/ssh"
)

type HostInfo struct {
	Name           string
	Address        string
	Port           int
	Username       string
	Password       string
	PrivateKeyFile string
	PrivateKeyText string
}

var (
	hostsFile = flag.String("f", "", "host data file (required)")
	hostsIDX  = flag.Int("i", 0, "index of host data (required)")

	hosts []HostInfo // All host data
	host  HostInfo   // Selected host
	key   []byte     // AES key (32byte = 256bit)
)

func openSession() (err error) {
	config := &ssh.ClientConfig{
		User:            host.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(host.Password)},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if host.PrivateKeyFile != "" {
		signer, err := setSigner(host.PrivateKeyFile)
		if err != nil {
			panic(err)
		}

		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	}

	hostport := fmt.Sprintf("%s:%d", host.Address, host.Port)
	conn, err := ssh.Dial("tcp", hostport, config)
	if err != nil {
		return fmt.Errorf("ssh.Dial %v: %v", hostport, err)
	}
	defer conn.Close()

	sess, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("conn.NewSession: %v", err)
	}
	defer sess.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
	}

	tty, err := gotty.Open()
	if err != nil {
		return fmt.Errorf("tty.Open: %v", err)
	}
	defer tty.Close()

	// termType := "xterm-256color" // Cursor on vi/vim not work
	// termType := "vt100" // Cursor on vi/vim not work
	// termType := "vt220" // No color on vi/vim
	// termType := "vt320" // No color on shell
	termType := "linux"
	w, h, err := tty.Size() // term.GetSize do not work on windows so, use mattn/go-tty
	if err != nil {
		w, h = 0, 0
	}

	setResizeControl(sess, tty, w, h)

	clean, err := tty.Raw()
	if err != nil {
		log.Fatal(err)
	}
	defer clean()

	err = sess.RequestPty(termType, h, w, modes)
	if err != nil {
		return fmt.Errorf("sess.RequestPty: %s", err)
	}

	pw, err := sess.StdinPipe()
	if err != nil {
		return fmt.Errorf("sess.StdinPipe: %v", err)
	}
	// sess.Stdin = os.Stdin
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	err = sess.Shell()
	if err != nil {
		return fmt.Errorf("sess.Shell: %v", err)
	}

	setEventControl(pw, tty)

	sess.Wait()

	return nil
}

func main() {
	flag.Parse()
	if flag.NArg() > 0 || *hostsFile == "" || *hostsIDX == 0 {
		binaryName := filepath.Base(os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", binaryName)
		flag.PrintDefaults()
		return
	}

	key = []byte("0123456789!#$%^&*()abcdefghijklm")
	err := loadHostData(*hostsFile, key, &hosts)
	if err != nil {
		fmt.Println("error loading host data file")
		return
	}

	if *hostsIDX > len(hosts) {
		fmt.Printf("index not exist. max index is %d\n", len(hosts))
		return
	}

	*hostsIDX--
	host = hosts[*hostsIDX]

	err = openSession()
	if err != nil {
		fmt.Println(err)
	}
}
