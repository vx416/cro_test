package main

import (
	"cro_test/cmd"
	"fmt"
	"os"

	cobra "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "server"}

func main() {
	rootCmd.AddCommand(cmd.ServerCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
