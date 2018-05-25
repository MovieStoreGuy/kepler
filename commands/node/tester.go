// Package node is used to be able to interface with node projects inside a
// git repo with submodules to automate potentially repeative or mundane tasks.
package node

import (
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"

	blast "github.com/LGUG2Z/blastradius"
)

// TestedProject contains all the data required to
// send back to main process and report weather or not it failed
type TestedProject struct {
	Name     string
	ExitCode int
	Output   []byte
}

// RunTestsOn will test the given project
// and all projects that use the given project
func RunTestsOn(project string, command ...string) (chan TestedProject, error) {
	if _, err := os.Stat(project); os.IsNotExist(err) {
		return nil, err
	}
	projects, err := blast.Calculate(".", project)
	if err != nil {
		return nil, err
	}
	results := make(chan TestedProject, runtime.NumCPU())
	wg := sync.WaitGroup{}
	// detaching the dispatcher thread from the main as
	// as to avoid blocking the main thread
	go func(ch chan TestedProject, projects []string) {
		wg.Add(len(projects))
		for _, p := range projects {
			go func(p string, wg *sync.WaitGroup, ch chan TestedProject) {
				ch <- executeTests(p, "npm", "test")
				wg.Done()
			}(p, &wg, ch)
		}
		wg.Wait()
		close(ch)
	}(results, append(projects, project))
	return results, nil
}

func executeTests(project string, cmd ...string) TestedProject {
	if len(cmd) < 2 {
		return TestedProject{project, 1, []byte("Not enough args pased to execute")}
	}
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Dir = project
	buff, err := c.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exiter, ok := err.(*exec.ExitError); ok {
			if status, ok := exiter.Sys().(syscall.WaitStatus); ok {
				exitCode = int(status.ExitStatus())
			}
		}
	}
	return TestedProject{
		Name:     project,
		ExitCode: exitCode,
		Output:   buff,
	}
}
