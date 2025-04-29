# 紫微斗数排盘服务

## 项目简介
基于Go语言开发的紫微斗数排盘微服务，提供阳历/农历两种排盘方式。

## 功能特性
- 阳历排盘接口 `/ziwei_yangli`
- 农历排盘接口 `/ziwei_nongli` 
- 支持时辰自动换算
- 返回标准JSON格式数据

## 快速开始
### 环境要求
- Go 1.23+ 

### 编译运行
```bash
go mod tidy
go run main.go
```

## API文档
### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|-----|-----|
| year   | 数字 | 是  | 年份（如2024） |
| month  | 数字 | 是  | 月份（1-12） |
| day    | 数字 | 是  | 日期（1-31） |
| hour   | 数字 | 是  | 小时（0-23） |
| minute | 数字 | 是  | 分钟（0-59） |
| sex    | 字符串 | 是  | 性别（男/女） |

### 时辰对照表
服务自动将小时转换为传统时辰：
```
0: 子时 (23:00-01:00)
1: 丑时 (01:00-03:00)
...
11: 亥时 (21:00-23:00)
```

## 服务配置
通过命令行参数配置：
```bash
# 指定监听地址和端口
go run main.go --host 127.0.0.1 --port 3000
```

## 依赖库
- [mcp-go](https://github.com/mark3labs/mcp-go) v0.20.1
- [goja](https://github.com/dop251/goja) JavaScript引擎
- [iztro](https://github.com/SylarLong/iztro) iztro.min.js 紫微斗数排盘核心逻辑文件 (位于项目根目录)
