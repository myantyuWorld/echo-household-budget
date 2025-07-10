package llm

type Tool interface {
	Name() string
	Description() string
	Execute(params map[string]interface{}) (interface{}, error)
}