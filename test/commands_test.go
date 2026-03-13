package test

import (
	"testing"

	"github.com/XungungoMarkets/xgg/cmd"
	"github.com/spf13/cobra"
)

func TestRootCommandStructure(t *testing.T) {
	root := cmd.GetRootCommand()

	// Verify root command exists and has correct metadata
	if root.Use == "" {
		t.Error("Root command should have a Use field")
	}
	if root.Short == "" {
		t.Error("Root command should have a Short field")
	}
	if root.Long == "" {
		t.Error("Root command should have a Long field")
	}
}

func TestSubCommandsExist(t *testing.T) {
	root := cmd.GetRootCommand()

	// Check that expected subcommands exist
	expectedCommands := []string{
		"stock",
		"search",
		"history",
		"technical",
		"version",
		"update",
	}

	for _, cmdName := range expectedCommands {
		found := false
		for _, subCmd := range root.Commands() {
			if subCmd.Name() == cmdName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand %q not found", cmdName)
		}
	}
}

func TestStockCommandValidation(t *testing.T) {
	root := cmd.GetRootCommand()

	// Find stock command
	var stockCmd *cobra.Command
	for _, subCmd := range root.Commands() {
		if subCmd.Name() == "stock" {
			stockCmd = subCmd
			break
		}
	}

	if stockCmd == nil {
		t.Fatal("Stock command not found")
	}

	// Validate stock command metadata
	if stockCmd.Use == "" {
		t.Error("Stock command should have Use field")
	}
	if stockCmd.Short == "" {
		t.Error("Stock command should have Short field")
	}
	if stockCmd.Example == "" {
		t.Error("Stock command should have Example field")
	}

	// Check argument validation
	if stockCmd.Args == nil {
		t.Error("Stock command should have argument validation")
	}

	// Test with no arguments should fail
	if err := stockCmd.Args(stockCmd, []string{}); err == nil {
		t.Error("Stock command should require at least one argument")
	}

	// Test with one argument should pass
	if err := stockCmd.Args(stockCmd, []string{"NVDA"}); err != nil {
		t.Errorf("Stock command should accept one argument, got error: %v", err)
	}

	// Test with multiple arguments should pass
	if err := stockCmd.Args(stockCmd, []string{"NVDA", "AAPL", "TSLA"}); err != nil {
		t.Errorf("Stock command should accept multiple arguments, got error: %v", err)
	}
}

func TestSearchCommandValidation(t *testing.T) {
	root := cmd.GetRootCommand()

	// Find search command
	var searchCmd *cobra.Command
	for _, subCmd := range root.Commands() {
		if subCmd.Name() == "search" {
			searchCmd = subCmd
			break
		}
	}

	if searchCmd == nil {
		t.Fatal("Search command not found")
	}

	// Validate search command metadata
	if searchCmd.Use == "" {
		t.Error("Search command should have Use field")
	}
	if searchCmd.Short == "" {
		t.Error("Search command should have Short field")
	}

	// Check argument validation
	if searchCmd.Args == nil {
		t.Error("Search command should have argument validation")
	}

	// Test with no arguments should fail
	if err := searchCmd.Args(searchCmd, []string{}); err == nil {
		t.Error("Search command should require exactly one argument")
	}

	// Test with one argument should pass
	if err := searchCmd.Args(searchCmd, []string{"NVDA"}); err != nil {
		t.Errorf("Search command should accept one argument, got error: %v", err)
	}

	// Test with multiple arguments should fail
	if err := searchCmd.Args(searchCmd, []string{"NVDA", "AAPL"}); err == nil {
		t.Error("Search command should not accept multiple arguments")
	}

	// Check that limit flag exists
	if searchCmd.Flags().Lookup("limit") == nil {
		t.Error("Search command should have limit flag")
	}

	// Check that market-data flag exists
	if searchCmd.Flags().Lookup("market-data") == nil {
		t.Error("Search command should have market-data flag")
	}
}

func TestHistoryCommandValidation(t *testing.T) {
	root := cmd.GetRootCommand()

	// Find history command
	var historyCmd *cobra.Command
	for _, subCmd := range root.Commands() {
		if subCmd.Name() == "history" {
			historyCmd = subCmd
			break
		}
	}

	if historyCmd == nil {
		t.Fatal("History command not found")
	}

	// Validate history command metadata
	if historyCmd.Use == "" {
		t.Error("History command should have Use field")
	}
	if historyCmd.Short == "" {
		t.Error("History command should have Short field")
	}

	// Check argument validation
	if historyCmd.Args == nil {
		t.Error("History command should have argument validation")
	}

	// Test with no arguments should fail
	if err := historyCmd.Args(historyCmd, []string{}); err == nil {
		t.Error("History command should require exactly one argument")
	}

	// Test with one argument should pass
	if err := historyCmd.Args(historyCmd, []string{"NVDA"}); err != nil {
		t.Errorf("History command should accept one argument, got error: %v", err)
	}

	// Test with multiple arguments should fail
	if err := historyCmd.Args(historyCmd, []string{"NVDA", "AAPL"}); err == nil {
		t.Error("History command should not accept multiple arguments")
	}

	// Check that period flag exists
	if historyCmd.Flags().Lookup("period") == nil {
		t.Error("History command should have period flag")
	}
}

