# Isolate

A Linux command to run other commands in its own namespaces.

This program is inspired by some popular and more powerful tools,
like unshare, lxc or docker.

Most isolation levels require some capabilities(7),
in most cases CAP_SYS_ADMIN.
User namespaces are the  exception:
since  Linux 3.8, no privilege is required to create a user namespace;
Note that some Linux distributions have user namespaces disabled by default,
you can enable it with `echo 1 > /proc/sys/kernel/unprivileged_userns_clone`.

Running a command in its own new user namespace
allows to run the command in an unprivileged namespace environment
where the process runs with a full set of capabilities,
meaning you no longer need to execute `isolate` as root
to get the CAP_SYS_ROOT capability required by the other isolation levels.
See the `userns` example below.

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
  Requires CAP_SYS_ROOT.
  ```
  ; sudo isolate -dir /tmp/foo pwd
  /tmp/foo
  ```

- Run a command in a new user namespace,
  getting a full set of capabilities in the new namespace.
  ```
  ; isolate -dir /tmp/foo pwd
  fork/exec /bin/pwd: operation not permitted
  ;
  ; isolate -userns -dir /tmp/foo pwd
  /tmp/foo
  ```

- Run a shell with isolated system identifiers.
  Requires CAP_SYS_ADMIN.

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

