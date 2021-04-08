package cmd

type Param struct {
	Name      string
	Shorthand string
	Value     interface{}
	Usage     string
	Required  bool
}
