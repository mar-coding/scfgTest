package main

import (
	"log"
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
	log.Printf("Running %v \n", args[2:])
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

		// Cloneflags:
		// syscall.CLONE_NEWUTS flag creates a new UTS (Unix Timesharing System) namespace for the child process.
		// This namespace provides isolation for the hostname and domain name,
		// allowing the container's hostname to be different from the host system.

		// syscall.CLONE_NEWPID flag creates a new PID (Process ID) namespace for the child process.
		// This namespace provides process isolation, meaning the child process will have a separate PID namespace,
		// making it appear as if it is running in its own isolated environment.

		// syscall.CLONE_NEWNS flag creates a new mount namespace for the child process.
		// This namespace provides isolation for the mount points,
		// allowing the child process to have its own separate filesystem view.

		// Unshareflags
		// syscall.CLONE_NEWNS flag indicates that the child process should have its own mount namespace.
		// This is used in conjunction with the Cloneflags mentioned above to
		//  ensure that the child process has a separate mount namespace.

		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	// set the hostname to show in the cmd

	err := cmd.Run()

	return err
}

func child(args ...string) error {
	log.Printf("Running %v as %d\n", args[2:], os.Getegid())

	must(syscall.Sethostname([]byte("container")))

	// '/tmp/alpine-rootfs/' is where the container's root filesystem is mounted.
	// where the mini root alpine was extracted.
	must(syscall.Chroot("/tmp/alpine-rootfs/"))
	must(syscall.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	cmd := exec.Command(args[2], args[3:]...)

	// adjust the os standard input,output,error to my cmd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// set the hostname to show in the cmd

	err := cmd.Run()

	must(syscall.Unmount("proc", 0))

	return err
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
