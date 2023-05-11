package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/mitchellh/colorstring"
	"github.com/urfave/cli/v2"
)

const version = "0.0.2"

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
				Name:      strings.Replace(*param.Name, path, "", 1),
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

func compParametersAction(c *cli.Context) error {
	tPath := c.String("target")
	bPath := c.String("base")
	region := c.String("region")

	tParams, err := fetchParametersByPath(tPath, region)
	if err != nil {
		return err
	}
	bParams, err := fetchParametersByPath(bPath, region)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, fmt.Sprintf("%s\t%-10s\t%-10s", "", "Name", "Value"))
	for _, t := range tParams {
		exists, value := findParamByPathFromParams(t.ValueFrom, bParams)
		if !exists {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("%-10s\t%-10s\t%-10s", colorstring.Color("[green][+]"), trimPath(t.ValueFrom), colorstring.Color("[green]"+t.Value)))
		} else if *&t.Value != *value {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("%-10s\t%-10s\t%-10s\t%-10s", colorstring.Color("[magenta][~]"), trimPath(t.ValueFrom), colorstring.Color("[green]"+t.Value), colorstring.Color("[red]"+*value)))
		}
	}

	for _, b := range bParams {
		exists, _ := findParamByPathFromParams(b.ValueFrom, tParams)
		if !exists {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("%-10s\t%-10s\t%-10s", colorstring.Color("[red][-]"), trimPath(b.ValueFrom), colorstring.Color("[red]"+b.Value)))
		}
	}

	return nil
}

func showParametersAction(c *cli.Context) error {
	path := c.String("path")
	region := c.String("region")

	params, err := fetchParametersByPath(path, region)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	fmt.Fprintln(os.Stdout, fmt.Sprintf("%s\t%-10s", "Name", "Value"))
	for _, p := range params {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("%-10s\t%-10s", p.Name, p.Value))
	}

	return nil
}

func main() {

	defaultRegion := "us-west-1"

	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "paramstore"
	app.Usage = "paramstore is a command line tool for AWS Parameter Store"
	app.Action = showParametersAction
	app.Version = version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "region",
			Aliases: []string{"r"},
			Usage:   "aws region",
			Value:   defaultRegion,
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "show",
			Usage:  "show parameters",
			Action: showParametersAction,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "path",
					Aliases: []string{"p"},
					Usage:   "path to show like /service/development_3rd/",
					Value:   "/",
				},
				&cli.StringFlag{
					Name:    "region",
					Aliases: []string{"r"},
					Usage:   "aws region",
					Value:   defaultRegion,
				},
			},
		},
		{
			Name:   "comp",
			Usage:  "compare parameters",
			Action: compParametersAction,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "target",
					Aliases: []string{"t"},
					Usage:   "target path to compare like /service/development_3rd/",
					Value:   "/",
				},
				&cli.StringFlag{
					Name:    "base",
					Aliases: []string{"b"},
					Usage:   "base path to compare like /service/development_3rd/",
					Value:   "",
				},
				&cli.StringFlag{
					Name:    "region",
					Aliases: []string{"r"},
					Usage:   "aws region",
					Value:   defaultRegion,
				},
			},
		},
	}
	app.Run(os.Args)
}
