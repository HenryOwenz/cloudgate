package update

import (
	"context"
	"fmt"
	"strconv"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	tea "github.com/charmbracelet/bubbletea"
)

// HandlePaginationKeyPress handles pagination key presses (h for previous, l for next)
func HandlePaginationKeyPress(m *model.Model, key string) (tea.Model, tea.Cmd) {
	if !view.IsPaginatedView(m.CurrentView) || m.Pagination.Type == model.PaginationTypeNone {
		return WrapModel(m), nil
	}

	newModel := m.Clone()

	switch key {
	case constants.KeyNextPage: // Next page
		if newModel.Pagination.HasMorePages && !newModel.Pagination.IsLoading {
			// Don't return the command directly, just set it up to be executed
			cmd := FetchNextPage(newModel)
			return WrapModel(newModel), cmd
		}
	case constants.KeyPreviousPage: // Previous page
		if newModel.Pagination.CurrentPage > 1 && !newModel.Pagination.IsLoading {
			// Don't return the command directly, just set it up to be executed
			cmd := FetchPreviousPage(newModel)
			return WrapModel(newModel), cmd
		}
	}

	return WrapModel(newModel), nil
}

// FetchNextPage fetches the next page based on the current view
func FetchNextPage(m *model.Model) tea.Cmd {
	switch m.CurrentView {
	case constants.ViewFunctionStatus:
		return FetchNextFunctionsPage(m)
	case constants.ViewPipelineStatus:
		return FetchNextPipelinesPage(m)
	case constants.ViewApprovals:
		return FetchNextApprovalsPage(m)
	default:
		return nil
	}
}

// FetchPreviousPage fetches the previous page based on the current view
func FetchPreviousPage(m *model.Model) tea.Cmd {
	// Don't fetch if already loading or if we're on the first page
	if m.Pagination.IsLoading || m.Pagination.CurrentPage <= 1 {
		return nil
	}

	// Calculate the previous page number
	prevPage := m.Pagination.CurrentPage - 1
	if prevPage < 1 {
		prevPage = 1
	}

	// If we have cached pages, use them for client-side pagination
	if m.Pagination.Type == model.PaginationTypeClientSide && len(m.Pagination.AllItems) > 0 {
		// Use cached data based on the view
		switch m.CurrentView {
		case constants.ViewFunctionStatus:
			// Calculate start and end indices for the page
			pageSize := m.Pagination.PageSize
			start := (prevPage - 1) * pageSize
			end := start + pageSize
			if end > len(m.Pagination.AllItems) {
				end = len(m.Pagination.AllItems)
			}

			// Extract functions for this page
			var pageFunctions []model.FunctionStatus
			for i := start; i < end && i < len(m.Pagination.AllItems); i++ {
				if function, ok := m.Pagination.AllItems[i].(model.FunctionStatus); ok {
					pageFunctions = append(pageFunctions, function)
				}
			}

			// Determine if there are more pages
			hasMorePages := true // There are always more pages when going back

			return func() tea.Msg {
				return model.FunctionsPageMsg{
					Functions:     pageFunctions,
					HasMorePages:  hasMorePages,
					NextPageToken: strconv.Itoa(prevPage), // Set the token to the previous page number
				}
			}
		case constants.ViewPipelineStatus:
			// Similar logic for pipelines
			pageSize := m.Pagination.PageSize
			start := (prevPage - 1) * pageSize
			end := start + pageSize
			if end > len(m.Pagination.AllItems) {
				end = len(m.Pagination.AllItems)
			}

			// Extract pipelines for this page
			var pagePipelines []model.PipelineStatus
			for i := start; i < end && i < len(m.Pagination.AllItems); i++ {
				if pipeline, ok := m.Pagination.AllItems[i].(model.PipelineStatus); ok {
					pagePipelines = append(pagePipelines, pipeline)
				}
			}

			// Determine if there are more pages
			hasMorePages := true // There are always more pages when going back

			return func() tea.Msg {
				return model.PipelinesPageMsg{
					Pipelines:     pagePipelines,
					HasMorePages:  hasMorePages,
					NextPageToken: strconv.Itoa(prevPage), // Set the token to the previous page number
				}
			}
		case constants.ViewApprovals:
			// Similar logic for approvals
			pageSize := m.Pagination.PageSize
			start := (prevPage - 1) * pageSize
			end := start + pageSize
			if end > len(m.Pagination.AllItems) {
				end = len(m.Pagination.AllItems)
			}

			// Extract approvals for this page
			var pageApprovals []model.ApprovalAction
			for i := start; i < end && i < len(m.Pagination.AllItems); i++ {
				if approval, ok := m.Pagination.AllItems[i].(model.ApprovalAction); ok {
					pageApprovals = append(pageApprovals, approval)
				}
			}

			// Determine if there are more pages
			hasMorePages := true // There are always more pages when going back

			return func() tea.Msg {
				return model.ApprovalsPageMsg{
					Approvals:     pageApprovals,
					HasMorePages:  hasMorePages,
					NextPageToken: strconv.Itoa(prevPage), // Set the token to the previous page number
				}
			}
		}
	}

	return nil
}

