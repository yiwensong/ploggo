package main

import (
	goflag "flag"
	fmt "fmt"
	os "os"
	path "path"

	glog "github.com/golang/glog"
	errors "github.com/pkg/errors"
	ksuid "github.com/segmentio/ksuid"
	cobra "github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "db-migration",
	Short: "helper for creating db migrations",
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new db migration",
	RunE:  rootCmdExec,
}

var migrationDir *string
var migrationName *string

func init() {
	rootCmd.AddCommand(createCmd)

	migrationDir = createCmd.Flags().StringP(
		"path",
		"p",
		"",
		"the directory in which the new migration file will be in",
	)
	migrationName = createCmd.Flags().StringP(
		"name",
		"n",
		"",
		"the name of the migration",
	)

	rootCmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		glog.Errorf(err.Error())
	}
}

func rootCmdExec(cmd *cobra.Command, args []string) error {
	goflag.Parse()

	glog.Infof("current path: %s\n", os.Getenv("PWD"))
	glog.Infof("migration dir: %s\n", *migrationDir)

	migrationId := ksuid.New()
	migrationPath := path.Join(
		*migrationDir,
		fmt.Sprintf(
			"%s-%s.sql",
			migrationId.String(),
			*migrationName,
		),
	)

	if _, err := os.Stat(migrationPath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(*migrationDir, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "MkdirAll(%q)", *migrationDir)
		}
	}

	glog.Infof("writing file to: %s\n", migrationPath)

	f, err := os.OpenFile(
		migrationPath,
		os.O_CREATE,
		0666,
	)
	if err != nil {
		return errors.Wrapf(err, "OpenFile(%q)", migrationPath)
	}
	defer f.Close()

	fmt.Println(migrationPath)

	return nil
}
