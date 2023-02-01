package gsenv

import (
	"github.com/jfy0o0/goStealer/errors/gserror"
	"os"
	"strings"
)

// All returns a copy of strings representing the environment,
// in the form "key=value".
func All() []string {
	return os.Environ()
}

// Map returns a copy of strings representing the environment as a map.
func Map() map[string]string {
	m := make(map[string]string)
	i := 0
	for _, s := range os.Environ() {
		i = strings.IndexByte(s, '=')
		m[s[0:i]] = s[i+1:]
	}
	return m
}

// Get creates and returns a Var with the value of the environment variable
// named by the `key`. It uses the given `def` if the variable does not exist
// in the environment.
func Get(key string, def ...string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		if len(def) > 0 {
			return def[0]
		}
		return ""
	}
	return v
}

// Set sets the value of the environment variable named by the `key`.
// It returns an error, if any.
func Set(key, value string) (err error) {
	err = os.Setenv(key, value)
	if err != nil {
		err = gserror.Wrapf(err, `set environment key-value failed with key "%s", value "%s"`, key, value)
	}
	return
}

// SetMap sets the environment variables using map.
func SetMap(m map[string]string) (err error) {
	for k, v := range m {
		if err = Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

// Contains checks whether the environment variable named `key` exists.
func Contains(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

// Remove deletes one or more environment variables.
func Remove(key ...string) (err error) {
	for _, v := range key {
		if err = os.Unsetenv(v); err != nil {
			err = gserror.Wrapf(err, `delete environment key failed with key "%s"`, v)
			return err
		}
	}
	return nil
}

// Build builds a map to an environment variable slice.
func Build(m map[string]string) []string {
	array := make([]string, len(m))
	index := 0
	for k, v := range m {
		array[index] = k + "=" + v
		index++
	}
	return array
}