// FetchNextFunctionsPage fetches the next page of Lambda functions
func FetchNextFunctionsPage(m *model.Model) tea.Cmd {
	return func() tea.Msg {
		// Set loading state
		newModel := m.Clone()
		newModel.Pagination.IsLoading = true

		// For client-side pagination, we already have all the data
		if newModel.Pagination.Type == model.PaginationTypeClientSide {
			// Calculate the next page
			nextPage := newModel.Pagination.CurrentPage + 1
			pageSize := newModel.Pagination.PageSize

			// Calculate start and end indices for the page
			startIdx := (nextPage - 1) * pageSize
			endIdx := startIdx + pageSize

			// Make sure we don't go out of bounds
			if startIdx >= len(newModel.Pagination.AllItems) {
				// No more items
				return model.FunctionsPageMsg{
					Functions:     []model.FunctionStatus{},
					NextPageToken: "",
					HasMorePages:  false,
				}
			}

			if endIdx > len(newModel.Pagination.AllItems) {
				endIdx = len(newModel.Pagination.AllItems)
			}

			// Extract functions for this page
			var pageFunctions []model.FunctionStatus
			for i := startIdx; i < endIdx; i++ {
				if function, ok := newModel.Pagination.AllItems[i].(model.FunctionStatus); ok {
					pageFunctions = append(pageFunctions, function)
				}
			}

			// Determine if there are more pages
			hasMorePages := endIdx < len(newModel.Pagination.AllItems)

			return model.FunctionsPageMsg{
				Functions:     pageFunctions,
				NextPageToken: fmt.Sprintf("%d", nextPage),
				HasMorePages:  hasMorePages,
			}
		}

		// For API-level pagination, we need to fetch from the API
		// Get the provider
		provider, err := m.Registry.Get(m.ProviderState.ProviderName)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get the function status operation
		functionOperation, err := provider.GetFunctionStatusOperation()
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get the next page of functions
		ctx := context.Background()
		functions, err := functionOperation.GetFunctionStatus(ctx)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// For now, we're simulating pagination since the actual API doesn't support it yet
		// In a real implementation, we would pass the nextMarker and pageSize to GetFunctionStatus

		// Determine if there are more pages (simulated for now)
		hasMorePages := false
		var nextPageToken string

		return model.FunctionsPageMsg{
			Functions:     functions,
			NextPageToken: nextPageToken,
			HasMorePages:  hasMorePages,
		}
	}
}

