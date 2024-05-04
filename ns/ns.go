package ns

import (
	"os"
)

func getenv(name, dflt string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		return dflt
	} else {
		return value
	}
}

func Setup() {
	ns := os.Getenv("NAMESPACE")
	if len(ns) == 0 {
		ns = "/tmp/ns." + getenv("USER", "") + "." + getenv("DISPLAY", ":0")
		os.Setenv("NAMESPACE", ns)
	}
	err := os.MkdirAll(ns, 0750)
	if err != nil {
		os.Stderr.WriteString("error creating namespace: " + err.Error() + "\n")
	}
}
