package kvconfig

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// MapStrStr is a very simple implementation satisfying the Getter and Setter interfaces
type MapStrStr map[string]string

func NewMap() *MapStrStr {
	newMap := make(MapStrStr)
	return &newMap
}

func (m *MapStrStr) Set(k, v string) {
	(*m)[k] = v
}

func (m *MapStrStr) Get(k string) string {
	return (*m)[k]
}

func (m *MapStrStr) Exists(k string) bool {
	_, ok := (*m)[k]
	return ok
}

func (m *MapStrStr) WriteEnvFile(filename string) error {
	f, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer f.Close()

	for k, v := range *m {
		_, err := f.WriteString(fmt.Sprintf("CFG_%s=%s\n", strings.ToUpper(k), v))
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MapStrStr) ReadEnvFile(filename string) error {
	f, err := os.Open(filename)

	// it's okay if our file doesn't exist, we can treat that as no/zero config
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		eqPos := strings.Index(line, "=")

		// primitive blank line, comment, and non-variable exclusion
		if len(line) < 1 || line[0] == '#' || eqPos == -1 {
			continue
		}

		// keys must start with "CFG_" (like envvars)
		if len(line) < 5 || line[0:4] != "CFG_" {
			continue
		}

		m.Set(strings.ToLower(line[4:eqPos]), line[eqPos+1:])
	}

	return scanner.Err()
}
