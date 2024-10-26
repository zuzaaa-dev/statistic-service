package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"statistic-service/config"
	"statistic-service/internal/infrastructure/database"
)

type PostgresConnectable interface {
	database.DatabaseConnectable
	CreateTables(db *gorm.DB, models ...interface{}) error
}

type postgresConnect struct {
	PostgresConnectable
	cfg *config.Config
}

func NewPostgresConnect(cfg *config.Config) PostgresConnectable {
	return &postgresConnect{cfg: cfg}
}

// Connect возвращает тип *gorm.DB, но мне лень писать дженерики, поэтому возвращаю интерфейс!
func (p *postgresConnect) Connect() (interface{}, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", p.cfg.DatabaseUser,
		p.cfg.DatabasePassword, p.cfg.DatabaseHost, p.cfg.DatabasePort, p.cfg.DatabaseName)
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CloseConnection принимает на вход *gorm.DB и закрывает подключение к БД
func (p *postgresConnect) CloseConnection(inputedDb interface{}) error {
	client := inputedDb.(*gorm.DB)
	dbClient, err := client.DB()
	if err != nil {
		return err
	}
	return dbClient.Close()
}

// CreateTables создает таблицы в базе данных
func (p *postgresConnect) CreateTables(db *gorm.DB, models ...interface{}) error {
	for _, model := range models {
		err := db.Migrator().CreateTable(&model)
		if err != nil {
			return err
		}
	}
	return nil
}
