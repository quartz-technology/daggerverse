# Dagger

A dagger module to execute Dagger inside Dagger.

## Features

| Command                               | Done |
|---------------------------------------|------|
| Setup Dagger CLI from scratch         | ✅    |
| Install Dagger CLI inside a container | ✅    |
| Execute `run`,`query` and `call`      | ✅    |
| Publish module                        | ✅    |
| Manage module                         | ⏳    |
| Run with debug version                | ⏳    |

## Example

### Run

Execute a simple command

```shell
dagger call cli --version 0.9.2 run --args "ps"
✔ dagger call cli run [1.05s]
┃ PID   USER     TIME  COMMAND                                                                                                                                                        
┃     1 root      0:00 /_shim /bin/dagger run ps                                                                                                                                      
┃    15 root      0:00 /bin/dagger run ps                                                                                                                                             
┃    20 root      0:00 ps     
```

Execute with multiple arguments

```
dagger call cli --version 0.9.2 run --args "echo" --args "hello world"
✔ dagger call cli run [1.94s]
┃ hello world
```



Made with ❤️ by Quartz.