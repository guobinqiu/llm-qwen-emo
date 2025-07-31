package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	DashApiKey   = "sk-3c26bc48e75044dd810a0838f18d75f9"
	DetectAPIURL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/image2video/face-detect"
	VideoAPIURL  = "https://dashscope.aliyuncs.com/api/v1/services/aigc/image2video/video-synthesis/"
	TaskAPIURL   = "https://dashscope.aliyuncs.com/api/v1/tasks/%s"
)

func CheckImage(imageURL string) ([]int, []int, error) {
	data := map[string]interface{}{
		"model":      "emo-detect-v1",
		"input":      map[string]interface{}{"image_url": imageURL},
		"parameters": map[string]interface{}{"ratio": "1:1"},
	}
	buf, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", DetectAPIURL, bytes.NewBuffer(buf))
	req.Header.Set("Authorization", "Bearer "+DashApiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("face-detect请求失败: %v", err)
	}
	defer resp.Body.Close()
	var respData struct {
		Output struct {
			CheckPass bool   `json:"check_pass"`
			FaceBBox  []int  `json:"face_bbox"`
			ExtBBox   []int  `json:"ext_bbox"`
			Code      string `json:"code"`
			Message   string `json:"message"`
		} `json:"output"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, nil, err
	}
	if !respData.Output.CheckPass {
		return nil, nil, fmt.Errorf("图片检测失败: code=%s message=%s", respData.Output.Code, respData.Output.Message)
	}
	return respData.Output.FaceBBox, respData.Output.ExtBBox, nil
}

func GenerateVideo(imageURL, audioURL string, faceBBox, extBBox []int) (string, error) {
	data := map[string]interface{}{
		"model": "emo-v1",
		"input": map[string]interface{}{
			"image_url": imageURL,
			"audio_url": audioURL,
			"face_bbox": faceBBox,
			"ext_bbox":  extBBox,
		},
		"parameters": map[string]any{"style_level": "normal"},
	}
	buf, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", VideoAPIURL, bytes.NewBuffer(buf))
	req.Header.Set("Authorization", "Bearer "+DashApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DashScope-Async", "enable")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("视频合成请求失败: %v", err)
	}
	defer resp.Body.Close()
	var respData struct {
		Output struct {
			TaskID string `json:"task_id"`
		} `json:"output"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return "", err
	}
	return respData.Output.TaskID, nil
}

func QueryTaskStatus(taskID string) (status, videoURL, code, message string, err error) {
	url := fmt.Sprintf(TaskAPIURL, taskID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+DashApiKey)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", "", "", fmt.Errorf("状态查询失败，HTTP状态码: %d", resp.StatusCode)
	}
	var respData struct {
		Output struct {
			TaskID     string `json:"task_id"`
			TaskStatus string `json:"task_status"`
			Code       string `json:"code,omitempty"`
			Message    string `json:"message,omitempty"`
			Results    struct {
				VideoURL string `json:"video_url"`
			} `json:"results"`
		} `json:"output"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return "", "", "", "", err
	}
	return respData.Output.TaskStatus, respData.Output.Results.VideoURL, respData.Output.Code, respData.Output.Message, nil
}
