package main

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Cashier struct {
	BankNoteOrCoin float64
	Amount         int
}

type PaidDetail struct {
	Paid         int `json:"paid"`
	ProductPrice int `json:"productPrice"`
}

type ChangeDetail struct {
	BankNoteOrCoin float64
	Count          int
}

type IndexXyz struct {
	X int `json:"indexX"`
	Y int `json:"indexY"`
	Z int `json:"indexZ"`
}

func cal(i int) float64 {
	var retult float64 = 1

	for j := 1; j <= i; j++ {
		a := 2 + (float64(j)-1)*0.5
		b := float64(j)
		diff := a * b
		retult += diff
	}

	return retult
}

func main() {
	cachiers := []Cashier{
		{BankNoteOrCoin: 1000, Amount: 10},
		{BankNoteOrCoin: 500, Amount: 20},
		{BankNoteOrCoin: 100, Amount: 15},
		{BankNoteOrCoin: 50, Amount: 20},
		{BankNoteOrCoin: 20, Amount: 30},
		{BankNoteOrCoin: 10, Amount: 20},
		{BankNoteOrCoin: 5, Amount: 20},
		{BankNoteOrCoin: 1, Amount: 20},
		{BankNoteOrCoin: 0.5, Amount: 50},
	}

	r := gin.Default()

	r.POST("/find-xyz", func(c *gin.Context) {
		var indexXyz IndexXyz
		if err := c.ShouldBindJSON(&indexXyz); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"result-x": cal(indexXyz.X),
			"result-y": cal(indexXyz.Y),
			"result-z": cal(indexXyz.Z),
		})
	})

	r.POST("/change", func(c *gin.Context) {
		var paidDetail PaidDetail
		if err := c.ShouldBindJSON(&paidDetail); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if paidDetail.ProductPrice > paidDetail.Paid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		result := []ChangeDetail{}
		var change float64 = float64(paidDetail.Paid - paidDetail.ProductPrice)
		var tempChange = change

		for index, cashier := range cachiers {
			if tempChange >= cashier.BankNoteOrCoin && cashier.Amount != 0 {
				bankOrCoinCount := int(math.Floor(float64(tempChange) / cashier.BankNoteOrCoin))
				if cashier.Amount-bankOrCoinCount < 0 {
					result = append(result, ChangeDetail{
						BankNoteOrCoin: cashier.BankNoteOrCoin,
						Count:          cashier.Amount,
					})
					cachiers[index].Amount = 0
					tempChange = tempChange - (cashier.BankNoteOrCoin * float64(cashier.Amount))
				} else {
					result = append(result, ChangeDetail{
						BankNoteOrCoin: cashier.BankNoteOrCoin,
						Count:          bankOrCoinCount,
					})
					cachiers[index].Amount -= bankOrCoinCount
					tempChange = math.Mod(float64(tempChange), cashier.BankNoteOrCoin)
				}
			}
		}

		if tempChange > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No bank not or coin cover of change."})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"change": change,
			"detail": result,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
