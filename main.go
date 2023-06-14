package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {

	case "run":
		run("parent")
	case "child":
		run("child")

	default:
		panic("help")
	}
}

func run(name string) {
	// fmt.Printf("Running %v\n", os.Args[2:])
	if name == "parent" {
		err := parent(os.Args...)
		if err != nil {
			panic(err)
		}
	} else if name == "child" {
		err := child(os.Args...)
		if err != nil {
			panic(err)
		}
	}
}

func parent(args ...string) error {
	// it is use the process that is running now
	// (the process that our golang code is running on.)
	// and it is run the child process.
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, args[2:]...)...)

	// adjust the os standard input,output,error to my cmd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		// change the uts, so our hostname in container
		// will be different from the host
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	// set the hostname to show in the cmd

	err := cmd.Run()

	return err
}

func child(args ...string) error {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getegid())

	syscall.Sethostname([]byte("container"))

	cmd := exec.Command(args[2], args[3:]...)

	// adjust the os standard input,output,error to my cmd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		// change the uts, so our hostname in container
		// will be different from the host
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	// set the hostname to show in the cmd

	err := cmd.Run()

	return err
}
