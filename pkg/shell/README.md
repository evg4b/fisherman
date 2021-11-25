# fisherman/shell

### Simple example
```golang
ctx := context.TODO()
host := shell.NewHost(ctx, shell.DefaultCmd())
if err := host.Run("echo 'this is demo shell' >> log.txt"); err != nil {
  panic(err)
}
```

### Simple example
```golang
ctx := context.TODO()
host := shell.NewHost(ctx, shell.DefaultCmd())


_, err = fmt.Fprintln(host, "echo 'this is demo shell' >> log.txt")
if err != nil {
  panic(err)
}

_, err = fmt.Fprintln(host, "ping -n 10 google.com")
if err != nil {
  panic(err)
}

if err := host.Wait(); err != nil {
  panic(err)
}
```

### Collect shell output
```golang
buffer := bytes.NewBufferString("")

ctx := context.TODO()
host := shell.NewHost(ctx, shell.DefaultCmd(), shell.WithStdout(buffer))
if err := host.Run("ping -n 10 google.com"); err != nil {
  panic(err)
}

fmt.Print(buffer.String())
```
