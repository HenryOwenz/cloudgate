package view

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/cloudproviders"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
)

// UpdateTableForView updates the table model based on the current view
func UpdateTableForView(m *model.Model) {
	columns := getColumnsForView(m)
	rows := getRowsForView(m)

	// Set the table height based on the current view
	tableHeight := constants.TableHeight
	if m.CurrentView == constants.ViewPipelineStages {
		tableHeight = constants.TableHeightLarge
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)

	t.SetStyles(m.Styles.Table)
	m.Table = t
}

// getColumnsForView returns the appropriate columns for the current view
func getColumnsForView(m *model.Model) []table.Column {
	switch m.CurrentView {
	case constants.ViewProviders:
		return []table.Column{
			{Title: "Provider", Width: constants.TableDefaultWidth},
			{Title: "Description", Width: constants.TableDescWidth},
		}
	case constants.ViewAuthMethodSelect:
		return []table.Column{
			{Title: "Authentication Method", Width: constants.TableDefaultWidth},
			{Title: "Description", Width: constants.TableDescWidth},
		}
	case constants.ViewAuthConfig:
		return []table.Column{
			{Title: m.ProviderState.AuthState.CurrentAuthConfigKey, Width: constants.TableDefaultWidth},
		}
	case constants.ViewProviderConfig:
		return []table.Column{
			{Title: m.ProviderState.CurrentConfigKey, Width: constants.TableDefaultWidth},
		}
	case constants.ViewAWSConfig:
		if m.AwsProfile == "" {
			return []table.Column{{Title: "Profile", Width: constants.TableDefaultWidth}}
		}
		return []table.Column{{Title: "Region", Width: constants.TableDefaultWidth}}
	case constants.ViewSelectService:
		return []table.Column{
			{Title: "Service", Width: constants.TableDefaultWidth},
			{Title: "Description", Width: constants.TableDescWidth},
		}
	case constants.ViewSelectCategory:
		return []table.Column{
			{Title: "Category", Width: constants.TableDefaultWidth},
			{Title: "Description", Width: constants.TableDescWidth},
		}
	case constants.ViewSelectOperation:
		return []table.Column{
			{Title: "Operation", Width: constants.TableDefaultWidth},
			{Title: "Description", Width: constants.TableDescWidth},
		}
	case constants.ViewApprovals:
		return []table.Column{
			{Title: "Pipeline", Width: constants.TableWideWidth},
			{Title: "Stage", Width: constants.TableDefaultWidth},
			{Title: "Action", Width: constants.TableNarrowWidth},
		}
	case constants.ViewConfirmation:
		return []table.Column{
			{Title: "Action", Width: constants.TableDefaultWidth},
			{Title: "Description", Width: constants.TableDescWidth},
		}
	case constants.ViewExecutingAction:
		return []table.Column{
			{Title: "Action", Width: constants.TableDefaultWidth},
			{Title: "Description", Width: constants.TableDescWidth},
		}
	case constants.ViewPipelineStatus:
		return []table.Column{
			{Title: "Pipeline", Width: constants.TableWideWidth},
			{Title: "Description", Width: constants.TableDescWidth},
		}
	case constants.ViewPipelineStages:
		return []table.Column{
			{Title: "Stage", Width: constants.TableDefaultWidth},
			{Title: "Status", Width: constants.TableNarrowWidth},
			{Title: "Last Updated", Width: constants.TableNarrowWidth},
		}
	case constants.ViewFunctionStatus:
		return []table.Column{
			{Title: "Function", Width: constants.TableWideWidth},
			{Title: "Runtime", Width: constants.TableNarrowWidth},
			{Title: "Last Updated", Width: constants.TableDefaultWidth},
		}
	case constants.ViewFunctionDetails:
		return []table.Column{
			{Title: "Property", Width: constants.TableDefaultWidth},
			{Title: "Value", Width: constants.TableWideWidth},
		}
	case constants.ViewSummary:
		return []table.Column{
			{Title: "Type", Width: constants.TableDefaultWidth},
			{Title: "Value", Width: constants.TableDescWidth},
		}
	default:
		return []table.Column{}
	}
}

