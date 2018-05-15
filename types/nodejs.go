package types

import (
	"io/ioutil"

	"github.com/MovieStoreGuy/resources/marshal"
)

//PackageJSON structure of package.json
type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	Main            string            `json:"main"`
	Bugs            map[string]string `json:"bugs,omitempty"`
	Scripts         map[string]string `json:"scripts,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`
	Private         bool              `json:"private,omitempty"`
	License         string            `json:"license,omitempty"`
}

// WriteTo will write the current contents of the PackageJSON
// into the given directory
func (pack *PackageJSON) WriteTo(path string) error {
	o, err := marshal.PureMarshalIndent(pack, "", "    ")
	if err != nil {
		return err
	}
	o = append(o, []byte("\n")...)
	return ioutil.WriteFile(path, o, 0644)
}
