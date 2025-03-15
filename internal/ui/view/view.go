package view

import (
	"fmt"
	"strings"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/charmbracelet/lipgloss"
)

// IsPaginatedView returns true if the view supports pagination
func IsPaginatedView(view constants.View) bool {
	return view == constants.ViewFunctionStatus ||
		view == constants.ViewPipelineStatus ||
		view == constants.ViewApprovals
}

// Render renders the UI
func Render(m *model.Model) string {
	if m.Err != nil {
		return renderErrorView(m)
	}

	// Create content array with appropriate spacing
	content := make([]string, constants.AppContentLines)

	// Set the title and context
	content[0] = renderTitle(m)
	content[1] = renderContext(m)
	content[2] = renderLoadingSpinner(m)
	content[3] = renderMainContent(m)
	content[4] = renderHelpText(m)

	// Join all content vertically with consistent spacing
	return m.Styles.App.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			content...,
		),
	)
}

// renderErrorView renders the error view
func renderErrorView(m *model.Model) string {
	return m.Styles.App.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.Styles.Error.Render("Error: "+m.Err.Error()),
			"\n",
			m.Styles.Help.Render(fmt.Sprintf("%s: quit • %s: back", constants.KeyQ, constants.KeyAltBack)),
		),
	)
}

// renderTitle renders the title based on the current view
func renderTitle(m *model.Model) string {
	return m.Styles.Title.Render(getTitleText(m))
}

// renderContext renders the context based on the current view
func renderContext(m *model.Model) string {
	return m.Styles.Context.Render(getContextText(m))
}

// renderLoadingSpinner renders the loading spinner if needed
func renderLoadingSpinner(m *model.Model) string {
	if m.IsLoading {
		return m.Spinner.View()
	}
	return ""
}

// renderMainContent renders the main content area based on the current view
func renderMainContent(m *model.Model) string {
	switch m.CurrentView {
	case constants.ViewProviders:
		return renderTable(m)
	case constants.ViewAWSConfig:
		if m.ManualInput {
			return m.TextInput.View()
		}
		return renderTable(m)
	case constants.ViewSelectService:
		return renderTable(m)
	case constants.ViewSelectCategory:
		return renderTable(m)
	case constants.ViewSelectOperation:
		return renderTable(m)
	case constants.ViewApprovals:
		return renderTable(m)
	case constants.ViewConfirmation:
		return renderTable(m)
	case constants.ViewSummary:
		if m.ManualInput {
			return m.TextInput.View()
		}
		return m.Summary
	case constants.ViewPipelineStatus:
		return renderTable(m)
	case constants.ViewPipelineStages:
		return renderTable(m)
	case constants.ViewFunctionStatus:
		return renderTable(m)
	case constants.ViewFunctionDetails:
		return renderTable(m)
	case constants.ViewLambdaExecute:
		// Set fixed height to match standard table views
		height := constants.TableHeight

		// Set the dimensions for the TextArea
		m.TextArea.SetWidth(m.Width - constants.ViewportMarginX*2)
		m.TextArea.SetHeight(height)

		// Create a custom view with the title and TextArea
		title := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(constants.ColorTitle)).
			Render(constants.TitleLambdaExecute)

		// Create a header with a line extending to the width
		line := strings.Repeat("─", max(0, m.Width-constants.ViewportMarginX*2-lipgloss.Width(title)))
		header := lipgloss.JoinHorizontal(lipgloss.Center, title, line)

		// Get the main content (TextArea view)
		content := m.TextArea.View()

		// Add a footer line
		footerText := ""
		if m.IsLambdaInputMode {
			footerText = "INPUT MODE"
		} else {
			footerText = "COMMAND MODE"
		}
		footer := lipgloss.NewStyle().
			Foreground(lipgloss.Color(constants.ColorPrimary)).
			Render(footerText)
		footerLine := strings.Repeat("─", max(0, m.Width-constants.ViewportMarginX*2-lipgloss.Width(footerText)))
		footer = lipgloss.JoinHorizontal(lipgloss.Center, footerLine, footer)

		// Return the complete view
		return fmt.Sprintf("%s\n%s\n%s", header, content, footer)
	case constants.ViewLambdaResponse:
		// Create a custom view with the title and Viewport
		title := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(constants.ColorTitle)).
			Render(constants.TitleLambdaResponse)

		// Create a header with a line extending to the viewport width
		line := strings.Repeat("─", max(0, m.Viewport.Width-lipgloss.Width(title)))
		header := lipgloss.JoinHorizontal(lipgloss.Center, title, line)

		// Add a footer with scroll percentage
		footerText := fmt.Sprintf("%3.f%%", m.Viewport.ScrollPercent()*100)
		footer := lipgloss.NewStyle().
			Foreground(lipgloss.Color(constants.ColorPrimary)).
			Render(footerText)
		footerLine := strings.Repeat("─", max(0, m.Viewport.Width-lipgloss.Width(footerText)))
		footer = lipgloss.JoinHorizontal(lipgloss.Center, footerLine, footer)

		// Set viewport height to match table height
		m.Viewport.Height = constants.TableHeight

		// Return the complete view
		return fmt.Sprintf("%s\n%s\n%s", header, m.Viewport.View(), footer)
	case constants.ViewExecutingAction:
		// Show the table instead of just the loading message
		return renderTable(m)
	default:
		return ""
	}
}

