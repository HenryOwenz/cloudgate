package update

import (
	"strings"
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
)

// TestActivateSearch tests the activation of search mode
func TestActivateSearch(t *testing.T) {
	// Create a model with pagination
	m := model.New()
	m.CurrentView = constants.ViewPipelineStatus
	m.Pagination.Type = model.PaginationTypeClientSide
	m.Pagination.AllItems = createTestItems(10)

	// Activate search
	newModel := ActivateSearch(m)

	// Check that search is activated
	if !newModel.Search.IsActive {
		t.Errorf("Expected search to be activated, got inactive")
	}

	// Check that search query is empty
	if newModel.Search.Query != "" {
		t.Errorf("Expected search query to be empty, got '%s'", newModel.Search.Query)
	}
}

// TestDeactivateSearch tests the deactivation of search mode
func TestDeactivateSearch(t *testing.T) {
	// Create a model with active search
	m := model.New()
	m.CurrentView = constants.ViewPipelineStatus
	m.Pagination.Type = model.PaginationTypeClientSide
	m.Pagination.AllItems = createTestItems(10)
	m.Search.IsActive = true
	m.Search.Query = "test"
	m.Search.FilteredItems = createTestItems(3)

	// Deactivate search
	newModel := DeactivateSearch(m)

	// Check that search is deactivated
	if newModel.Search.IsActive {
		t.Errorf("Expected search to be deactivated, got active")
	}

	// Check that search query is cleared
	if newModel.Search.Query != "" {
		t.Errorf("Expected search query to be cleared, got '%s'", newModel.Search.Query)
	}

	// Check that filtered items are cleared
	if len(newModel.Search.FilteredItems) != 0 {
		t.Errorf("Expected filtered items to be cleared, got %d items", len(newModel.Search.FilteredItems))
	}

	// Check that pagination is reset to page 1
	if newModel.Pagination.CurrentPage != 1 {
		t.Errorf("Expected pagination to be reset to page 1, got page %d", newModel.Pagination.CurrentPage)
	}
}

// TestUpdateSearchQuery tests updating the search query and filtering items
func TestUpdateSearchQuery(t *testing.T) {
	tests := []struct {
		name           string
		setupModel     func() *model.Model
		query          string
		expectedCount  int
		expectedFilter func(string) bool
	}{
		{
			name: "Empty query shows all items",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewPipelineStatus
				m.Pagination.Type = model.PaginationTypeClientSide

				// Create test items with different names
				for i := 0; i < 10; i++ {
					pipeline := model.PipelineStatus{
						Name: "pipeline-" + string(rune('A'+i%3)) + "-" + string(rune('0'+i)),
					}
					m.Pagination.AllItems = append(m.Pagination.AllItems, pipeline)
				}

				return m
			},
			query:         "",
			expectedCount: 0, // No filtered items when query is empty
			expectedFilter: func(name string) bool {
				return true // All items should pass
			},
		},
		{
			name: "Filter by letter A",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewPipelineStatus
				m.Pagination.Type = model.PaginationTypeClientSide

				// Create test items with different names
				for i := 0; i < 10; i++ {
					pipeline := model.PipelineStatus{
						Name: "pipeline-" + string(rune('A'+i%3)) + "-" + string(rune('0'+i)),
					}
					m.Pagination.AllItems = append(m.Pagination.AllItems, pipeline)
				}

				return m
			},
			query:         "A",
			expectedCount: 4, // Should match items with 'A' in the name
			expectedFilter: func(name string) bool {
				return strings.Contains(strings.ToLower(name), "a")
			},
		},
		{
			name: "Filter by number",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewPipelineStatus
				m.Pagination.Type = model.PaginationTypeClientSide

				// Create test items with different names
				for i := 0; i < 10; i++ {
					pipeline := model.PipelineStatus{
						Name: "pipeline-" + string(rune('A'+i%3)) + "-" + string(rune('0'+i)),
					}
					m.Pagination.AllItems = append(m.Pagination.AllItems, pipeline)
				}

				return m
			},
			query:         "5",
			expectedCount: 1, // Should match only pipeline-X-5
			expectedFilter: func(name string) bool {
				return strings.Contains(name, "5")
			},
		},
		{
			name: "No matches",
			setupModel: func() *model.Model {
				m := model.New()
				m.CurrentView = constants.ViewPipelineStatus
				m.Pagination.Type = model.PaginationTypeClientSide

				// Create test items with different names
				for i := 0; i < 10; i++ {
					pipeline := model.PipelineStatus{
						Name: "pipeline-" + string(rune('A'+i%3)) + "-" + string(rune('0'+i)),
					}
					m.Pagination.AllItems = append(m.Pagination.AllItems, pipeline)
				}

				return m
			},
			query:         "XYZ",
			expectedCount: 0, // Should match no items
			expectedFilter: func(name string) bool {
				return false // No items should pass
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup the model
			m := tt.setupModel()

			// Update search query
			newModel := UpdateSearchQuery(m, tt.query)

			// Check that search query is updated
			if newModel.Search.Query != tt.query {
				t.Errorf("Expected search query to be '%s', got '%s'", tt.query, newModel.Search.Query)
			}

			// Check filtered items count
			if tt.query == "" {
				// Empty query should not filter items
				if len(newModel.Search.FilteredItems) != tt.expectedCount {
					t.Errorf("Expected %d filtered items for empty query, got %d",
						tt.expectedCount, len(newModel.Search.FilteredItems))
				}
			} else {
				// Check that filtered items match the expected count
				if len(newModel.Search.FilteredItems) != tt.expectedCount {
					t.Errorf("Expected %d filtered items, got %d",
						tt.expectedCount, len(newModel.Search.FilteredItems))
				}

				// Check that all filtered items match the filter
				for _, item := range newModel.Search.FilteredItems {
					if pipeline, ok := item.(model.PipelineStatus); ok {
						if !tt.expectedFilter(pipeline.Name) {
							t.Errorf("Item '%s' should not have passed the filter", pipeline.Name)
						}
					}
				}
			}

			// Check that pagination is reset to page 1
			if newModel.Pagination.CurrentPage != 1 {
				t.Errorf("Expected pagination to be reset to page 1, got page %d",
					newModel.Pagination.CurrentPage)
			}
		})
	}
}

