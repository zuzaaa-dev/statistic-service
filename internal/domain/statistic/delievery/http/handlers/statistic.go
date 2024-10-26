package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"statistic-service/config"
	_ "statistic-service/internal/domain/statistic/usecases/repository_interface"
	interfaceRepo "statistic-service/internal/domain/statistic/usecases/repository_interface"
	"strconv"
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

type LeaderBoardStatistics struct {
	Place int    `json:"place"`
	FIO   string `json:"fio"`
	Score int    `json:"score"`
}

// GetLeaderBoard godoc
// @Summary      Получить топ активных участников за указанный период
// @Description  Возвращает список топ-N самых активных участников за последние X дней, отсортированный по количеству завершенных задач.
// @Tags         statistics
// @Param        period   path      int  true  "Период в днях для определения активности (например, 7 для последней недели)"
// @Param        count    path      int  true  "Количество записей в списке топа (например, 10 для топ-10)"
// @Success      200      {array}   LeaderBoardStatistics  "Список участников с их местом, ФИО и количеством завершенных задач"
// @Failure      400      {object}  ErrorResponse "Описание ошибки в случае некорректного запроса"
// @Router       /get-leaderboard/{period}/{count} [get]
func (h *StatisticHandlers) GetLeaderBoard(c *gin.Context) {
	// Получаем параметры из запроса
	periodParam := c.Param("period")
	countParam := c.Param("count")

	// Конвертируем параметры в int
	period, err := strconv.Atoi(periodParam)
	if err != nil {
		period = 7 // Если параметр отсутствует или некорректен, используем значение по умолчанию
	}

	count, err := strconv.Atoi(countParam)
	if err != nil {
		count = 10 // Если параметр отсутствует или некорректен, используем значение по умолчанию
	}

	// Выполняем запрос
	var leaderboard []interfaceRepo.LeaderBoardStatistics
	query := `
        SELECT 
            ROW_NUMBER() OVER(ORDER BY COUNT(ap.challenge_id) DESC) AS place,
            u.full_name AS fio,
            COUNT(ap.challenge_id) AS score
        FROM 
            authentication_participant AS ap
        JOIN 
            authentication_challenge AS ac ON ap.challenge_id = ac.id
        JOIN 
            authentication_user AS u ON ap.user_id = u.id
        WHERE 
            ac.end_date >= CURRENT_DATE - ($1 * INTERVAL '1 day')
        GROUP BY 
            u.id, u.full_name
        ORDER BY 
            score DESC
        LIMIT $2;
    `

	if err := h.db.Raw(query, period, count).Scan(&leaderboard).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}
