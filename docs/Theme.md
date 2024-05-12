## create custom theme 

example: [default theme](../internal/theme/default.json)

## apply your theme

### through command line flag
```bash
g -theme=path/to/theme [other options] path
```

### set in the config file
add `Theme: $location` to your profile. `$location` can be found by `g --help` like this:

```bash
> g --help
# ...
CONFIG:
  Configuration: /Users/equationzhao/Library/Application Support/g/g.yaml
  See More at: g.equationzhao.space
# ...
```

here is an example of setting custom theme in the config file:

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

## advanced

### build with custom theme

the custom theme must be placed in the theme directory and named custom_builtin.json (theme/custom_builtin.json)

```bash
go build -tags=custom .
```
