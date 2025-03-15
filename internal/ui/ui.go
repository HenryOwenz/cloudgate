package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/update"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
)

// Model is the main UI model that implements the tea.Model interface
type Model struct {
	core *model.Model
}

// New creates a new UI model
func New() Model {
	m := Model{
		core: model.New(),
	}
	// Initialize the table for the current view
	view.UpdateTableForView(m.core)
	return m
}

// Init initializes the UI model
func (m Model) Init() tea.Cmd {
	// Make sure to initialize the table before returning
	view.UpdateTableForView(m.core)
	return m.core.Init()
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		newModel := m.Clone()
		newModel.core.Width = msg.Width
		newModel.core.Height = msg.Height

		// Update viewport dimensions if we're in the Lambda response view
		if newModel.core.CurrentView == constants.ViewLambdaResponse {
			// Create temporary header and footer to calculate their heights
			title := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(constants.ColorTitle)).
				Render(constants.TitleLambdaResponse)
			line := strings.Repeat("─", max(0, msg.Width-lipgloss.Width(title)))
			header := lipgloss.JoinHorizontal(lipgloss.Center, title, line)
			headerHeight := lipgloss.Height(header)

			footerText := fmt.Sprintf("%3.f%%", 0.0)
			footer := lipgloss.NewStyle().
				Foreground(lipgloss.Color(constants.ColorPrimary)).
				Render(footerText)
			footerLine := strings.Repeat("─", max(0, msg.Width-lipgloss.Width(footerText)))
			footer = lipgloss.JoinHorizontal(lipgloss.Center, footerLine, footer)
			footerHeight := lipgloss.Height(footer)

			verticalMarginHeight := headerHeight + footerHeight

			if !newModel.core.ViewportReady {
				// First time initialization with proper window dimensions
				newModel.core.Viewport.Width = msg.Width - constants.ViewportMarginX*2
				newModel.core.Viewport.Height = msg.Height - verticalMarginHeight
				newModel.core.Viewport.YPosition = headerHeight
				newModel.core.ViewportReady = true
			} else {
				// Just update dimensions for subsequent window size changes
				newModel.core.Viewport.Width = msg.Width - constants.ViewportMarginX*2
				newModel.core.Viewport.Height = msg.Height - verticalMarginHeight
			}
		}

		view.UpdateTableForView(newModel.core)
		return newModel, nil
	case model.ErrMsg:
		newModel := m.Clone()
		newModel.core.Err = msg.Err
		newModel.core.IsLoading = false
		return newModel, nil
	case model.ApprovalsMsg:
		newModel := m.Clone()
		newModel.core.Approvals = msg.Approvals
		newModel.core.Provider = msg.Provider
		newModel.core.CurrentView = constants.ViewApprovals
		newModel.core.IsLoading = false

		// Only show the first page of approvals based on page size
		pageSize := newModel.core.PageSize
		initialApprovals := newModel.core.Approvals

		// If we have more approvals than the page size, only show the first page
		if len(initialApprovals) > pageSize {
			initialApprovals = initialApprovals[:pageSize]
		}

		// Convert to an ApprovalsPageMsg for consistent pagination handling
		approvalsPageMsg := model.ApprovalsPageMsg{
			Approvals:     initialApprovals,
			NextPageToken: "2", // Set to page 2 for next page
			HasMorePages:  len(newModel.core.Approvals) > pageSize,
		}

		// Use the pagination handler
		newModel.core = update.HandleApprovalsPagination(newModel.core, approvalsPageMsg)

		// Store all approvals in the pagination state for client-side pagination
		if len(newModel.core.Approvals) > 0 {
			for _, approval := range msg.Approvals {
				// We need to add all approvals to AllItems, not just the ones displayed
				found := false
				for _, item := range newModel.core.Pagination.AllItems {
					if a, ok := item.(model.ApprovalAction); ok &&
						a.PipelineName == approval.PipelineName &&
						a.StageName == approval.StageName &&
						a.ActionName == approval.ActionName {
						found = true
						break
					}
				}
				if !found {
					newModel.core.Pagination.AllItems = append(newModel.core.Pagination.AllItems, approval)
				}
			}
			// Update total items count
			newModel.core.Pagination.TotalItems = int64(len(newModel.core.Pagination.AllItems))
		}

		return newModel, nil
	case model.ApprovalResultMsg:
		newModel := m.Clone()
		newModel.core.IsLoading = false // Ensure loading is turned off
		update.HandleApprovalResult(newModel.core, msg.Err)
		view.UpdateTableForView(newModel.core)
		return newModel, nil
	case model.PipelineExecutionMsg:
		newModel := m.Clone()
		newModel.core.IsLoading = false // Ensure loading is turned off
		update.HandlePipelineExecution(newModel.core, msg.Err)
		view.UpdateTableForView(newModel.core)
		return newModel, nil
	case model.FunctionStatusMsg:
		newModel := m.Clone()
		newModel.core.Functions = msg.Functions
		newModel.core.Provider = msg.Provider
		newModel.core.CurrentView = constants.ViewFunctionStatus
		newModel.core.IsLoading = false

		// Sort functions by name in ascending order (case-insensitive)
		// This preserves the original case of function names in the display
		// while providing a consistent sorting order regardless of casing.
		// The lowercase conversion is used only for comparison during sorting.
		sort.Slice(newModel.core.Functions, func(i, j int) bool {
			return strings.ToLower(newModel.core.Functions[i].Name) < strings.ToLower(newModel.core.Functions[j].Name)
		})

		// Convert to a FunctionsPageMsg for consistent pagination handling
		// Only show the first page of functions based on page size
		pageSize := newModel.core.PageSize
		initialFunctions := newModel.core.Functions

		// If we have more functions than the page size, only show the first page
		if len(initialFunctions) > pageSize {
			initialFunctions = initialFunctions[:pageSize]
		}

		functionsPageMsg := model.FunctionsPageMsg{
			Functions:     initialFunctions,
			NextPageToken: "2", // Set to page 2 for next page
			HasMorePages:  len(newModel.core.Functions) > pageSize,
		}

		// Use the pagination handler
		newModel.core = update.HandleFunctionStatusPagination(newModel.core, functionsPageMsg)

		// Store all functions in the pagination state for client-side pagination
		if len(newModel.core.Functions) > 0 {
			for _, function := range msg.Functions {
				// We need to add all functions to AllItems, not just the ones displayed
				found := false
				for _, item := range newModel.core.Pagination.AllItems {
					if f, ok := item.(model.FunctionStatus); ok && f.Name == function.Name {
						found = true
						break
					}
				}
				if !found {
					newModel.core.Pagination.AllItems = append(newModel.core.Pagination.AllItems, function)
				}
			}
			// Update total items count
			newModel.core.Pagination.TotalItems = int64(len(newModel.core.Pagination.AllItems))
		}

		return newModel, nil
	case model.LambdaExecuteResultMsg:
		newModel := m.Clone()
		newModel.core = update.HandleLambdaExecuteResult(newModel.core, &msg)
		return newModel, nil
	case spinner.TickMsg:
		newModel := m.Clone()
		var cmd tea.Cmd
		newModel.core.Spinner, cmd = newModel.core.Spinner.Update(msg)
		if newModel.core.IsLoading {
			cmds = append(cmds, cmd)
		}
		return newModel, tea.Batch(cmds...)
	case tea.KeyMsg:
		// Ignore navigation key presses when loading
		if m.core.IsLoading {
			// Only allow quit commands during loading
			switch msg.String() {
			case constants.KeyCtrlC, constants.KeyQ:
				return m, tea.Quit
			default:
				// Ignore all other key presses during loading
				return m, nil
			}
		}

		// Special handling for Lambda response view
		if m.core.CurrentView == constants.ViewLambdaResponse {
			// Handle quit and back navigation
			switch msg.String() {
			case constants.KeyQ, constants.KeyCtrlC:
				return m, tea.Quit
			case constants.KeyEsc, constants.KeyAltBack:
				// Navigate back to the Lambda execution view
				newModel := m.Clone()
				newModel.core.CurrentView = constants.ViewLambdaExecute
				return newModel, nil
			default:
				// Pass ALL other keys to the viewport
				newModel := m.Clone()
				var cmd tea.Cmd
				newModel.core.Viewport, cmd = newModel.core.Viewport.Update(msg)
				return newModel, cmd
			}
		}

		// Special handling for Lambda execution view
		if m.core.CurrentView == constants.ViewLambdaExecute {
			// Handle special keys
			switch msg.String() {
			case constants.KeyQ:
				return m, tea.Quit
			case constants.KeyCtrlC:
				// If in input mode, exit to command mode
				if m.core.IsLambdaInputMode {
					newModel := m.Clone()
					newModel.core.IsLambdaInputMode = false
					return newModel, nil
				}
				// Otherwise, quit the application
				return m, tea.Quit
			case constants.KeyEsc, constants.KeyAltBack:
				// If in input mode, exit to command mode
				if m.core.IsLambdaInputMode {
					newModel := m.Clone()
					newModel.core.IsLambdaInputMode = false
					return newModel, nil
				}
				// Otherwise, navigate back
				newCore := update.NavigateBack(m.core)
				view.UpdateTableForView(newCore)
				return Model{core: newCore}, nil
			case constants.KeyShiftEnter, constants.KeyCtrlEnter, constants.KeyF5:
				// Execute the Lambda function regardless of mode
				modelWrapper, cmd := update.HandleLambdaExecute(m.core)
				if wrapper, ok := modelWrapper.(update.ModelWrapper); ok {
					newModel := Model{core: wrapper.Model}
					if newModel.core.IsLoading {
						return newModel, tea.Batch(cmd, newModel.core.Spinner.Tick)
					}
					return newModel, cmd
				}
				return modelWrapper, cmd
			case constants.KeyTab:
				// Tab key is now used for text input only
				if m.core.IsLambdaInputMode {
					// Pass the tab key to the text area for indentation
					var cmd tea.Cmd
					newModel := m.Clone()
					newModel.core.TextArea, cmd = newModel.core.TextArea.Update(msg)
					return newModel, cmd
				}
				return m, nil
			case "i":
				// Enter input mode (only if not already in input mode)
				if !m.core.IsLambdaInputMode {
					newModel := m.Clone()
					newModel.core.IsLambdaInputMode = true
					return newModel, nil
				}
				// If already in input mode, pass the 'i' key to the TextArea
				if m.core.IsLambdaInputMode {
					var cmd tea.Cmd
					newModel := m.Clone()
					newModel.core.TextArea, cmd = newModel.core.TextArea.Update(msg)
					newModel.core.LambdaPayload = newModel.core.TextArea.Value()
					return newModel, cmd
				}
				return m, nil
			case constants.KeyEnter:
				// In input mode, Enter adds new lines
				// In command mode, Enter executes the Lambda function
				if !m.core.IsLambdaInputMode {
					// Command mode - run the Lambda function
					modelWrapper, cmd := update.HandleLambdaExecute(m.core)
					if wrapper, ok := modelWrapper.(update.ModelWrapper); ok {
						newModel := Model{core: wrapper.Model}
						if newModel.core.IsLoading {
							return newModel, tea.Batch(cmd, newModel.core.Spinner.Tick)
						}
						return newModel, cmd
					}
					return modelWrapper, cmd
				}
				// If in input mode, let the default handler process the Enter key
				// (which will add a new line in the TextArea)
				if m.core.IsLambdaInputMode {
					var cmd tea.Cmd
					newModel := m.Clone()
					newModel.core.TextArea, cmd = newModel.core.TextArea.Update(msg)
					newModel.core.LambdaPayload = newModel.core.TextArea.Value()
					return newModel, cmd
				}
				return m, nil
			// Add viewport navigation keys
			case constants.KeyUp, constants.KeyAltUp:
				// If in input mode, pass 'k' to the text input
				if m.core.IsLambdaInputMode && msg.String() == constants.KeyAltUp {
					var cmd tea.Cmd
					newModel := m.Clone()
					newModel.core.TextArea, cmd = newModel.core.TextArea.Update(msg)
					newModel.core.LambdaPayload = newModel.core.TextArea.Value()
					return newModel, cmd
				}
				// Scroll viewport up
				newModel := m.Clone()
				newModel.core.Viewport.LineUp(1)
				return newModel, nil
			case constants.KeyDown, constants.KeyAltDown:
				// If in input mode, pass 'j' to the text input
				if m.core.IsLambdaInputMode && msg.String() == constants.KeyAltDown {
					var cmd tea.Cmd
					newModel := m.Clone()
					newModel.core.TextArea, cmd = newModel.core.TextArea.Update(msg)
					newModel.core.LambdaPayload = newModel.core.TextArea.Value()
					return newModel, cmd
				}
				// Scroll viewport down
				newModel := m.Clone()
				newModel.core.Viewport.LineDown(1)
				return newModel, nil
			case constants.KeyPageUp, constants.KeyAltPageUp:
				// If in input mode, pass 'b' to the text input
				if m.core.IsLambdaInputMode && msg.String() == constants.KeyAltPageUp {
					var cmd tea.Cmd
					newModel := m.Clone()
					newModel.core.TextArea, cmd = newModel.core.TextArea.Update(msg)
					newModel.core.LambdaPayload = newModel.core.TextArea.Value()
					return newModel, cmd
				}
				// Page up in viewport
				newModel := m.Clone()
				newModel.core.Viewport.HalfViewUp()
				return newModel, nil
			case constants.KeyPageDown, constants.KeyAltPageDown:
				// If in input mode, pass 'f' to the text input
				if m.core.IsLambdaInputMode && msg.String() == constants.KeyAltPageDown {
					var cmd tea.Cmd
					newModel := m.Clone()
					newModel.core.TextArea, cmd = newModel.core.TextArea.Update(msg)
					newModel.core.LambdaPayload = newModel.core.TextArea.Value()
					return newModel, cmd
				}
				// Page down in viewport
				newModel := m.Clone()
				newModel.core.Viewport.HalfViewDown()
				return newModel, nil
			case constants.KeyHome, constants.KeyGotoTop:
				// If in input mode, pass 'g' to the text input
				if m.core.IsLambdaInputMode && msg.String() == constants.KeyGotoTop {
					var cmd tea.Cmd
					newModel := m.Clone()
					newModel.core.TextArea, cmd = newModel.core.TextArea.Update(msg)
					newModel.core.LambdaPayload = newModel.core.TextArea.Value()
					return newModel, cmd
				}
				// Go to top of viewport
				newModel := m.Clone()
				newModel.core.Viewport.GotoTop()
				return newModel, nil
			case constants.KeyEnd, constants.KeyGotoBottom:
				// If in input mode, pass 'G' to the text input
				if m.core.IsLambdaInputMode && msg.String() == constants.KeyGotoBottom {
					var cmd tea.Cmd
					newModel := m.Clone()
					newModel.core.TextArea, cmd = newModel.core.TextArea.Update(msg)
					newModel.core.LambdaPayload = newModel.core.TextArea.Value()
					return newModel, cmd
				}
				// Go to bottom of viewport
				newModel := m.Clone()
				newModel.core.Viewport.GotoBottom()
				return newModel, nil
			default:
				// Pass keys to TextArea only if in input mode
				if m.core.IsLambdaInputMode {
					// Always pass these navigation keys to the TextArea
					key := msg.String()
					if key == "f" || key == "g" || key == "j" || key == "k" || key == "b" {
						var cmd tea.Cmd
						newModel := m.Clone()
						newModel.core.TextArea, cmd = newModel.core.TextArea.Update(msg)
						newModel.core.LambdaPayload = newModel.core.TextArea.Value()
						return newModel, cmd
					}
					var cmd tea.Cmd
					newModel := m.Clone()
					newModel.core.TextArea, cmd = newModel.core.TextArea.Update(msg)
					newModel.core.LambdaPayload = newModel.core.TextArea.Value()
					return newModel, cmd
				}
				return m, nil
			}
		}

		// Handle key presses when not loading
		switch msg.String() {
		case constants.KeyCtrlC, constants.KeyQ:
			return m, tea.Quit
		case constants.KeyEnter:
			// If there's an error, clear it and allow navigation
			if m.core.Err != nil {
				newModel := m.Clone()
				newModel.core.Err = nil
				return newModel, nil
			}

			modelWrapper, cmd := update.HandleEnter(m.core)
			if wrapper, ok := modelWrapper.(update.ModelWrapper); ok {
				// Since ModelWrapper embeds *model.Model, we can create a new Model with it
				newModel := Model{core: wrapper.Model}
				if newModel.core.IsLoading {
					return newModel, tea.Batch(cmd, newModel.core.Spinner.Tick)
				}
				return newModel, cmd
			}
			return modelWrapper, cmd
		case constants.KeyEsc, constants.KeyAltBack:
			// If there's an error, clear it and navigate back
			if m.core.Err != nil {
				newCore := update.NavigateBack(m.core)
				newCore.Err = nil // Clear the error
				view.UpdateTableForView(newCore)
				return Model{core: newCore}, nil
			}

			// Only use '-' for back navigation if not in text input mode
			if msg.String() == constants.KeyAltBack && m.core.ManualInput {
				// If in text input mode, '-' should be treated as a character
				newModel := m.Clone()
				var cmd tea.Cmd
				newModel.core.TextInput, cmd = newModel.core.TextInput.Update(msg)
				return newModel, cmd
			}

			// Handle back navigation
			if m.core.ManualInput {
				newModel := m.Clone()
				newModel.core.ManualInput = false
				newModel.core.ResetTextInput()
				view.UpdateTableForView(newModel.core)
				return newModel, nil
			}
			newCore := update.NavigateBack(m.core)
			view.UpdateTableForView(newCore)
			return Model{core: newCore}, nil
		case constants.KeyUp, constants.KeyAltUp:
			// If in text input mode, pass 'k' to the text input
			if m.core.ManualInput && msg.String() == constants.KeyAltUp {
				newModel := m.Clone()
				var cmd tea.Cmd
				newModel.core.TextInput, cmd = newModel.core.TextInput.Update(msg)
				return newModel, cmd
			}
			newModel := m.Clone()
			newModel.core.Table.MoveUp(1)
			return newModel, nil
		case constants.KeyDown, constants.KeyAltDown:
			// If in text input mode, pass 'j' to the text input
			if m.core.ManualInput && msg.String() == constants.KeyAltDown {
				newModel := m.Clone()
				var cmd tea.Cmd
				newModel.core.TextInput, cmd = newModel.core.TextInput.Update(msg)
				return newModel, cmd
			}
			newModel := m.Clone()
			newModel.core.Table.MoveDown(1)
			return newModel, nil
		// Add vim-like navigation keys
		case constants.KeyGotoTop, constants.KeyHome:
			// If in text input mode, pass the key to the text input
			if m.core.ManualInput {
				newModel := m.Clone()
				var cmd tea.Cmd
				newModel.core.TextInput, cmd = newModel.core.TextInput.Update(msg)
				return newModel, cmd
			}
			newModel := m.Clone()
			newModel.core.Table.GotoTop()
			return newModel, nil
		case constants.KeyGotoBottom, constants.KeyEnd:
			// If in text input mode, pass the key to the text input
			if m.core.ManualInput {
				newModel := m.Clone()
				var cmd tea.Cmd
				newModel.core.TextInput, cmd = newModel.core.TextInput.Update(msg)
				return newModel, cmd
			}
			newModel := m.Clone()
			newModel.core.Table.GotoBottom()
			return newModel, nil
		case constants.KeyHalfPageUp, constants.KeyAltHalfPageUp:
			// If in text input mode, pass the key to the text input
			if m.core.ManualInput {
				newModel := m.Clone()
				var cmd tea.Cmd
				newModel.core.TextInput, cmd = newModel.core.TextInput.Update(msg)
				return newModel, cmd
			}
			newModel := m.Clone()
			newModel.core.Table.MoveUp(newModel.core.Table.Height() / 2)
			return newModel, nil
		case constants.KeyHalfPageDown, constants.KeyAltHalfPageDown:
			// If in text input mode, pass the key to the text input
			if m.core.ManualInput {
				newModel := m.Clone()
				var cmd tea.Cmd
				newModel.core.TextInput, cmd = newModel.core.TextInput.Update(msg)
				return newModel, cmd
			}
			newModel := m.Clone()
			newModel.core.Table.MoveDown(newModel.core.Table.Height() / 2)
			return newModel, nil
		case constants.KeyPageUp, constants.KeyAltPageUp:
			// If in text input mode, pass the key to the text input
			if m.core.ManualInput {
				newModel := m.Clone()
				var cmd tea.Cmd
				newModel.core.TextInput, cmd = newModel.core.TextInput.Update(msg)
				return newModel, cmd
			}
			newModel := m.Clone()
			newModel.core.Table.MoveUp(newModel.core.Table.Height())
			return newModel, nil
		case constants.KeyPageDown, constants.KeyAltPageDown, constants.KeySpace:
			// If in text input mode, pass the key to the text input
			if m.core.ManualInput {
				newModel := m.Clone()
				var cmd tea.Cmd
				newModel.core.TextInput, cmd = newModel.core.TextInput.Update(msg)
				return newModel, cmd
			}
			newModel := m.Clone()
			newModel.core.Table.MoveDown(newModel.core.Table.Height())
			return newModel, nil
		case constants.KeyTab:
			// Tab key is now used for text input only
			if m.core.CurrentView == constants.ViewLambdaExecute && m.core.IsLambdaInputMode {
				// Pass the tab key to the text area for indentation
				newModel := m.Clone()
				var cmd tea.Cmd
				newModel.core.TextArea, cmd = newModel.core.TextArea.Update(msg)
				return newModel, cmd
			}
			return m, nil
		// Add pagination key handlers
		case constants.KeyPreviousPage, constants.KeyNextPage:
			// Handle pagination key presses if in a paginated view
			if view.IsPaginatedView(m.core.CurrentView) {
				// Pass the exact key string that was pressed
				modelWrapper, cmd := update.HandlePaginationKeyPress(m.core, msg.String())
				if wrapper, ok := modelWrapper.(update.ModelWrapper); ok {
					newModel := Model{core: wrapper.Model}
					if cmd != nil {
						return newModel, cmd
					}
					return newModel, nil
				}
				return modelWrapper, cmd
			}
			return m, nil
		default:
			if m.core.ManualInput {
				newModel := m.Clone()
				var cmd tea.Cmd
				newModel.core.TextInput, cmd = newModel.core.TextInput.Update(msg)

				// If we're in the summary view with manual commit ID
				if newModel.core.CurrentView == constants.ViewSummary && newModel.core.SelectedOperation != nil &&
					newModel.core.SelectedOperation.Name == "Start Pipeline" && newModel.core.ManualInput {
					newModel.core.CommitID = newModel.core.TextInput.Value()
					newModel.core.ManualCommitID = true
				}

				// If we're in the summary view with approval comment
				if newModel.core.CurrentView == constants.ViewSummary && newModel.core.SelectedApproval != nil {
					newModel.core.ApprovalComment = newModel.core.TextInput.Value()
				}

				// For AWS config view, the actual setting happens when Enter is pressed in HandleEnter
				return newModel, cmd
			}
		}
	case model.PipelineStatusMsg:
		newModel := m.Clone()
		newModel.core.Pipelines = msg.Pipelines
		newModel.core.Provider = msg.Provider
		newModel.core.CurrentView = constants.ViewPipelineStatus
		newModel.core.IsLoading = false

		// Only show the first page of pipelines based on page size
		pageSize := newModel.core.PageSize
		initialPipelines := newModel.core.Pipelines

		// If we have more pipelines than the page size, only show the first page
		if len(initialPipelines) > pageSize {
			initialPipelines = initialPipelines[:pageSize]
		}

		// Convert to a PipelinesPageMsg for consistent pagination handling
		pipelinesPageMsg := model.PipelinesPageMsg{
			Pipelines:     initialPipelines,
			NextPageToken: "2", // Set to page 2 for next page
			HasMorePages:  len(newModel.core.Pipelines) > pageSize,
		}

		// Use the pagination handler
		newModel.core = update.HandlePipelineStatusPagination(newModel.core, pipelinesPageMsg)

		// Store all pipelines in the pagination state for client-side pagination
		if len(newModel.core.Pipelines) > 0 {
			for _, pipeline := range msg.Pipelines {
				// We need to add all pipelines to AllItems, not just the ones displayed
				found := false
				for _, item := range newModel.core.Pagination.AllItems {
					if p, ok := item.(model.PipelineStatus); ok && p.Name == pipeline.Name {
						found = true
						break
					}
				}
				if !found {
					newModel.core.Pagination.AllItems = append(newModel.core.Pagination.AllItems, pipeline)
				}
			}
			// Update total items count
			newModel.core.Pagination.TotalItems = int64(len(newModel.core.Pagination.AllItems))
		}

		return newModel, nil
	// Add handlers for pagination messages
	case model.FunctionsPageMsg:
		newModel := m.Clone()
		newModel.core = update.HandleFunctionStatusPagination(newModel.core, msg)
		return newModel, nil
	case model.PipelinesPageMsg:
		newModel := m.Clone()
		newModel.core = update.HandlePipelineStatusPagination(newModel.core, msg)
		return newModel, nil
	case model.ApprovalsPageMsg:
		newModel := m.Clone()
		newModel.core = update.HandleApprovalsPagination(newModel.core, msg)
		return newModel, nil
	case tea.MouseMsg:
		// If we're in the Lambda response view, pass mouse events to the viewport
		if m.core.CurrentView == constants.ViewLambdaResponse {
			newModel := m.Clone()
			var cmd tea.Cmd
			newModel.core.Viewport, cmd = newModel.core.Viewport.Update(msg)
			return newModel, cmd
		}
		return m, nil
	}

	// If we're loading, make sure to keep the spinner spinning
	if m.core.IsLoading {
		return m, m.core.Spinner.Tick
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	return view.Render(m.core)
}

// Clone creates a deep copy of the model
func (m Model) Clone() Model {
	return Model{
		core: m.core.Clone(),
	}
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
