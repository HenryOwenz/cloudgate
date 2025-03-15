package model

import (
	"github.com/HenryOwenz/cloudgate/internal/cloud"
)

// Service represents a cloud service
type Service struct {
	ID          string
	Name        string
	Description string
	Available   bool
}

// Category represents a service category
type Category struct {
	ID          string
	Name        string
	Description string
	Available   bool
}

// Operation represents a service operation
type Operation struct {
	ID          string
	Name        string
	Description string
}

// ErrMsg represents an error message
type ErrMsg struct {
	Err error
}

// ApprovalAction is an alias for cloud.ApprovalAction
type ApprovalAction = cloud.ApprovalAction

// StageStatus is an alias for cloud.StageStatus
type StageStatus = cloud.StageStatus

// PipelineStatus is an alias for cloud.PipelineStatus
type PipelineStatus = cloud.PipelineStatus

// FunctionStatus is an alias for cloud.FunctionStatus
type FunctionStatus = cloud.FunctionStatus

// PaginationType represents the type of pagination to use
type PaginationType int

const (
	// PaginationTypeNone indicates no pagination is needed
	PaginationTypeNone PaginationType = iota
	// PaginationTypeAPI indicates API-level pagination (e.g., Lambda)
	PaginationTypeAPI
	// PaginationTypeClientSide indicates client-side pagination (e.g., S3)
	PaginationTypeClientSide
)

// Pagination represents a generic pagination state
type Pagination struct {
	Type          PaginationType
	CurrentPage   int
	PageSize      int
	TotalItems    int64 // May be unknown (-1)
	HasMorePages  bool
	NextPageToken string // Service-specific token
	IsLoading     bool

	// For client-side pagination
	AllItems      []interface{} // All items fetched so far
	FilteredItems []interface{} // Items after filtering
}

// ApprovalsMsg represents a message containing approvals
type ApprovalsMsg struct {
	Approvals []ApprovalAction
	Provider  cloud.Provider
}

// ApprovalResultMsg represents the result of an approval action
type ApprovalResultMsg struct {
	Err error
}

// PipelineStatusMsg represents a message containing pipeline status
type PipelineStatusMsg struct {
	Pipelines []PipelineStatus
	Provider  cloud.Provider
}

// PipelineExecutionMsg represents the result of a pipeline execution
type PipelineExecutionMsg struct {
	Err error
}

// FunctionStatusMsg represents a message containing function status
type FunctionStatusMsg struct {
	Functions []FunctionStatus
	Provider  cloud.Provider
}

// LambdaExecuteResultMsg represents the result of a Lambda execution
type LambdaExecuteResultMsg struct {
	Result *cloud.LambdaExecuteResult
	Err    error
}

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
