package main

import (
	"log"
	"os"

	"github.com/oidc-proxy-ecosystem/proxy-server/cmd"
	"github.com/oidc-proxy-ecosystem/proxy-server/config"
	"github.com/oidc-proxy-ecosystem/proxy-server/database"
	"github.com/oidc-proxy-ecosystem/proxy-server/plugins/storage"
)

var (
	Version  = ""
	Revision = ""
)

func main() {
	config.NewSettingFile()
	storageConfig := config.NewStorageConfig(config.File.Storage)
	err := storage.New(storageConfig.Type, database.Config{
		UserName:  storageConfig.Username,
		Password:  storageConfig.Password,
		Endpoints: storageConfig.Endpoints,
		TLSConfig: storageConfig.TLSConfig,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer storage.Client.Close()
	cmd := cmd.Command()
	cmd.Version = Version + " - " + Revision
	cmd.Usage = "opendid connect proxy server"
	if err := cmd.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
