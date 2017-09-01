# Isolate

A Linux command to run other commands in its own namespaces.

This program is inspired by some popular and more powerful tools,
like unshare, lxc or docker.

# Instalation

```
; go get github.com/alcortesm/isolate
```

# Examples

- Run a command without any isolation whatsoever:
  ```
  ; isolate echo no isolation at all
  no isolation at all
  ```

- Run a command and show its exit code:
  ```
  ; isolate -exitCode echo show me your exit code
  show me your exit code
  Exit Code 0
  ```

- Run a shell with isolated system identifiers:
  ```
  ; hostname
  cherry
  ;
  ;
  ; isolate -uts bash
  $ hostname
  cherry
  $ hostname foo
  $ hostname
  foo
  $ exit
  ;
  ;
  ; hostname
  cherry
  ```

- Run a command in a root jail:
  ```
  ; isolate -dir /tmp/foo ls /bla
  [will show the list of files at /tmp/foo/bla]
  ```
