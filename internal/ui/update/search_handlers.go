package update

import (
	"strings"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
)

// ActivateSearch activates search mode
func ActivateSearch(m *model.Model) *model.Model {
	newModel := m.Clone()
	newModel.Search.IsActive = true
	newModel.Search.Query = ""
	return newModel
}

// DeactivateSearch deactivates search mode
func DeactivateSearch(m *model.Model) *model.Model {
	newModel := m.Clone()
	newModel.Search.IsActive = false
	newModel.Search.Query = ""

	// Reset filtered items
	newModel.Search.FilteredItems = make([]interface{}, 0)

	// Reset pagination to show all items
	newModel.Pagination.CurrentPage = 1

	// Update the table with all items
	refreshTableWithAllItems(newModel)

	return newModel
}

// UpdateSearchQuery updates the search query and filters items
func UpdateSearchQuery(m *model.Model, query string) *model.Model {
	newModel := m.Clone()
	newModel.Search.Query = query

	// If query is empty, show all items
	if query == "" {
		newModel.Search.FilteredItems = make([]interface{}, 0)
		newModel.Pagination.CurrentPage = 1
		refreshTableWithAllItems(newModel)
		return newModel
	}

	// Filter items based on the query
	filteredItems := filterItemsByQuery(newModel.Pagination.AllItems, query)
	newModel.Search.FilteredItems = filteredItems
	newModel.Pagination.CurrentPage = 1

	// Update the table with filtered items
	updateTableWithItems(newModel, filteredItems)

	return newModel
}

// filterItemsByQuery filters items based on the query using simple substring matching
func filterItemsByQuery(items []interface{}, query string) []interface{} {
	if query == "" {
		return items
	}

	// Convert query to lowercase for case-insensitive matching
	lowerQuery := strings.ToLower(query)

	var filteredItems []interface{}

	for _, item := range items {
		// Extract searchable text based on item type
		var searchText string

		switch v := item.(type) {
		case model.FunctionStatus:
			searchText = strings.ToLower(v.Name)
			// Include other fields for more comprehensive search
			if v.Runtime != "" {
				searchText += " " + strings.ToLower(v.Runtime)
			}
		case model.PipelineStatus:
			searchText = strings.ToLower(v.Name)
			// Include stages for more comprehensive search
			for _, stage := range v.Stages {
				searchText += " " + strings.ToLower(stage.Name)
				searchText += " " + strings.ToLower(stage.Status)
			}
		case model.ApprovalAction:
			searchText = strings.ToLower(v.PipelineName + " " + v.StageName + " " + v.ActionName)
		default:
			// Skip unknown types
			continue
		}

		// Simple substring match
		if strings.Contains(searchText, lowerQuery) {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

// refreshTableWithAllItems updates the table with all items
func refreshTableWithAllItems(m *model.Model) {
	updateTableWithItems(m, m.Pagination.AllItems)
}

// updateTableWithItems updates the table with the given items
func updateTableWithItems(m *model.Model, items []interface{}) {
	// Calculate the current page of items
	pageSize := m.PageSize
	startIdx := (m.Pagination.CurrentPage - 1) * pageSize
	endIdx := startIdx + pageSize

	if endIdx > len(items) {
		endIdx = len(items)
	}

	// Update the table based on the current view
	switch m.CurrentView {
	case constants.ViewFunctionStatus:
		updateFunctionsTable(m, items, startIdx, endIdx)
	case constants.ViewPipelineStatus:
		updatePipelinesTable(m, items, startIdx, endIdx)
	case constants.ViewApprovals:
		updateApprovalsTable(m, items, startIdx, endIdx)
	}

	// Update pagination info
	m.Pagination.TotalItems = int64(len(items))
	m.Pagination.HasMorePages = endIdx < len(items)

	// Update the table view to reflect the changes
	view.UpdateTableForView(m)
}

// updateFunctionsTable updates the functions table with the given items
func updateFunctionsTable(m *model.Model, items []interface{}, startIdx, endIdx int) {
	var functions []model.FunctionStatus

	// Extract the current page of functions
	for i := startIdx; i < endIdx && i < len(items); i++ {
		if function, ok := items[i].(model.FunctionStatus); ok {
			functions = append(functions, function)
		}
	}

	m.Functions = functions
}

// updatePipelinesTable updates the pipelines table with the given items
func updatePipelinesTable(m *model.Model, items []interface{}, startIdx, endIdx int) {
	var pipelines []model.PipelineStatus

	// Extract the current page of pipelines
	for i := startIdx; i < endIdx && i < len(items); i++ {
		if pipeline, ok := items[i].(model.PipelineStatus); ok {
			pipelines = append(pipelines, pipeline)
		}
	}

	m.Pipelines = pipelines
}

// updateApprovalsTable updates the approvals table with the given items
func updateApprovalsTable(m *model.Model, items []interface{}, startIdx, endIdx int) {
	var approvals []model.ApprovalAction

	// Extract the current page of approvals
	for i := startIdx; i < endIdx && i < len(items); i++ {
		if approval, ok := items[i].(model.ApprovalAction); ok {
			approvals = append(approvals, approval)
		}
	}

	m.Approvals = approvals
}

// IsPrintableChar checks if a character is printable
func IsPrintableChar(r rune) bool {
	return r >= 32 && r < 127
}
