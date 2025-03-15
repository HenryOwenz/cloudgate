package ui

import (
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// TestArrowKeyPagination tests that arrow keys are correctly handled for pagination
func TestArrowKeyPagination(t *testing.T) {
	// Create a model with pagination
	m := Model{
		core: &model.Model{
			CurrentView: constants.ViewPipelineStatus,
			Pagination: model.Pagination{
				Type:         model.PaginationTypeClientSide,
				CurrentPage:  1,
				HasMorePages: true,
				PageSize:     5,
				AllItems:     createMockItems(10),
			},
		},
	}
	view.UpdateTableForView(m.core)

	tests := []struct {
		name           string
		key            tea.KeyMsg
		expectedPage   int
		expectedChange bool
	}{
		{
			name:           "Right arrow key for next page",
			key:            tea.KeyMsg{Type: tea.KeyRight},
			expectedPage:   2,
			expectedChange: true,
		},
		{
			name:           "Left arrow key for previous page (no change on page 1)",
			key:            tea.KeyMsg{Type: tea.KeyLeft},
			expectedPage:   1,
			expectedChange: false,
		},
		{
			name:           "l key for next page",
			key:            tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(constants.KeyNextPage)},
			expectedPage:   2,
			expectedChange: true,
		},
		{
			name:           "h key for previous page (no change on page 1)",
			key:            tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(constants.KeyPreviousPage)},
			expectedPage:   1,
			expectedChange: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the model for each test
			testModel := Model{
				core: &model.Model{
					CurrentView: constants.ViewPipelineStatus,
					Pagination: model.Pagination{
						Type:         model.PaginationTypeClientSide,
						CurrentPage:  1,
						HasMorePages: true,
						PageSize:     5,
						AllItems:     createMockItems(10),
					},
				},
			}
			view.UpdateTableForView(testModel.core)

			// Initial page
			initialPage := testModel.core.Pagination.CurrentPage

			// Call the Update method with the key message
			newModel, cmd := testModel.Update(tt.key)

			// Check if we got a command (indicates a page change request)
			if tt.expectedChange {
				if cmd == nil {
					t.Errorf("Expected a command for page change, got nil")
				}
			} else {
				// If we don't expect a change, the page should remain the same
				updatedModel, ok := newModel.(Model)
				if !ok {
					t.Fatalf("Expected updated model to be of type Model, got %T", newModel)
				}

				if updatedModel.core.Pagination.CurrentPage != initialPage {
					t.Errorf("Expected page to remain %d, got %d",
						initialPage, updatedModel.core.Pagination.CurrentPage)
				}
			}
		})
	}
}

// TestSearchKeyHandling tests that the search key is correctly handled
func TestSearchKeyHandling(t *testing.T) {
	// Create a model with pagination
	m := Model{
		core: &model.Model{
			CurrentView: constants.ViewPipelineStatus,
			Pagination: model.Pagination{
				Type:         model.PaginationTypeClientSide,
				CurrentPage:  1,
				HasMorePages: true,
				PageSize:     5,
				AllItems:     createMockItems(10),
			},
		},
	}
	view.UpdateTableForView(m.core)

	// Create a search key message
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(constants.KeySearch)}

	// Call the Update method with the search key message
	newModel, _ := m.Update(keyMsg)

	// Verify the model was updated
	updatedModel, ok := newModel.(Model)
	if !ok {
		t.Fatalf("Expected updated model to be of type Model, got %T", newModel)
	}

	// Check that search is active
	if !updatedModel.core.Search.IsActive {
		t.Errorf("Expected search to be active in the updated model, got inactive")
	}
}

