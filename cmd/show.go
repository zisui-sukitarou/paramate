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

func showParametersAction(cmd *cobra.Command, args []string) {
	path := args[0]
	region := cmd.Flag("region").Value.String()

	fmt.Fprintln(os.Stdout, colorstring.Color(fmt.Sprint("[blue]loading parameters from AWS SSM Parameter Store in region:", region, " path:", path, " ...")))

	params, err := fetchParametersByPath(path, region)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Println("")
	fmt.Fprintln(os.Stdout, fmt.Sprintf("%s\t%-10s", "Name", "Value"))
	for _, p := range params {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("%-10s\t%-10s", p.Name, p.Value))
	}
}


// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show [path]",
	Short: "show envs of the path",
	Long: "show envs of the path",
	Args: cobra.ExactArgs(1),
	Run: showParametersAction,
}

func init() {
	showCmd.Flags().StringP("region", "r", "ap-northeast-1", "AWS region")
	rootCmd.AddCommand(showCmd)
}
