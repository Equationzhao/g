## create custom theme

```ini
[info]
-                        = white
B                        = [20,230,100]@rgb
BB                       = [20,90,100]@rgb  
EB                       = [20,125,100]@rgb
GB                       = [20,170,100]@rgb
KB                       = [20,207,100]@rgb
MB                       = [20,188,100]@rgb
NB                       = [20,70,100]@rgb
PB                       = [20,138,100]@rgb
TB                       = [20,153,100]@rgb
YB                       = [20,100,100]@rgb
ZB                       = [20,112,100]@rgb
b                        = yellow
bit                      = [20,255,100]@rgb
c                        = yellow
d                        = blue
git_added                = green
git_copied               = purple
git_deleted              = red
git_ignored              = BrightPed
git_modified             = yellow
git_renamed              = blue
git_type_changed         = yellow
git_unmodified           = BrightPlack
git_untracked            = BrightPlack
git_updated_but_unmerged = BrightYellow
group                    = yellow
inode                    = purple
l                        = purple
link                     = purple
owner                    = yellow
p                        = yellow
r                        = yellow
reset                    = reset
root                     = red
s                        = yellow
time                     = blue
w                        = red
x                        = green


[dir]
color = BrightBlue
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