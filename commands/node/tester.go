// Package node is used to be able to interface with node projects inside a
// git repo with submodules to automate potentially repeative or mundane tasks.
package node

import (
	"os/exec"
	"runtime"
	"sync"

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
	projects, err := blast.Calculate(".", project)
	if err != nil {
		return nil, err
	}
	results := make(chan TestedProject, runtime.NumCPU())
	wg := sync.WaitGroup{}
	// detaching the dispatcher thread from the main as
	// as to avoid blocking the main thread
	// THIS CODE IS GROSS AND I REALLY DON'T LIKE IT
	go func(ch chan TestedProject, projects []string) {
		wg.Add(len(projects) + 1)
		for _, p := range projects {
			go func(p string, wg *sync.WaitGroup, ch chan TestedProject) {
				c := exec.Command(command[0], command[1:]...)
				c.Dir = p
				exitCode := 0
				if err := c.Wait(); err != nil {
					// Figure out why the thing failed
				}
				buff, err := c.CombinedOutput()
				if err != nil {
					// Complain on some degree
				}
				ch <- TestedProject{
					Name:     p,
					ExitCode: exitCode,
					Output:   buff,
				}
				wg.Done()
			}(p, &wg, ch)
		}
		wg.Wait()
		close(ch)
	}(results, append(projects, project))
	return results, nil
}

func executeTests(project string, cmd ...string) (TestedProject, error) {
	return TestedProject{}, nil
}
