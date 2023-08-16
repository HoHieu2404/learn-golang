package repositories

import (
	"database/sql"
	"learn-golang/back-end/models"
	"time"
)

type RatesByDateResp struct {
	Date  string             `json:"date"`
	Rates map[string]float32 `json:"rates"`
}

type RatesAnalyzeResp struct {
	Rates map[string]RatesAnalyze `json:"rates"`
}

type RatesAnalyze struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
	Avg float32 `json:"avg"`
}

type RepositoryInterface interface {
	GetRatesLatest() (RatesByDateResp, string)
	GetRatesByDate(date string) (RatesByDateResp, string)
	GetRatesAnalyze() (RatesAnalyzeResp, string)
	ImportDataInit(data *models.Data) error
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) RepositoryInterface {
	return &Repository{db}
}

func (r *Repository) GetRatesLatest() (RatesByDateResp, string) {
	yesterday := time.Now().Format("2006-01-02")
	res, err := r.db.Query("SELECT currency, rate FROM Rates WHERE date=?", yesterday)
	if err != nil {
		error := "Could not get currency at" + yesterday
		return RatesByDateResp{}, error
	}

	var currency string
	var rate float32
	rates := make(map[string]float32)

	for res.Next() {
		err = res.Scan(&currency, &rate)
		if err != nil {
			error := "Something went wrong"
			return RatesByDateResp{}, error
		}
		rates[currency] = rate
	}
	data := RatesByDateResp{
		Date:  yesterday,
		Rates: rates,
	}
	return data, ""
}

func (r *Repository) GetRatesByDate(date string) (RatesByDateResp, string) {
	res, err := r.db.Query("SELECT currency, rate FROM Rates WHERE date=?", date)
	if err != nil {
		error := "Could not get rates on " + date
		return RatesByDateResp{}, error
	}

	var currency string
	var rate float32
	rates := make(map[string]float32)

	for res.Next() {
		err = res.Scan(&currency, &rate)
		if err != nil {
			error := "Something went wrong"
			return RatesByDateResp{}, error
		}
		rates[currency] = rate
	}
	data := RatesByDateResp{
		Date:  date,
		Rates: rates,
	}
	if len(rates) == 0 {
		error := "Could not get rates on " + date
		return RatesByDateResp{}, error
	}
	return data, ""
}

func (r *Repository) GetRatesAnalyze() (RatesAnalyzeResp, string) {
	res, err := r.db.Query("SELECT currency, max(rate) as max, min(rate) as min, (sum(rate)/count(currency)) as avg FROM Rates GROUP BY currency")
	if err != nil {
		error := "Could not get rates analyze"
		return RatesAnalyzeResp{}, error
	}

	var currency string
	var max float32
	var min float32
	var avg float32

	currencies := map[string]RatesAnalyze{}

	for res.Next() {
		err = res.Scan(&currency, &max, &min, &avg)
		if err != nil {
			error := "Something went wrong"
			return RatesAnalyzeResp{}, error
		}
		currencies[currency] = RatesAnalyze{
			Min: min,
			Max: max,
			Avg: avg,
		}
	}
	data := RatesAnalyzeResp{
		Rates: currencies,
	}
	return data, ""
}

func (r *Repository) ImportDataInit(data *models.Data) error {
	r.db.Query("DELETE * FROM Rates")
	for date, currencies := range data.Rates {
		for currency, rate := range currencies {
			_, err := r.db.Exec(`INSERT INTO Rates(date, currency, rate) VALUES( ?, ?, ?)`, date, currency, rate)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
