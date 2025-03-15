# Cloudgate Pagination Implementation Plan

This document outlines the plan for implementing pagination in Cloudgate for AWS resource listing operations, specifically for Lambda Functions, CodePipeline, and Pipeline Approvals views.

## Current Implementation Status

As of now, we have implemented client-side pagination with the following features:
- Core pagination infrastructure with models and state management
- UI components for pagination display and navigation
- Keyboard handlers for pagination ('h' for previous, 'l' for next)
- Client-side pagination for Lambda Functions, CodePipeline, and Pipeline Approvals

The current implementation works by:
1. Fetching all data at once when a view is first loaded
2. Storing all items in memory in the `AllItems` cache
3. Displaying pages by slicing this cache based on the current page and page size
4. Providing navigation between pages using 'h' and 'l' keys

## 1. Core Pagination Infrastructure

### 1.1 Pagination Models (Implemented)

We've added pagination models to the `internal/ui/model/types.go` file:

```go
// PaginationType represents the type of pagination to use
type PaginationType string

const (
    PaginationTypeNone       PaginationType = "none"
    PaginationTypeClientSide PaginationType = "client-side"
    PaginationTypeAPILevel   PaginationType = "api-level"
)

// Pagination represents a generic pagination state
type Pagination struct {
    Type            PaginationType
    CurrentPage     int
    PageSize        int
    TotalItems      int64
    HasMorePages    bool
    NextPageToken   string
    IsLoading       bool
    
    // For client-side pagination
    AllItems        []interface{}
    FilteredItems   []interface{}
}
```

### 1.2 Pagination Message Types (Implemented)

We've added message types for pagination in `internal/ui/model/types.go`:

```go
// FunctionsPageMsg represents a message containing a page of functions
type FunctionsPageMsg struct {
    Functions     []FunctionStatus
    NextPageToken string
    HasMorePages  bool
}

// PipelinesPageMsg represents a message containing a page of pipelines
type PipelinesPageMsg struct {
    Pipelines     []PipelineStatus
    NextPageToken string
    HasMorePages  bool
}

// ApprovalsPageMsg represents a message containing a page of approvals
type ApprovalsPageMsg struct {
    Approvals     []ApprovalAction
    NextPageToken string
    HasMorePages  bool
}
```

### 1.3 Pagination Handlers (Implemented)

We've implemented pagination handlers in `internal/ui/update/pagination_handlers.go`:

```go
// HandlePaginationKeyPress handles pagination key presses (h for previous, l for next)
func HandlePaginationKeyPress(m *model.Model, key string) (tea.Model, tea.Cmd) {
    // Implementation details...
}

// FetchNextPage fetches the next page based on the current view
func FetchNextPage(m *model.Model) tea.Cmd {
    // Implementation details...
}

// FetchPreviousPage fetches the previous page based on the current view
func FetchPreviousPage(m *model.Model) tea.Cmd {
    // Implementation details...
}
```

## 2. UI Components

### 2.1 Pagination View Identification (Implemented)

We've added a function to identify paginated views in `internal/ui/view/view.go`:

```go
// IsPaginatedView returns true if the view supports pagination
func IsPaginatedView(view constants.View) bool {
    return view == constants.ViewFunctionStatus ||
        view == constants.ViewPipelineStatus ||
        view == constants.ViewApprovals
}
```

### 2.2 Pagination Help Text (Implemented)

We've updated the `renderHelpText` function in `internal/ui/view/view.go` to include pagination information:

```go
// renderHelpText renders the help text
func renderHelpText(m *model.Model) string {
    baseHelpText := getHelpText(m)

    // Add pagination help text if the view supports pagination
    if IsPaginatedView(m.CurrentView) && m.Pagination.Type != model.PaginationTypeNone {
        // Calculate total pages
        totalPages := 1
        if m.Pagination.TotalItems > 0 && m.Pagination.PageSize > 0 {
            totalPages = int((m.Pagination.TotalItems + int64(m.Pagination.PageSize) - 1) / int64(m.Pagination.PageSize))
        }

        // Add pagination info and controls
        paginationInfo := fmt.Sprintf("Page %d of %d", m.Pagination.CurrentPage, totalPages)

        // Add item count if available
        if m.Pagination.TotalItems > 0 {
            paginationInfo += fmt.Sprintf(" (%d items)", m.Pagination.TotalItems)
        }

        // Create pagination controls based on current state
        var paginationControls string
        if m.Pagination.CurrentPage > 1 {
            paginationControls += "h: prev page"
        }
        if m.Pagination.HasMorePages {
            if paginationControls != "" {
                paginationControls += " • "
            }
            paginationControls += "l: next page"
        }

        // Combine base help text with pagination info
        if paginationControls != "" {
            baseHelpText = paginationControls + " • " + baseHelpText
        }

        // Align pagination info to the right
        return lipgloss.JoinHorizontal(
            lipgloss.Left,
            m.Styles.Help.Render(baseHelpText),
            strings.Repeat(" ", max(0, m.Width-lipgloss.Width(baseHelpText)-lipgloss.Width(paginationInfo)-4)), // 4 for padding
            m.Styles.Help.Render(paginationInfo),
        )
    }

    return m.Styles.Help.Render(baseHelpText)
}
```