// FetchNextPipelinesPage fetches the next page of CodePipeline pipelines
func FetchNextPipelinesPage(m *model.Model) tea.Cmd {
	return func() tea.Msg {
		// Set loading state
		newModel := m.Clone()
		newModel.Pagination.IsLoading = true

		// For client-side pagination, we already have all the data
		if newModel.Pagination.Type == model.PaginationTypeClientSide {
			// Calculate the next page
			nextPage := newModel.Pagination.CurrentPage + 1
			pageSize := newModel.Pagination.PageSize

			// Calculate start and end indices for the page
			startIdx := (nextPage - 1) * pageSize
			endIdx := startIdx + pageSize

			// Make sure we don't go out of bounds
			if startIdx >= len(newModel.Pagination.AllItems) {
				// No more items
				return model.PipelinesPageMsg{
					Pipelines:     []model.PipelineStatus{},
					NextPageToken: "",
					HasMorePages:  false,
				}
			}

			if endIdx > len(newModel.Pagination.AllItems) {
				endIdx = len(newModel.Pagination.AllItems)
			}

			// Extract pipelines for this page
			var pagePipelines []model.PipelineStatus
			for i := startIdx; i < endIdx; i++ {
				if pipeline, ok := newModel.Pagination.AllItems[i].(model.PipelineStatus); ok {
					pagePipelines = append(pagePipelines, pipeline)
				}
			}

			// Determine if there are more pages
			hasMorePages := endIdx < len(newModel.Pagination.AllItems)

			return model.PipelinesPageMsg{
				Pipelines:     pagePipelines,
				NextPageToken: fmt.Sprintf("%d", nextPage),
				HasMorePages:  hasMorePages,
			}
		}

		// For API-level pagination, we need to fetch from the API
		// Get the provider
		provider, err := m.Registry.Get(m.ProviderState.ProviderName)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get the pipeline status operation
		pipelineOperation, err := provider.GetPipelineStatusOperation()
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get the next page of pipelines
		ctx := context.Background()
		pipelines, err := pipelineOperation.GetPipelineStatus(ctx)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// For now, we're simulating pagination since the actual API doesn't support it yet
		// In a real implementation, we would pass the nextToken and pageSize to GetPipelineStatus

		// Determine if there are more pages (simulated for now)
		hasMorePages := false
		var nextPageToken string

		return model.PipelinesPageMsg{
			Pipelines:     pipelines,
			NextPageToken: nextPageToken,
			HasMorePages:  hasMorePages,
		}
	}
}

// FetchNextApprovalsPage fetches the next page of pipeline approvals
func FetchNextApprovalsPage(m *model.Model) tea.Cmd {
	return func() tea.Msg {
		// Set loading state
		newModel := m.Clone()
		newModel.Pagination.IsLoading = true

		// For client-side pagination, we already have all the data
		if newModel.Pagination.Type == model.PaginationTypeClientSide {
			// Calculate the next page
			nextPage := newModel.Pagination.CurrentPage + 1
			pageSize := newModel.Pagination.PageSize

			// Calculate start and end indices for the page
			startIdx := (nextPage - 1) * pageSize
			endIdx := startIdx + pageSize

			// Make sure we don't go out of bounds
			if startIdx >= len(newModel.Pagination.AllItems) {
				// No more items
				return model.ApprovalsPageMsg{
					Approvals:     []model.ApprovalAction{},
					NextPageToken: "",
					HasMorePages:  false,
				}
			}

			if endIdx > len(newModel.Pagination.AllItems) {
				endIdx = len(newModel.Pagination.AllItems)
			}

			// Extract approvals for this page
			var pageApprovals []model.ApprovalAction
			for i := startIdx; i < endIdx; i++ {
				if approval, ok := newModel.Pagination.AllItems[i].(model.ApprovalAction); ok {
					pageApprovals = append(pageApprovals, approval)
				}
			}

			// Determine if there are more pages
			hasMorePages := endIdx < len(newModel.Pagination.AllItems)

			return model.ApprovalsPageMsg{
				Approvals:     pageApprovals,
				NextPageToken: fmt.Sprintf("%d", nextPage),
				HasMorePages:  hasMorePages,
			}
		}

		// For API-level pagination, we need to fetch from the API
		// Get the provider
		provider, err := m.Registry.Get(m.ProviderState.ProviderName)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get the approval operation
		approvalOperation, err := provider.GetCodePipelineManualApprovalOperation()
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get the next page of approvals
		ctx := context.Background()
		approvals, err := approvalOperation.GetPendingApprovals(ctx)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// For now, we're simulating pagination since the actual API doesn't support it yet
		// In a real implementation, we would pass the nextToken and pageSize to GetPendingApprovals

		// Determine if there are more pages (simulated for now)
		hasMorePages := false
		var nextPageToken string

		return model.ApprovalsPageMsg{
			Approvals:     approvals,
			NextPageToken: nextPageToken,
			HasMorePages:  hasMorePages,
		}
	}
}

