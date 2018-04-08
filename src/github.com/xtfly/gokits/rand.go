package gokits

import (
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/rand"
	"time"

	crand "crypto/rand"
)

func NewRandWithPrefix(prefix string, len int) string {
	return hex.EncodeToString([]byte(prefix)) + NewRand(len)
}

// Create a rang string
func NewRand(len int) string {
	u := make([]byte, len/2)
	// Reader is a global, shared instance of a cryptographically strong pseudo-random generator.
	// On Unix-like systems, Reader reads from /dev/urandom.
	// On Windows systems, Reader uses the CryptGenRandom API.
	_, err := io.ReadFull(crand.Reader, u)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(u)
}

// NewUUID generates a new UUID based on version 4.
func NewUUID() string {
	u := make([]byte, 16)
	// Reader is a global, shared instance of a cryptographically strong pseudo-random generator.
	// On Unix-like systems, Reader reads from /dev/urandom.
	// On Windows systems, Reader uses the CryptGenRandom API.
	_, err := io.ReadFull(crand.Reader, u)
	if err != nil {
		panic(err)
	}

	// Set version (4) and variant (2).
	var version byte = 4 << 4
	var variant byte = 2 << 4
	u[6] = version | (u[6] & 15)
	u[8] = variant | (u[8] & 15)

	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

//生成规定范围内的整数
//设置起始数字范围，0开始,n截止
func RangeRand(n int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(n)
}

//标准正态分布随机整数，n为随机个数,从0开始
func NormRand(n int64) float64 {
	//sample = NormFloat64() * desiredStdDev + desiredMean
	// 默认位置参数(期望desiredMean)为0,尺度参数(标准差desiredStdDev)为1.

	var i, sample int64 = 0, 0
	desiredMean := 0.0
	desiredStdDev := 100.0

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i < n {
		rn := int64(r.NormFloat64()*desiredStdDev + desiredMean)
		sample = rn % n
		i += 1
	}

	return math.Abs(float64(sample))
}
