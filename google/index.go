package google

import (
	"github.com/dgryski/dgoogauth"
	"math"
	"strconv"
	"time"
)

func Index(secretKey string) string {
	t := int64(math.Floor(float64(time.Now().Unix() / 30)))
	return strconv.Itoa(dgoogauth.ComputeCode(secretKey,t))
}