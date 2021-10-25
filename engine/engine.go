package engine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	vegeta "github.com/tsenart/vegeta/lib"
	"math"
	"net/http"
	"strconv"
	"time"
)

type attackRequest struct {
	Method         string `json:"method" binding:"required"`
	Url            string `json:"url" binding:"required"`
	AttackDuration string `json:"attack_duration" binding:"required"`
	AttackRate     string `json:"attack_rate" binding:"required"`
	PayLoad        string `json:"pay_load"`
	PassRate       int64  `json:"pass_rate"`
}

type latencyResponse struct {
	Total float64 `json:"total"`
	Mean  float64 `json:"mean"`
}

type attackResponse struct {
	Latencies   latencyResponse `json:"latencies"`
	Duration    float64         `json:"duration"`
	Wait        float64         `json:"wait"`
	Requests    float64         `json:"requests"`
	Throughput  float64         `json:"throughput"`
	Success     float64         `json:"success"`
	StatusCodes map[string]int  `json:"status_codes"`
	Pass        bool            `json:"pass"`
}

func Attack(ctx *gin.Context) {

	var req attackRequest
	var response attackResponse

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "Error")
	}

	response = processAttack(req)

	ctx.JSON(http.StatusOK, response)
}

func processAttack(request attackRequest) attackResponse {
	var targeter vegeta.Targeter

	if request.Method == "GET" {
		targeter = vegeta.NewStaticTargeter(vegeta.Target{
			Method: request.Method,
			URL:    request.Url,
		})
	} else {
		targeter = func(tgt *vegeta.Target) error {
			if tgt == nil {
				return vegeta.ErrNilTarget
			}

			tgt.Method = request.Method

			tgt.URL = request.Url

			payloadString := fmt.Sprintf("%v", request.PayLoad)

			tgt.Body = []byte(payloadString)
			return nil
		}
	}

	attackRate, err := strconv.Atoi(request.AttackRate)

	if err != nil {
		fmt.Println("Error converting rate :" + request.AttackRate)
	}
	attackDuration, err := strconv.Atoi(request.AttackRate)

	if err != nil {
		fmt.Println("Error converting rate :" + request.AttackRate)
	}

	rate := vegeta.Rate{Freq: attackRate, Per: time.Second} // change the rate here
	duration := time.Duration(attackDuration) * time.Second
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Load Test") {
		metrics.Add(res)
	}
	metrics.Close()
	success := metrics.Success * 100
	var pass bool

	if int64(success) >= request.PassRate {
		pass = true
	} else {
		pass = false
	}

	return attackResponse{
		Latencies: latencyResponse{
			Total: float64(metrics.Latencies.Total) / 1000000000,
			Mean:  float64(metrics.Latencies.Mean) / 1000000000,
		},
		Duration:    math.Ceil(float64(metrics.Duration) / 1000000000),
		Wait:        float64(metrics.Wait) / 1000000000,
		Requests:    float64(metrics.Requests),
		Throughput:  metrics.Throughput,
		Success:     success,
		StatusCodes: metrics.StatusCodes,
		Pass:        pass,
	}

}
