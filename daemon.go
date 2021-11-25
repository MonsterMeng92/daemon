package daemon

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var (
	h bool
	d bool
	s bool
)

func init() {
	flag.BoolVar(&h, "h", false, "usage")
	flag.BoolVar(&d, "d", false, "run apiserver as a daemon with -d=true")
	flag.BoolVar(&s, "s", false, "shutdown apiserver")
	flag.Usage = usage
	flag.Parse()

	if d {
		cmd := exec.Command(os.Args[0], flag.Args()...)
		if err := cmd.Start(); err != nil {
			log.Printf("start %s failed with error: %s", os.Args[0], err.Error())
			os.Exit(-1)
		}
		f, _ := os.OpenFile("apiserver.lock", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		fmt.Fprintf(f, "%d", cmd.Process.Pid)
		log.Printf("%s [PID] %d running...\n", os.Args[0], cmd.Process.Pid)
		f.Close()
		os.Exit(0)
	}
	if h {
		flag.Usage()
	}
	if s {
		data, err := ioutil.ReadFile("apiserver.lock")
		if err != nil {
			log.Fatal(err.Error())
		}
		cmd := exec.Command("kill", "-9", string(data))
		if err := cmd.Start(); err != nil {
			log.Printf("shutdown apiserver error: %s", err.Error())
			os.Exit(-1)
		}
		log.Println("apiserver is down")
		os.Exit(0)
	}
}

func usage() {
	fmt.Fprintf(os.Stdout, `apiserver usage:
Options:
`)
	flag.PrintDefaults()
	os.Exit(0)
}
