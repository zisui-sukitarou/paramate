/*
Copyright © 2023 NAME HERE zisuisukitarou@gmail.com

*/
package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"
)

const version = "1.0.0"

type Secret struct {
	Name      string `json:"name"`
	ValueFrom string `json:"valueFrom"`
	Value     string `json:"value"`
}

func fetchParametersByPath(path string, region string) ([]Secret, error) {
	ctx := context.Background()
	var secrets []Secret
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	p := ssm.NewGetParametersByPathPaginator(ssm.NewFromConfig(cfg), &ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		MaxResults:     aws.Int32(10),
		WithDecryption: aws.Bool(true),
		Recursive:      aws.Bool(true),
	})

	for p.HasMorePages() {
		params, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, param := range params.Parameters {
			secrets = append(secrets, Secret{
				Name:      strings.Replace(*param.Name, path+"/", "", 1),
				ValueFrom: *param.Name,
				Value:     *param.Value,
			})
		}
	}
	return secrets, nil
}

func findParamByPathFromParams(path string, params []Secret) (bool, *string) {
	for _, p := range params {
		if trimPath(p.ValueFrom) == trimPath(path) {
			return true, &p.Value
		}
	}
	return false, nil
}

/* /service/development_3rd/XXX -> XXX */
func trimPath(path string) string {
	return path[strings.LastIndex(path, "/")+1:]
}


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "paramate",
	Version: version,
	Short: "paramstore is a command line tool for AWS Parameter Store",
	Long: "paramstore is a command line tool for AWS Parameter Store",
}


func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("region", "r", "ap-northeast-1", "AWS Region")
}
