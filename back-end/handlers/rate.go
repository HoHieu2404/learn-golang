package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	"golang/backend/models"
	"golang/backend/services"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mitchellh/mapstructure"
	gomail "gopkg.in/mail.v2"
	"github.com/gorilla/mux"
)

type ListEmail struct {
	Email []string
}

type RatesByDateResp struct {
	Date  string             `json:"date"`
	Rates map[string]float32 `json:"rates"`
}

type HandlerInterface interface {
	GetRatesLatest(response http.ResponseWriter, request *http.Request)
	GetRatesByDate(response http.ResponseWriter, request *http.Request)
	GetAnalysisRates(response http.ResponseWriter, request *http.Request)
	SyncData()
}

type Handler struct {
	Service services.ServiceInterface
}

func NewHandler(service services.ServiceInterface) HandlerInterface {
	return &Handler{Service: service}
}

func (h *Handler) GetRatesLatest(response http.ResponseWriter, request *http.Request) {
	resp, err := h.Service.GetRatesLatest()
	if err != "" {
		// responseWithError(response, http.StatusBadRequest, err)
		json.NewEncoder(response).Encode(err)
		return
	}
	// responseWithJson(response, http.StatusOK, resp)
	json.NewEncoder(response).Encode(resp)
	return

}

func (h *Handler) GetRatesByDate(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	date := vars["date"]
	resp, err := h.Service.GetRatesByDate(date)
	if err != "" {
		json.NewEncoder(response).Encode(err)
		// responseWithError(response, http.StatusBadRequest, err)
		return
	}
	// responseWithJson(response, http.StatusOK, resp)
	json.NewEncoder(response).Encode(resp)

}

func (h *Handler) GetAnalysisRates(response http.ResponseWriter, request *http.Request) {
	resp, err := h.Service.GetRatesAnalyze()
	if err != "" {
		// responseWithError(response, http.StatusBadRequest, err)
		json.NewEncoder(response).Encode(err)
		return
	}
	// responseWithJson(response, http.StatusOK, resp)
	json.NewEncoder(response).Encode(resp)
	return


}

func (h *Handler) SyncData() {
	startDate := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")
	url := "https://api.apilayer.com/exchangerates_data/timeseries?start_date="+startDate+"&end_date=" + endDate
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic("Error: Failed to get api")
	}
	req.Header.Set("apikey", "Jj8J8G1gN02VNhLnjAAwk4BsJrzDdHp1")
	res, _ := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)

	var data models.Data
	json.Unmarshal(body, &data)
	err = h.Service.ImportDataInit(&data)
	if err != nil {
		log.Println("Sync Data: ", err)
	}
	fmt.Println("Updated data at ", time.Now())
}

// func responseWithJson(response http.ResponseWriter, statusCode int, data interface{}) {
// 	result, _ := json.Marshal(data)
// 	response.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
// 	// response.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
// 	response.WriteHeader(statusCode)
// 	response.Write(result)
// }

// func responseWithError(response http.ResponseWriter, statusCode int, msg string) {
// 	responseWithJson(response, statusCode, map[string]string{
// 		"error": msg,
// 	})
// }