// HandleFunctionStatusPagination handles the pagination of Lambda functions
func HandleFunctionStatusPagination(m *model.Model, msg model.FunctionsPageMsg) *model.Model {
	newModel := m.Clone()
	newModel.IsLoading = false
	newModel.Pagination.IsLoading = false

	// Update pagination state
	if newModel.Pagination.Type == model.PaginationTypeNone {
		// First time initialization
		newModel.Pagination.Type = model.PaginationTypeClientSide
		newModel.Pagination.CurrentPage = 1
		newModel.Pagination.AllItems = make([]interface{}, 0)
		newModel.Pagination.FilteredItems = make([]interface{}, 0)

		// Store all functions for client-side pagination
		for _, function := range msg.Functions {
			newModel.Pagination.AllItems = append(newModel.Pagination.AllItems, function)
		}

		// Update total items count
		newModel.Pagination.TotalItems = int64(len(newModel.Pagination.AllItems))
	} else if msg.NextPageToken != "" {
		// Only increment page if this is a new page (not a refresh of the same page)
		// And make sure we don't increment beyond the total number of pages
		nextPage, err := strconv.Atoi(msg.NextPageToken)
		if err == nil && nextPage > 0 {
			newModel.Pagination.CurrentPage = nextPage
		}
	}

	newModel.Pagination.HasMorePages = msg.HasMorePages

	// Store functions for the current view
	newModel.Functions = msg.Functions

	// Update the table for the current view
	view.UpdateTableForView(newModel)

	return newModel
}

// HandlePipelineStatusPagination handles the pagination of CodePipeline pipelines
func HandlePipelineStatusPagination(m *model.Model, msg model.PipelinesPageMsg) *model.Model {
	newModel := m.Clone()
	newModel.IsLoading = false
	newModel.Pagination.IsLoading = false

	// Update pagination state
	if newModel.Pagination.Type == model.PaginationTypeNone {
		// First time initialization
		newModel.Pagination.Type = model.PaginationTypeClientSide
		newModel.Pagination.CurrentPage = 1
		newModel.Pagination.AllItems = make([]interface{}, 0)
		newModel.Pagination.FilteredItems = make([]interface{}, 0)

		// Store all pipelines for client-side pagination
		for _, pipeline := range msg.Pipelines {
			newModel.Pagination.AllItems = append(newModel.Pagination.AllItems, pipeline)
		}

		// Update total items count
		newModel.Pagination.TotalItems = int64(len(newModel.Pagination.AllItems))
	} else if msg.NextPageToken != "" {
		// Only increment page if this is a new page (not a refresh of the same page)
		// And make sure we don't increment beyond the total number of pages
		nextPage, err := strconv.Atoi(msg.NextPageToken)
		if err == nil && nextPage > 0 {
			newModel.Pagination.CurrentPage = nextPage
		}
	}

	newModel.Pagination.HasMorePages = msg.HasMorePages

	// Store pipelines for the current view
	newModel.Pipelines = msg.Pipelines

	// Update the table for the current view
	view.UpdateTableForView(newModel)

	return newModel
}

// HandleApprovalsPagination handles the pagination of pipeline approvals
func HandleApprovalsPagination(m *model.Model, msg model.ApprovalsPageMsg) *model.Model {
	newModel := m.Clone()
	newModel.IsLoading = false
	newModel.Pagination.IsLoading = false

	// Update pagination state
	if newModel.Pagination.Type == model.PaginationTypeNone {
		// First time initialization
		newModel.Pagination.Type = model.PaginationTypeClientSide
		newModel.Pagination.CurrentPage = 1
		newModel.Pagination.AllItems = make([]interface{}, 0)
		newModel.Pagination.FilteredItems = make([]interface{}, 0)

		// Store all approvals for client-side pagination
		for _, approval := range msg.Approvals {
			newModel.Pagination.AllItems = append(newModel.Pagination.AllItems, approval)
		}

		// Update total items count
		newModel.Pagination.TotalItems = int64(len(newModel.Pagination.AllItems))
	} else if msg.NextPageToken != "" {
		// Only increment page if this is a new page (not a refresh of the same page)
		// And make sure we don't increment beyond the total number of pages
		nextPage, err := strconv.Atoi(msg.NextPageToken)
		if err == nil && nextPage > 0 {
			newModel.Pagination.CurrentPage = nextPage
		}
	}

	newModel.Pagination.HasMorePages = msg.HasMorePages

	// Store approvals for the current view
	newModel.Approvals = msg.Approvals

	// Update the table for the current view
	view.UpdateTableForView(newModel)

	return newModel
}
