package cmd

import (
    "database/sql"
    "fmt"
    "github.com/golang-migrate/migrate"
    "github.com/golang-migrate/migrate/database/mysql"
    "github.com/neotroops/go-realworld/configs"
    "github.com/neotroops/go-realworld/constant"
    "github.com/neotroops/go-realworld/i18n"
    "github.com/neotroops/go-realworld/pkg"
    "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/golang-migrate/migrate/source/file"
    "os"
)

var initCmd = &cobra.Command{
    Use:   "gorealworld",
    Short: "Go lang real world implementation",
    Run: func(cmd *cobra.Command, args []string) {
        appConfig := configs.Config()
        i18n.Init()
        logrus.Info(fmt.Sprintf(constant.INIT_MESSAGE, appConfig.AppName, appConfig.AppPort))
        pkg.StartAPIServer(appConfig)
    },
}

var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Go lang real world implementation",
    Run: func(cmd *cobra.Command, args []string) {
        logrus.SetReportCaller(true)
        appConfig := configs.Config()
        dbUser := appConfig.DbConfig.DbUser
        dbPassword := appConfig.DbConfig.DbPassword
        dbHost := appConfig.DbConfig.DbHost
        dbPort := appConfig.DbConfig.DbPort
        db, _ := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/golang?multiStatements=true")
        driver, _ := mysql.WithInstance(db, &mysql.Config{})
        m, _ := migrate.NewWithDatabaseInstance(
            "file://db/",
            "mysql",
            driver,
        )
        if err := m.Steps(1); err != nil {
            logrus.Panic(err)
            os.Exit(1)
        }

    },
}

func Exec() {
    initCmd.AddCommand(migrateCmd)
    if err := initCmd.Execute(); err != nil {
        logrus.Panic(err)
        os.Exit(1)
    }
}
