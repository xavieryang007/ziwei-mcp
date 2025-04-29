package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"log"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type DateService struct {
	server *server.MCPServer
	addr   string
}

func NewDateService(addr string) *DateService {
	return &DateService{
		server: server.NewMCPServer(
			"date-mcp",
			"0.1.0",
		),
		addr: addr,
	}
}

func (s *DateService) Start() error {
	// Create date tool
	yangliTool := mcp.NewTool("ziwei_yangli",
		mcp.WithDescription("通过阳历/公历年月日时分获取紫微斗数排盘数据信息"),
		mcp.WithNumber("year",
			mcp.Required(),
			mcp.Description("Year (e.g., 2024)"),
		),
		mcp.WithNumber("month",
			mcp.Required(),
			mcp.Description("Month (1-12)"),
		),
		mcp.WithNumber("day",
			mcp.Required(),
			mcp.Description("Day (1-31)"),
		),
		mcp.WithNumber("hour",
			mcp.Required(),
			mcp.Description("Hour (0-23)"),
		),
		mcp.WithNumber("minute",
			mcp.Required(),
			mcp.Description("Minute (0-59)"),
		),
		mcp.WithString("sex",
			mcp.Required(),
			mcp.Description("性别 值为(男/女)"),
		),
	)

	// Add date tool handler
	s.server.AddTool(yangliTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get parameters from request
		year, ok := request.Params.Arguments["year"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid year")
		}
		month, ok := request.Params.Arguments["month"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid month")
		}
		day, ok := request.Params.Arguments["day"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid day")
		}
		hour, ok := request.Params.Arguments["hour"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid hour")
		}
		minute, ok := request.Params.Arguments["minute"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid minute")
		}
		hourInt := int(hour)
		hourInt = hourToTimePeriod(hourInt)
		// Create time object
		date := time.Date(
			int(year),
			time.Month(month),
			int(day),
			int(hour),
			int(minute),
			0, 0, time.UTC,
		)
		//获取这样格式的时间 2000-01-01
		dateStr := date.Format("2006-01-02")
		data := jsRuntime(dateStr, hourInt, request.Params.Arguments["sex"].(string), 1)
		return mcp.NewToolResultText(data), nil
	})

	nongliTool := mcp.NewTool("ziwei_nongli",
		mcp.WithDescription("通过农历/阴历年月日时分获取紫微斗数排盘数据信息"),
		mcp.WithNumber("year",
			mcp.Required(),
			mcp.Description("Year (e.g., 2024)"),
		),
		mcp.WithNumber("month",
			mcp.Required(),
			mcp.Description("Month (1-12)"),
		),
		mcp.WithNumber("day",
			mcp.Required(),
			mcp.Description("Day (1-31)"),
		),
		mcp.WithNumber("hour",
			mcp.Required(),
			mcp.Description("Hour (0-23)"),
		),
		mcp.WithNumber("minute",
			mcp.Required(),
			mcp.Description("Minute (0-59)"),
		),
		mcp.WithString("sex",
			mcp.Required(),
			mcp.Description("性别 值为(男/女)"),
		),
	)

	// Add date tool handler
	s.server.AddTool(nongliTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get parameters from request
		year, ok := request.Params.Arguments["year"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid year")
		}
		month, ok := request.Params.Arguments["month"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid month")
		}
		day, ok := request.Params.Arguments["day"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid day")
		}
		hour, ok := request.Params.Arguments["hour"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid hour")
		}
		minute, ok := request.Params.Arguments["minute"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid minute")
		}
		hourInt := int(hour)
		hourInt = hourToTimePeriod(hourInt)
		// Create time object
		date := time.Date(
			int(year),
			time.Month(month),
			int(day),
			int(hour),
			int(minute),
			0, 0, time.UTC,
		)
		//获取这样格式的时间 2000-01-01
		dateStr := date.Format("2006-01-02")
		data := jsRuntime(dateStr, hourInt, request.Params.Arguments["sex"].(string), 2)
		return mcp.NewToolResultText(data), nil
	})
	// Start the server
	return server.NewSSEServer(s.server).Start(s.addr)
}

