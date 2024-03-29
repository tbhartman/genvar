package genvar

import (
	"os"
	"strings"
)

// EnvRO encapsulates read only functions from the os package pertaining to
// environment variables.
type EnvRO interface {
	Environ() []string
	ExpandEnv(s string) string
	Getenv(key string) string
	LookupEnv(key string) (string, bool)
}

// Env encapsulates functions from the os package pertaining to environment
// variables.
type Env interface {
	EnvRO
	ClearEnv()
	Setenv(key, value string) error
	Unsetenv(key string) error
}

// Update updates Env with an existing map
func Update(a Env, m map[string]string) (err error) {
	for k, v := range m {
		err = a.Setenv(k, v)
		if err != nil {
			return
		}
	}
	return
}

// UpdateFromEnviron updates Env with string slice (a la Environ)
func UpdateFromEnviron(a Env, m []string) (err error) {
	for _, v := range m {
		sp := strings.SplitN(v, "=", 2)
		k := sp[0]
		if len(sp) > 1 {
			v = sp[1]
		} else {
			v = ""
		}
		err = a.Setenv(k, v)
		if err != nil {
			return
		}
	}
	return
}

type osVars struct {
}

func (o osVars) Environ() []string                   { return os.Environ() }
func (o osVars) Setenv(key, value string) error      { return os.Setenv(key, value) }
func (o osVars) Getenv(key string) string            { return os.Getenv(key) }
func (o osVars) LookupEnv(key string) (string, bool) { return os.LookupEnv(key) }
func (o osVars) ClearEnv()                           { os.Clearenv() }
func (o osVars) Unsetenv(key string) error           { return os.Unsetenv(key) }
func (o osVars) ExpandEnv(s string) string           { return os.ExpandEnv(s) }

// NewOs returns a VarUser that wraps the os package functions
func NewOs() Env {
	return &osVars{}
}

type mapVars struct {
	m map[string]string
}

func (m mapVars) Environ() []string {
	ret := make([]string, 0, len(m.m))
	for k, v := range m.m {
		ret = append(ret, k+"="+v)
	}
	return ret
}
func (m mapVars) Setenv(key, value string) error {
	m.m[key] = value
	return nil
}
func (m mapVars) Getenv(key string) string {
	return m.m[key]
}
func (m mapVars) LookupEnv(key string) (string, bool) {
	value, ok := m.m[key]
	return value, ok
}
func (m mapVars) ClearEnv() {
	for k := range m.m {
		delete(m.m, k)
	}
}
func (m mapVars) Unsetenv(key string) error {
	delete(m.m, key)
	return nil
}
func (m mapVars) ExpandEnv(s string) string {
	return os.Expand(s, m.Getenv)
}

// NewMap returns a VarUser that uses an in-memory map
func NewMap() Env {
	return mapVars{
		m: make(map[string]string),
	}
}
