package services

import (
	"learn-golang/back-end/models"
	repo "learn-golang/back-end/repositories"
)

type ServiceInterface interface {
	GetRatesLatest() (interface{}, string)
	GetRatesByDate(date string) (interface{}, string)
	GetRatesAnalyze() (interface{}, string)
	ImportDataInit(data *models.Data) error
}

type Service struct {
	repo repo.RepositoryInterface
}

func NewService(repo repo.RepositoryInterface) ServiceInterface {
	return &Service{repo: repo}
}

func (s *Service) GetRatesLatest() (interface{}, string) {
	return s.repo.GetRatesLatest()
}

func (s *Service) GetRatesByDate(date string) (interface{}, string) {
	return s.repo.GetRatesByDate(date)
}

func (s *Service) GetRatesAnalyze() (interface{}, string) {
	return s.repo.GetRatesAnalyze()
}

func (s *Service) ImportDataInit(data *models.Data) error {
	return s.repo.ImportDataInit(data)
}