func TestTechnicalCommandValidation(t *testing.T) {
	root := cmd.GetRootCommand()

	// Find technical command
	var technicalCmd *cobra.Command
	for _, subCmd := range root.Commands() {
		if subCmd.Name() == "technical" {
			technicalCmd = subCmd
			break
		}
	}

	if technicalCmd == nil {
		t.Fatal("Technical command not found")
	}

	// Validate technical command metadata
	if technicalCmd.Use == "" {
		t.Error("Technical command should have Use field")
	}
	if technicalCmd.Short == "" {
		t.Error("Technical command should have Short field")
	}
	if technicalCmd.Example == "" {
		t.Error("Technical command should have Example field")
	}

	// Check argument validation
	if technicalCmd.Args == nil {
		t.Error("Technical command should have argument validation")
	}

	// Test with no arguments should fail
	if err := technicalCmd.Args(technicalCmd, []string{}); err == nil {
		t.Error("Technical command should require exactly one argument")
	}

	// Test with one argument should pass
	if err := technicalCmd.Args(technicalCmd, []string{"NVDA"}); err != nil {
		t.Errorf("Technical command should accept one argument, got error: %v", err)
	}

	// Test with multiple arguments should fail
	if err := technicalCmd.Args(technicalCmd, []string{"NVDA", "AAPL"}); err == nil {
		t.Error("Technical command should not accept multiple arguments")
	}

	// Check that indicator flag exists
	if technicalCmd.Flags().Lookup("indicator") == nil {
		t.Error("Technical command should have indicator flag")
	}

	// Check that period flag exists
	if technicalCmd.Flags().Lookup("period") == nil {
		t.Error("Technical command should have period flag")
	}
}

func TestVersionCommandValidation(t *testing.T) {
	root := cmd.GetRootCommand()

	// Find version command
	var versionCmd *cobra.Command
	for _, subCmd := range root.Commands() {
		if subCmd.Name() == "version" {
			versionCmd = subCmd
			break
		}
	}

	if versionCmd == nil {
		t.Fatal("Version command not found")
	}

	// Validate version command metadata
	if versionCmd.Use == "" {
		t.Error("Version command should have Use field")
	}
	if versionCmd.Short == "" {
		t.Error("Version command should have Short field")
	}

	// Version command may not have argument validation (Args could be nil)
	if versionCmd.Args != nil {
		// Version command should not require arguments
		if err := versionCmd.Args(versionCmd, []string{}); err != nil {
			t.Errorf("Version command should not require arguments, got error: %v", err)
		}
	}
}

func TestUpdateCommandValidation(t *testing.T) {
	root := cmd.GetRootCommand()

	// Find update command
	var updateCmd *cobra.Command
	for _, subCmd := range root.Commands() {
		if subCmd.Name() == "update" {
			updateCmd = subCmd
			break
		}
	}

	if updateCmd == nil {
		t.Fatal("Update command not found")
	}

	// Validate update command metadata
	if updateCmd.Use == "" {
		t.Error("Update command should have Use field")
	}
	if updateCmd.Short == "" {
		t.Error("Update command should have Short field")
	}

	// Update command may not have argument validation (Args could be nil)
	if updateCmd.Args != nil {
		// Update command should not require arguments
		if err := updateCmd.Args(updateCmd, []string{}); err != nil {
			t.Errorf("Update command should not require arguments, got error: %v", err)
		}
	}
}

func TestRootFlags(t *testing.T) {
	root := cmd.GetRootCommand()

	// Check that persistent flags exist
	expectedFlags := []string{
		"provider",
		"rate-limit",
		"max-retries",
		"retry-delay",
		"timeout",
		"watchlist-type",
		"json",
	}

	for _, flagName := range expectedFlags {
		if root.PersistentFlags().Lookup(flagName) == nil {
			t.Errorf("Root command should have %s flag", flagName)
		}
	}
}

func TestCommandHelpText(t *testing.T) {
	root := cmd.GetRootCommand()

	// Test that all commands have help text
	for _, subCmd := range root.Commands() {
		if subCmd.Long == "" && subCmd.Short == "" {
			t.Errorf("Command %q should have help text", subCmd.Name())
		}
	}
}

func TestCommandExamples(t *testing.T) {
	root := cmd.GetRootCommand()

	// Commands that should have examples
	commandsWithExamples := []string{
		"stock",
		"search",
		"history",
		"technical",
	}

	for _, cmdName := range commandsWithExamples {
		var foundCmd *cobra.Command
		for _, subCmd := range root.Commands() {
			if subCmd.Name() == cmdName {
				foundCmd = subCmd
				break
			}
		}
		if foundCmd == nil {
			t.Errorf("Command %q not found", cmdName)
			continue
		}
		if foundCmd.Example == "" {
			t.Errorf("Command %q should have an example", cmdName)
		}
	}
}
