package update

import (
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
)

func TestLambdaNavigationFlow(t *testing.T) {
	testCases := []struct {
		name                 string
		isExecuteLambdaFlow  bool
		startView            constants.View
		expectedView         constants.View
		shouldClearSelection bool
	}{
		{
			name:                 "From_ViewLambdaExecute_to_ViewFunctionStatus_in_Execute_Function_flow",
			isExecuteLambdaFlow:  true,
			startView:            constants.ViewLambdaExecute,
			expectedView:         constants.ViewFunctionStatus,
			shouldClearSelection: true,
		},
		{
			name:                 "From_ViewLambdaExecute_to_ViewFunctionDetails_in_Function_Status_flow",
			isExecuteLambdaFlow:  false,
			startView:            constants.ViewLambdaExecute,
			expectedView:         constants.ViewFunctionDetails,
			shouldClearSelection: false,
		},
		{
			name:                 "From_ViewLambdaResponse_to_ViewLambdaExecute",
			isExecuteLambdaFlow:  true, // This should be the same for both flows
			startView:            constants.ViewLambdaResponse,
			expectedView:         constants.ViewLambdaExecute,
			shouldClearSelection: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a model with the test case configuration
			m := &model.Model{
				CurrentView:         tc.startView,
				IsExecuteLambdaFlow: tc.isExecuteLambdaFlow,
				SelectedFunction: &cloud.FunctionStatus{
					Name:    "test-function",
					Runtime: "nodejs14.x",
				},
			}

			// Navigate back
			result := NavigateBack(m)

			// Check that the view changed as expected
			if result.CurrentView != tc.expectedView {
				t.Errorf("Expected view to be %v, got %v", tc.expectedView, result.CurrentView)
			}

			// Check if the function selection was cleared when expected
			if tc.shouldClearSelection && result.SelectedFunction != nil {
				t.Errorf("Expected function selection to be cleared, but it wasn't")
			} else if !tc.shouldClearSelection && result.SelectedFunction == nil {
				t.Errorf("Expected function selection to be preserved, but it was cleared")
			}

			// For the Execute Function flow, verify the table is updated
			if tc.isExecuteLambdaFlow && tc.startView == constants.ViewLambdaExecute {
				// We can't directly test that UpdateTableForView was called,
				// but we can check that the model was properly prepared for it
				if result.CurrentView != constants.ViewFunctionStatus {
					t.Errorf("Expected view to be ViewFunctionStatus for Execute Function flow")
				}
			}
		})
	}
}

// TestNavigateBackChain tests a complete navigation chain for the Lambda execution flow
func TestLambdaNavigationChain(t *testing.T) {
	// Create a model for the Lambda execution flow
	m := &model.Model{
		CurrentView:         constants.ViewLambdaResponse,
		IsExecuteLambdaFlow: true,
		SelectedFunction:    &cloud.FunctionStatus{Name: "test-function"},
	}

	// First back navigation: ViewLambdaResponse -> ViewLambdaExecute
	result := NavigateBack(m)
	if result.CurrentView != constants.ViewLambdaExecute {
		t.Errorf("First navigation: Expected view to be ViewLambdaExecute, got %v", result.CurrentView)
	}
	if result.SelectedFunction == nil {
		t.Errorf("First navigation: Expected function selection to be preserved")
	}

	// Second back navigation: ViewLambdaExecute -> ViewFunctionStatus
	result = NavigateBack(result)
	if result.CurrentView != constants.ViewFunctionStatus {
		t.Errorf("Second navigation: Expected view to be ViewFunctionStatus, got %v", result.CurrentView)
	}
	if result.SelectedFunction != nil {
		t.Errorf("Second navigation: Expected function selection to be cleared")
	}

	// Third back navigation: ViewFunctionStatus -> ViewSelectOperation
	result = NavigateBack(result)
	if result.CurrentView != constants.ViewSelectOperation {
		t.Errorf("Third navigation: Expected view to be ViewSelectOperation, got %v", result.CurrentView)
	}
}

// TestStateResetOnNavigation tests that search and pagination states are properly reset when navigating
func TestStateResetOnNavigation(t *testing.T) {
	// Create a model with active search and pagination
	m := model.New()
	m.CurrentView = constants.ViewPipelineStatus
	m.Pagination.Type = model.PaginationTypeClientSide
	m.Pagination.CurrentPage = 3
	m.Pagination.HasMorePages = true
	m.Pagination.PageSize = 5
	m.Pagination.AllItems = createMockItems(20)
	m.Pagination.TotalItems = 20

	// Set up search state
	m.Search.IsActive = true
	m.Search.Query = "test-query"
	m.Search.FilteredItems = createMockItems(5)

	// Test navigation back to operation selection
	t.Run("Reset state when navigating back to operation selection", func(t *testing.T) {
		// Navigate back to operation selection
		newModel := NavigateBack(m)

		// Check that search state is reset
		if newModel.Search.IsActive {
			t.Errorf("Expected search to be deactivated after navigation, got active")
		}

		if newModel.Search.Query != "" {
			t.Errorf("Expected search query to be reset after navigation, got '%s'", newModel.Search.Query)
		}

		if len(newModel.Search.FilteredItems) > 0 {
			t.Errorf("Expected filtered items to be cleared after navigation, got %d items",
				len(newModel.Search.FilteredItems))
		}

		// Check that pagination state is reset
		if newModel.Pagination.CurrentPage != 1 {
			t.Errorf("Expected pagination to be reset to page 1, got page %d",
				newModel.Pagination.CurrentPage)
		}

		if newModel.Pagination.Type != model.PaginationTypeNone {
			t.Errorf("Expected pagination type to be reset to None, got %v",
				newModel.Pagination.Type)
		}

		if len(newModel.Pagination.AllItems) > 0 {
			t.Errorf("Expected all items to be cleared after navigation, got %d items",
				len(newModel.Pagination.AllItems))
		}
	})

	// Test navigation to a different view
	t.Run("Reset state when navigating to a different view", func(t *testing.T) {
		// Set up a model with a different view
		m.CurrentView = constants.ViewApprovals

		// Navigate to a different view (simulate by directly changing the view)
		newModel := m.Clone()
		newModel.CurrentView = constants.ViewPipelineStatus
		newModel = NavigateBack(newModel) // Use NavigateBack to reset state

		// Check that search state is reset
		if newModel.Search.IsActive {
			t.Errorf("Expected search to be deactivated after navigation, got active")
		}

		if newModel.Search.Query != "" {
			t.Errorf("Expected search query to be reset after navigation, got '%s'", newModel.Search.Query)
		}

		// Check that pagination state is reset
		if newModel.Pagination.CurrentPage != 1 {
			t.Errorf("Expected pagination to be reset to page 1, got page %d",
				newModel.Pagination.CurrentPage)
		}
	})
}
