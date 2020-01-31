package main

import (
	"errors"

	"github.com/ne2blink/antenna/pkg/storage"
	"github.com/spf13/cobra"
)

var (
	errMissingName = errors.New("missing name")
)

func newAppCommand() *cobra.Command {
	appCmd := &cobra.Command{
		Use:   "app",
		Short: "manages applications",
		Args:  cobra.NoArgs,
	}
	{
		cmd := &cobra.Command{
			Use:   "create",
			Short: "create application",
			Args:  cobra.NoArgs,
			RunE:  createApp,
		}
		flags := cmd.Flags()
		flags.StringP("name", "n", "", "app name (required)")
		flags.StringP("secret", "s", "", "app secret")
		cmd.MarkFlagRequired("name")
		appCmd.AddCommand(cmd)
	}
	{
		cmd := &cobra.Command{
			Use:   "update",
			Short: "update application",
			Args:  cobra.NoArgs,
			RunE:  updateApp,
		}
		flags := cmd.Flags()
		flags.StringP("id", "i", "", "app ID (required)")
		cmd.MarkFlagRequired("id")
		flags.StringP("name", "n", "", "new app name")
		flags.StringP("secret", "s", "", "new app secret (combined with -r)")
		flags.BoolP("rotate-secret", "r", false, "rotate app secret")
		appCmd.AddCommand(cmd)
	}
	{
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "delete application",
			Args:  cobra.NoArgs,
			RunE:  deleteApp,
		}
		flags := cmd.Flags()
		flags.StringP("id", "i", "", "app ID (required)")
		cmd.MarkFlagRequired("id")
		appCmd.AddCommand(cmd)
	}
	{
		cmd := &cobra.Command{
			Use:   "list",
			Short: "list applications",
			Args:  cobra.NoArgs,
			RunE:  listApps,
		}
		appCmd.AddCommand(cmd)
	}
	return appCmd
}

func createApp(cmd *cobra.Command, _ []string) error {
	name := flagString(cmd, "name")
	if name == "" {
		return errMissingName
	}

	// Construct App
	app := storage.App{
		Name: name,
	}
	secret, err := app.SetSecret(flagString(cmd, "secret"))
	if err != nil {
		return err
	}

	// Create App
	store, err := getStore(cmd)
	if err != nil {
		return err
	}
	defer store.Close()
	id, err := store.CreateApp(app)
	if err != nil {
		return err
	}

	// Reconstruct App with raw secret
	app = storage.App{
		ID:     id,
		Name:   name,
		Secret: secret,
	}

	return printJSON(app)
}

func updateApp(cmd *cobra.Command, _ []string) error {
	// Fetch App
	store, err := getStore(cmd)
	if err != nil {
		return err
	}
	defer store.Close()
	app, err := store.GetApp(flagString(cmd, "id"))
	if err != nil {
		return err
	}

	// Update App
	updated := false
	var secret string
	if name := flagString(cmd, "name"); name != "" {
		app.Name = name
		updated = true
	}
	if rotate := flagBool(cmd, "rotate-secret"); rotate {
		secret, err = app.SetSecret(flagString(cmd, "secret"))
		if err != nil {
			return err
		}
		updated = true
	}
	if updated {
		if err := store.UpdateApp(app); err != nil {
			return err
		}
	}
	app.Secret = secret

	return printJSON(app)
}

func deleteApp(cmd *cobra.Command, _ []string) error {
	// Fetch App
	store, err := getStore(cmd)
	if err != nil {
		return err
	}
	defer store.Close()
	app, err := store.GetApp(flagString(cmd, "id"))
	if err != nil {
		return err
	}

	// Delete App
	if err := store.DeleteApp(app.ID); err != nil {
		return err
	}

	app.Secret = ""
	return printJSON(app)
}

func listApps(cmd *cobra.Command, _ []string) error {
	// Fetch App
	store, err := getStore(cmd)
	if err != nil {
		return err
	}
	defer store.Close()
	apps, err := store.ListApps()
	if err != nil {
		return err
	}

	// Sanitize apps
	for i := range apps {
		apps[i].Secret = ""
	}

	return printJSON(apps)
}
