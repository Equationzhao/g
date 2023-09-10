# get the latest git tag like v0.14.0 and remove the prefix v
latest := `git describe --abbrev=0 --tags | sed 's/v//'`
ldflags := "-ldflags='-s -w'"

build: # build binaries for all platforms
    # build the binary in build/
    # Linux macOS Windows
    # 386 amd64 arm arm64
    mkdir -p build
    CGO_ENABLED=0 GOOS=linux GOARCH=386     go build {{ldflags}} -o build/g-linux-386
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64   go build {{ldflags}} -o build/g-linux-amd64
    CGO_ENABLED=0 GOOS=linux GOARCH=arm     go build {{ldflags}} -o build/g-linux-arm
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64   go build {{ldflags}} -o build/g-linux-arm64

    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64  go build {{ldflags}} -o build/g-darwin-amd64
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64  go build {{ldflags}} -o build/g-darwin-arm64

    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build {{ldflags}} -o build/g-windows-amd64.exe
    CGO_ENABLED=0 GOOS=windows GOARCH=386   go build {{ldflags}} -o build/g-windows-386.exe
    CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build {{ldflags}} -o build/g-windows-arm64.exe
    CGO_ENABLED=0 GOOS=windows GOARCH=arm   go build {{ldflags}} -o build/g-windows-arm.exe

    upx build/g-linux-386
    upx build/g-linux-amd64
    upx build/g-linux-arm
    upx build/g-linux-arm64
    upx build/g-darwin-amd64
    upx build/g-windows-amd64.exe
    upx build/g-windows-386.exe

compress: # compress the binaries for all platforms
    # g_OS_ARCH.tar.gz or g_Windows_ARCH.zip
    tar -zcvf build/g-Linux-386.tar.gz build/g-linux-386
    tar -zcvf build/g-Linux-amd64.tar.gz build/g-linux-amd64
    tar -zcvf build/g-Linux-arm.tar.gz build/g-linux-arm
    tar -zcvf build/g-Linux-arm64.tar.gz build/g-linux-arm64

    tar -zcvf build/g-Darwin-amd64.tar.gz build/g-darwin-amd64
    tar -zcvf build/g-Darwin-arm64.tar.gz build/g-darwin-arm64

    zip -r build/g-Windows-amd64.zip build/g-windows-amd64.exe
    zip -r build/g-Windows-386.zip build/g-windows-386.exe
    zip -r build/g-Windows-arm64.zip build/g-windows-arm64.exe
    zip -r build/g-Windows-arm.zip build/g-windows-arm.exe

