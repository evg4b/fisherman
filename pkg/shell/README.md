# fisherman/pkg/shell

A package providing a simple cross-platform wrapper over the system shell.

Main features:
- Expandable (implements `shell.ShellStrategy` interface to integrate )
- Shell host implements `io.Writer` interface to simple put commands into `stdin`.
- The commands are executed after entering a new line. This can be used to execute
  commands over time using a single terminal process.

# Available shell strategies:

- **shell.Cmd()** - Strategy to use [cmd.exe](https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/cmd?ranMID=46131&ranEAID=a1LgFw09t88&ranSiteID=a1LgFw09t88-GnIhyB.n9js83pQZuKa7oQ&epi=a1LgFw09t88-GnIhyB.n9js83pQZuKa7oQ&irgwc=1&OCID=AID2200057_aff_7806_1243925&tduid=%28ir__tx9o062soskf6jneo33yzzo9bn2xoglw11jzryfe00%29%287806%29%281243925%29%28a1LgFw09t88-GnIhyB.n9js83pQZuKa7oQ%29%28%29&irclickid=_tx9o062soskf6jneo33yzzo9bn2xoglw11jzryfe00) as shell host. Can be used only on `windows`. It is also the default shell for windows.
- **shell.Bash()** - Strategy to use [bash](https://www.gnu.org/software/bash/) as shell host. Available on  `windows`, `linux` and `darwin`. It is the default shell for `linux` and `darwin`.   **Note**: for windows `bash` available only in package [MSYS2](https://www.msys2.org/).
- **shell.PowerShell()** - Strategy to use [powershell](https://docs.microsoft.com/en-us/powershell/scripting/overview?view=powershell-7.2) as shell host. Can be used on `windows`, `linux` and `darwin`.

# Examples:

### Basic
```golang
ctx := context.Background()
host := shell.NewHost(ctx, shell.Cmd())
if err := host.Run("echo 'this is demo shell' >> log.txt"); err != nil {
  panic(err)
}
```

### Advanced script input
```golang
ctx := context.Background()
host := shell.NewHost(ctx, shell.Cmd())

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

ctx := context.Background()
host := shell.NewHost(ctx, shell.Cmd(), shell.WithStdout(buffer))
if err := host.Run("ping -n 10 google.com"); err != nil {
  panic(err)
}

fmt.Print(buffer.String())
```
