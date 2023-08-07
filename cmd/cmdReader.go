package cmd

import (
	"github.com/NickLovera/KzStats/db"
	"github.com/NickLovera/KzStats/mgr"
	"github.com/NickLovera/KzStats/routes"
	"github.com/NickLovera/KzStats/utils"
	mongo "github.com/NickLovera/mongo-utils"
	rest "github.com/NickLovera/rest-utils-go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	v "github.com/spf13/viper"
	"log"
)

//TODO update swagger yml to correct data types and field name formats
var viper = v.New()
var rootCmd = &cobra.Command{
	Use:   "Kz Stats",
	Short: "Kz Stat retriever",
	Long:  "Retrieve Kz stats from kzstats.com for fun",
	Run: func(cmd *cobra.Command, args []string) {

		dbConn, dbErr := mongo.ConnectAndMigrate(nil, "KzStats")
		if dbErr != nil {
			log.Panicln("Cannot connect to mongo db: ", dbErr)
		}

		log.Println("Migration Done")

		mongoClient := mongo.NewMongoClient(dbConn)
		log.Println("Connected to mongo db")

		kzRepo := db.NewKzRepo("KzStats", mongoClient)
		utils := utils.NewUtils()
		requestService := mgr.NewKzService(kzRepo, utils)
		log.Println("Created services and repos")

		//Initialize routin env vars

		rest.InitWebConfig(nil)

		routes.InitServer(requestService)
	},
}

func init() {
	//get env vars
	BindFlagWithCmd(rootCmd, mongo.InitMongoWithViper)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

//TODO move to library
func BindFlagWithCmd(cmd *cobra.Command, inits ...func() (*v.Viper, *pflag.FlagSet)) {
	for _, init := range inits {
		b, f := init()
		cmd.Flags().AddFlagSet(f)
		_ = b.BindPFlags(cmd.Flags())
	}
}