// Convert hour to Chinese time period (时辰)
// Returns 0-11 where:
// 0: 子时 (23:00-01:00)
// 1: 丑时 (01:00-03:00)
// 2: 寅时 (03:00-05:00)
// 3: 卯时 (05:00-07:00)
// 4: 辰时 (07:00-09:00)
// 5: 巳时 (09:00-11:00)
// 6: 午时 (11:00-13:00)
// 7: 未时 (13:00-15:00)
// 8: 申时 (15:00-17:00)
// 9: 酉时 (17:00-19:00)
// 10: 戌时 (19:00-21:00)
// 11: 亥时 (21:00-23:00)
func hourToTimePeriod(hour int) int {
	return ((hour + 1) % 24) / 2
}

func jsRuntime(nianyueri string, shichen int, sex string, _type int) string {
	vm := goja.New()
	registry := new(require.Registry) // this can be shared by multiple runtimes
	registry.Enable(vm)

	vm.Set("self", vm.GlobalObject())
	// 定义一个完整的 `console` 对象
	console := make(map[string]interface{})

	// log 方法
	console["log"] = func(call goja.FunctionCall) goja.Value {
		fmt.Print("[LOG] ")
		for _, arg := range call.Arguments {
			fmt.Print(arg.String(), " ")
		}
		fmt.Println()
		return goja.Undefined()
	}

	// warn 方法
	console["warn"] = func(call goja.FunctionCall) goja.Value {
		fmt.Print("[WARN] ")
		for _, arg := range call.Arguments {
			fmt.Print(arg.String(), " ")
		}
		fmt.Println()
		return goja.Undefined()
	}

	// error 方法
	console["error"] = func(call goja.FunctionCall) goja.Value {
		fmt.Fprint(os.Stderr, "[ERROR] ")
		for _, arg := range call.Arguments {
			fmt.Fprint(os.Stderr, arg.String(), " ")
		}
		fmt.Fprintln(os.Stderr)
		return goja.Undefined()
	}

	// info 方法
	console["info"] = func(call goja.FunctionCall) goja.Value {
		fmt.Print("[INFO] ")
		for _, arg := range call.Arguments {
			fmt.Print(arg.String(), " ")
		}
		fmt.Println()
		return goja.Undefined()
	}

	// debug 方法
	console["debug"] = func(call goja.FunctionCall) goja.Value {
		fmt.Print("[DEBUG] ")
		for _, arg := range call.Arguments {
			fmt.Print(arg.String(), " ")
		}
		fmt.Println()
		return goja.Undefined()
	}
	SCRIPT := ""
	if _type == 1 {
		yangli := fmt.Sprintf(`
    	var {astro} = require("./iztro.min.js");
   		const astrolabe = astro.bySolar("%s", %d, "%s");
var outdata = JSON.stringify(astrolabe);
    `, nianyueri, shichen, sex)
		SCRIPT = yangli
	} else {
		nongli := fmt.Sprintf(`
    	var {astro} = require("./iztro.min.js");
   		const astrolabe = astro.bySolar("%s", %d, "%s");
		var outdata = JSON.stringify(astrolabe);
    `, nianyueri, shichen, sex)
		SCRIPT = nongli
	}

	// 将 `console` 注册为全局变量
	vm.Set("console", console)
	_, err := vm.RunString(SCRIPT)

	if err != nil {
		fmt.Println(err)
	}
	data := vm.Get("outdata").String()
	return data
}
func main() {

	// Define command line flags
	host := flag.String("host", "0.0.0.0", "Server host address")
	port := flag.Int("port", 8080, "Server port number")
	flag.Parse()

	// Construct server address
	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Starting Date Service on %s", addr)

	// Create and start service
	service := NewDateService(addr)
	if err := service.Start(); err != nil {
		log.Fatalf("Failed to start date service: %v", err)
	}
}