// getRowsForView returns the appropriate rows for the current view
func getRowsForView(m *model.Model) []table.Row {
	switch m.CurrentView {
	case constants.ViewProviders:
		// Initialize providers if not already done
		if len(m.Registry.GetProviderNames()) == 0 {
			cloudproviders.InitializeProviders(m.Registry)
		}

		// Get all providers from the registry
		providerList := m.Registry.List()
		rows := make([]table.Row, len(providerList))

		for i, provider := range providerList {
			rows[i] = table.Row{provider.Name(), provider.Description()}
		}

		return rows
	case constants.ViewAuthMethodSelect:
		rows := make([]table.Row, len(m.ProviderState.AuthState.AvailableMethods))
		for i, method := range m.ProviderState.AuthState.AvailableMethods {
			description := getAuthMethodDescription(m.ProviderState.ProviderName, method)
			rows[i] = table.Row{method, description}
		}
		return rows
	case constants.ViewAuthConfig:
		key := m.ProviderState.AuthState.CurrentAuthConfigKey
		options, ok := m.ProviderState.ConfigOptions[key]
		if !ok {
			return []table.Row{}
		}

		rows := make([]table.Row, len(options)+1)
		rows[0] = table.Row{"Manual Entry"}
		for i, option := range options {
			rows[i+1] = table.Row{option}
		}
		return rows
	case constants.ViewProviderConfig:
		key := m.ProviderState.CurrentConfigKey
		options, ok := m.ProviderState.ConfigOptions[key]
		if !ok {
			return []table.Row{}
		}

		rows := make([]table.Row, len(options)+1)
		rows[0] = table.Row{"Manual Entry"}
		for i, option := range options {
			rows[i+1] = table.Row{option}
		}
		return rows
	case constants.ViewAWSConfig:
		if m.AwsProfile == "" {
			rows := make([]table.Row, len(m.Profiles)+1)
			rows[0] = table.Row{"Manual Entry"}
			for i, profile := range m.Profiles {
				rows[i+1] = table.Row{profile}
			}
			return rows
		}
		rows := make([]table.Row, len(m.Regions)+1)
		rows[0] = table.Row{"Manual Entry"}
		for i, region := range m.Regions {
			rows[i+1] = table.Row{region}
		}
		return rows
	case constants.ViewSelectService:
		// Get the AWS provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return []table.Row{}
		}

		// Get all services from the provider
		services := provider.Services()
		rows := make([]table.Row, len(services))

		for i, service := range services {
			rows[i] = table.Row{service.Name(), service.Description()}
		}

		// Sort services by name in ascending order
		sort.Slice(rows, func(i, j int) bool {
			return rows[i][0] < rows[j][0]
		})

		return rows
	case constants.ViewSelectCategory:
		if m.SelectedService == nil {
			return []table.Row{}
		}

		// Get the AWS provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return []table.Row{}
		}

		// Find the selected service
		var selectedService cloud.Service
		for _, service := range provider.Services() {
			if service.Name() == m.SelectedService.Name {
				selectedService = service
				break
			}
		}

		if selectedService == nil {
			return []table.Row{}
		}

		// Get all categories from the service
		categories := selectedService.Categories()

		// Filter out non-UI visible categories
		var visibleCategories []cloud.Category
		for _, category := range categories {
			if category.IsUIVisible() {
				visibleCategories = append(visibleCategories, category)
			}
		}

		rows := make([]table.Row, len(visibleCategories))
		for i, category := range visibleCategories {
			rows[i] = table.Row{category.Name(), category.Description()}
		}

		return rows
	case constants.ViewSelectOperation:
		if m.SelectedService == nil || m.SelectedCategory == nil {
			return []table.Row{}
		}

		// Get the AWS provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return []table.Row{}
		}

		// Find the selected service
		var selectedService cloud.Service
		for _, service := range provider.Services() {
			if service.Name() == m.SelectedService.Name {
				selectedService = service
				break
			}
		}

		if selectedService == nil {
			return []table.Row{}
		}

		// Find the selected category
		var selectedCategory cloud.Category
		for _, category := range selectedService.Categories() {
			if category.Name() == m.SelectedCategory.Name {
				selectedCategory = category
				break
			}
		}

		if selectedCategory == nil {
			return []table.Row{}
		}

		// Get all operations from the category
		operations := selectedCategory.Operations()

		// Filter out non-UI visible operations
		var visibleOperations []cloud.Operation
		for _, operation := range operations {
			if operation.IsUIVisible() {
				visibleOperations = append(visibleOperations, operation)
			}
		}

		// Sort operations alphabetically by name
		sort.Slice(visibleOperations, func(i, j int) bool {
			return visibleOperations[i].Name() < visibleOperations[j].Name()
		})

		rows := make([]table.Row, len(visibleOperations))
		for i, operation := range visibleOperations {
			rows[i] = table.Row{operation.Name(), operation.Description()}
		}

		return rows
	case constants.ViewApprovals:
		rows := make([]table.Row, len(m.Approvals))
		for i, approval := range m.Approvals {
			rows[i] = table.Row{
				approval.PipelineName,
				approval.StageName,
				approval.ActionName,
			}
		}
		return rows
	case constants.ViewConfirmation:
		return []table.Row{
			{"Approve", "Approve the pipeline stage"},
			{"Reject", "Reject the pipeline stage"},
		}
	case constants.ViewExecutingAction:
		if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
			return []table.Row{
				{"Execute", "Start pipeline with latest commit"},
				{"Cancel", "Cancel and return to main menu"},
			}
		}
		action := "approve"
		if !m.ApproveAction {
			action = "reject"
		}
		return []table.Row{
			{"Execute", fmt.Sprintf("Execute %s action", action)},
			{"Cancel", "Cancel and return to main menu"},
		}
	case constants.ViewPipelineStatus:
		if m.Pipelines == nil {
			return []table.Row{}
		}
		rows := make([]table.Row, len(m.Pipelines))
		for i, pipeline := range m.Pipelines {
			rows[i] = table.Row{
				pipeline.Name,
				fmt.Sprintf("%d stages", len(pipeline.Stages)),
			}
		}
		return rows
	case constants.ViewPipelineStages:
		if m.SelectedPipeline == nil {
			return []table.Row{}
		}
		rows := make([]table.Row, len(m.SelectedPipeline.Stages))
		for i, stage := range m.SelectedPipeline.Stages {
			rows[i] = table.Row{
				stage.Name,
				stage.Status,
				stage.LastUpdated,
			}
		}
		return rows
	case constants.ViewFunctionStatus:
		if m.Functions == nil {
			return []table.Row{}
		}
		rows := make([]table.Row, len(m.Functions))
		for i, function := range m.Functions {
			// Clean up timestamp by removing the milliseconds and timezone offset
			lastUpdate := function.LastUpdate
			if len(lastUpdate) > 19 { // Format: "2024-06-29T07:10:02.331+0000"
				lastUpdate = lastUpdate[:19] // Keep only "2024-06-29T07:10:02"
				// Replace 'T' with a space for better readability
				lastUpdate = strings.Replace(lastUpdate, "T", " ", 1)
			}

			rows[i] = table.Row{
				function.Name,
				function.Runtime,
				lastUpdate,
			}
		}
		return rows
	case constants.ViewFunctionDetails:
		if m.SelectedFunction == nil {
			return []table.Row{}
		}
		function := m.SelectedFunction

		// Format code size to be more readable (KB or MB)
		var codeSizeFormatted string
		if function.CodeSize < 1024 {
			codeSizeFormatted = fmt.Sprintf("%d bytes", function.CodeSize)
		} else if function.CodeSize < 1024*1024 {
			codeSizeFormatted = fmt.Sprintf("%.2f KB", float64(function.CodeSize)/1024)
		} else {
			codeSizeFormatted = fmt.Sprintf("%.2f MB", float64(function.CodeSize)/(1024*1024))
		}

		// Clean up timestamp by removing the milliseconds and timezone offset
		lastUpdate := function.LastUpdate
		if len(lastUpdate) > 19 { // Format: "2024-06-29T07:10:02.331+0000"
			lastUpdate = lastUpdate[:19] // Keep only "2024-06-29T07:10:02"
			// Replace 'T' with a space for better readability
			lastUpdate = strings.Replace(lastUpdate, "T", " ", 1)
		}

		rows := []table.Row{
			{"Name", function.Name},
			{"Description", function.Description},
			{"ARN", function.FunctionArn},
			{"Runtime", function.Runtime},
			{"Handler", function.Handler},
			{"Memory", fmt.Sprintf("%d MB", function.Memory)},
			{"Timeout", fmt.Sprintf("%d seconds", function.Timeout)},
			{"Code Size", codeSizeFormatted},
			{"Last Updated", lastUpdate},
			{"Version", function.Version},
			{"Package Type", function.PackageType},
			{"Architecture", function.Architecture},
			{"Role", function.Role},
		}

		// Only add log group if it's available
		if function.LogGroup != "" {
			rows = append(rows, table.Row{"Log Group", function.LogGroup})
		}

		return rows
	case constants.ViewSummary:
		if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
			if m.SelectedPipeline == nil {
				return []table.Row{}
			}
			return []table.Row{
				{"Latest Commit", "Use latest commit from source"},
				{"Manual Input", "Enter specific commit ID"},
			}
		}
		// For approval summary, don't show any rows since we're showing text input
		return []table.Row{}
	default:
		return []table.Row{}
	}
}

// getAuthMethodDescription returns a description for an authentication method
func getAuthMethodDescription(providerName, method string) string {
	descriptions := map[string]map[string]string{
		"AWS": {
			"profile": "Use AWS profile from ~/.aws/credentials",
		},
		"Azure": {
			"cli":        "Use Azure CLI authentication",
			"config-dir": "Use Azure configuration directory",
		},
		"GCP": {
			"service-account": "Use GCP service account key file",
			"adc":             "Use Application Default Credentials",
		},
	}

	if providerDescriptions, ok := descriptions[providerName]; ok {
		if description, ok := providerDescriptions[method]; ok {
			return description
		}
	}

	return ""
}
