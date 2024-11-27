package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func main() {
	// 静态文件服务
	http.Handle("/", http.FileServer(http.Dir("static")))

	// 抽奖API
	http.HandleFunc("/draw", drawHandler)

	// 启动服务器并打开浏览器
	serverURL := "http://localhost:8080"
	log.Printf("服务启动在 %s", serverURL)
	log.Println("正在尝试打开浏览器...")
	log.Println("如果浏览器没有自动打开，请手动访问上述链接")

	// 延迟一小段时间后打开浏览器，确保服务器已经启动
	go func() {
		time.Sleep(100 * time.Millisecond)
		if err := openBrowser(serverURL); err != nil {
			log.Printf("无法自动打开浏览器: %v", err)
			log.Println("请手动复制链接到浏览器访问")
		} else {
			log.Println("浏览器已打开")
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func drawHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持 POST 请求", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Names string `json:"names"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	names := strings.Split(strings.TrimSpace(request.Names), "\n")
	names = filterEmptyStrings(names)

	winner, err := drawWinner(names)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := struct {
		Winner string `json:"winner"`
	}{
		Winner: winner,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func filterEmptyStrings(strs []string) []string {
	var result []string

	for _, str := range strs {
		trimmed := strings.TrimSpace(str)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}

	return result
}

// openBrowser 根据操作系统打开默认浏览器
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func drawWinner(names []string) (string, error) {
	if len(names) == 0 {
		return "", fmt.Errorf("名单为空")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(names))))
	if err != nil {
		return "", fmt.Errorf("随机数生成失败: %v", err)
	}

	return names[n.Int64()], nil
}
