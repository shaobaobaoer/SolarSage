package main

import (
	"fmt"
	"path/filepath"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/sweph"
)

func main() {
	// 使用绝对路径
	absPath, _ := filepath.Abs("third_party/swisseph/ephe")
	fmt.Printf("Absolute path: %s\n", absPath)
	
	sweph.Init(absPath)
	defer sweph.Close()

	// 测试Chiron
	jd := sweph.JulDay(2026, 3, 1, 0, true)
	res, err := sweph.CalcUT(jd, sweph.SE_CHIRON)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Printf("Chiron: %.4f°\n", res.Longitude)
	}
}