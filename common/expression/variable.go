package expression

type variable struct {
	name string
}

func (v variable) Eval(vars map[string]int64) int64 {
	if val, exists := vars[v.name]; exists {
		return val
	}
	return 0
}
