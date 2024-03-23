package cmd

import (
	"fmt"
	"log"
	"time"
	"totoro1302/goparser/services"

	"github.com/spf13/cobra"
)

var (
	filepath   string
	entityType string
	rootCmd    = &cobra.Command{
		Use:   "fastload",
		Short: "fastload is a discogs xml dump importer",
		Long: `fastload can import discogs dumps (monthly dumps
			for artists, albums, labels`,
		Run: func(cmd *cobra.Command, args []string) {
			switch entityType {
			case "label":
				services.ParseLabel(filepath)
			case "artist":
				services.ParseArtist(filepath)
			default:
				log.Fatalln("Invalid entity-type parameter")
			}
		},
		Version: "0.0.1",
	}
)

func Execute() {
	defer timer("parser")()
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&filepath, "filepath", "f", "", "Path to xml file")
	rootCmd.MarkPersistentFlagRequired("filepath")
	rootCmd.PersistentFlags().StringVarP(&entityType, "entity-type", "t", "", "Entity type to import (label, artist, album)")
	rootCmd.MarkPersistentFlagRequired("entity-type")
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
