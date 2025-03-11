package update

import (
	"context"
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleFunctionStatus handles the function status operation
func HandleFunctionStatus(m *model.Model) (tea.Model, tea.Cmd) {
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

// HandleFunctionSelection handles the selection of a function
func HandleFunctionSelection(m *model.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		// Clone the model to avoid modifying the original
		newModel := m.Clone()

		functionName := selected[0]

		// Find the selected function
		var selectedFunction *cloud.FunctionStatus
		for _, function := range m.Functions {
			if function.Name == functionName {
				selectedFunction = &function
				break
			}
		}

		if selectedFunction == nil {
			return WrapModel(m), func() tea.Msg {
				return model.ErrMsg{Err: fmt.Errorf(constants.MsgErrorNoFunction)}
			}
		}

		// Update the model
		newModel.SetSelectedFunction(selectedFunction)

		// If we're in the Lambda execution flow and a function is already selected,
		// go directly to the Lambda execution view
		if newModel.IsExecuteLambdaFlow && newModel.SelectedFunction != nil {
			return HandleLambdaExecuteSelection(newModel)
		}

		newModel.CurrentView = constants.ViewFunctionDetails
		view.UpdateTableForView(newModel)
		return WrapModel(newModel), nil
	}
	return WrapModel(m), nil
}

// HandleLambdaExecuteSelection handles the selection to execute a Lambda function
func HandleLambdaExecuteSelection(m *model.Model) (tea.Model, tea.Cmd) {
	if m.SelectedFunction == nil {
		return WrapModel(m), func() tea.Msg {
			return model.ErrMsg{Err: fmt.Errorf(constants.MsgErrorNoFunction)}
		}
	}

	// Clone the model to avoid modifying the original
	newModel := m.Clone()

	// Initialize the TextArea for JSON input
	ta := textarea.New()
	ta.Placeholder = constants.MsgEnterLambdaPayload
	ta.Focus()
	ta.SetWidth(m.Width - 4)
	ta.SetHeight(10)
	ta.ShowLineNumbers = true
	ta.CharLimit = 0

	// Initialize the Viewport for output display
	vp := viewport.New(m.Width-4, m.Height-20)
	vp.SetContent("")

	// Update the model
	newModel.TextArea = ta
	newModel.Viewport = vp
	newModel.CurrentView = constants.ViewLambdaExecute
	newModel.SetLambdaPayload("{}")
	newModel.SetLambdaResult(nil)
	newModel.IsLambdaInputMode = false // Start in command mode (not input mode)

	return WrapModel(newModel), nil
}

// HandleLambdaExecute handles the execution of a Lambda function
func HandleLambdaExecute(m *model.Model) (tea.Model, tea.Cmd) {
	if m.SelectedFunction == nil {
		return WrapModel(m), func() tea.Msg {
			return model.ErrMsg{Err: fmt.Errorf(constants.MsgErrorNoFunction)}
		}
	}

	// Get the payload from the TextArea
	payload := m.TextArea.Value()
	if payload == "" {
		payload = "{}"
	}

	// Set the payload in the model
	newModel := m.Clone()
	newModel.SetLambdaPayload(payload)
	newModel.IsLoading = true
	newModel.LoadingMsg = constants.MsgExecutingLambda

	return WrapModel(newModel), func() tea.Msg {
		// Get the provider
		provider, err := m.Registry.Get(m.ProviderState.ProviderName)
		if err != nil {
			return model.LambdaExecuteResultMsg{Err: err}
		}

		// Get the LambdaExecuteOperation from the provider
		lambdaOperation, err := provider.GetLambdaExecuteOperation()
		if err != nil {
			return model.LambdaExecuteResultMsg{Err: err}
		}

		// Execute the Lambda function
		ctx := context.Background()
		result, err := lambdaOperation.ExecuteFunction(ctx, m.SelectedFunction.Name, payload)
		if err != nil {
			return model.LambdaExecuteResultMsg{Err: err}
		}

		return model.LambdaExecuteResultMsg{
			Result: result,
		}
	}
}
