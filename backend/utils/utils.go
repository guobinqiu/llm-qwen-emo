package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func DownloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("下载文件失败，状态码: %d", resp.StatusCode)
	}

	outFile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}

// 按数字排序的排序器
func SortFilesByNumber(files []string) {
	// 用正则提取数字，假设文件名中包含数字部分，比如 part_123.mp3
	re := regexp.MustCompile(`\d+`)

	sort.Slice(files, func(i, j int) bool {
		// 提取第i个文件的数字
		numI := extractNumber(files[i], re)
		numJ := extractNumber(files[j], re)
		return numI < numJ
	})
}

// 从字符串中提取第一个数字，没找到返回0
func extractNumber(s string, re *regexp.Regexp) int {
	numStr := re.FindString(s)
	if numStr == "" {
		return 0
	}
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0
	}
	return num
}
