# fisherman/pkg/shell

A package providing a lightweight cross-platform wrapper over the system shell.

Features:

- Expandable (implements `shell.ShellStrategy` interface to integrate)
- Shell host implements `io.Writer` interface to simple put commands into `stdin`.
- The commands are executed after entering a new line. This can be used to execute
  commands over time using a single terminal process.

## Available shell strategies

- **shell.Cmd()** - Strategy to
  use [cmd.exe](https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/cmd?ranMID=46131&ranEAID=a1LgFw09t88&ranSiteID=a1LgFw09t88-GnIhyB.n9js83pQZuKa7oQ&epi=a1LgFw09t88-GnIhyB.n9js83pQZuKa7oQ&irgwc=1&OCID=AID2200057_aff_7806_1243925&tduid=%28ir__tx9o062soskf6jneo33yzzo9bn2xoglw11jzryfe00%29%287806%29%281243925%29%28a1LgFw09t88-GnIhyB.n9js83pQZuKa7oQ%29%28%29&irclickid=_tx9o062soskf6jneo33yzzo9bn2xoglw11jzryfe00)
  as shell host. Can be used only on `windows`. It is also the default
  shell for windows.
- **shell.Bash()** - Strategy to use [bash](https://www.gnu.org/software/bash/)
  as shell host. Available on  `windows`, `linux` and `darwin`. It is the default
  shell for `linux` and `darwin`. **Note**: for windows `bash` available only
  in package [MSYS2](https://www.msys2.org/).
- **shell.PowerShell()** - Strategy to
  use [powershell](https://docs.microsoft.com/en-us/powershell/scripting/overview?view=powershell-7.2)
  as shell host. Can be used on `windows`, `linux` and `darwin`.

## Options

- **shell.WithCwd(dir)** - specifies the working directory of the shell.
  If it is the empty ot not set, host runs the shell in the calling process's
  current directory.
- **shell.WithStdout(output)** - specify the process's standard output.
- **shell.WithStderr(output)** - specify the process's standard error.
- **shell.WithEnv(env)** - specifies the environment of the shell process
  (including the shell strategy provided variables).
- **shell.WithRawEnv(env)** - specifies the environment of the shell process
  (without the shell strategy provided variables).
- **shell.WithArgs(args)** - holds command line arguments (including the shell
  strategy provided args)
- **shell.WithRawArgs(args)** - holds command line arguments (without the shell
  strategy provided args)
- **shell.WithEncoding(encoding)** - overwrite encoding for shell. (To setup utf-8
  encoding pass `encoding.Nop`)

## Examples

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

### Custom encoding

``` golang
buffer := bytes.NewBuffer([]byte{})

host := shell.NewHost(
  ctx,
  shell.Cmd(),
  shell.WithStdout(buffer),
  shell.WithEncoding(charmap.Windows1251),
)

fmt.Fprintln(host, "chcp 1251 > nul")
fmt.Fprintln(host, "echo проверка русского текста")

if err := host.Wait(); err != nil {
  panic(err)
}

fmt.Print(buffer.String())
```

## Helpers

Also, the package contains helpers functions in `fisherman/pkg/shell/helpers`.
These functions can be useful when working with shell host.

### MergeEnv

- **helpers.MergeEnv** - Merges a slice of environment variables
  (for example `os.Environ()`) with a map of custom variables.
  In case when variable already defined in slice, it will replaced from map.

## Known issues

There is a problem with rendering utf8 output on windows. This is related to the
problem with UTF-8 support in Windows. ReadConsoleA and ReadFile just silently
fail with the code page set to 65001. Read
more [here](https://social.msdn.microsoft.com/Forums/vstudio/en-US/6db367e1-6b39-4c91-bd08-e3779ae5fc23/problems-with-readingwriting-utf8-characters-to-console?forum=vcgeneral).
