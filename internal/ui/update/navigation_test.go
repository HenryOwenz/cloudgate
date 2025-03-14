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
