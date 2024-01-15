## 创建自定义主题 create custom theme 

example: [default theme](internal/theme/default.json)

## 应用主题 apply theme
```bash
g -theme=path/to/theme [other options] path
```

## 高级 advanced

### 嵌入自定义主题 build with custom theme
自定义主题必须在 theme 目录下，命名为 custom_builtin.json (internal/theme/custom_builtin.json)

the custom theme must be placed in the theme directory and named custom_builtin.json (theme/custom_builtin.json)

```bash
go build -tags=custom .
```

### 在配置文件中指定主题位置 specify theme location in config file
在配置文件中添加 `Theme: $location`，`$location` 可以通过 `g --help` 查看

add `Theme: $location` to your profile. `$location` can be found by `g --help` like this:
```bash
> g --help
# ...
CONFIG:
  Configuration: /Users/equationzhao/Library/Application Support/g/g.yaml
  See More at: g.equationzhao.space
# ...
```
下面是一个使用自定义主题的配置文件示例
this is an example of profile using custom theme
```yaml
Args:
  - hyperlink=never
  - icons
  - disable-index

CustomTreeStyle:
  Child: "├── "
  LastChild: "╰── "
  Mid: "│   "
  Empty: "    "
  
Theme: /Users/equationzhao/g/internal/theme/your_custom_theme.json
```