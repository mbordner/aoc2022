package expression

import "github.com/pkg/errors"

type variable struct {
	name string
}

func (v variable) Eval(vars map[string]int64) int64 {
	if val, exists := vars[v.name]; exists {
		return val
	}
	return 0
}

func (v variable) EvalKnown(vars map[string]int64) (int64, error) {
	if val, exists := vars[v.name]; exists {
		return val, nil
	}
	return 0, errors.Errorf("unknown var %s", v.name)
}

func (v variable) String() string {
	return v.name
}