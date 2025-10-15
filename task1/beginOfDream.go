package main

import (
	"errors"
	"fmt"
	"strconv"
)

func main() {
	// 查找只出现一次的元素
	fmt.Println("---------------------------查找只出现一次的元素-------------------------")
	var test1_1 = []int{1, 1, 2, 3, 3}
	var test1_2 = []int{1, 1, 2, 2, 3, 3}
	testFindOnceNumber(test1_1)
	testFindOnceNumber(test1_2)

	// 判断是否回文数
	fmt.Println("---------------------------判断是否回文数---------------------------")
	test2 := []int{121, -121, 10, 5, 12321}
	for _, num := range test2 {
		if isPalindrome(num) {
			fmt.Printf("%d 是回文数\n", num)
		} else {
			fmt.Printf("%d 不是回文数\n", num)
		}
	}

	// 判断字符串是否有效
	fmt.Println("---------------------------判断字符串是否有效---------------------------")
	test3 := []string{"()", "()[]{}", "(]", "([)]", "{[]}", "{([])}"}
	for _, str := range test3 {
		if validStr(str) {
			fmt.Printf("%s 是有效字符串\n", str)
		} else {
			fmt.Printf("%s 不是有效字符串\n", str)
		}
	}

	// 最长公共前缀
	fmt.Println("---------------------------最长公共前缀---------------------------")
	test4_1 := []string{"flower", "flow", "flight"}
	test4_2 := []string{"dog", "racecar", "car"}
	fmt.Printf("最长公共前缀：%s\n", longestCommonPrefix(test4_1))
	fmt.Printf("最长公共前缀：%s\n", longestCommonPrefix(test4_2))

	// 给定一个表示 大整数 的整数数组 digits，其中 digits[i] 是整数的第 i 位数字。将大整数加 1，并返回结果的数字数组。
	fmt.Println("---------------------------大整数加一---------------------------")
	test5_1 := []int{1, 2, 3}
	test5_2 := []int{4, 3, 2, 1}
	test5_3 := []int{9}
	test5_4 := []int{9, 9, 9}
	fmt.Printf("大整数加一：%v\n", plusOne(test5_1))
	fmt.Printf("大整数加一：%v\n", plusOne(test5_2))
	fmt.Printf("大整数加一：%v\n", plusOne(test5_3))
	fmt.Printf("大整数加一：%v\n", plusOne(test5_4))

	// 删除有序数组中的重复项
	fmt.Println("---------------------------删除有序数组中的重复项---------------------------")
	test6_1 := []int{1, 1, 2}
	test6_2 := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	fmt.Printf("删除有序数组中的重复项：%v\n", removeDuplicates(test6_1))
	fmt.Printf("删除有序数组中的重复项：%v\n", removeDuplicates(test6_2))

	// 合并区间
	fmt.Println("---------------------------合并区间---------------------------")
	test7_1 := [][]int{{1, 3}, {15, 18}, {8, 10}, {2, 6}}
	test7_2 := [][]int{{4, 7}, {1, 5}, {3, 6}}
	fmt.Printf("合并区间：%v\n", merge(test7_1))
	fmt.Printf("合并区间：%v\n", merge(test7_2))

	// 两数之和
	fmt.Println("---------------------------两数之和---------------------------")
	test8_1 := []int{2, 7, 11, 15}
	target1 := 9
	test8_2 := []int{3, 2, 4}
	target2 := 6
	fmt.Printf("两数之和对应数组下标：%v\n", twoSum(test8_1, target1))
	fmt.Printf("两数之和对应数组下标：%v\n", twoSum(test8_2, target2))
}

func testFindOnceNumber(arr []int) {
	if v, err := findOnceNumber(arr); err != nil {
		fmt.Println("异常错误：", err)
	} else {
		fmt.Println("findOnceNumber返回只出现一次的元素：", v)
	}
}

// 给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素
func findOnceNumber(arr []int) (int, error) {
	tempMap := make(map[int]int)
	for _, val := range arr {
		tempMap[val] += 1
	}
	for k, v := range tempMap {
		if v == 1 {
			return k, nil
		}
	}
	return 0, errors.New("不存在只出现一次的元素")
}

// 判断回文数（即正读和反读都相同的数，如121、1331 ）
func isPalindrome(a int) bool {
	if a < 0 || (a != 0 && a%10 == 0) {
		// 负数或末尾为0的特殊情况，直接返回不是
		return false
	} else if a < 10 {
		// 个位数直接返回是
		return true
	} else {
		str := strconv.Itoa(a)
		reverseStr := reverseASCII(str)
		if str == reverseStr {
			return true
		} else {
			return false
		}
	}
}

// 字符串反转
func reverseASCII(s string) string {
	b := []byte(s)
	len := len(b)
	n := len / 2
	for i := 0; i < n; i++ {
		j := len - 1 - i
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func validStr(s string) bool {
	if s == "" || len(s) == 0 || len(s)%2 != 0 {
		return false
	}

	// 用切片模拟栈（存储左括号）
	stack := []rune{}
	match := map[rune]rune{')': '(', ']': '[', '}': '{'}
	for _, v := range s {
		switch v {
		case '(', '[', '{':
			stack = append(stack, v)
		case ')', ']', '}':
			len := len(stack)
			if len == 0 || stack[len-1] != match[v] {
				// 栈为空，或不匹配
				return false
			}
			stack = stack[:len-1]
		default:
			// 非括号
			return false
		}
	}
	return len(stack) == 0
}

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	for i := 0; i < len(strs[0]); i++ {
		// 以第一个字符串为基准
		char := strs[0][i]
		for j := 1; j < len(strs); j++ {
			if i >= len(strs[j]) || strs[j][i] != char {
				return strs[0][:i]
			}
		}
	}
	// 第一个最短，且全匹配上了
	return strs[0]
}

func plusOne(digits []int) []int {
	for i := len(digits) - 1; i >= 0; i-- {
		// 如果当前位小于9，直接加1并返回结果
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}
	// 最高位也进位了
	digits = append([]int{1}, digits...) // 在切片开头添加1
	return digits
}

func removeDuplicates(nums []int) []int {
	if len(nums) == 0 {
		return nums
	}
	slow := 0
	for fast := 1; fast < len(nums); fast++ {
		if nums[fast] != nums[slow] {
			slow++
			nums[slow] = nums[fast]
		}
	}
	return nums[:slow+1]
}

func merge(intervals [][]int) [][]int {
	len := len(intervals)
	if len == 0 {
		return intervals
	}
	// 先对二维数组按第一个元素排序
	for i := 0; i < len-1; i++ {
		for j := 0; j < len-1-i; j++ {
			if intervals[j][0] > intervals[j+1][0] {
				intervals[j], intervals[j+1] = intervals[j+1], intervals[j]
			}
		}
	}
	fmt.Println("排序后的二维数组：", intervals)
	// 合并区间
	result := [][]int{}
	temp := intervals[0]
	for i := 1; i < len; i++ {
		if temp[1] >= intervals[i][0] {
			// 有交集，合并
			if temp[1] < intervals[i][1] {
				temp[1] = intervals[i][1]
			}
		} else {
			// 无交集，添加到结果集
			result = append(result, temp)
			temp = intervals[i]
		}
	}
	// 添加最后一个区间
	result = append(result, temp)
	return result
}

func twoSum(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	fmt.Println("没有满足条件的索引")
	return []int{}
}