// renderHelpText renders the help text based on the current view
func renderHelpText(m *model.Model) string {
	helpText := getHelpText(m)
	return m.Styles.Help.Render(helpText)
}

// getContextText returns the appropriate context text for the current view
func getContextText(m *model.Model) string {
	switch m.CurrentView {
	case constants.ViewProviders:
		return getProvidersContextText()
	case constants.ViewAWSConfig:
		return getAWSConfigContextText(m)
	case constants.ViewSelectService:
		return getSelectServiceContextText(m)
	case constants.ViewSelectCategory:
		return getSelectCategoryContextText(m)
	case constants.ViewSelectOperation:
		return getSelectOperationContextText(m)
	case constants.ViewApprovals:
		return getApprovalsContextText(m)
	case constants.ViewConfirmation, constants.ViewSummary:
		return getConfirmationSummaryContextText(m)
	case constants.ViewExecutingAction:
		return getExecutingActionContextText(m)
	case constants.ViewPipelineStatus:
		return getPipelineStatusContextText(m)
	case constants.ViewPipelineStages:
		return getPipelineStagesContextText(m)
	case constants.ViewFunctionStatus:
		return getFunctionStatusContextText(m)
	case constants.ViewFunctionDetails:
		return getFunctionDetailsContextText(m)
	case constants.ViewLambdaExecute:
		return getLambdaExecuteContextText(m)
	case constants.ViewLambdaResponse:
		return getLambdaResponseContextText(m)
	default:
		return ""
	}
}

// getProvidersContextText returns the context text for the providers view
func getProvidersContextText() string {
	return constants.MsgAppDescription
}

// getAWSConfigContextText returns the context text for the AWS config view
func getAWSConfigContextText(m *model.Model) string {
	if m.AwsProfile == "" {
		// If in manual entry mode for profile, show the text input in the context
		if m.ManualInput {
			return fmt.Sprintf("Amazon Web Services\n\nEnter AWS Profile: %s", m.TextInput.View())
		}
		return "Amazon Web Services"
	}
	// If in manual entry mode for region, show the text input in the context
	if m.ManualInput {
		return fmt.Sprintf("Profile: %s\n\nEnter AWS Region: %s", m.AwsProfile, m.TextInput.View())
	}
	return fmt.Sprintf("Profile: %s", m.AwsProfile)
}

// getSelectServiceContextText returns the context text for the select service view
func getSelectServiceContextText(m *model.Model) string {
	return fmt.Sprintf("Profile: %s\nRegion: %s",
		m.AwsProfile,
		m.AwsRegion)
}

// getSelectCategoryContextText returns the context text for the select category view
func getSelectCategoryContextText(m *model.Model) string {
	if m.SelectedService == nil {
		return ""
	}
	return fmt.Sprintf("Profile: %s\nRegion: %s\nService: %s",
		m.AwsProfile,
		m.AwsRegion,
		m.SelectedService.Name)
}

// getSelectOperationContextText returns the context text for the select operation view
func getSelectOperationContextText(m *model.Model) string {
	if m.SelectedService == nil || m.SelectedCategory == nil {
		return ""
	}
	return fmt.Sprintf("Profile: %s\nRegion: %s\nService: %s\nCategory: %s",
		m.AwsProfile,
		m.AwsRegion,
		m.SelectedService.Name,
		m.SelectedCategory.Name)
}

// getApprovalsContextText returns the context text for the approvals view
func getApprovalsContextText(m *model.Model) string {
	return fmt.Sprintf("Profile: %s\nRegion: %s",
		m.AwsProfile,
		m.AwsRegion)
}

