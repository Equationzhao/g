## 创建自定义主题 create custom theme 

example: [default theme](theme/default.json)

## 应用主题 apply theme
```bash
g -theme=path/to/theme [other options] path
```

## 高级 advanced

### 嵌入自定义主题 build with custom theme
自定义主题必须在 theme 目录下，命名为 custom_builtin.json (theme/custom_builtin.json)

the custom theme must be placed in the theme directory and named custom_builtin.json (theme/custom_builtin.json)

```bash
go build -tags=custom .
```