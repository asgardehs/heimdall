package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/asgardehs/heimdall"
	"github.com/spf13/cobra"
)

func openHeimdall() (*heimdall.Heimdall, error) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))
	return heimdall.Open(logger)
}

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
	}
	cmd.AddCommand(
		configGetCmd(),
		configSetCmd(),
		configListCmd(),
		configResetCmd(),
	)
	return cmd
}

func configGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <namespace> <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := openHeimdall()
			if err != nil {
				return err
			}
			defer h.Close()

			entry, err := h.Get(args[0], args[1])
			if err != nil {
				return err
			}

			displayValue := maskSecret(entry.Type, entry.Value)
			fmt.Printf("%s.%s = %s  (source: %s, type: %s)\n", args[0], args[1], displayValue, entry.Source, entry.Type)
			return nil
		},
	}
}

func configSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <namespace> <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := openHeimdall()
			if err != nil {
				return err
			}
			defer h.Close()

			if err := h.Set(args[0], args[1], args[2]); err != nil {
				return err
			}

			fmt.Printf("%s.%s = %s\n", args[0], args[1], args[2])
			return nil
		},
	}
}

func configListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <namespace>",
		Short: "List configuration values for a namespace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := openHeimdall()
			if err != nil {
				return err
			}
			defer h.Close()

			entries, err := h.List(args[0])
			if err != nil {
				return err
			}

			if len(entries) == 0 {
				fmt.Printf("No config entries for namespace %q.\n", args[0])
				return nil
			}

			for _, e := range entries {
				displayValue := maskSecret(e.Type, e.Value)
				fmt.Printf("  %-30s = %-30s  (%s)\n", e.Key, displayValue, e.Source)
			}
			return nil
		},
	}
}

func configResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset <namespace> <key>",
		Short: "Reset a configuration value to its default",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := openHeimdall()
			if err != nil {
				return err
			}
			defer h.Close()

			if err := h.Reset(args[0], args[1]); err != nil {
				return err
			}

			fmt.Printf("%s.%s reset to default.\n", args[0], args[1])
			return nil
		},
	}
}

// maskSecret redacts secret values for display, showing only the first and
// last 4 characters for values long enough to be meaningful.
func maskSecret(valueType, value string) string {
	if valueType != "secret" || value == "" {
		return value
	}
	if len(value) > 8 {
		return value[:4] + "****" + value[len(value)-4:]
	}
	return "****"
}
