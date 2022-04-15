package actions

import "fmt"

var (
	TemporaryVariableCounter int
)

// generateTempVarName generates tmp name used for temporary variables.
func generateTempVarName(name string) string {
	tempVariableName := fmt.Sprintf("%s_%d", name, TemporaryVariableCounter)
	TemporaryVariableCounter++

	return tempVariableName
}
