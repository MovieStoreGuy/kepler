// Package node is used to be able to interface with node projects inside a
// git repo with submodules to automate potentially repeative or mundane tasks.
package node

import blast "github.com/LGUG2Z/blastradius"

// TestedProject contains all the data required to
// send back to main process and report weather or not it failed
type TestedProject struct {
	Name     string
	ExitCode int
}

// RunTestsOn will test the given project
// and all projects that use the given project
func RunTestsOn(project string, command ...string) (chan TestedProject, error) {
	projects, err := blast.Calculate(".", project)
	if err != nil {
		return nil, err
	}
	results := make(chan TestedProject)
	_ = projects
	return results, nil
}
