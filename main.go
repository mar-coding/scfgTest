package main

import (
	"os"
	"os/exec"
)

func main() {
	switch os.Args[1] {

	case "run":
		run()

	default:
		panic("help")
	}
}

func run() {
	// fmt.Printf("Running %v\n", os.Args[2:])
	err := execution(os.Args...)
	if err != nil {
		panic(err)
	}
}

func execution(args ...string) error {
	cmd := exec.Command(args[2], args[3:]...)

	// adjust the os standard input,output,error to my cmd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err
}
