package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/urfave/cli/v2"
	"github.com/mitchellh/colorstring"
)

func fetchParametersByPath(path string) ([]*ssm.Parameter, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String("ap-northeast-1"),
		},
	})
	if err != nil {
		return []*ssm.Parameter{}, err
	}

	svc := ssm.New(sess)
	res, err := svc.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return []*ssm.Parameter{}, err
	}

	return res.Parameters, nil
}

func findParamByPathFromParams(path string, params []*ssm.Parameter) (bool, *string) {
	for _, p := range params {
		if trimPath(*p.Name) == trimPath(path) {
			return true, p.Value
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
	tParams, err := fetchParametersByPath(tPath)
	if err != nil {
		return err
	}
	bParams, err := fetchParametersByPath(bPath)
	if err != nil {
		return err
	}

	for _, t := range tParams {
		exists, value := findParamByPathFromParams(*t.Name, bParams)
		if !exists {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("%s\t%s:\t%s", colorstring.Color("[green] [+]"), trimPath(*t.Name), colorstring.Color("[green]" + *t.Value)))
		} else if *t.Value != *value {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("%s\t%s:\t%s %s", colorstring.Color("[magenta] [~]"), trimPath(*t.Name), colorstring.Color("[green]" + *t.Value), colorstring.Color("[red]" + *value)))
		}
	}

	for _, b := range bParams {
		exists, _ := findParamByPathFromParams(*b.Name, tParams)
		if !exists {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("%s\t%s:\t%s", colorstring.Color("[red] [-]"), trimPath(*b.Name), colorstring.Color("[red]" + *b.Value)))
		}
	}

	return nil
}

func showParametersAction(c *cli.Context) error {
	path := c.String("path")
	params, err := fetchParametersByPath(path)
	if err != nil {
		return err
	}

	for _, p := range params {
		log.Printf("%s: %s", *p.Name, *p.Value)
	}

	return nil
}

func main() {

	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "paramstore"
	app.Usage = "paramstore is a command line tool for AWS Parameter Store"
	app.Action = showParametersAction
	app.Commands = []*cli.Command{
		{
			Name: "show",
			Usage: "show parameters",
			Action: showParametersAction,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "path, p",
					Usage: "path to show like /service/development_3rd/",
					Value: "",
				},
			},
		},
		{
			Name: "comp",
			Usage: "compare parameters",
			Action: compParametersAction,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "target, t",
					Usage: "target path to compare like /service/development_3rd/",
					Value: "",
				},
				&cli.StringFlag{
					Name: "base, b",
					Usage: "base path to compare like /service/development_3rd/",
					Value: "",
				},
			},
		},
	}
	app.Run(os.Args)
}