// TestFilterItemsByQuery tests the filtering of items based on a query
func TestFilterItemsByQuery(t *testing.T) {
	// Create test items
	items := []interface{}{
		model.FunctionStatus{Name: "function-A", Runtime: "nodejs14.x"},
		model.FunctionStatus{Name: "function-B", Runtime: "python3.8"},
		model.FunctionStatus{Name: "function-C", Runtime: "java11"},
		model.PipelineStatus{
			Name: "pipeline-A",
			Stages: []cloud.StageStatus{
				{Name: "Build", Status: "Succeeded"},
				{Name: "Test", Status: "Failed"},
			},
		},
		model.PipelineStatus{
			Name: "pipeline-B",
			Stages: []cloud.StageStatus{
				{Name: "Build", Status: "Succeeded"},
				{Name: "Deploy", Status: "InProgress"},
			},
		},
		model.ApprovalAction{
			PipelineName: "pipeline-C",
			StageName:    "Approval",
			ActionName:   "ManualApproval",
		},
	}

	tests := []struct {
		name          string
		query         string
		expectedCount int
		expectedItems []string // Names of expected items
	}{
		{
			name:          "Empty query returns all items",
			query:         "",
			expectedCount: 6,
			expectedItems: []string{
				"function-A", "function-B", "function-C",
				"pipeline-A", "pipeline-B", "pipeline-C",
			},
		},
		{
			name:          "Filter by 'function'",
			query:         "function",
			expectedCount: 3,
			expectedItems: []string{"function-A", "function-B", "function-C"},
		},
		{
			name:          "Filter by 'pipeline'",
			query:         "pipeline",
			expectedCount: 3,
			expectedItems: []string{"pipeline-A", "pipeline-B", "pipeline-C"},
		},
		{
			name:          "Filter by 'A'",
			query:         "a",
			expectedCount: 4, // function-A, pipeline-A, ManualApproval, Approval
			expectedItems: []string{"function-A", "pipeline-A", "pipeline-C"},
		},
		{
			name:          "Filter by runtime",
			query:         "python",
			expectedCount: 1,
			expectedItems: []string{"function-B"},
		},
		{
			name:          "Filter by stage name",
			query:         "deploy",
			expectedCount: 1,
			expectedItems: []string{"pipeline-B"},
		},
		{
			name:          "Filter by stage status",
			query:         "failed",
			expectedCount: 1,
			expectedItems: []string{"pipeline-A"},
		},
		{
			name:          "No matches",
			query:         "xyz123",
			expectedCount: 0,
			expectedItems: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Filter items
			filteredItems := filterItemsByQuery(items, tt.query)

			// Check count
			if len(filteredItems) != tt.expectedCount {
				t.Errorf("Expected %d filtered items, got %d", tt.expectedCount, len(filteredItems))
			}

			// Check that all expected items are present
			for _, expectedName := range tt.expectedItems {
				found := false
				for _, item := range filteredItems {
					var name string
					switch v := item.(type) {
					case model.FunctionStatus:
						name = v.Name
					case model.PipelineStatus:
						name = v.Name
					case model.ApprovalAction:
						name = v.PipelineName
					}

					if name == expectedName {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("Expected item '%s' not found in filtered results", expectedName)
				}
			}
		})
	}
}

