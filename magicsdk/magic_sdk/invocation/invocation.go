package invocation

type Invocation struct {
	ParentJSON []byte
	ParentName string
	FnName     string
	InputArgs  map[string][]byte
}