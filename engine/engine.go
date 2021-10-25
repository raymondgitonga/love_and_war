package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	vegeta "github.com/tsenart/vegeta/lib"
	"log"
	"math"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

type attackRequest struct {
	Method         string `json:"method" binding:"required"`
	Url            string `json:"url" binding:"required"`
	AttackDuration string `json:"attack_duration" binding:"required"`
	AttackRate     string `json:"attack_rate" binding:"required"`
	PayLoad        string `json:"pay_load"`
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
}

func Attack(ctx *gin.Context) {

	var req attackRequest
	var response attackResponse

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "Error")
	}

	if req.Method == "GET" {
		response = processGet(req)
	} else {
		response = processPost(req)
	}

	ctx.JSON(http.StatusOK, response)
}

func processGet(req attackRequest) attackResponse {
	cmdStr := fmt.Sprintf("echo %s %s | vegeta attack -duration=%ss -rate=%s | vegeta report --type=json",
		req.Method, req.Url, req.AttackDuration, req.AttackRate)

	cmd := exec.Command("zsh", "-c", cmdStr)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()

	if err != nil {
		fmt.Println("The error is " + err.Error())
		log.Fatal(err)
	}
	var attackResp attackResponse

	json.Unmarshal(outb.Bytes(), &attackResp)

	latencyResp := &attackResp.Latencies

	return attackResponse{
		Latencies: latencyResponse{
			Total: latencyResp.Total / 1000000000,
			Mean:  latencyResp.Mean / 1000000000,
		},
		Duration:    math.Ceil(attackResp.Duration / 1000000000),
		Wait:        attackResp.Wait / 1000000000,
		Requests:    attackResp.Requests,
		Throughput:  attackResp.Throughput,
		Success:     attackResp.Success * 100,
		StatusCodes: attackResp.StatusCodes,
	}
}

func processPost(request attackRequest) attackResponse {
	targeter := func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = request.Method

		tgt.URL = request.Url

		payloadString := fmt.Sprintf("%v", request.PayLoad)

		tgt.Body = []byte(payloadString)
		return nil
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

	return attackResponse{
		Latencies: latencyResponse{
			Total: float64(metrics.Latencies.Total) / 1000000000,
			Mean:  float64(metrics.Latencies.Mean) / 1000000000,
		},
		Duration:    math.Ceil(float64(metrics.Duration) / 1000000000),
		Wait:        float64(metrics.Wait) / 1000000000,
		Requests:    float64(metrics.Requests),
		Throughput:  metrics.Throughput,
		Success:     metrics.Success * 100,
		StatusCodes: metrics.StatusCodes,
	}

}
