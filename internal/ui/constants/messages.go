package constants

// Message constants for UI text
const (
	// Description messages
	MsgAppDescription = "A simple tool to manage your cloud resources"

	// Loading messages
	MsgLoadingApprovals  = "Loading approvals..."
	MsgLoadingPipelines  = "Loading pipelines..."
	MsgLoadingFunctions  = "Loading functions..."
	MsgStartingPipeline  = "Starting pipeline..."
	MsgExecutingApproval = "Executing approval action..."
	MsgExecutingLambda   = "Executing Lambda function..."

	// Input placeholders
	MsgEnterProfile          = "Enter AWS profile name..."
	MsgEnterRegion           = "Enter AWS region..."
	MsgEnterComment          = "Enter comment..."
	MsgEnterApprovalComment  = "Enter approval comment..."
	MsgEnterRejectionComment = "Enter rejection comment..."
	MsgEnterCommitID         = "Enter commit ID..."
	MsgEnterLambdaPayload    = "Enter Lambda JSON payload..."

	// Success messages
	MsgApprovalSuccess      = "Successfully approved pipeline: %s, stage: %s, action: %s"
	MsgRejectionSuccess     = "Successfully rejected pipeline: %s, stage: %s, action: %s"
	MsgPipelineStartSuccess = "Successfully started pipeline: %s"
	MsgLambdaExecuteSuccess = "Successfully executed Lambda function: %s"

	// Error messages
	MsgErrorGeneric       = "Error: %s"
	MsgErrorNoApproval    = "No approval selected"
	MsgErrorNoPipeline    = "No pipeline selected"
	MsgErrorNoFunction    = "No function selected"
	MsgNoFunctionSelected = "No function selected"
	MsgLambdaExecuteError = "Error executing Lambda function %s: %v"
	MsgErrorEmptyCommitID = "Commit ID cannot be empty"
	MsgErrorEmptyComment  = "Comment cannot be empty"
	MsgErrorInvalidJSON   = "Invalid JSON payload"
)
