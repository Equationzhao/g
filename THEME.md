## create custom theme

```ini
[info]
d 		= blue
l 		= purple
b 		= yellow
c 		= yellow
p 		= yellow
s 		= yellow
r 		= yellow
w 		= red
x 		= green
-       = white
time 	= blue
size 	= green
owner 	= yellow
group 	= yellow
reset 	= reset
root 	= red

[dir]
color = blue
icon = 📁

[exec,exe]
color = green
icon = 🚀

[file]
color = white
icon = 📄

...
```

## apply theme
```bash
g -theme=path/to/theme [other options] path
```

## default theme

see [default](theme/default.ini)