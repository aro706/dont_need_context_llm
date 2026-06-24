package models

// Node represents a structural element (Function, Route, Import, Variable)
type Node struct {
	ProjectName  string
	FilePath     string
	FunctionName string
	Type         string
}

// Call represents a dependency relationship between functions
type Call struct {
	ProjectName string
	Caller      string
	Callee      string
}
