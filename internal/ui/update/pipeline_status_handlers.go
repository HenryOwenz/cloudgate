package update

import (
	"context"
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
)

// HandlePipelineExecution handles the result of a pipeline execution
func HandlePipelineExecution(m *model.Model, err error) {
	if err != nil {
		m.Error = fmt.Sprintf(constants.MsgErrorGeneric, err.Error())
		m.CurrentView = constants.ViewError
		return
	}

	m.Success = fmt.Sprintf(constants.MsgPipelineStartSuccess, m.SelectedPipeline.Name)

	// Reset pipeline state
	m.SelectedPipeline = nil
	m.CommitID = ""
	m.ManualCommitID = false

	// Completely reset the text input
	m.ResetTextInput()
	m.TextInput.Placeholder = constants.MsgEnterComment
	m.ManualInput = false

	// Reset pagination state
	m.Pagination.Type = model.PaginationTypeNone
	m.Pagination.CurrentPage = 1
	m.Pagination.HasMorePages = false
	m.Pagination.AllItems = make([]interface{}, 0)
	m.Pagination.FilteredItems = make([]interface{}, 0)
	m.Pagination.TotalItems = 0

	// Reset search state
	m.Search.IsActive = false
	m.Search.Query = ""
	m.Search.FilteredItems = make([]interface{}, 0)

	// Navigate back to the operation selection view
	m.CurrentView = constants.ViewSelectOperation

	// Clear all lists to force a refresh next time
	m.Pipelines = nil
	m.Functions = nil
	m.Approvals = nil

	// Update the table for the current view
	view.UpdateTableForView(m)
}

// FetchPipelineStatus fetches pipeline status from the provider
func FetchPipelineStatus(m *model.Model) tea.Cmd {
	// Initialize pagination state in the model
	newModel := m.Clone()
	newModel.Pagination.Type = model.PaginationTypeClientSide
	newModel.Pagination.CurrentPage = 1
	newModel.Pagination.PageSize = newModel.PageSize
	newModel.Pagination.TotalItems = -1 // Unknown until we fetch the data
	newModel.Pagination.HasMorePages = false
	newModel.Pagination.IsLoading = true
	newModel.Pagination.AllItems = make([]interface{}, 0)
	newModel.Pagination.FilteredItems = make([]interface{}, 0)

	return func() tea.Msg {
		// Get the provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get the PipelineStatusOperation from the provider
		statusOperation, err := provider.GetPipelineStatusOperation()
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get pipeline status using the operation
		ctx := context.Background()
		pipelines, err := statusOperation.GetPipelineStatus(ctx)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Sort pipelines by name in ascending order (case-insensitive)
		sort.Slice(pipelines, func(i, j int) bool {
			return strings.ToLower(pipelines[i].Name) < strings.ToLower(pipelines[j].Name)
		})

		// Implement client-side pagination
		totalItems := int64(len(pipelines))
		pageSize := m.PageSize

		// Determine if there are more pages
		hasMorePages := totalItems > int64(pageSize)

		// Get the first page of pipelines
		endIdx := pageSize
		if endIdx > len(pipelines) {
			endIdx = len(pipelines)
		}
		firstPagePipelines := pipelines[:endIdx]

		return model.PipelinesPageMsg{
			Pipelines:     firstPagePipelines,
			NextPageToken: "1", // Use page number as token for client-side pagination
			HasMorePages:  hasMorePages,
		}
	}
}

// ExecutePipeline executes a pipeline
func ExecutePipeline(m *model.Model) tea.Cmd {
	return func() tea.Msg {
		if m.SelectedPipeline == nil {
			return model.ErrMsg{Err: fmt.Errorf("no pipeline selected")}
		}

		// Get the provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get the StartPipelineOperation from the provider
		startOperation, err := provider.GetStartPipelineOperation()
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Execute the pipeline using the operation
		ctx := context.Background()
		err = startOperation.StartPipelineExecution(ctx, m.SelectedPipeline.Name, m.CommitID)
		if err != nil {
			return model.PipelineExecutionMsg{Err: err}
		}

		return model.PipelineExecutionMsg{Err: nil}
	}
}
