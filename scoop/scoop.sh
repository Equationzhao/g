#! /bin/sh

latest=$(git describe --abbrev=0 --tags | sed 's/v//')
hash64=$(shasum -a 256 ../build/g-windows-amd64.exe | cut -d ' ' -f 1)
hash386=$(shasum -a 256 ../build/g-windows-386.exe | cut -d ' ' -f 1)
hasharm64=$(shasum -a 256 ../build/g-windows-arm64.exe | cut -d ' ' -f 1)

echo "{ \
 \"homepage\": \"g.equationzhao.space\", \
 \"bin\": \"bin/g.exe\", \
 \"architecture\": { \
     \"64bit\": { \
     \"url\": \"https://github.com/Equationzhao/g/releases/download/v$latest/g-windows-amd64.exe\", \
     \"hash\": \"$hash64\", \
     \"bin\": \"g-windows-amd64.exe\", \
     \"post_install\":[ \
         \"cd \$scoopdir/shims\", \
         \"mv g-windows-amd64.exe g.exe\", \
         \"mv g-windows-amd64.shim g.shim\" \
     ], \
     \"shortcuts\":[ \
         [ \
         \"g-windows-amd64.exe\", \
         \"g\" \
         ] \
     ] \
     }, \
     \"32bit\": { \
     \"url\": \"https://github.com/Equationzhao/g/releases/download/v$latest/g-windows-386.exe\", \
     \"hash\": \"$hash386\", \
     \"bin\": \"g-windows-386.exe\", \
     \"post_install\":[ \
         \"cd \$scoopdir/shims\", \
         \"mv g-windows-386.exe g.exe\", \
         \"mv g-windows-386.shim g.shim\" \
     ], \
     \"shortcuts\":[ \
         [ \
         \"g-windows-386.exe\", \
         \"g\" \
         ] \
     ] \
     }, \
     \"arm64\": { \
     \"url\": \"https://github.com/Equationzhao/g/releases/download/v$latest/g-windows-arm64.exe\", \
     \"hash\": \"$hasharm64\", \
     \"bin\": \"g-windows-arm64.exe\", \
     \"post_install\":[ \
         \"cd \$scoopdir/shims\", \
         \"mv g-windows-arm64.exe g.exe\", \
         \"mv g-windows-arm64.shim g.shim\" \
     ], \
     \"shortcuts\":[ \
         [ \
         \"g-windows-arm64.exe\", \
         \"g\" \
         ] \
     ] \
     } \
 }, \
 \"license\": \"MIT\", \
 \"version\": \"v$latest\" \
}" > g.json

prettier -w g.json