### 2.3 Keyboard Handlers for Pagination (Implemented)

We've added keyboard handlers for pagination in `internal/ui/ui.go`:

```go
// Add pagination key handlers
case constants.KeyPreviousPage, constants.KeyNextPage:
    // Handle pagination key presses if in a paginated view
    if view.IsPaginatedView(m.core.CurrentView) {
        // Pass the exact key string that was pressed
        modelWrapper, cmd := update.HandlePaginationKeyPress(m.core, msg.String())
        if wrapper, ok := modelWrapper.(update.ModelWrapper); ok {
            newModel := Model{core: wrapper.Model}
            if cmd != nil {
                return newModel, cmd
            }
            return newModel, nil
        }
        return modelWrapper, cmd
    }
    return m, nil
```

## 3. API-Level Pagination (To Be Implemented)

To implement API-level pagination, we need to make the following changes:

### 3.1 Update Provider Interfaces

Modify the provider interfaces to support pagination parameters:

```go
// In internal/cloud/provider.go
type FunctionStatusOperation interface {
    GetFunctionStatus(ctx context.Context, nextToken string, maxItems int) ([]FunctionStatus, string, error)
}

type PipelineStatusOperation interface {
    GetPipelineStatus(ctx context.Context, nextToken string, maxItems int) ([]PipelineStatus, string, error)
}

type CodePipelineManualApprovalOperation interface {
    GetPendingApprovals(ctx context.Context, nextToken string, maxItems int) ([]ApprovalAction, string, error)
}
```

### 3.2 Update AWS Provider Implementations

#### 3.2.1 Lambda Functions

Update the Lambda functions implementation in `internal/cloud/aws/lambda/operations.go`:

```go
func (o *FunctionStatusOperation) GetFunctionStatus(ctx context.Context, nextToken string, maxItems int) ([]cloud.FunctionStatus, string, error) {
    // Create a new AWS SDK client
    client, err := getClient(ctx, o.profile, o.region)
    if err != nil {
        return nil, "", err
    }

    // Create the input with pagination parameters
    input := &lambda.ListFunctionsInput{}
    if maxItems > 0 {
        input.MaxItems = aws.Int32(int32(maxItems))
    }
    if nextToken != "" {
        input.Marker = aws.String(nextToken)
    }

    // Call the AWS API
    result, err := client.ListFunctions(ctx, input)
    if err != nil {
        return nil, "", fmt.Errorf("failed to list functions: %w", err)
    }

    // Convert to cloud.FunctionStatus
    functionStatuses := make([]cloud.FunctionStatus, len(result.Functions))
    for i, function := range result.Functions {
        // Existing conversion code...
    }

    // Return the next marker for pagination
    var newNextToken string
    if result.NextMarker != nil {
        newNextToken = *result.NextMarker
    }

    return functionStatuses, newNextToken, nil
}
```

#### 3.2.2 CodePipeline

Update the CodePipeline implementation in `internal/cloud/aws/codepipeline/operations.go`:

