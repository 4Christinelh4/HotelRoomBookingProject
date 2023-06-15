package forms

type errors map[string][]string

func (e errors) Add(field string, message string) {
	e[field] = append(e[field], message)
}

// Get returns error message
func (e errors) Get(field string) string {
	estrs := e[field]
	if len(estrs) == 0 {
		return ""
	}
	return estrs[0]
}
