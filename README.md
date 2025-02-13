# WXapkg_API_Extractor

## 介绍

`WXapkg_API_Extractor` 是一个 Go 程序，用于扫描指定目录下的 JavaScript 文件，提取 API 路径，并将结果保存到一个 JSON 文件中。该工具使用正则表达式匹配 API 路径，并通过一系列过滤条件来排除无关路径。

## 功能

- 扫描指定目录中的 JavaScript 文件
- 提取 API 路径
- 排除特定的无效路径
- 将结果保存为 JSON 文件

## Usage

```
go run main.go -directory <path> -output <file>
```
排除情况
```azure
排除以非字母字符开头的路径
排除包含特定关键词的路径
排除包含特定关键词的路径，包括包含 "js" 的路径
排除路径中包含任何数字的情况
排除两个斜杠之间只有一个字符的路径
排除斜杠后只有一个字符的路径
排除以 /wxb 开头的路径
排除以单个字母开头的路径，例如 /t/
排除包含双斜杠的路径
排除长度不符合要求的路径
```
该工具只是辅助查找小程序反编译的API，规则不满足所有情况，还是有一部分误报，但是很明显，提取到的接口可能是分开写的，所以接口不一定是完整的接口，建议结合代码寻找完整的API