package update

import (
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	tea "github.com/charmbracelet/bubbletea"
)

// TestHandlePaginationKeyPress tests the pagination key handling functionality
func TestHandlePaginationKeyPress(t *testing.T) {
	tests := []struct {
		name                string
		setupModel          func() *model.Model
		key                 string
		expectedPageChange  bool
		expectedPageNumber  int
		expectedHasMoreFlag bool
	}{
		{
			name: "Next page with 'l' key",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewPipelineStatus
				m.Pagination.Type = model.PaginationTypeClientSide
				m.Pagination.CurrentPage = 1
				m.Pagination.HasMorePages = true
				m.Pagination.PageSize = 5
				m.Pagination.AllItems = createMockItems(12) // 3 pages total
				return m
			},
			key:                 constants.KeyNextPage, // "l"
			expectedPageChange:  true,
			expectedPageNumber:  2,
			expectedHasMoreFlag: true,
		},
		{
			name: "Next page with right arrow key",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewPipelineStatus
				m.Pagination.Type = model.PaginationTypeClientSide
				m.Pagination.CurrentPage = 1
				m.Pagination.HasMorePages = true
				m.Pagination.PageSize = 5
				m.Pagination.AllItems = createMockItems(12) // 3 pages total
				return m
			},
			key:                 constants.KeyNextPage, // Use the mapped key, not the arrow key directly
			expectedPageChange:  true,
			expectedPageNumber:  2,
			expectedHasMoreFlag: true,
		},
		{
			name: "Previous page with 'h' key",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewPipelineStatus
				m.Pagination.Type = model.PaginationTypeClientSide
				m.Pagination.CurrentPage = 2
				m.Pagination.HasMorePages = true
				m.Pagination.PageSize = 5
				m.Pagination.AllItems = createMockItems(12) // 3 pages total
				return m
			},
			key:                 constants.KeyPreviousPage, // "h"
			expectedPageChange:  true,
			expectedPageNumber:  1,
			expectedHasMoreFlag: true,
		},
		{
			name: "Previous page with left arrow key",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewPipelineStatus
				m.Pagination.Type = model.PaginationTypeClientSide
				m.Pagination.CurrentPage = 2
				m.Pagination.HasMorePages = true
				m.Pagination.PageSize = 5
				m.Pagination.AllItems = createMockItems(12) // 3 pages total
				return m
			},
			key:                 constants.KeyPreviousPage, // Use the mapped key, not the arrow key directly
			expectedPageChange:  true,
			expectedPageNumber:  1,
			expectedHasMoreFlag: true,
		},
		{
			name: "Next page on last page",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewPipelineStatus
				m.Pagination.Type = model.PaginationTypeClientSide
				m.Pagination.CurrentPage = 3
				m.Pagination.HasMorePages = false
				m.Pagination.PageSize = 5
				m.Pagination.AllItems = createMockItems(12) // 3 pages total
				return m
			},
			key:                 constants.KeyNextPage,
			expectedPageChange:  false,
			expectedPageNumber:  3,
			expectedHasMoreFlag: false,
		},
		{
			name: "Previous page on first page",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewPipelineStatus
				m.Pagination.Type = model.PaginationTypeClientSide
				m.Pagination.CurrentPage = 1
				m.Pagination.HasMorePages = true
				m.Pagination.PageSize = 5
				m.Pagination.AllItems = createMockItems(12) // 3 pages total
				return m
			},
			key:                 constants.KeyPreviousPage,
			expectedPageChange:  false,
			expectedPageNumber:  1,
			expectedHasMoreFlag: true,
		},
		{
			name: "Pagination in non-paginated view",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewProviders // Not a paginated view
				return m
			},
			key:                 constants.KeyNextPage,
			expectedPageChange:  false,
			expectedPageNumber:  0, // Not relevant
			expectedHasMoreFlag: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup the model
			m := tt.setupModel()
			initialPage := m.Pagination.CurrentPage

			// Call the function with the key
			result, cmd := HandlePaginationKeyPress(m, tt.key)

			// Check the result type
			wrapper, ok := result.(ModelWrapper)
			if !ok {
				t.Fatalf("Expected HandlePaginationKeyPress to return a ModelWrapper, got %T", result)
			}

			// Check if page changed as expected
			if tt.expectedPageChange {
				// If we expect a page change, we should have a command
				if cmd == nil {
					t.Errorf("Expected a command for page change, got nil")
				}

				// Execute the command to simulate the page change
				if cmd != nil {
					msg := executeCommand(t, cmd)

					// Handle the message based on the view
					var newModel *model.Model
					switch m.CurrentView {
					case constants.ViewFunctionStatus:
						if funcMsg, ok := msg.(model.FunctionsPageMsg); ok {
							newModel = HandleFunctionStatusPagination(wrapper.Model, funcMsg)
						}
					case constants.ViewPipelineStatus:
						if pipeMsg, ok := msg.(model.PipelinesPageMsg); ok {
							newModel = HandlePipelineStatusPagination(wrapper.Model, pipeMsg)
						}
					case constants.ViewApprovals:
						if approvalMsg, ok := msg.(model.ApprovalsPageMsg); ok {
							newModel = HandleApprovalsPagination(wrapper.Model, approvalMsg)
						}
					}

					// Check the page number after handling the message
					if newModel != nil {
						if newModel.Pagination.CurrentPage != tt.expectedPageNumber {
							t.Errorf("Expected page number %d after handling message, got %d",
								tt.expectedPageNumber, newModel.Pagination.CurrentPage)
						}
						if newModel.Pagination.HasMorePages != tt.expectedHasMoreFlag {
							t.Errorf("Expected HasMorePages to be %v after handling message, got %v",
								tt.expectedHasMoreFlag, newModel.Pagination.HasMorePages)
						}
					}
				}
			} else {
				// If we don't expect a page change, the page should remain the same
				if wrapper.Model.Pagination.CurrentPage != initialPage {
					t.Errorf("Expected page to remain %d, got %d",
						initialPage, wrapper.Model.Pagination.CurrentPage)
				}
			}
		})
	}
}

