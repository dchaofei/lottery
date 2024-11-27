package main

import (
	"fmt"
	"testing"
)

func TestDrawWinner_Distribution(t *testing.T) {
	// 测试数据
	names := []string{
		"张三", "李四", "王五", "赵六", "钱七",
		"孙八", "周九", "吴十",
	}

	// 抽奖次数
	iterations := 100000

	// 统计每个名字被抽中的次数
	results := make(map[string]int)
	for _, name := range names {
		results[name] = 0
	}

	// 进行多次抽奖
	for i := 0; i < iterations; i++ {
		winner, err := drawWinner(names)
		if err != nil {
			t.Fatalf("抽奖失败: %v", err)
		}
		results[winner]++
	}

	// 计算期望值 (每人12.5%)
	expected := iterations / len(names)

	// 打印统计信息
	fmt.Printf("\n期望每人抽中次数: %d (12.5%%)\n", expected)
	fmt.Println("实际抽奖结果:")
	for name, count := range results {
		percentage := float64(count) / float64(iterations) * 100
		fmt.Printf("%s: %d 次 (%.2f%%)\n", name, count, percentage)
	}
}

// 辅助函数：计算绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// 测试空输入情况
func TestDrawWinner_EmptyInput(t *testing.T) {
	_, err := drawWinner([]string{})
	if err == nil {
		t.Error("期望空输入时返回错误，但没有")
	}
}

// 测试单个名字的情况
func TestDrawWinner_SingleName(t *testing.T) {
	names := []string{"张三"}
	winner, err := drawWinner(names)
	if err != nil {
		t.Errorf("单个名字抽奖失败: %v", err)
	}
	if winner != "张三" {
		t.Errorf("期望抽中 '张三'，实际抽中 '%s'", winner)
	}
}
