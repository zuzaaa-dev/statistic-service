package repository

import (
	"gorm.io/gorm"
	"log/slog"
	"statistic-service/config"
	interfaceRepo "statistic-service/internal/domain/statistic/usecases/repository_interface"
)

type statisticsRepository struct {
	interfaceRepo.StatisticsRepositoryInterface
	cfg *config.Config
	log *slog.Logger
	db  *gorm.DB
}

func NewStatisticsRepository(cfg *config.Config, log *slog.Logger, db *gorm.DB) *statisticsRepository {
	return &statisticsRepository{cfg: cfg, log: log, db: db}
}

func (s *statisticsRepository) GetByUserId(id int) ([]interfaceRepo.ProfileResponseStatistics, error) {
	var result []interfaceRepo.ProfileResponseStatistics
	if err := s.db.Where("user_id = ?", id).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