// getConfirmationSummaryContextText returns the context text for the confirmation and summary views
func getConfirmationSummaryContextText(m *model.Model) string {
	if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
		if m.SelectedPipeline == nil {
			return ""
		}
		return fmt.Sprintf("Profile: %s\nRegion: %s\nPipeline: %s",
			m.AwsProfile,
			m.AwsRegion,
			m.SelectedPipeline.Name)
	}
	if m.SelectedApproval == nil {
		return ""
	}
	return fmt.Sprintf("Profile: %s\nRegion: %s\nPipeline: %s\nStage: %s\nAction: %s",
		m.AwsProfile,
		m.AwsRegion,
		m.SelectedApproval.PipelineName,
		m.SelectedApproval.StageName,
		m.SelectedApproval.ActionName)
}

// getExecutingActionContextText returns the context text for the executing action view
func getExecutingActionContextText(m *model.Model) string {
	if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
		if m.SelectedPipeline == nil {
			return ""
		}

		revisionID := "Latest commit"
		if m.ManualCommitID && m.CommitID != "" {
			revisionID = m.CommitID
		}

		return fmt.Sprintf("Profile: %s\nRegion: %s\nPipeline: %s\nRevisionID: %s",
			m.AwsProfile,
			m.AwsRegion,
			m.SelectedPipeline.Name,
			revisionID)
	}
	if m.SelectedApproval == nil {
		return ""
	}
	return fmt.Sprintf("Profile: %s\nRegion: %s\nPipeline: %s\nStage: %s\nAction: %s\nComment: %s",
		m.AwsProfile,
		m.AwsRegion,
		m.SelectedApproval.PipelineName,
		m.SelectedApproval.StageName,
		m.SelectedApproval.ActionName,
		m.ApprovalComment)
}

// getPipelineStatusContextText returns the context text for the pipeline status view
func getPipelineStatusContextText(m *model.Model) string {
	return fmt.Sprintf("Profile: %s\nRegion: %s",
		m.AwsProfile,
		m.AwsRegion)
}

// getPipelineStagesContextText returns the context text for the pipeline stages view
func getPipelineStagesContextText(m *model.Model) string {
	if m.SelectedPipeline == nil {
		return ""
	}
	return fmt.Sprintf("Profile: %s\nRegion: %s\nPipeline: %s",
		m.AwsProfile,
		m.AwsRegion,
		m.SelectedPipeline.Name)
}

// getFunctionStatusContextText returns the context text for the function status view
func getFunctionStatusContextText(m *model.Model) string {
	return fmt.Sprintf("Profile: %s\nRegion: %s\nService: %s\nCategory: %s",
		m.AwsProfile,
		m.AwsRegion,
		m.SelectedService.Name,
		m.SelectedCategory.Name)
}

// getFunctionDetailsContextText returns the context text for the function details view
func getFunctionDetailsContextText(m *model.Model) string {
	if m.SelectedFunction == nil {
		return ""
	}
	return fmt.Sprintf("Profile: %s\nRegion: %s\nFunction: %s",
		m.AwsProfile,
		m.AwsRegion,
		m.SelectedFunction.Name)
}

// getLambdaExecuteContextText returns the context text for the Lambda execution view
func getLambdaExecuteContextText(m *model.Model) string {
	if m.SelectedFunction == nil {
		return ""
	}

	return fmt.Sprintf(
		"Profile: %s\nRegion: %s\nService: Lambda\nFunction: %s\nRuntime: %s",
		m.AwsProfile,
		m.AwsRegion,
		m.SelectedFunction.Name,
		m.SelectedFunction.Runtime,
	)
}

// getLambdaResponseContextText returns the context text for the Lambda response view
func getLambdaResponseContextText(m *model.Model) string {
	if m.SelectedFunction == nil {
		return ""
	}

	return fmt.Sprintf(
		"Profile: %s\nRegion: %s\nService: Lambda\nFunction: %s\nRuntime: %s",
		m.AwsProfile,
		m.AwsRegion,
		m.SelectedFunction.Name,
		m.SelectedFunction.Runtime,
	)
}

