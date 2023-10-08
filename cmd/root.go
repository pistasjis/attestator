package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var OutputAsJson bool

var rootCmd = &cobra.Command{
	Use:   "attestator",
	Short: "Attestator is an open-source application that allows you to check how apps on your computer respect privacy and security",
	Long:  "Attestator is an that allows you to check how apps on your computer respect privacy and security\nDeveloped by pistasjis with Go.",
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS != "windows" {
			panic("pistasjis Attestator is not supported on " + runtime.GOOS)
		}

		if runtime.GOARCH == "arm64" {
			fmt.Println("The arm64 architecture is not exactly supported (I can't test it), so please report your findings on GitHub at https://github.com/pistasjis/attestator.")
		}

		if runtime.GOARCH == "386" && runtime.GOOS == "windows" {
			fmt.Println("The binary for the 386 (x86) architecture cannot get 64-bit apps, please keep that in mind.")
		}

		RunAttestator(OutputAsJson)
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&OutputAsJson, "json", "j", false, "Output attestation as JSON instead of HTML")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
