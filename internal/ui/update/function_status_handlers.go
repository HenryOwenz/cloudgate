package update

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleFunctionStatusOperation handles the function status operation
func HandleFunctionStatusOperation(m *model.Model) (tea.Model, tea.Cmd) {
	newModel := m.Clone()
	newModel.IsLoading = true
	newModel.LoadingMsg = constants.MsgLoadingFunctions

	return WrapModel(newModel), func() tea.Msg {
		// Get the provider
		provider, err := m.Registry.Get(m.ProviderState.ProviderName)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get the FunctionStatusOperation from the provider
		functionOperation, err := provider.GetFunctionStatusOperation()
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get function status using the operation
		ctx := context.Background()
		functions, err := functionOperation.GetFunctionStatus(ctx)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		return model.FunctionStatusMsg{
			Functions: functions,
			Provider:  provider,
		}
	}
}

// HandleLambdaExecuteResult handles the result of a Lambda execution
func HandleLambdaExecuteResult(m *model.Model, result *model.LambdaExecuteResultMsg) *model.Model {
	newModel := m.Clone()
	newModel.IsLoading = false

	if result.Err != nil {
		newModel.Err = result.Err
		return newModel
	}

	if newModel.SelectedFunction == nil {
		newModel.Error = constants.MsgNoFunctionSelected
		return newModel
	}

	// Store the result in the model
	newModel.SetLambdaResult(result.Result)

	// Format the payload for display
	var formattedPayload string
	if result.Result.Payload != "" {
		// Try to pretty-print the JSON
		var jsonObj interface{}
		if err := json.Unmarshal([]byte(result.Result.Payload), &jsonObj); err == nil {
			if prettyJSON, err := json.MarshalIndent(jsonObj, "", "  "); err == nil {
				formattedPayload = string(prettyJSON)
			} else {
				formattedPayload = result.Result.Payload
			}
		} else {
			formattedPayload = result.Result.Payload
		}
	} else {
		formattedPayload = "(empty response)"
	}

	// Create the content for the viewport
	content := fmt.Sprintf("Status Code: %d\nExecuted Version: %s\n\nResponse:\n%s\n\nLogs:\n%s",
		result.Result.StatusCode,
		result.Result.ExecutedVersion,
		formattedPayload,
		result.Result.LogResult)

	// Initialize a new viewport (following the example pattern)
	// We'll set initial dimensions, but these will be updated on WindowSizeMsg
	newModel.Viewport = viewport.New(newModel.Width-constants.ViewportMarginX*2, constants.TableHeight)
	newModel.Viewport.YPosition = constants.HeaderHeight // Position below the title
	newModel.Viewport.SetContent(content)
	newModel.Viewport.GotoTop()

	// Mark that the viewport needs proper initialization with window dimensions
	newModel.ViewportReady = false

	// Set success message
	newModel.Success = fmt.Sprintf(constants.MsgLambdaExecuteSuccess, newModel.SelectedFunction.Name)

	// Transition to the response view
	newModel.CurrentView = constants.ViewLambdaResponse

	return newModel
}