```go
func (o *CloudPipelineStatusOperation) GetPipelineStatus(ctx context.Context, nextToken string, maxItems int) ([]cloud.PipelineStatus, string, error) {
    // Create a new AWS SDK client
    client, err := getClient(ctx, o.profile, o.region)
    if err != nil {
        return nil, "", err
    }

    // Create the input with pagination parameters
    input := &codepipeline.ListPipelinesInput{}
    if maxItems > 0 {
        input.MaxResults = aws.Int32(int32(maxItems))
    }
    if nextToken != "" {
        input.NextToken = aws.String(nextToken)
    }

    // Call the AWS API
    result, err := client.ListPipelines(ctx, input)
    if err != nil {
        return nil, "", fmt.Errorf("failed to list pipelines: %w", err)
    }

    // Convert to cloud.PipelineStatus
    pipelineStatuses := make([]cloud.PipelineStatus, len(result.Pipelines))
    for i, pipeline := range result.Pipelines {
        // Existing conversion code...
    }

    // Return the next token for pagination
    var newNextToken string
    if result.NextToken != nil {
        newNextToken = *result.NextToken
    }

    return pipelineStatuses, newNextToken, nil
}
```

#### 3.2.3 Pipeline Approvals

Update the Pipeline Approvals implementation in `internal/cloud/aws/codepipeline/approvals.go`:

```go
func (o *CloudManualApprovalOperation) GetPendingApprovals(ctx context.Context, nextToken string, maxItems int) ([]cloud.ApprovalAction, string, error) {
    // Create a new AWS SDK client
    client, err := getClient(ctx, o.profile, o.region)
    if err != nil {
        return nil, "", err
    }

    // Create the input with pagination parameters
    input := &codepipeline.ListPipelinesInput{}
    if maxItems > 0 {
        input.MaxResults = aws.Int32(int32(maxItems))
    }
    if nextToken != "" {
        input.NextToken = aws.String(nextToken)
    }

    // Call the AWS API
    result, err := client.ListPipelines(ctx, input)
    if err != nil {
        return nil, "", fmt.Errorf("failed to list pipelines: %w", err)
    }

    // Existing code to process approvals...

    // Return the next token for pagination
    var newNextToken string
    if result.NextToken != nil {
        newNextToken = *result.NextToken
    }

    return approvals, newNextToken, nil
}
```

### 3.3 Update Pagination Handlers for API-Level Pagination

Modify the pagination handlers in `internal/ui/update/pagination_handlers.go` to support API-level pagination:

```go
// FetchNextFunctionsPage fetches the next page of Lambda functions
func FetchNextFunctionsPage(m *model.Model) tea.Cmd {
    return func() tea.Msg {
        // Set loading state
        newModel := m.Clone()
        newModel.Pagination.IsLoading = true

        // Determine which pagination type to use
        if newModel.Pagination.Type == model.PaginationTypeClientSide {
            // Existing client-side pagination code...
        } else if newModel.Pagination.Type == model.PaginationTypeAPILevel {
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
            nextToken := m.Pagination.NextPageToken
            pageSize := m.Pagination.PageSize
            functions, newNextToken, err := functionOperation.GetFunctionStatus(ctx, nextToken, pageSize)
            if err != nil {
                return model.ErrMsg{Err: err}
            }

            // Determine if there are more pages
            hasMorePages := newNextToken != ""

            return model.FunctionsPageMsg{
                Functions:     functions,
                NextPageToken: newNextToken,
                HasMorePages:  hasMorePages,
            }
        }

        // Default case
        return model.FunctionsPageMsg{
            Functions:     []model.FunctionStatus{},
            NextPageToken: "",
            HasMorePages:  false,
        }
    }
}
```

### 3.4 Implement Token History for Previous Page Navigation

For API-level pagination, we need to keep track of previous tokens to support backward navigation:

```go
// Add to internal/ui/model/model.go
type Pagination struct {
    // Existing fields...
    
    // For API-level pagination
    TokenHistory    []string  // History of tokens for backward navigation
}

// Add to internal/ui/update/pagination_handlers.go
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
    if m.Pagination.Type == model.PaginationTypeClientSide {
        // Existing client-side pagination code...
    } else if m.Pagination.Type == model.PaginationTypeAPILevel {
        // Get the previous token from the token history
        if len(m.Pagination.TokenHistory) >= prevPage {
            prevToken := m.Pagination.TokenHistory[prevPage-1]
            
            // Use the previous token to fetch the page
            return fetchWithToken(m, prevToken, prevPage)
        } else {
            // If we don't have the token history, we need to start from the beginning
            return fetchFromBeginning(m, prevPage)
        }
    }

    return nil
}

// Helper function to fetch a page with a specific token
func fetchWithToken(m *model.Model, token string, targetPage int) tea.Cmd {
    // Implementation depends on the current view
    switch m.CurrentView {
    case constants.ViewFunctionStatus:
        return fetchFunctionsWithToken(m, token, targetPage)
    case constants.ViewPipelineStatus:
        return fetchPipelinesWithToken(m, token, targetPage)
    case constants.ViewApprovals:
        return fetchApprovalsWithToken(m, token, targetPage)
    default:
        return nil
    }
}

// Helper function to fetch from the beginning up to a specific page
func fetchFromBeginning(m *model.Model, targetPage int) tea.Cmd {
    // Implementation depends on the current view
    switch m.CurrentView {
    case constants.ViewFunctionStatus:
        return fetchFunctionsFromBeginning(m, targetPage)
    case constants.ViewPipelineStatus:
        return fetchPipelinesFromBeginning(m, targetPage)
    case constants.ViewApprovals:
        return fetchApprovalsFromBeginning(m, targetPage)
    default:
        return nil
    }
}
```

