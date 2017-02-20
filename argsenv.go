package kvconfig

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func normalizeArgumentName(arg string) string {
	name := strings.Replace(strings.TrimPrefix(strings.ToLower(arg), "-"), "-", "_", -1)
	pos := strings.LastIndex(name, "_")
	if pos == -1 {
		return fmt.Sprintf("%s_0", name)
	} else if len(name) >= pos+2 {
		trialingNum := name[strings.LastIndex(name, "_")+1:]
		if _, err := strconv.Atoi(trialingNum); err != nil {
			return fmt.Sprintf("%s_0", name)
		}
	}
	return name
}

// Parse command-line arguments into the key/value store.
// Note that argument names may be transformed.
func ParseArgs(kv Setter) error {
	for i := 1; i < len(os.Args); i++ {
		if !strings.HasPrefix(os.Args[i], "-") {
			return errors.New(fmt.Sprintf("invalid argument: \"%s\"", os.Args[i]))
		}
		if strings.Index(os.Args[i], "=") == -1 {
			if len(os.Args) <= i+1 {
				return errors.New(fmt.Sprintf("missing value to argument: %s", os.Args[i]))
			}
			if strings.HasPrefix(os.Args[i+1], "-") {
				return errors.New(fmt.Sprintf("value following argument cannot start with \"-\": %s", os.Args[i+1]))
			}
			kv.Set(normalizeArgumentName(os.Args[i]), os.Args[i+1])
			i++
		} else {
			split := strings.SplitN(os.Args[i], "=", 2)
			kv.Set(normalizeArgumentName(split[0]), split[1])
		}
	}

	return nil
}

// Parse environment variables starting with "CFG_" into the key/value store.
// Note that environment variable names may be transformed.
func ParseEnv(kv Setter) {
	for _, arg := range os.Environ() {
		if !strings.HasPrefix(arg, "CFG_") {
			continue
		}
		split := strings.SplitN(arg[4:], "=", 2)
		kv.Set(normalizeArgumentName(split[0]), split[1])
	}
}
