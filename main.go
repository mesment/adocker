// +build linux

package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
	"github.com/docker/docker-ce/components/engine/pkg/reexec"
)

func init()  {
	reexec.Register("nsInitialisation", nsInitialisation)
	if reexec.Init() {
		os.Exit(0)
	}
}

func nsInitialisation() {
	log.Printf("\n>> namespace setup code goes here <<\n\n")
	nsRun()
}
func nsRun() {
	cmd := exec.Command("/bin/sh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"PS1=-[adocker-process]- # "}

	if err := cmd.Run(); err != nil {
		log.Printf("Error running the /bin/sh command - %s\n", err)
		os.Exit(1)
	}
}

func main()  {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:syscall.CLONE_NEWUTS|
			syscall.CLONE_NEWIPC|
			syscall.CLONE_NEWPID|
			syscall.CLONE_NEWNS|
			syscall.CLONE_NEWNET|
			syscall.CLONE_NEWUSER,
			UidMappings:[]syscall.SysProcIDMap{
				{
					ContainerID: 0,
					HostID:      os.Getuid(),
					Size:        1,
				},
			},
			GidMappings: []syscall.SysProcIDMap{
					{
						ContainerID: 0,
						HostID:      os.Getgid(),
						Size:        1,
					},
			},
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error running the /bin/sh command - %s\n", err)
		panic(err)
	}
}
