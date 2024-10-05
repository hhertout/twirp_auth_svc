# Installing protobuf

`brew install protobuf`

`go install github.com/twitchtv/twirp/protoc-gen-twirp@latest`
`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

To add in your zshrc : 

`export PATH="~/go/bin:$PATH"`
`source ~/.zshrc`

# Generate protobuf

run `make generate file=<path_to_file_in_rpc_folder`
