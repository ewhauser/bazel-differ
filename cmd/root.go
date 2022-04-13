package cmd

import (
	"fmt"
	"github.com/ewhauser/bazel-differ/internal"
	"github.com/spf13/cobra"
	"os"
)

var WorkspacePath string
var BazelPath string
var StartingHashes string
var FinalHashes string
var Output string
var BazelStartupOptions string
var BazelCommandOptions string
var KeepGoing bool
var Verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bazel-differ",
	Short: "bazel-differ is a CLI tool to assist with doing differential Bazel builds",
	Long:  `bazel-differ is a CLI tool to assist with doing differential Bazel builds`,
}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	WorkspacePath = cwd
	rootCmd.PersistentFlags().StringVarP(&WorkspacePath, "workspacePath", "w", "", "Path to Bazel workspace directory.")
	rootCmd.PersistentFlags().StringVarP(&BazelPath, "bazelPath", "b", "", "Path to Bazel binary")
	rootCmd.PersistentFlags().StringVarP(&BazelStartupOptions, "bazelStartupOptions", "y", "",
		"Additional space separated Bazel client startup options used when invoking Bazel")
	rootCmd.PersistentFlags().StringVarP(&BazelCommandOptions, "bazelCommandOptions", "z", "",
		"Additional space separated Bazel command options used when invoking Bazel")
	rootCmd.PersistentFlags().BoolVarP(&KeepGoing, "keep_going", "k", true,
		"This flag controls if `bazel query` will be executed with the `--keep_going` flag or not. Disabling this flag allows you to catch configuration issues in your Bazel graph, but may not work for some Bazel setups. Defaults to `true`")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "enables verbose output")
}

func GetBazelClient() internal.BazelClient {
	return internal.NewBazelClient(internal.Filesystem, WorkspacePath, BazelPath, Verbose, KeepGoing,
		displayElapsedTime,
		BazelStartupOptions, BazelCommandOptions)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func ExitIfError(err error, msg string) {
	if msg == "" {
		msg = "Error: "
	}

	if err != nil {
		fmt.Printf("%s: %s \n", msg, err)
		os.Exit(1)
	}
}
