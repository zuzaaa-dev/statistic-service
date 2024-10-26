package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"statistic-service/config"
	_ "statistic-service/internal/domain/statistic/usecases/repository_interface"
	interfaceRepo "statistic-service/internal/domain/statistic/usecases/repository_interface"
)

type StatisticHandlers struct {
	cfg *config.Config
	log *slog.Logger
	db  *gorm.DB
}

func NewStatisticHandlers(cfg *config.Config, log *slog.Logger, db *gorm.DB) *StatisticHandlers {
	return &StatisticHandlers{cfg: cfg, log: log, db: db}
}

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
}

// ProfileResponseStatisticsDoc используется для генерации документации
type ProfileResponseStatisticsDoc struct {
	Date            string `json:"date" example:"2024-10-26T00:00:00Z"`
	SuccessCount    int64  `json:"success_count" example:"5"`
	InProgressCount int64  `json:"in_progress_count" example:"3"`
}

// GetByUserId godoc
// @Summary      Получить статистику пользователя по его ID
// @Description  Возвращает статистику задач для пользователя: завершенные и находящиеся в процессе задачи для каждой даты
// @Tags         statistics
// @Param        user_id   path      string  true  "User ID"
// @Success      200       {array}   ProfileResponseStatisticsDoc
// @Failure      400       {object}  ErrorResponse
// @Router       /getStatistic/user/{user_id} [get]
func (h *StatisticHandlers) GetByUserId(c *gin.Context) {
	id := c.Param("user_id")

	var result []interfaceRepo.ProfileResponseStatistics
	query := `
		SELECT 
			DATE(end_date) AS date,
			COUNT(CASE WHEN status = 'завершён' THEN 1 END) AS success_count,
			COUNT(CASE WHEN status = 'в_процессе' THEN 1 END) AS in_progress_count
			FROM 
				authentication_participant AS ap
			JOIN 
				authentication_challenge AS ac ON ap.challenge_id = ac.id AND user_id = ?
			GROUP BY 
				DATE(end_date)
			ORDER BY 
				DATE(end_date);
	`

	if err := h.db.Raw(query, id).Scan(&result).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
