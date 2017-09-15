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

- Run a command in a root jail.
  Requires capability CAP_SYS_ROOT:
  ```
  ; sudo isolate -dir /tmp/foo pwd
  /tmp/foo
  ```

- Run a command in a new user namespace.
  This allows to run the command in an unprivileged namespace environment
  where the process runs with a full set of capabilities
  (meaning you no longer need to execute isolate as root).
  ```
  ; isolate -dir /tmp/foo pwd
  fork/exec /bin/pwd: operation not permitted
  ;
  ; isolate -userns -dir /tmp/foo pwd
  /tmp/foo
  ```

  Note that user namespaces are disabled by default in some Linux distributions,
  you can enable it with `echo 1 > /proc/sys/kernel/unprivileged_userns_clone`.
  

- Run a shell with isolated system identifiers.
  Requires capability CAP_SYS_ADMIN:
  ```
  ; sudo isolate -uts bash
  $ hostname
  cherry
  $ hostname foo
  $ hostname
  foo
  $ exit
  ; hostname
  cherry
  ```

