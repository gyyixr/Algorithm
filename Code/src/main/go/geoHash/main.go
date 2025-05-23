package main

import (
	"fmt"
)

const (
	base32 = "0123456789bcdefghjkmnpqrstuvwxyz" // Base32 编码字符集 (移除了 a, i, l, o)
)

// Encode 将经纬度编码为指定精度的 GeoHash 字符串 (不使用 bytes.Buffer)
func Encode(latitude, longitude float64, precision int) (string, error) {
	if precision <= 0 {
		return "", fmt.Errorf("precision must be greater than 0")
	}

	geohashChars := make([]byte, 0, precision) // 初始化一个 byte 切片，预设容量
	var bits uint                              // 当前处理的二进制位
	// var bit uint  // 当前二进制位的值 (0 或 1) // 这个变量实际上在之前的逻辑中没有直接使用来赋值，而是通过 ch |= 来设置
	var ch byte // Base32 编码字符的索引

	// 经纬度范围
	latInterval := []float64{-90.0, 90.0}
	lonInterval := []float64{-180.0, 180.0}

	isEven := true // true 表示处理经度, false 表示处理纬度

	for len(geohashChars) < precision {
		if isEven { // 处理经度
			mid := (lonInterval[0] + lonInterval[1]) / 2
			if longitude > mid {
				ch |= (1 << (4 - bits)) // 设置当前二进制位为 1
				lonInterval[0] = mid
			} else {
				lonInterval[1] = mid
			}
		} else { // 处理纬度
			mid := (latInterval[0] + latInterval[1]) / 2
			if latitude > mid {
				ch |= (1 << (4 - bits)) // 设置当前二进制位为 1
				latInterval[0] = mid
			} else {
				latInterval[1] = mid
			}
		}

		isEven = !isEven // 切换处理经纬度
		bits++

		if bits == 5 { // 每 5 个二进制位生成一个 Base32 字符
			geohashChars = append(geohashChars, base32[ch])
			bits = 0
			ch = 0
		}
	}

	return string(geohashChars), nil
}

func main() {
	// 示例：编码旧金山的经纬度
	latitude := 37.7749
	longitude := -122.4194
	precision := 9 // GeoHash 字符串长度

	geohash, err := Encode(latitude, longitude, precision)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Latitude: %f, Longitude: %f\n", latitude, longitude)
	fmt.Printf("GeoHash (precision %d): %s\n", precision, geohash)

	// 示例：编码北京的经纬度
	latitudeBeijing := 39.9042
	longitudeBeijing := 116.4074
	precisionBeijing := 7

	geohashBeijing, err := Encode(latitudeBeijing, longitudeBeijing, precisionBeijing)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("\nLatitude: %f, Longitude: %f\n", latitudeBeijing, longitudeBeijing)
	fmt.Printf("GeoHash (precision %d): %s\n", precisionBeijing, geohashBeijing)
}
