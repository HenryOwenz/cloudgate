# Lazy String Search Implementation

## Overview

This document outlines the implementation of a simple, efficient lazy string search feature for the CloudGate application. The search functionality allows users to filter items in paginated views using a case-insensitive substring matching approach.

## Design Principles

1. **Simplicity**: Straightforward implementation with minimal complexity
2. **Efficiency**: Optimized for performance with potentially large datasets
3. **Case-insensitivity**: Search is case-insensitive by default
4. **Substring matching**: Matches items that contain the search query as a substring
5. **Modularity**: Search logic is separate from UI logic
6. **Extensibility**: Works with different item types (Functions, Pipelines, Approvals)

## User Experience

- Press `/` to activate search mode in any paginated view
- Type search query to filter items in real-time
- Press `Enter` to confirm search and exit search mode
- Press `Esc` to cancel search and restore original view
- Empty search query shows all items

## Implementation Details

### 1. Model Updates

The model is extended to support search functionality:

```go
// Add to model.Model struct
type Model struct {
    // ... existing fields
    
    // Search state
    Search struct {
        IsActive      bool         // Whether search mode is active
        Query         string       // Current search query
        FilteredItems []interface{} // Items that match the search query
    }
}
```

### 2. Search Logic

The search logic is implemented in a new file `internal/ui/update/search_handlers.go`:

```go
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
    
    // Reset pagination to show all items
    newModel.Pagination.FilteredItems = nil
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
        newModel.Pagination.FilteredItems = nil
        newModel.Pagination.CurrentPage = 1
        refreshTableWithAllItems(newModel)
        return newModel
    }
    
    // Filter items based on the query
    filteredItems := filterItemsByQuery(newModel.Pagination.AllItems, query)
    newModel.Pagination.FilteredItems = filteredItems
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
            if v.Status != "" {
                searchText += " " + strings.ToLower(v.Status)
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
```

### 3. UI Integration

The UI is updated to handle search-related key presses:

```go
// Add to ui.go in the Update function's key handling section
case "/":
    // Only activate search in paginated views with data
    if view.IsPaginatedView(m.core.CurrentView) && len(m.core.Pagination.AllItems) > 0 {
        newModel := m.Clone()
        newModel.core = update.ActivateSearch(newModel.core)
        return newModel, nil
    }
    return m, nil

// Add handling for when search is active
if m.core.Search.IsActive {
    switch msg.String() {
    case "esc", "ctrl+c":
        // Exit search mode
        newModel := m.Clone()
        newModel.core = update.DeactivateSearch(newModel.core)
        return newModel, nil
    case "enter":
        // Confirm search and exit search mode
        newModel := m.Clone()
        newModel.core.Search.IsActive = false
        return newModel, nil
    case "backspace":
        // Handle backspace in search query
        if len(m.core.Search.Query) > 0 {
            newModel := m.Clone()
            newModel.core.Search.Query = newModel.core.Search.Query[:len(newModel.core.Search.Query)-1]
            newModel.core = update.UpdateSearchQuery(newModel.core, newModel.core.Search.Query)
            return newModel, nil
        }
        return m, nil
    default:
        // Add character to search query if it's a printable character
        r := msg.Runes
        if len(r) == 1 && isPrintableChar(r[0]) {
            newModel := m.Clone()
            newModel.core.Search.Query += string(r)
            newModel.core = update.UpdateSearchQuery(newModel.core, newModel.core.Search.Query)
            return newModel, nil
        }
        return m, nil
    }
}
```

### 4. View Updates

The view is updated to display the search input:

```go
// Add to view.go in the renderHelpText function
func renderHelpText(m *model.Model) string {
    // If search is active, show search help text
    if m.Search.IsActive {
        searchPrompt := "Search: "
        searchText := m.Search.Query
        if len(searchText) == 0 {
            searchText = "_" // Show cursor placeholder when empty
        } else {
            searchText += "_" // Add cursor at the end
        }
        
        // Style the search prompt and text
        styledSearchPrompt := m.Styles.SearchPrompt.Render(searchPrompt)
        styledSearchText := m.Styles.SearchText.Render(searchText)
        
        // Add help text
        helpText := m.Styles.Help.Render(" | Enter: confirm • Esc: cancel")
        
        return lipgloss.JoinHorizontal(lipgloss.Left, styledSearchPrompt, styledSearchText, helpText)
    }
    
    // Add search hint to regular help text for paginated views
    helpText := getHelpText(m)
    if view.IsPaginatedView(m.CurrentView) && len(m.Pagination.AllItems) > 0 {
        helpText += " • /: search"
    }
    
    return m.Styles.Help.Render(helpText)
}
```

## Performance Optimizations

The implementation includes several optimizations for better performance:

1. **Lazy Evaluation**: Items are only filtered when the query changes
2. **Early Termination**: The function returns early if the query is empty
3. **Efficient String Operations**: Uses `strings.Contains` which is optimized in Go
4. **Pagination Awareness**: Only processes items for the current page when displaying

## Integration with Pagination

The search functionality integrates with the existing pagination system:

```go
// Modify pagination handlers to use FilteredItems when available
func FetchNextPage(m *model.Model) tea.Cmd {
    // Use FilteredItems if available, otherwise use AllItems
    sourceItems := m.Pagination.AllItems
    if len(m.Pagination.FilteredItems) > 0 {
        sourceItems = m.Pagination.FilteredItems
    }
    
    // Rest of the function remains the same, but uses sourceItems
    // ...
}
```

## Testing Strategy

1. **Unit Tests**: Test the search logic with different item types and queries
2. **Edge Cases**: Test empty queries, case sensitivity, and special characters
3. **Performance Tests**: Benchmark search performance with large datasets
4. **Integration Tests**: Test the search UI with the rest of the application

## Future Enhancements

Potential future enhancements to the search functionality:

1. **Fuzzy Matching**: Implement fuzzy matching for more forgiving search
2. **Highlighting**: Highlight matching text in search results
3. **Search History**: Add support for search history with up/down arrows
4. **Field-Specific Search**: Allow searching in specific fields (e.g., name:value)
5. **Regular Expressions**: Support for regex-based search

## Conclusion

This lazy string search implementation provides a simple, efficient way to filter items in paginated views. It enhances the user experience by allowing quick access to specific items without requiring complex external dependencies. 