// TestSearchInputHandling tests that search input is correctly handled
func TestSearchInputHandling(t *testing.T) {
	// Test adding a character to the search query
	t.Run("Add character to search query", func(t *testing.T) {
		// Create a model with active search
		m := Model{
			core: &model.Model{
				CurrentView: constants.ViewPipelineStatus,
				Pagination: model.Pagination{
					Type:         model.PaginationTypeClientSide,
					CurrentPage:  1,
					HasMorePages: true,
					PageSize:     5,
					AllItems:     createMockItems(10),
				},
				Search: model.SearchState{
					IsActive: true,
					Query:    "",
				},
			},
		}
		view.UpdateTableForView(m.core)

		// Create a key message for the letter 'a'
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}

		// Call the Update method with the key message
		newModel, _ := m.Update(keyMsg)

		// Verify the model was updated
		updatedModel, ok := newModel.(Model)
		if !ok {
			t.Fatalf("Expected updated model to be of type Model, got %T", newModel)
		}

		// Check that the query was updated
		if updatedModel.core.Search.Query != "a" {
			t.Errorf("Expected search query to be 'a' in the updated model, got '%s'",
				updatedModel.core.Search.Query)
		}
	})

	// Test handling backspace in search query
	t.Run("Backspace in search query", func(t *testing.T) {
		// Create a model with active search and existing query
		m := Model{
			core: &model.Model{
				CurrentView: constants.ViewPipelineStatus,
				Pagination: model.Pagination{
					Type:         model.PaginationTypeClientSide,
					CurrentPage:  1,
					HasMorePages: true,
					PageSize:     5,
					AllItems:     createMockItems(10),
				},
				Search: model.SearchState{
					IsActive: true,
					Query:    "test",
				},
			},
		}
		view.UpdateTableForView(m.core)

		// Create a backspace key message
		keyMsg := tea.KeyMsg{Type: tea.KeyBackspace}

		// Call the Update method with the key message
		newModel, _ := m.Update(keyMsg)

		// Verify the model was updated
		updatedModel, ok := newModel.(Model)
		if !ok {
			t.Fatalf("Expected updated model to be of type Model, got %T", newModel)
		}

		// Check that the query was updated
		if updatedModel.core.Search.Query != "tes" {
			t.Errorf("Expected search query to be 'tes' in the updated model, got '%s'",
				updatedModel.core.Search.Query)
		}
	})

	// Test exiting search with Escape key
	t.Run("Exit search with Escape", func(t *testing.T) {
		// Create a model with active search
		m := Model{
			core: &model.Model{
				CurrentView: constants.ViewPipelineStatus,
				Pagination: model.Pagination{
					Type:         model.PaginationTypeClientSide,
					CurrentPage:  1,
					HasMorePages: true,
					PageSize:     5,
					AllItems:     createMockItems(10),
				},
				Search: model.SearchState{
					IsActive: true,
					Query:    "test",
				},
			},
		}
		view.UpdateTableForView(m.core)

		// Create an Escape key message
		keyMsg := tea.KeyMsg{Type: tea.KeyEsc}

		// Call the Update method with the key message
		newModel, _ := m.Update(keyMsg)

		// Verify the model was updated
		updatedModel, ok := newModel.(Model)
		if !ok {
			t.Fatalf("Expected updated model to be of type Model, got %T", newModel)
		}

		// Check that search is inactive
		if updatedModel.core.Search.IsActive {
			t.Errorf("Expected search to be inactive in the updated model, got active")
		}

		// Check that the query was cleared
		if updatedModel.core.Search.Query != "" {
			t.Errorf("Expected search query to be empty in the updated model, got '%s'",
				updatedModel.core.Search.Query)
		}
	})
}

// TestFunctionalKeysInTextInputMode tests that functional keys are passed through in text input mode
func TestFunctionalKeysInTextInputMode(t *testing.T) {
	// Create a model with text input active
	m := Model{
		core: &model.Model{
			CurrentView: constants.ViewConfirmation,
			ManualInput: true,
			TextInput:   textinput.New(),
		},
	}

	// Test pagination keys in text input mode
	t.Run("Pagination keys in text input mode", func(t *testing.T) {
		// Create a key message for 'h' (previous page)
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(constants.KeyPreviousPage)}

		// Call the Update method with the key message
		_, cmd := m.Update(keyMsg)

		// Verify the key was passed through (no pagination command)
		if cmd != nil {
			t.Errorf("Expected pagination key to be passed through in text input mode, got a command")
		}

		// Create a key message for 'l' (next page)
		keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(constants.KeyNextPage)}

		// Call the Update method with the key message
		_, cmd = m.Update(keyMsg)

		// Verify the key was passed through (no pagination command)
		if cmd != nil {
			t.Errorf("Expected pagination key to be passed through in text input mode, got a command")
		}
	})

	// Test search key in text input mode
	t.Run("Search key in text input mode", func(t *testing.T) {
		// Create a key message for '/' (search)
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(constants.KeySearch)}

		// Call the Update method with the key message
		newModel, _ := m.Update(keyMsg)

		// Verify the key was passed through (no search activation)
		updatedModel, ok := newModel.(Model)
		if !ok {
			t.Fatalf("Expected updated model to be of type Model, got %T", newModel)
		}

		if updatedModel.core.Search.IsActive {
			t.Errorf("Expected search key to be passed through in text input mode, got search activated")
		}
	})
}

// Helper function to create mock items for tests
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
