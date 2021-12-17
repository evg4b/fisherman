# fisherman/pkg/shell

A package providing a lightweight cross-platform wrapper over the system shell.

Features:
- Expandable (implements `shell.ShellStrategy` interface to integrate)
- Shell host implements `io.Writer` interface to simple put commands into `stdin`.
- The commands are executed after entering a new line. This can be used to execute
  commands over time using a single terminal process.

# Available shell strategies:

- **shell.Cmd()** - Strategy to use [cmd.exe](https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/cmd?ranMID=46131&ranEAID=a1LgFw09t88&ranSiteID=a1LgFw09t88-GnIhyB.n9js83pQZuKa7oQ&epi=a1LgFw09t88-GnIhyB.n9js83pQZuKa7oQ&irgwc=1&OCID=AID2200057_aff_7806_1243925&tduid=%28ir__tx9o062soskf6jneo33yzzo9bn2xoglw11jzryfe00%29%287806%29%281243925%29%28a1LgFw09t88-GnIhyB.n9js83pQZuKa7oQ%29%28%29&irclickid=_tx9o062soskf6jneo33yzzo9bn2xoglw11jzryfe00) as shell host. Can be used only on `windows`. It is also the default shell for windows.
- **shell.Bash()** - Strategy to use [bash](https://www.gnu.org/software/bash/) as shell host. Available on  `windows`, `linux` and `darwin`. It is the default shell for `linux` and `darwin`.   **Note**: for windows `bash` available only in package [MSYS2](https://www.msys2.org/).
- **shell.PowerShell()** - Strategy to use [powershell](https://docs.microsoft.com/en-us/powershell/scripting/overview?view=powershell-7.2) as shell host. Can be used on `windows`, `linux` and `darwin`.

# Options:

- **shell.WithCwd(dir)** - specifies the working directory of the shell. If it is the empty ot not set, host runs the shell in the calling process's current directory.
- **shell.WithStdout(output)** - specify the process's standard output.
- **shell.WithStderr(output)** - specify the process's standard error.
- **shell.WithEnv(env)** - specifies the environment of the shell process (including the shell strategy provided variables).
- **shell.WithRawEnv(env)** - specifies the environment of the shell process (without the shell strategy provided variables).
- **shell.WithArgs(args)** - holds command line arguments (including the shell strategy provided args)
- **shell.WithRawArgs(args)** - holds command line arguments (without the shell strategy provided args)

# Encoding:
This wrapper does not work with the shell encodings. To execute commands in a different encoding, you can wrap the host in Encoder and Stdout in Decoder from packages `golang.org/x/text/encoding/charmap` and `golang.org/x/text/transform` like this:

``` GO
cp := charmap.CodePage866 // Your encoding
buffer := bytes.NewBufferString("")

host := shell.NewHost(ctx, shell.Cmd(), shell.WithStdout(buffer))
convertedHost := transform.NewWriter(host, cp.NewEncoder())

fmt.Fprintln(convertedHost, "chcp")

if err := host.Wait(); err != nil {
  panic(err)
}

output, _, err := transform.String(cp.NewDecoder(), buff.String())
if err != nil {
  panic(err)
}

fmt.Print(output)
```

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

## Helpers

Also, the package contains helpers functions in `fisherman/pkg/shell/helpers`. These functions can be useful when working with shell host.

### MergeEnv

- **helpers.MergeEnv** - Merges a slice of environment variables (for example `os.Environ()`) with a map of custom variables.
  In case when variable already defined in slice, it will replaced from map.