// getTitleText returns the appropriate title for the current view
func getTitleText(m *model.Model) string {
	// Map of view types to their corresponding titles
	titleMap := map[constants.View]string{
		constants.ViewProviders:       constants.TitleProviders,
		constants.ViewSelectService:   constants.TitleSelectService,
		constants.ViewSelectCategory:  constants.TitleSelectCategory,
		constants.ViewSelectOperation: constants.TitleSelectOperation,
		constants.ViewApprovals:       constants.TitleApprovals,
		constants.ViewConfirmation:    constants.TitleConfirmation,
		constants.ViewSummary:         constants.TitleSummary,
		constants.ViewExecutingAction: constants.TitleExecutingAction,
		constants.ViewPipelineStatus:  constants.TitlePipelineStatus,
		constants.ViewPipelineStages:  constants.TitlePipelineStages,
		constants.ViewError:           constants.TitleError,
		constants.ViewSuccess:         constants.TitleSuccess,
		constants.ViewHelp:            constants.TitleHelp,
		constants.ViewFunctionStatus:  constants.TitleFunctionStatus,
		constants.ViewFunctionDetails: constants.TitleFunctionDetails,
		constants.ViewLambdaExecute:   constants.TitleLambdaExecute,
		constants.ViewLambdaResponse:  constants.TitleLambdaResponse,
	}

	// Special case for AWS config view
	if m.CurrentView == constants.ViewAWSConfig {
		if m.AwsProfile == "" {
			return constants.TitleSelectProfile
		}
		return constants.TitleSelectRegion
	}

	// Get the base title
	var title string
	if t, ok := titleMap[m.CurrentView]; ok {
		title = t
	} else {
		return ""
	}

	// Add pagination information for paginated views
	if IsPaginatedView(m.CurrentView) && m.Pagination.Type != model.PaginationTypeNone {
		// Calculate total pages
		totalPages := 1
		if m.Pagination.TotalItems > 0 && m.Pagination.PageSize > 0 {
			totalPages = int((m.Pagination.TotalItems + int64(m.Pagination.PageSize) - 1) / int64(m.Pagination.PageSize))
		}

		// Create pagination text
		paginationText := fmt.Sprintf(" - Page %d of %d", m.Pagination.CurrentPage, totalPages)

		if m.Pagination.TotalItems >= 0 {
			paginationText += fmt.Sprintf(" (%d items)", m.Pagination.TotalItems)
		}

		// Add pagination text to the title
		title += paginationText
	}

	return title
}

// getHelpText returns the appropriate help text for the current view
func getHelpText(m *model.Model) string {
	// Define common help text patterns
	const (
		defaultHelpText        = "↑/↓: navigate • %s: select • %s: back • %s: quit"
		manualInputHelpText    = "%s: confirm • %s: cancel • %s: quit"
		summaryHelpText        = "↑/↓: navigate • %s: select • %s: back • %s: quit"
		providersHelpText      = "↑/↓: navigate • %s: select • %s: quit"
		lambdaCommandModeText  = "-- COMMAND MODE -- • i: enter input mode • enter: execute • %s: back • %s: quit"
		lambdaInputModeText    = "-- INPUT MODE -- • enter: new line • ctrl+c/esc: exit input mode • %s: back • %s: quit"
		lambdaResponseHelpText = "↑/↓: scroll • pgup/pgdn: page • home/end: top/bottom • %s: back to editor • %s: quit"
	)

	// Special cases based on view and state
	switch {
	case m.CurrentView == constants.ViewProviders:
		return fmt.Sprintf(providersHelpText, constants.KeyEnter, constants.KeyQ)
	case m.CurrentView == constants.ViewAWSConfig && m.ManualInput:
		return fmt.Sprintf(manualInputHelpText, constants.KeyEnter, constants.KeyEsc, constants.KeyCtrlC)
	case m.CurrentView == constants.ViewSummary && m.ManualInput:
		return fmt.Sprintf(manualInputHelpText, constants.KeyEnter, constants.KeyEsc, constants.KeyCtrlC)
	case m.CurrentView == constants.ViewSummary:
		return fmt.Sprintf(summaryHelpText, constants.KeyEnter, constants.KeyEsc, constants.KeyQ)
	case m.CurrentView == constants.ViewLambdaExecute:
		if m.IsLambdaInputMode {
			return fmt.Sprintf(lambdaInputModeText, constants.KeyEsc, constants.KeyQ)
		}
		return fmt.Sprintf(lambdaCommandModeText, constants.KeyEsc, constants.KeyQ)
	case m.CurrentView == constants.ViewLambdaResponse:
		return fmt.Sprintf(lambdaResponseHelpText, constants.KeyEsc, constants.KeyQ)
	default:
		return fmt.Sprintf(defaultHelpText, constants.KeyEnter, constants.KeyEsc, constants.KeyQ)
	}
}

// renderTable renders the table for the current view
func renderTable(m *model.Model) string {
	if m.Table.Rows() == nil {
		return ""
	}

	// Create a table style with fixed height to match viewport
	tableStyle := lipgloss.NewStyle().
		Height(constants.TableHeight).
		PaddingTop(1).
		PaddingRight(2).
		PaddingBottom(0).
		PaddingLeft(0)

	// Render the table with the appropriate styles
	return tableStyle.Render(m.Table.View())
}

// Add the max helper function at the end of the file
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
