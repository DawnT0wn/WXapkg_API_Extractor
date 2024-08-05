package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type API struct {
	File string `json:"file"`
	Path string `json:"path"`
}

func findAPIs(directory string) ([]API, error) {
	// 定义正则表达式以匹配路径
	apiPattern1 := regexp.MustCompile(`(?:(?:[\w/-]+)?[\w/-]+(?:[\w/-]+)?(?:[\w/-]+)?)/([\w/-]+)[\w/-]*`)
	apiPattern2 := regexp.MustCompile(`/[\w-]+(?:/[\w-]+)*`)
	var apis []API
	apiMap := make(map[string]string) // 使用map来存储更具体的路径

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".js" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			contentStr := string(content)

			// 使用第一个正则表达式进行匹配
			matches1 := apiPattern1.FindAllString(contentStr, -1)
			for _, match := range matches1 {
				if !hasFileExtension(match) && !isFalsePositive(match) {
					updateMap(apiMap, match)
				}
			}

			// 使用第二个正则表达式进行匹配
			matches2 := apiPattern2.FindAllString(contentStr, -1)
			for _, match := range matches2 {
				if !hasFileExtension(match) && !isFalsePositive(match) {
					updateMap(apiMap, match)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 将去重后的路径转换为 API 列表
	for path := range apiMap {
		apis = append(apis, API{File: "", Path: path})
	}

	// 重新读取文件以填充 File 字段
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".js" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			contentStr := string(content)

			for i := range apis {
				if strings.Contains(contentStr, apis[i].Path) {
					apis[i].File = path
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return apis, nil
}

// updateMap 更新 map，如果存在相同路径但较长，则替换
func updateMap(apiMap map[string]string, path string) {
	for existingPath := range apiMap {
		if strings.HasSuffix(path, existingPath) && len(path) > len(existingPath) {
			apiMap[path] = path
			return
		}
		if strings.HasSuffix(existingPath, path) && len(existingPath) > len(path) {
			return
		}
	}
	apiMap[path] = path
}

func isFalsePositive(path string) bool {
	// 排除以非字母字符开头的路径
	if matched, _ := regexp.MatchString(`^[^a-zA-Z/]`, path); matched {
		return true
	}
	// 排除包含特定关键词的路径
	if matched, _ := regexp.MatchString(`^(?:text/|image/|font|form-data)`, path); matched {
		return true
	}
	// 排除包含特定关键词的路径，包括包含 "js" 的路径
	if matched, _ := regexp.MatchString(`/.*js.*|.*weapp.*`, path); matched {
		return true
	}
	// 排除路径中包含任何数字的情况
	if matched, _ := regexp.MatchString(`\d`, path); matched {
		return true
	}
	// 排除两个斜杠之间只有一个字符的路径
	if matched, _ := regexp.MatchString(`/[^/]/`, path); matched {
		return true
	}
	// 排除斜杠后只有一个字符的路径
	if matched, _ := regexp.MatchString(`/[^/]/[^/]`, path); matched {
		return true
	}
	// 排除以 /wxb 开头的路径
	if matched, _ := regexp.MatchString(`^/wxb`, path); matched {
		return true
	}
	// 排除以单个字母开头的路径，例如 /t/
	if matched, _ := regexp.MatchString(`^/[a-zA-Z]/`, path); matched {
		return true
	}
	// 排除包含双斜杠的路径
	if matched, _ := regexp.MatchString(`//`, path); matched {
		return true
	}
	// 排除长度不符合要求的路径
	if len(path) < 5 {
		return true
	}
	return false
}

func hasFileExtension(path string) bool {
	// 检查路径是否包含文件扩展名
	return strings.Contains(path, ".")
}

func saveResults(results []API, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(results)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// 定义命令行标志
	directory := flag.String("d", "", "Directory to scan")
	directoryLong := flag.String("directory", "", "Directory to scan")
	outputFile := flag.String("o", "", "Output file for results")
	outputFileLong := flag.String("output", "", "Output file for results")
	flag.Parse()

	// 处理长短标志的兼容性
	if *directory == "" {
		*directory = *directoryLong
	}
	if *outputFile == "" {
		*outputFile = *outputFileLong
	}

	// 检查命令行参数
	if *directory == "" || *outputFile == "" {
		fmt.Println("Usage: go run main.go -d <directory> -o <output>")
		return
	}

	apis, err := findAPIs(*directory)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = saveResults(apis, *outputFile)
	if err != nil {
		fmt.Println("Error saving results:", err)
		return
	}

	fmt.Printf("Results saved to %s\n", *outputFile)
	// 打印 API 路径的总数
	fmt.Printf("一个匹配到: %d个API\n\n", len(apis))

	// 打印每个 API 的路径
	fmt.Println("Matched API Paths: ")
	for _, api := range apis {
		fmt.Println(api.Path)
	}
}
