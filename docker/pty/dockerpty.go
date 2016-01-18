package pty

import (
	"errors"
	"io"
	"os"
	gosignal "os/signal"
	"time"

	"github.com/docker/docker/pkg/signal"
	"github.com/docker/docker/pkg/term"
	"github.com/fsouza/go-dockerclient"
)

// Start starts and attach a pty (Pseudo-TTY) to the specified container.
// It will take over the current process' TTY until the container's PTY is closed.
func Start(client *docker.Client, container *docker.Container, hostConfig *docker.HostConfig) (err error) {
	var (
		terminalFd uintptr
		oldState   *term.State
		out        io.Writer = os.Stdout
	)

	if file, ok := out.(*os.File); ok {
		terminalFd = file.Fd()
	} else {
		return errors.New("Not a terminal!")
	}

	// Set up the pseudo terminal
	oldState, err = term.SetRawTerminal(terminalFd)
	if err != nil {
		return
	}

	// Clean up after the container has exited
	defer term.RestoreTerminal(terminalFd, oldState)

	// Attach to the container on a separate thread
	attachChan := make(chan error)
	go attachToContainer(client, container.ID, attachChan)

	// Start it
	err = client.StartContainer(container.ID, hostConfig)
	if err != nil {
		return
	}

	// Make sure terminal resizes are passed on to the container
	monitorTty(client, container.ID, terminalFd)

	return <-attachChan
}

func attachToContainer(client *docker.Client, containerID string, errorChan chan error) {
	r, w := io.Pipe()
	go io.Copy(w, os.Stdin)
	err := client.AttachToContainer(docker.AttachToContainerOptions{
		Container:    containerID,
		InputStream:  r,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Stdin:        true,
		Stdout:       true,
		Stderr:       true,
		Stream:       true,
		RawTerminal:  true,
	})
	errorChan <- err
}

// From https://github.com/docker/docker/blob/0d70706b4b6bf9d5a5daf46dd147ca71270d0ab7/api/client/utils.go#L222-L233
func monitorTty(client *docker.Client, containerID string, terminalFd uintptr) {
	resizeTty(client, containerID, terminalFd)

	sigchan := make(chan os.Signal, 1)
	gosignal.Notify(sigchan, signal.SIGWINCH)
	go func() {
		for range sigchan {
			resizeTty(client, containerID, terminalFd)
		}
	}()
}

// From https://github.com/docker/docker/blob/0d70706b4b6bf9d5a5daf46dd147ca71270d0ab7/api/client/utils.go#L222-L233
func monitorExecTty(client *docker.Client, execID string, terminalFd uintptr) {
	// HACK: For some weird reason on Docker 1.4.1 this resize is being triggered
	//       before the Exec instance is running resulting in an error on the
	//       Docker server. So we wait a little bit before triggering this first
	//       resize
	time.Sleep(50 * time.Millisecond)
	resizeExecTty(client, execID, terminalFd)

	sigchan := make(chan os.Signal, 1)
	gosignal.Notify(sigchan, signal.SIGWINCH)
	go func() {
		for range sigchan {
			resizeExecTty(client, execID, terminalFd)
		}
	}()
}

func resizeTty(client *docker.Client, containerID string, terminalFd uintptr) error {
	height, width := getTtySize(terminalFd)
	if height == 0 && width == 0 {
		return nil
	}
	return client.ResizeContainerTTY(containerID, height, width)
}

func resizeExecTty(client *docker.Client, containerID string, terminalFd uintptr) error {
	height, width := getTtySize(terminalFd)
	if height == 0 && width == 0 {
		return nil
	}
	return client.ResizeExecTTY(containerID, height, width)
}

// From https://github.com/docker/docker/blob/0d70706b4b6bf9d5a5daf46dd147ca71270d0ab7/api/client/utils.go#L235-L247
func getTtySize(terminalFd uintptr) (int, int) {
	ws, err := term.GetWinsize(terminalFd)
	if err != nil {
		if ws == nil {
			return 0, 0
		}
	}
	return int(ws.Height), int(ws.Width)
}
