package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sonalgaud12/grater/internal/compare"
	"github.com/sonalgaud12/grater/internal/gather"
	"github.com/sonalgaud12/grater/internal/test"
	"github.com/sonalgaud12/grater/pkg/client"
	"github.com/sonalgaud12/grater/pkg/report"
	"github.com/spf13/cobra"
)

var (
	outputFile     string
	maxDependents  int
	timeout        int
	workDir        string
	dependentsFile string
)

var rootCmd = &cobra.Command{
	Use:   "impact",
	Short: "Analyze the impact of changes on downstream Go dependencies",
	Long: `A tool to help Go library maintainers understand the impact of changes
on downstream consumers by running test suites and comparing results.`,
}

var gatherCmd = &cobra.Command{
	Use:   "gather [module-path]",
	Short: "Gather list of dependents for a module",
	Args:  cobra.ExactArgs(1),
	RunE:  runGather,
}

var testCmd = &cobra.Command{
	Use:   "test [module-path] [version]",
	Short: "Run test suites of dependents",
	Args:  cobra.ExactArgs(2),
	RunE:  runTest,
}

var compareCmd = &cobra.Command{
	Use:   "compare [results-file-1] [results-file-2]",
	Short: "Compare two test suite results",
	Args:  cobra.ExactArgs(2),
	RunE:  runCompare,
}

func init() {
	gatherCmd.Flags().StringVarP(&outputFile, "output", "o", "dependents.json", "Output file for dependents list")
	gatherCmd.Flags().IntVarP(&maxDependents, "max", "m", 50, "Maximum number of dependents to gather")

	testCmd.Flags().StringVarP(&outputFile, "output", "o", "results.json", "Output file for test results")
	testCmd.Flags().IntVarP(&timeout, "timeout", "t", 300, "Timeout per test in seconds")
	testCmd.Flags().StringVarP(&workDir, "workdir", "w", "/tmp/impact", "Working directory for cloning repos")
	testCmd.Flags().StringVarP(&dependentsFile, "dependents", "d", "dependents.json", "File containing list of dependents")

	compareCmd.Flags().StringVarP(&outputFile, "output", "o", "comparison.json", "Output file for comparison report")

	rootCmd.AddCommand(gatherCmd, testCmd, compareCmd)
}

func runGather(cmd *cobra.Command, args []string) error {
	modulePath := args[0]

	pkgClient := client.NewPkgGoDevClient()
	svc := gather.NewService(pkgClient)

	opts := gather.Options{
		ModulePath: modulePath,
		MaxResults: maxDependents,
		OutputFile: outputFile,
		Client:     pkgClient,
	}

	_, err := svc.Gather(opts)
	return err
}

func runTest(cmd *cobra.Command, args []string) error {
	modulePath := args[0]
	version := args[1]

	dependents, err := gather.LoadDependents(dependentsFile)
	if err != nil {
		return err
	}

	runner := test.NewRunner(workDir, time.Duration(timeout)*time.Second)

	opts := test.Options{
		ModulePath: modulePath,
		Version:    version,
		Dependents: dependents,
		WorkDir:    workDir,
		Timeout:    time.Duration(timeout) * time.Second,
		OutputFile: outputFile,
	}

	results, err := runner.RunAll(opts)
	if err != nil {
		return err
	}

	report.PrintTestSummary(results)
	return nil
}

func runCompare(cmd *cobra.Command, args []string) error {
	baseFile := args[0]
	newFile := args[1]

	fmt.Printf("Comparing test results:\n")
	fmt.Printf("  Base: %s\n", baseFile)
	fmt.Printf("  New:  %s\n", newFile)

	baseResults, err := test.LoadResults(baseFile)
	if err != nil {
		return err
	}

	newResults, err := test.LoadResults(newFile)
	if err != nil {
		return err
	}

	svc := compare.NewService()
	compReport := svc.Compare(baseResults, newResults)

	if err := svc.SaveReport(compReport, outputFile); err != nil {
		return err
	}

	report.PrintComparisonSummary(compReport)
	fmt.Printf("\nDetailed report written to: %s\n", outputFile)
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
