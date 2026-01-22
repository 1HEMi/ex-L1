package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	for {
		fmt.Print("minishell> ")

		line, err := reader.ReadString('\n')
		if err != nil {

			fmt.Println()
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, "|") {
			runPipeline(line)
			continue
		}

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		if runBuiltin(args) {
			continue
		}

		runCommand(args)

		select {
		case <-sigCh:
			fmt.Println()
		default:
		}
	}
}

func runBuiltin(args []string) bool {
	switch args[0] {

	case "cd":
		if len(args) < 2 {
			fmt.Println("cd: missing operand")
			return true
		}
		if err := os.Chdir(args[1]); err != nil {
			fmt.Println("cd:", err)
		}
		return true

	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(dir)
		}
		return true

	case "echo":
		fmt.Println(strings.Join(args[1:], " "))
		return true

	case "kill":
		if len(args) < 2 {
			fmt.Println("kill: missing pid")
			return true
		}
		pid, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("kill: invalid pid")
			return true
		}
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			fmt.Println("kill:", err)
		}
		return true

	case "ps":
		cmd := exec.Command("ps")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return true
	}

	return false
}

func runCommand(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("error:", err)
	}
}

func runPipeline(line string) {
	parts := strings.Split(line, "|")

	var cmds []*exec.Cmd

	for _, part := range parts {
		args := strings.Fields(strings.TrimSpace(part))
		if len(args) == 0 {
			return
		}
		cmds = append(cmds, exec.Command(args[0], args[1:]...))
	}

	for i := 0; i < len(cmds)-1; i++ {
		pipe, err := cmds[i].StdoutPipe()
		if err != nil {
			fmt.Println("pipe error:", err)
			return
		}
		cmds[i+1].Stdin = pipe
	}

	cmds[0].Stdin = os.Stdin
	cmds[len(cmds)-1].Stdout = os.Stdout
	cmds[len(cmds)-1].Stderr = os.Stderr

	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			fmt.Println("start error:", err)
			return
		}
	}

	for _, cmd := range cmds {
		cmd.Wait()
	}
}
