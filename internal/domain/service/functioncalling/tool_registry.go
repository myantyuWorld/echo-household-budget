package functioncalling

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/infrastructure/llm"
)

type ToolRegistry struct {
	tools map[string]llm.Tool
}

func NewToolRegistry(
	shoppingRepository domainmodel.ShoppingRepository,
	householdRepository domainmodel.HouseHoldRepository,
) *ToolRegistry {
	registry := &ToolRegistry{
		tools: make(map[string]llm.Tool),
	}

	registry.Register(NewExpenseSearchTool(shoppingRepository))
	registry.Register(NewLimitRetrievalTool(householdRepository))
	registry.Register(NewPredictionTool(shoppingRepository, householdRepository))

	return registry
}

func (r *ToolRegistry) Register(tool llm.Tool) {
	r.tools[tool.Name()] = tool
}

func (r *ToolRegistry) GetTool(name string) llm.Tool {
	return r.tools[name]
}

func (r *ToolRegistry) GetAllTools() []llm.Tool {
	tools := make([]llm.Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}