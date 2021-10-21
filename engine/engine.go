package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net/http"
	"os/exec"
)

type attackRequest struct {
	Method         string `json:"method" binding:"required"`
	Url            string `json:"url" binding:"required"`
	AttackDuration string `json:"attack_duration" binding:"required"`
	AttackRate     string `json:"attack_rate" binding:"required"`
}

type latencyResponse struct {
	Total float64 `json:"total"`
	Mean float64 `json:"mean"`
}

type attackResponse struct {
	Latencies latencyResponse `json:"latencies"`
	Duration float64 `json:"duration"`
	Wait float64 `json:"wait"`
	Requests float64 `json:"requests"`
	Throughput float64 `json:"throughput"`
	Success float64 `json:"success"`
	StatusCodes map[string]interface{} `json:"status_codes"`
}

func Attack(ctx *gin.Context) {

	var req attackRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "Error")
	}

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

	fmt.Println(outb.String())

	json.Unmarshal(outb.Bytes(), &attackResp)

	latencyResp := &attackResp.Latencies

	finResp := &attackResponse{
		Latencies:   latencyResponse{
			Total: latencyResp.Total / 1000000000,
			Mean : latencyResp.Mean / 1000000000,
		},
		Duration:    math.Ceil(attackResp.Duration / 1000000000),
		Wait:        attackResp.Wait / 1000000000,
		Requests:    attackResp.Requests,
		Throughput:  attackResp.Throughput,
		Success:     attackResp.Success * 100,
		StatusCodes: attackResp.StatusCodes,
	}

	ctx.JSON(http.StatusOK, finResp)
}
