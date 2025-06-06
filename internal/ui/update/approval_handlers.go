package update

import (
	"context"
	"fmt"
	"sort"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
)

// HandleApprovalResult handles the result of an approval action
func HandleApprovalResult(m *model.Model, err error) {
	if err != nil {
		m.Error = fmt.Sprintf(constants.MsgErrorGeneric, err.Error())
		m.CurrentView = constants.ViewError
		return
	}

	// Check if SelectedApproval is nil before accessing its fields
	if m.SelectedApproval == nil {
		m.Error = constants.MsgErrorNoApproval
		m.CurrentView = constants.ViewError
		return
	}

	// Use the appropriate message constant based on approval action
	if m.ApproveAction {
		m.Success = fmt.Sprintf(constants.MsgApprovalSuccess,
			m.SelectedApproval.PipelineName,
			m.SelectedApproval.StageName,
			m.SelectedApproval.ActionName)
	} else {
		m.Success = fmt.Sprintf(constants.MsgRejectionSuccess,
			m.SelectedApproval.PipelineName,
			m.SelectedApproval.StageName,
			m.SelectedApproval.ActionName)
	}

	// Reset approval state
	m.SelectedApproval = nil
	m.ApprovalComment = ""

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

	// Clear the approvals list to force a refresh next time
	m.Approvals = nil
	m.Pipelines = nil
	m.Functions = nil

	// Update the table for the current view
	view.UpdateTableForView(m)
}

// FetchApprovals fetches pipeline approvals from the provider
func FetchApprovals(m *model.Model) tea.Cmd {
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

		// Get the CodePipelineManualApprovalOperation from the provider
		approvalOperation, err := provider.GetCodePipelineManualApprovalOperation()
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get approvals using the operation
		ctx := context.Background()
		approvals, err := approvalOperation.GetPendingApprovals(ctx)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Sort approvals by pipeline name, stage name, and action name
		sort.Slice(approvals, func(i, j int) bool {
			if approvals[i].PipelineName != approvals[j].PipelineName {
				return approvals[i].PipelineName < approvals[j].PipelineName
			}
			if approvals[i].StageName != approvals[j].StageName {
				return approvals[i].StageName < approvals[j].StageName
			}
			return approvals[i].ActionName < approvals[j].ActionName
		})

		// Implement client-side pagination
		totalItems := int64(len(approvals))
		pageSize := m.PageSize

		// Determine if there are more pages
		hasMorePages := totalItems > int64(pageSize)

		// Get the first page of approvals
		endIdx := pageSize
		if endIdx > len(approvals) {
			endIdx = len(approvals)
		}
		firstPageApprovals := approvals[:endIdx]

		return model.ApprovalsPageMsg{
			Approvals:     firstPageApprovals,
			NextPageToken: "1", // Use page number as token for client-side pagination
			HasMorePages:  hasMorePages,
		}
	}
}

// ExecuteApproval executes an approval action
func ExecuteApproval(m *model.Model) tea.Cmd {
	return func() tea.Msg {
		if m.SelectedApproval == nil {
			return model.ErrMsg{Err: fmt.Errorf("no approval selected")}
		}

		// Get the provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get the CodePipelineManualApprovalOperation from the provider
		approvalOperation, err := provider.GetCodePipelineManualApprovalOperation()
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Execute the approval action using the operation
		ctx := context.Background()
		err = approvalOperation.ApproveAction(ctx, *m.SelectedApproval, m.ApproveAction, m.ApprovalComment)
		if err != nil {
			return model.ApprovalResultMsg{Err: err}
		}

		return model.ApprovalResultMsg{Err: nil}
	}
}
