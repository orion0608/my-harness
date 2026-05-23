# /proc-run

启动当前工作区的软件服务。

## Agent 行为

1. 确认当前工作区（Shell 操作以该路径为 cwd）
2. 检测项目技术栈；确认依赖已安装（未安装则提示或安装）
3. 按项目约定顺序启动服务（前后端、单体等）
4. 验证启动成功

## Web 前端

项目含 Web/H5 需在后端启动后在IDE浏览器访问协助用户打开前端界面：

- 须 **cursor-ide-browser**：`browser_navigate`  
- 汇报时给出完整 `http(s)://` URL
