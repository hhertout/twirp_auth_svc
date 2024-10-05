# Installing protobuf

`brew install protobuf`

add to your .zshrc : 

```bash
export GOBIN=$PWD/bin
export PATH=$GOBIN:$PATH
```

`make proto_install` to install proto bin locally.

It will create a bin directory, witch it gitignore by default.

# Generate protobuf

run `make generate file=<path_to_file_in_rpc_folder`
