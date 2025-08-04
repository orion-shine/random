package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"math/big"
	"testing"
)

var bitFull = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)

// RandomBigInt 生成伪随机 big.Int 值
func RandomBigInt(secret, seed string) *big.Int {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(seed))
	hash := h.Sum(nil)
	return new(big.Int).SetBytes(hash)
}

// Random 随机整数 [0, length]
func Random(length int64, seed, secret string) int64 {
	rnd := RandomBigInt(secret, seed)
	frac := new(big.Rat).SetFrac(rnd, bitFull)
	x, _ := frac.Float64()
	return int64(x * float64(length+1))
}

// RandomFloat 随机浮点数 [0.0, 1.0)
func RandomFloat(secret, seed string) float64 {
	rnd := RandomBigInt(secret, seed)
	frac := new(big.Rat).SetFrac(rnd, bitFull)
	f, _ := frac.Float64()
	return f
}

// RandomSort 数组排列（洗牌）
func RandomSort(arr []int64, secret, seed string) []int64 {
	ids := make([]int64, len(arr))
	for i := range arr {
		ids[i] = int64(i)
	}

	for i := len(arr) - 1; i > 0; i-- {
		curSeed := fmt.Sprintf("%s_%d", seed, i)
		j := Random(int64(i), curSeed, secret)
		ids[i], ids[j] = ids[j], ids[i]
	}

	res := make([]int64, len(arr))
	for i, idx := range ids {
		res[i] = arr[idx]
	}
	return res
}

// RandomArrRet 数组随机取一个元素
func RandomArrRet(arr []int64, secret, seed string) int64 {
	if len(arr) == 0 {
		return -1
	}
	idx := Random(int64(len(arr)-1), seed+"_pick", secret)
	return arr[idx]
}

// RandomPC28 彩票：PC28 模拟（0-9取3个求和）
func RandomPC28(secret, seed string) []int64 {
	result := make([]int64, 3)
	for i := 0; i < 3; i++ {
		subSeed := fmt.Sprintf("%s_pc28_%d", seed, i)
		result[i] = Random(9, subSeed, secret)
	}
	return result
}

// RandomMarkSix 彩票：Mark Six（香港六合彩）从1~49中随机选6个不重复
func RandomMarkSix(secret, seed string) []int64 {
	nums := make([]int64, 49)
	for i := 0; i < 49; i++ {
		nums[i] = int64(i + 1)
	}
	shuffled := RandomSort(nums, secret, seed+"_mark6")
	return shuffled[:6]
}

func TestRandom(t *testing.T) {
	fmt.Println("Random int [0-99]:", Random(99, "123", "123"))
	fmt.Println("Random float:", RandomFloat("132", "132"))
	fmt.Println("Random sorted [0-4]:", RandomSort([]int64{0, 1, 2, 3, 4}, "123", "132"))
	fmt.Println("Random pick from array:", RandomArrRet([]int64{10, 20, 30, 40, 50}, "key", "seed"))
	fmt.Println("PC28 result:", RandomPC28("secret", "pc28_seed"))
	fmt.Println("Mark Six result:", RandomMarkSix("secret", "marksix_seed"))
}