deb: # build deb package for all ARCH
    # g_VERSION_ARCH.deb
    # |__DEBIAN
    # |  |__control
    # |__usr
    #    |__local
    #       |__bin
    #       |  |__g
    #       |__share
    #           |__man
    #              |__man1
    #                 |__g.1.gz

    # amd64
    mkdir -vp build/g_{{latest}}_amd64/DEBIAN
    mkdir -vp build/g_{{latest}}_amd64/usr/local/bin
    mkdir -vp build/g_{{latest}}_amd64/usr/local/share/man/man1
    cp build/g-linux-amd64 build/g_{{latest}}_amd64/usr/local/bin/g
    cp man/g.1.gz build/g_{{latest}}_amd64/usr/local/share/man/man1/g.1.gz

    echo "Package: g" > build/g_{{latest}}_amd64/DEBIAN/control
    echo "Version: {{latest}}" >> build/g_{{latest}}_amd64/DEBIAN/control
    echo "Architecture: amd64" >> build/g_{{latest}}_amd64/DEBIAN/control
    echo "Maintainer: Equationzhao <equationzhao at foxmail.com>" >> build/g_{{latest}}_amd64/DEBIAN/control
    echo "Description: a powerful ls tool" >> build/g_{{latest}}_amd64/DEBIAN/control

    dpkg-deb -b build/g_{{latest}}_amd64
    rm -rf build/g_{{latest}}_amd64


    # arm64
    mkdir -vp build/g_{{latest}}_arm64/DEBIAN
    mkdir -vp build/g_{{latest}}_arm64/usr/local/bin
    mkdir -vp build/g_{{latest}}_arm64/usr/local/share/man/man1
    cp build/g-linux-arm64 build/g_{{latest}}_arm64/usr/local/bin/g
    cp man/g.1.gz build/g_{{latest}}_arm64/usr/local/share/man/man1/g.1.gz

    echo "Package: g" > build/g_{{latest}}_arm64/DEBIAN/control
    echo "Version: {{latest}}" >> build/g_{{latest}}_arm64/DEBIAN/control
    echo "Architecture: arm64" >> build/g_{{latest}}_arm64/DEBIAN/control
    echo "Maintainer: Equationzhao <equationzhao at foxmail.com>" >> build/g_{{latest}}_arm64/DEBIAN/control
    echo "Description: a powerful ls tool" >> build/g_{{latest}}_arm64/DEBIAN/control

    dpkg-deb -b build/g_{{latest}}_arm64
    rm -rf build/g_{{latest}}_arm64


    # 386
    mkdir -vp build/g_{{latest}}_386/DEBIAN
    mkdir -vp build/g_{{latest}}_386/usr/local/bin
    mkdir -vp build/g_{{latest}}_386/usr/local/share/man/man1
    cp build/g-linux-386 build/g_{{latest}}_386/usr/local/bin/g
    cp man/g.1.gz build/g_{{latest}}_386/usr/local/share/man/man1/g.1.gz

    echo "Package: g" > build/g_{{latest}}_386/DEBIAN/control
    echo "Version: {{latest}}" >> build/g_{{latest}}_386/DEBIAN/control
    echo "Architecture: 386" >> build/g_{{latest}}_386/DEBIAN/control
    echo "Maintainer: Equationzhao <equationzhao at foxmail.com>" >> build/g_{{latest}}_386/DEBIAN/control
    echo "Description: a powerful ls tool" >> build/g_{{latest}}_386/DEBIAN/control

    dpkg-deb -b build/g_{{latest}}_386
    rm -rf build/g_{{latest}}_386

    # arm
    mkdir -vp build/g_{{latest}}_arm/DEBIAN
    mkdir -vp build/g_{{latest}}_arm/usr/local/bin
    mkdir -vp build/g_{{latest}}_arm/usr/local/share/man/man1
    cp build/g-linux-arm build/g_{{latest}}_arm/usr/local/bin/g
    cp man/g.1.gz build/g_{{latest}}_arm/usr/local/share/man/man1/g.1.gz

    echo "Package: g" > build/g_{{latest}}_arm/DEBIAN/control
    echo "Version: {{latest}}" >> build/g_{{latest}}_arm/DEBIAN/control
    echo "Architecture: arm" >> build/g_{{latest}}_arm/DEBIAN/control
    echo "Maintainer: Equationzhao <equationzhao at foxmail.com>" >> build/g_{{latest}}_arm/DEBIAN/control
    echo "Description: a powerful ls tool" >> build/g_{{latest}}_arm/DEBIAN/control

    dpkg-deb -b build/g_{{latest}}_arm
    rm -rf build/g_{{latest}}_arm

checksum: # generate the checksum file for all build files
    shasum -a 256 build/* > build/checksum.txt

release: # release to github
    gh release create v{{latest}} build/*

url := "https://github.com/Equationzhao/g"

aur: # update aur 
    # download the latest source code from {{url}}/"archive/refs/tags/v{{latest}}.tar.gz"
    #!/usr/bin/env bash
    wget -c {{url}}/archive/refs/tags/v{{latest}}.tar.gz -O v{{latest}}.tar.gz

    # update PKGBUILD
    sed -i bak "s/sha256sums=.*/sha256sums=('$(shasum -a 256 v{{latest}}.tar.gz | choose 0)')/g" ../g-ls/PKGBUILD
    sed -i bak "s/pkgver=.*/pkgver={{latest}}/g" ../g-ls/PKGBUILD

    # update .SRCINFO
    sed -i bak "s/pkgver = .*/pkgver = {{latest}}/g" ../g-ls/.SRCINFO
    sed -i bak "s/sha256sums = .*/sha256sums = '$(shasum -a 256 v{{latest}}.tar.gz | choose 0)'/g" ../g-ls/.SRCINFO

    # input git commit message
    git add -u
    git commit
    git push

brew: # update homebrew
    # class GLs < Formula
    #   desc "a powerfull cross-platform ls"
    #   homepage "g.equationzhao.space"
    #   url "https://github.com/Equationzhao/g/archive/refs/tags/v0.13.2.tar.gz" , :tag => "v0.13.2"
    #   sha256 "be8afefb7952c2e74127a55fdc0c056fadbd58e652957af1b9f913d0e0a82123"  
    #   license "MIT"
    #   depends_on "go" => :build
    #   def install
    #     system "go build -ldflags='-s -w'"
    #     bin.install "g"
    #   end
    # end
    sed -i bak "s#url .*#url \"{{url}}/archive/refs/tags/v{{latest}}.tar.gz\" , :tag => \"v{{latest}}\"#g" ../homebrew-g/g-ls.rb
    sed -i bak "s/sha256 .*/sha256 \"$(shasum -a 256 v{{latest}}.tar.gz | choose 0)\"/g" ../homebrew-g/g-ls.rb
    git add -u
    git commit
    git push

all : build compress deb checksum 

clean: # clean the build directory
    rm -rf build

format: # format the code
    gofumpt -w -l .

doc: # generate the documentation
    go build -tags 'doc'
    ./g 
    rm g

theme: # generate the theme
    go build -tags 'theme'
    ./g 
    rm g