// TestSearchWithPagination tests the integration of search with pagination
func TestSearchWithPagination(t *testing.T) {
	// Create a model with pagination
	m := model.New()
	m.CurrentView = constants.ViewPipelineStatus
	m.Pagination.Type = model.PaginationTypeClientSide
	m.Pagination.CurrentPage = 1
	m.Pagination.PageSize = 5

	// Create 20 items with different names
	for i := 0; i < 20; i++ {
		pipeline := model.PipelineStatus{
			Name: "pipeline-" + string(rune('A'+i%3)) + "-" + string(rune('0'+i)),
		}
		m.Pagination.AllItems = append(m.Pagination.AllItems, pipeline)
	}

	// Activate search
	m = ActivateSearch(m)

	// Update search query to filter items
	m = UpdateSearchQuery(m, "A")

	// Check that search is active
	if !m.Search.IsActive {
		t.Errorf("Expected search to be active, got inactive")
	}

	// Check that search query is set
	if m.Search.Query != "A" {
		t.Errorf("Expected search query to be 'A', got '%s'", m.Search.Query)
	}

	// Check that filtered items are populated
	if len(m.Search.FilteredItems) == 0 {
		t.Errorf("Expected filtered items to be populated, got empty")
	}

	// Check that all filtered items match the query
	for _, item := range m.Search.FilteredItems {
		if pipeline, ok := item.(model.PipelineStatus); ok {
			if !strings.Contains(strings.ToLower(pipeline.Name), "a") {
				t.Errorf("Item '%s' should not have passed the filter", pipeline.Name)
			}
		}
	}

	// Check that pagination is reset to page 1
	if m.Pagination.CurrentPage != 1 {
		t.Errorf("Expected pagination to be reset to page 1, got page %d", m.Pagination.CurrentPage)
	}

	// Check that the table is updated with filtered items
	if len(m.Pipelines) == 0 {
		t.Errorf("Expected pipelines to be populated with filtered items")
	}

	// Check that all items in the table match the filter
	for _, pipeline := range m.Pipelines {
		if !strings.Contains(strings.ToLower(pipeline.Name), "a") {
			t.Errorf("Table item '%s' should not have passed the filter", pipeline.Name)
		}
	}
}

// TestIsPrintableChar tests the IsPrintableChar function
func TestIsPrintableChar(t *testing.T) {
	tests := []struct {
		name     string
		char     rune
		expected bool
	}{
		{"Lowercase letter", 'a', true},
		{"Uppercase letter", 'Z', true},
		{"Number", '5', true},
		{"Symbol", '@', true},
		{"Space", ' ', true},
		{"Tab", '\t', false},
		{"Newline", '\n', false},
		{"Null", 0, false},
		{"Delete", 127, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPrintableChar(tt.char)
			if result != tt.expected {
				t.Errorf("Expected IsPrintableChar('%c') to be %v, got %v",
					tt.char, tt.expected, result)
			}
		})
	}
}

// Helper function to create test items for search tests
func createTestItems(count int) []interface{} {
	items := make([]interface{}, 0, count)
	for i := 0; i < count; i++ {
		pipeline := model.PipelineStatus{
			Name: "pipeline-" + string(rune('A'+i%3)) + "-" + string(rune('0'+i)),
		}
		items = append(items, pipeline)
	}
	return items
}
