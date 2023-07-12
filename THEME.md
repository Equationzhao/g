## create custom theme

example: [default theme](theme/default.json)

## apply theme
```bash
g -theme=path/to/theme [other options] path
```

## advanced

### build with custom theme
the custom must be placed in the theme directory and named custom_builtin.json (theme/custom_builtin.json)

```bash
go build -tags=custom .
```