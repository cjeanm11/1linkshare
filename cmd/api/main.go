package main

import (
	"1linkshare/internal/network"
	srv "1linkshare/internal/server"
	"1linkshare/internal/utils"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

func init() {
	if err := utils.DeleteAllFilesInUploadDir(); err != nil {
		log.Println("Error cleaning up upload directory:", err)
	}
}

func commandLineInterface() {
	app := &cli.App{
		Name:  "1linkssh",
		Usage: "Generate secure URLs for file sharing via SSH tunnel",
		Commands: []*cli.Command{
			{
				Name:    "upload",
				Aliases: []string{"u"},
				Usage:   "Upload a file and generate a secure URL",
				Action: func(c *cli.Context) error {

					if c.Args().Len() != 1 {
						return fmt.Errorf("expected one argument: path to file")
					}
					filePath := c.Args().First()
					url := "http://localhost:8080/upload"

					header := lipgloss.NewStyle().
						Foreground(lipgloss.Color("#00FF00")).
						Background(lipgloss.Color("#000000")).
						Bold(true).
						Italic(true).
						Padding(0, 1).
						Render

					body := lipgloss.NewStyle().
						Foreground(lipgloss.Color("#FF8800")).
						Render

					fmt.Println()

					headerApp := lipgloss.NewStyle().
						Foreground(lipgloss.Color("#00008B")).
						Padding(0, 1).
						Bold(true).
						Render("1linkssh")

					fmt.Println(headerApp)

					fmt.Println()

					fmt.Println(header("Uploading file..."))
					if err := network.UploadFile(url, filePath); err != nil {
						return err
					}

					bigSwag := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff88")).
						Padding(0, 1).Render
					secureURL := filepath.Base(filePath)
					fmt.Print(bigSwag("Secure URL generated for file: "), body(secureURL))
					fmt.Println()

					fmt.Println()

					fmt.Println(header("File uploaded successfully!"))
					return nil
				},
			},
			{
				Name:    "share",
				Aliases: []string{"s"},
				Usage:   "Share an existing file using a secure URL",
				Action: func(c *cli.Context) error {
					id := c.Args().First()
					url := fmt.Sprintf("https://example.com/share/%s", id)

					header := lipgloss.NewStyle().
						Foreground(lipgloss.Color("#00FF00")).
						Background(lipgloss.Color("#000000")).
						Bold(true).
						Italic(true).
						Padding(0, 1).
						Render

					body := lipgloss.NewStyle().
						Foreground(lipgloss.Color("#FF8800")).
						Render

					fmt.Println(header("Secure URL for file:"))
					fmt.Println(body(url))

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	port := srv.GetPortOrDefault(8080)
	server := srv.NewServer(srv.WithPort(port))
	go commandLineInterface()
	server.Start()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
}
