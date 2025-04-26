package main

import (
	"copier/internal/copy"
	"copier/internal/log"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func main() {
	var logFile string
	var ignoreFile string

	var rootCmd = &cobra.Command{
		Use:   "copier <source> <destination>",
		Short: "Recursively copy a folder, skipping files/dirs via .copyignore, logging errors/skips.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			src := args[0]
			dst := args[1]

			if logFile == "" {
				logFile = "copy.log"
			}
			if ignoreFile == "" {
				ignoreFile = filepath.Join(src, ".copyignore")
			}

			logger, err := log.NewLogger(logFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
				os.Exit(1)
			}
			defer logger.Close()

			err = copy.CopyDirWithIgnore(src, dst, logger, ignoreFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Copy completed with errors. See log for details.\n")
			} else {
				fmt.Println("Copy completed successfully.")
			}
		},
	}

	rootCmd.Flags().StringVar(&logFile, "log", "", "Path to log file (default: copy.log)")
	rootCmd.Flags().StringVar(&ignoreFile, "ignore", "", "Path to ignore file (default: .copyignore in source)")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
