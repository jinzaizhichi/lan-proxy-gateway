// Package script runs a user-provided JavaScript "enhancer" over the mihomo
// config YAML, same shape as Clash Verge Rev 的 "Script" feature:
//
//	function main(config) { /* mutate */ return config; }
//
// 调用 Apply(scriptPath, configYAML) 得到被 JS 修改过后的 YAML。
// 失败返回原 YAML 前的 error，不污染 input。
//
// 实现思路：
//  1. YAML → Go generic value → JSON（避免 Go↔JS 复杂类型映射）
//  2. 在 goja VM 里 JSON.parse 成 JS 对象
//  3. 调 main(config)
//  4. JSON.stringify 结果 → Go generic value → YAML
package script

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dop251/goja"
	"gopkg.in/yaml.v3"
)

// Apply reads the JS script, runs main(config) against configYAML, returns
// the modified YAML. The script must export a `main(config) -> config` function.
func Apply(scriptPath string, configYAML []byte) ([]byte, error) {
	scriptCode, err := os.ReadFile(scriptPath)
	if err != nil {
		return nil, fmt.Errorf("读取脚本失败: %w", err)
	}

	// YAML → interface{} → JSON
	var configObj interface{}
	if err := yaml.Unmarshal(configYAML, &configObj); err != nil {
		return nil, fmt.Errorf("解析 config YAML 失败: %w", err)
	}
	jsonBytes, err := json.Marshal(configObj)
	if err != nil {
		return nil, fmt.Errorf("config → JSON 失败: %w", err)
	}

	vm := goja.New()
	if _, err := vm.RunScript(scriptPath, string(scriptCode)); err != nil {
		return nil, fmt.Errorf("脚本加载失败: %w", err)
	}

	mainFn, ok := goja.AssertFunction(vm.Get("main"))
	if !ok {
		return nil, fmt.Errorf("脚本里找不到 main 函数（签名应为 function main(config) { return config; }）")
	}

	// 在 JS 里 JSON.parse 出 config，避开 Go→JS 的深层类型转换坑
	parseFn, _ := goja.AssertFunction(vm.Get("JSON").ToObject(vm).Get("parse"))
	configVal, err := parseFn(goja.Undefined(), vm.ToValue(string(jsonBytes)))
	if err != nil {
		return nil, fmt.Errorf("JSON.parse 失败: %w", err)
	}

	result, err := mainFn(goja.Undefined(), configVal)
	if err != nil {
		return nil, fmt.Errorf("脚本 main() 报错: %w", err)
	}

	stringifyFn, _ := goja.AssertFunction(vm.Get("JSON").ToObject(vm).Get("stringify"))
	jsonResult, err := stringifyFn(goja.Undefined(), result)
	if err != nil {
		return nil, fmt.Errorf("JSON.stringify 失败: %w", err)
	}

	var resultObj interface{}
	if err := json.Unmarshal([]byte(jsonResult.String()), &resultObj); err != nil {
		return nil, fmt.Errorf("解析脚本返回值失败: %w", err)
	}
	output, err := yaml.Marshal(resultObj)
	if err != nil {
		return nil, fmt.Errorf("生成 YAML 失败: %w", err)
	}
	return output, nil
}