### 3.5 Pagination Type Detection

Add logic to detect when to use API-level pagination:

```go
// Add to internal/ui/model/model.go
func (m *Model) DeterminePaginationType() PaginationType {
    // Check if API-level pagination is enabled in configuration
    if m.Config.UseAPILevelPagination {
        return PaginationTypeAPILevel
    }
    
    // Check if the resource count is above the threshold
    if m.Config.PaginationThreshold > 0 {
        switch m.CurrentView {
        case constants.ViewFunctionStatus:
            if len(m.Functions) > m.Config.PaginationThreshold {
                return PaginationTypeAPILevel
            }
        case constants.ViewPipelineStatus:
            if len(m.Pipelines) > m.Config.PaginationThreshold {
                return PaginationTypeAPILevel
            }
        case constants.ViewApprovals:
            if len(m.Approvals) > m.Config.PaginationThreshold {
                return PaginationTypeAPILevel
            }
        }
    }
    
    // Default to client-side pagination
    return PaginationTypeClientSide
}
```

### 3.6 Configuration Options

Add configuration options for pagination:

```go
// Add to internal/ui/model/model.go
type Config struct {
    // Existing fields...
    
    // Pagination configuration
    PageSize                int
    UseAPILevelPagination   bool
    PaginationThreshold     int  // Number of items above which to use API-level pagination
}

// Add to internal/cmd/root.go
func init() {
    // Existing code...
    
    rootCmd.Flags().IntVar(&config.PageSize, "page-size", 10, "Number of items per page")
    rootCmd.Flags().BoolVar(&config.UseAPILevelPagination, "api-pagination", false, "Use API-level pagination")
    rootCmd.Flags().IntVar(&config.PaginationThreshold, "pagination-threshold", 100, "Number of items above which to use API-level pagination")
    
    // Existing code...
}
```

## 4. Hybrid Pagination Approach

For the best user experience, we'll implement a hybrid pagination approach:

### 4.1 Initial Loading Strategy

When a view is first loaded:
1. Determine whether to use client-side or API-level pagination based on configuration
2. If using client-side pagination, fetch all items at once
3. If using API-level pagination, fetch only the first page

### 4.2 Progressive Loading

For API-level pagination:
1. As the user navigates through pages, fetch each page as needed
2. Store fetched items in memory for potential reuse
3. Keep track of token history for backward navigation

### 4.3 Automatic Switching

Implement automatic switching between pagination types:
1. Start with API-level pagination for efficiency
2. If the total number of items is small, switch to client-side pagination
3. If memory usage becomes a concern, switch back to API-level pagination

## 5. Next Steps

1. **Implement API-Level Pagination**:
   - Update provider interfaces to support pagination parameters
   - Modify AWS provider implementations to use pagination
   - Update pagination handlers to support API-level pagination
   - Implement token history for previous page navigation

2. **Add Configuration Options**:
   - Add command-line flags for pagination configuration
   - Implement pagination type detection logic

3. **Enhance User Experience**:
   - Improve loading indicators during pagination
   - Add more detailed pagination information in the UI
   - Implement the hybrid pagination approach

4. **Testing and Optimization**:
   - Test with large datasets to ensure performance
   - Optimize memory usage for client-side pagination
   - Ensure smooth navigation experience in all scenarios

## Conclusion

The pagination implementation in Cloudgate has made significant progress with the client-side pagination implementation. The next phase will focus on implementing API-level pagination to support large enterprise environments with thousands of resources. The hybrid approach will provide the best of both worlds: efficient data fetching for large datasets and smooth navigation for all users. 