// TestPaginationWithArrowKeys tests that arrow keys are correctly mapped to vim-style keys
func TestPaginationWithArrowKeys(t *testing.T) {
	// Create a model with pagination
	m := model.New()
	m.CurrentView = constants.ViewPipelineStatus
	m.Pagination.Type = model.PaginationTypeClientSide
	m.Pagination.CurrentPage = 2
	m.Pagination.HasMorePages = true
	m.Pagination.PageSize = 5
	m.Pagination.AllItems = createMockItems(15) // 3 pages total

	// Test that left arrow key is mapped to 'h'
	t.Run("Left arrow key mapped to 'h'", func(t *testing.T) {
		// Call with left arrow key
		result, _ := HandlePaginationKeyPress(m, constants.KeyArrowPreviousPage)

		// Call with 'h' key
		expectedResult, _ := HandlePaginationKeyPress(m, constants.KeyPreviousPage)

		// Verify both calls produce the same result
		wrapper1, ok1 := result.(ModelWrapper)
		wrapper2, ok2 := expectedResult.(ModelWrapper)

		if !ok1 || !ok2 {
			t.Fatalf("Expected HandlePaginationKeyPress to return a ModelWrapper")
		}

		// Compare the models
		if wrapper1.Model.Pagination.CurrentPage != wrapper2.Model.Pagination.CurrentPage {
			t.Errorf("Expected left arrow key to be mapped to '%s', but got different results",
				constants.KeyPreviousPage)
		}
	})

	// Test that right arrow key is mapped to 'l'
	t.Run("Right arrow key mapped to 'l'", func(t *testing.T) {
		// Call with right arrow key
		result, _ := HandlePaginationKeyPress(m, constants.KeyArrowNextPage)

		// Call with 'l' key
		expectedResult, _ := HandlePaginationKeyPress(m, constants.KeyNextPage)

		// Verify both calls produce the same result
		wrapper1, ok1 := result.(ModelWrapper)
		wrapper2, ok2 := expectedResult.(ModelWrapper)

		if !ok1 || !ok2 {
			t.Fatalf("Expected HandlePaginationKeyPress to return a ModelWrapper")
		}

		// Compare the models
		if wrapper1.Model.Pagination.CurrentPage != wrapper2.Model.Pagination.CurrentPage {
			t.Errorf("Expected right arrow key to be mapped to '%s', but got different results",
				constants.KeyNextPage)
		}
	})
}

