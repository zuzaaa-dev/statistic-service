package repository_interface

import "time"

// ProfileResponseStatistics представляет статистику задач на определённую дату.
// @Description Информация о количестве задач для конкретной даты.
type ProfileResponseStatistics struct {
	Date            time.Time `json:"date" example:"2024-10-26T00:00:00Z"`
	SuccessCount    int64     `json:"success_count" example:"5"`
	InProgressCount int64     `json:"in_progress_count" example:"3"`
}

func (ProfileResponseStatistics) TableName() string {
	return "authentication_participant"
}

type LeaderBoardStatistics struct {
	Place int    `json:"place"`
	FIO   string `json:"fio"`
	Score int    `json:"score"`
}

type StatisticsRepositoryInterface interface {
	GetByUserId(id int) ([]ProfileResponseStatistics, error)
	GetLeaderBoard(period, count int) ([]LeaderBoardStatistics, error)
}
