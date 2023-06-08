/*
Copyright © 2023 NAME HERE zisuisukitarou@gmail.com

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/colorstring"
	"github.com/spf13/cobra"
)

func compParametersAction(cmd *cobra.Command, args []string) {
	path1 := args[0]
	path2 := args[1]
	region := cmd.Flag("region").Value.String()

	params1, err := fetchParametersByPath(path1, region)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	params2, err := fetchParametersByPath(path2, region)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, fmt.Sprintf("%s\t%-10s\t%-10s", "", "Name", "Value"))
	for _, t := range params1 {
		exists, value := findParamByPathFromParams(t.ValueFrom, params2)
		if !exists {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("%-10s\t%-10s\t%-10s", colorstring.Color("[green][+]"), trimPath(t.ValueFrom), colorstring.Color("[green]"+t.Value)))
		} else if *&t.Value != *value {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("%-10s\t%-10s\t%-10s\t%-10s", colorstring.Color("[magenta][~]"), trimPath(t.ValueFrom), colorstring.Color("[green]"+t.Value), colorstring.Color("[red]"+*value)))
		}
	}

	for _, b := range params2 {
		exists, _ := findParamByPathFromParams(b.ValueFrom, params1)
		if !exists {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("%-10s\t%-10s\t%-10s", colorstring.Color("[red][-]"), trimPath(b.ValueFrom), colorstring.Color("[red]"+b.Value)))
		}
	}
}

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "show difference in envs of the two paths",
	Long: `show difference in envs of the two paths`,
	Run: compParametersAction,
}

func init() {
	diffCmd.Flags().StringP("region", "r", "ap-northeast-1", "AWS region")
	rootCmd.AddCommand(diffCmd)
}
