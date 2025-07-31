package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 使用 ffmpeg 切割音频，返回排序后的切片文件列表
func SliceAudio(audioPath, outDir string, segmentSecond int) ([]string, error) {
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		return nil, err
	}

	if segmentSecond <= 0 {
		return nil, fmt.Errorf("segmentSecond must be greater than 0")
	}

	ext := filepath.Ext(audioPath)

	pattern := filepath.Join(outDir, "part_%d"+ext)
	cmd := exec.Command("ffmpeg", "-i", audioPath, "-f", "segment", "-segment_time", fmt.Sprintf("%d", segmentSecond), "-c", "copy", pattern)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg执行失败: %v, output: %s", err, string(out))
	}

	files, err := os.ReadDir(outDir)
	if err != nil {
		return nil, err
	}

	var audioFiles []string
	prefix := "part_"
	for _, f := range files {
		if strings.HasPrefix(f.Name(), prefix) && strings.HasSuffix(f.Name(), ext) {
			audioFiles = append(audioFiles, filepath.Join(outDir, f.Name()))
		}
	}

	// sort.Strings(audioFiles)
	SortFilesByNumber(audioFiles)
	return audioFiles, nil
}

// 合并视频片段
func MergeVideos(videoPaths []string, outFile string) error {
	tmpDir := filepath.Dir(outFile)
	os.MkdirAll(tmpDir, 0755)

	listFile := filepath.Join(tmpDir, "files.txt")
	f, err := os.Create(listFile)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, lf := range videoPaths {
		absPath, _ := filepath.Abs(lf)
		_, err := fmt.Fprintf(f, "file '%s'\n", absPath)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", listFile, "-c", "copy", outFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg合并失败: %v, 输出: %s", err, string(out))
	}

	return nil
}
