package cmd

import (
	"errors"

	"github.com/sawadashota/fpick"
	"github.com/spf13/cobra"
)

// Cmd is Exported for entry point
var Cmd *cobra.Command

func init() {
	Cmd = RootCmd()
}

// RootCmd .
func RootCmd() *cobra.Command {
	var src string
	var dst string
	var filename string
	var regex string
	var flat bool

	var err error
	var matcher fpick.FileMatcher

	cmd := &cobra.Command{
		Use:     "fpick --src <PATH> --dst <PATH>",
		Short:   "search and pick files",
		Version: Version,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if src == "" {
				return errors.New("argument error: --src is empty")
			}
			if dst == "" {
				return errors.New("argument error: --dst is empty")
			}
			if filename != "" && regex == "" {
				matcher = fpick.FilenameExtractMatch(filename)
				return nil
			}
			if filename == "" && regex != "" {
				matcher, err = fpick.FilenameRegexMatch(regex)
				return err
			}
			return errors.New("argument error: either --filename or --regex should be present")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := fpick.New(src, dst)
			if err != nil {
				return err
			}

			var opts []fpick.OutputOption
			if flat {
				opts = append(opts, fpick.OutputFlatDirOption)
			}

			return c.Pick(matcher, opts...)
		},
	}

	cmd.Flags().StringVarP(&src, "src", "s", "", "required. source directory")
	cmd.Flags().StringVarP(&dst, "dst", "d", "", "required. district directory")
	cmd.Flags().StringVarP(&filename, "name", "n", "", "search extract match filename")
	cmd.Flags().StringVarP(&regex, "regex", "r", "", "search regex match filename")
	cmd.Flags().BoolVar(&flat, "flat", false, "output flat directory")

	return cmd
}
