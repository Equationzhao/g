FROM golang
RUN git clone https://github.com/Equationzhao/g.git
RUN go install github.com/goreleaser/goreleaser@latest \
    && go install mvdan.cc/gofumpt@latest \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
    && go install github.com/daixiang0/gci@latest \
    && cd g \
    && go install