// TestPaginationWithSearch tests that pagination works correctly with search filtering
func TestPaginationWithSearch(t *testing.T) {
	// Create a model with pagination and search
	m := model.New()
	m.CurrentView = constants.ViewPipelineStatus
	m.Pagination.Type = model.PaginationTypeClientSide
	m.Pagination.CurrentPage = 1
	m.Pagination.HasMorePages = true
	m.Pagination.PageSize = 5

	// Create 20 items with different names
	for i := 0; i < 20; i++ {
		pipeline := model.PipelineStatus{
			Name: "pipeline-" + string(rune('A'+i%3)) + "-" + string(rune('0'+i)),
		}
		m.Pagination.AllItems = append(m.Pagination.AllItems, pipeline)
	}

	// Filter to only items with 'A' in the name (should be ~7 items)
	m.Search.Query = "A"
	m.Search.IsActive = true

	// Manually filter items to simulate search
	for _, item := range m.Pagination.AllItems {
		if pipeline, ok := item.(model.PipelineStatus); ok {
			if pipeline.Name[9] == 'A' {
				m.Search.FilteredItems = append(m.Search.FilteredItems, pipeline)
			}
		}
	}

	// Test next page with filtered items
	t.Run("Next page with filtered items", func(t *testing.T) {
		// Call the function with next page key
		_, cmd := HandlePaginationKeyPress(m, constants.KeyNextPage)

		// Execute the command
		if cmd == nil {
			t.Fatalf("Expected a command for page change, got nil")
		}

		msg := executeCommand(t, cmd)

		// Handle the message
		if pipeMsg, ok := msg.(model.PipelinesPageMsg); ok {
			newModel := HandlePipelineStatusPagination(m, pipeMsg)

			// Check that we're on page 2
			if newModel.Pagination.CurrentPage != 2 {
				t.Errorf("Expected to be on page 2, got %d", newModel.Pagination.CurrentPage)
			}

			// Check that we're still using filtered items
			if len(newModel.Pipelines) == 0 {
				t.Errorf("Expected pipelines to be populated with filtered items")
			}

			// Check that all items on this page match the filter
			for _, pipeline := range newModel.Pipelines {
				if pipeline.Name[9] != 'A' {
					t.Errorf("Expected all pipelines to match filter 'A', got %s", pipeline.Name)
				}
			}
		} else {
			t.Errorf("Expected PipelinesPageMsg, got %T", msg)
		}
	})
}

