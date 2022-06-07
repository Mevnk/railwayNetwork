package cmd

import (
	"github.com/spf13/viper"
	"os"
	"railwayNetwork/pkg/repo"
)

//import (
//)

func main() {
	db, err := repo.NewMysqlDB(repo.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		Password: os.Getenv("DB_PASSWORD"),
	})
}
