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

# Migrations

To create a new migration: 

```bash
make migration-generate:
```

It will create a new migration in `/migrations`.

Migrations are automatcly loaded in the database at the start of the server.

To disable migration, you can set the env variable 

```bash
MIGRATION_ENABLE=false
```