// TestPaginationReset tests that pagination state is reset when navigating away
func TestPaginationReset(t *testing.T) {
	// Create a model with pagination
	m := model.New()
	m.CurrentView = constants.ViewPipelineStatus
	m.Pagination.Type = model.PaginationTypeClientSide
	m.Pagination.CurrentPage = 2
	m.Pagination.HasMorePages = true
	m.Pagination.PageSize = 5
	m.Pagination.AllItems = createMockItems(15)
	m.Pagination.TotalItems = 15

	// Also set up search state
	m.Search.IsActive = true
	m.Search.Query = "test"
	m.Search.FilteredItems = createMockItems(5)

	// Navigate back
	newModel := NavigateBack(m)

	// Check that pagination state is reset
	if newModel.Pagination.Type != model.PaginationTypeNone {
		t.Errorf("Expected pagination type to be reset to None, got %v", newModel.Pagination.Type)
	}

	if newModel.Pagination.CurrentPage != 1 {
		t.Errorf("Expected current page to be reset to 1, got %d", newModel.Pagination.CurrentPage)
	}

	if newModel.Pagination.HasMorePages {
		t.Errorf("Expected HasMorePages to be reset to false, got true")
	}

	if len(newModel.Pagination.AllItems) != 0 {
		t.Errorf("Expected AllItems to be reset to empty, got %d items", len(newModel.Pagination.AllItems))
	}

	if len(newModel.Pagination.FilteredItems) != 0 {
		t.Errorf("Expected FilteredItems to be reset to empty, got %d items", len(newModel.Pagination.FilteredItems))
	}

	if newModel.Pagination.TotalItems != -1 {
		t.Errorf("Expected TotalItems to be reset to -1, got %d", newModel.Pagination.TotalItems)
	}

	// Check that search state is also reset
	if newModel.Search.IsActive {
		t.Errorf("Expected search to be deactivated, got active")
	}

	if newModel.Search.Query != "" {
		t.Errorf("Expected search query to be reset to empty, got '%s'", newModel.Search.Query)
	}

	if len(newModel.Search.FilteredItems) != 0 {
		t.Errorf("Expected search filtered items to be reset to empty, got %d items", len(newModel.Search.FilteredItems))
	}
}

// TestPaginationBoundaries tests that pagination stops at boundaries rather than wrapping around
func TestPaginationBoundaries(t *testing.T) {
	// Create a model with pagination
	m := model.New()
	m.CurrentView = constants.ViewPipelineStatus
	m.Pagination.Type = model.PaginationTypeClientSide
	m.Pagination.PageSize = 5
	m.Pagination.AllItems = createMockItems(15) // 3 pages total

	// Test first page boundary
	t.Run("First page boundary", func(t *testing.T) {
		// Set up model at first page
		m.Pagination.CurrentPage = 1
		m.Pagination.HasMorePages = true

		// Try to go to previous page
		result, cmd := HandlePaginationKeyPress(m, constants.KeyPreviousPage)

		// Check that we didn't get a command (no page change)
		if cmd != nil {
			t.Errorf("Expected no command when trying to go before first page, got a command")
		}

		// Check that we're still on page 1
		wrapper, ok := result.(ModelWrapper)
		if !ok {
			t.Fatalf("Expected HandlePaginationKeyPress to return a ModelWrapper, got %T", result)
		}

		if wrapper.Model.Pagination.CurrentPage != 1 {
			t.Errorf("Expected to stay on page 1, got page %d", wrapper.Model.Pagination.CurrentPage)
		}
	})

	// Test last page boundary
	t.Run("Last page boundary", func(t *testing.T) {
		// Set up model at last page
		m.Pagination.CurrentPage = 3
		m.Pagination.HasMorePages = false

		// Try to go to next page
		result, cmd := HandlePaginationKeyPress(m, constants.KeyNextPage)

		// Check that we didn't get a command (no page change)
		if cmd != nil {
			t.Errorf("Expected no command when trying to go past last page, got a command")
		}

		// Check that we're still on page 3
		wrapper, ok := result.(ModelWrapper)
		if !ok {
			t.Fatalf("Expected HandlePaginationKeyPress to return a ModelWrapper, got %T", result)
		}

		if wrapper.Model.Pagination.CurrentPage != 3 {
			t.Errorf("Expected to stay on page 3, got page %d", wrapper.Model.Pagination.CurrentPage)
		}
	})
}

// Helper function to create mock items for pagination tests
func createMockItems(count int) []interface{} {
	items := make([]interface{}, 0, count)
	for i := 0; i < count; i++ {
		pipeline := model.PipelineStatus{
			Name: "pipeline-" + string(rune('0'+i)),
		}
		items = append(items, pipeline)
	}
	return items
}

// Helper function to execute a command and return the message
func executeCommand(t *testing.T, cmd tea.Cmd) tea.Msg {
	msg := cmd()
	if msg == nil {
		t.Fatalf("Command returned nil message")
	}
	return msg
}
