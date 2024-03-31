package main

import (
	"fmt"
	"os"

	imgtoico "github.com/erancihan/img-to-ico/pkg/img-to-ico"
	"github.com/spf13/cobra"
)

var rootcmd = &cobra.Command{
	Use:   "img-to-ico [SRC]",
	Short: "Convert image to ico",
	Long:  `Convert image to ico`,

	Args: func(cmd *cobra.Command, args []string) error {
		// Optionally run one of the validators provided by cobra
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		converter := &imgtoico.Converter{
			From: args[0],
		}
		converter.Convert()
		converter.Write()
	},
}

func init() {
}

func Execute() {
	if err := rootcmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
