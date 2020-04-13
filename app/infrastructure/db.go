package infrastructure

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

func NewPgPool(config Configuration, logger *zap.Logger) *pgxpool.Pool {
	dbConf := config.DB
	databaseUrl :=
		fmt.Sprintf("user=%s "+
			"password=%s "+
			"host=%s "+
			"port=%d "+
			"dbname=%s "+
			"pool_max_conns=%d", dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Dbname, dbConf.PoolMaxConns)
	poolConfig, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		logger.Fatal("Error creating directory for logs", zap.String("database_url", ""), zap.Error(err))
	}
	poolConfig.ConnConfig.Logger = zapadapter.NewLogger(logger)

	db, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Fatal("Unable to create connection pool", zap.String("database_url", ""), zap.Error(err))
	}

	return db
}
