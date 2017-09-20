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
  See chroot(2).
  ```
  ; sudo isolate -dir /tmp/foo pwd
  /tmp/foo
  ```

- Run a command in a new user namespace,
  getting a full set of capabilities in the new namespace.
  See user_namespaces(7).
  ```
  ; isolate -dir /tmp/foo pwd
  fork/exec /bin/pwd: operation not permitted
  ;
  ; isolate -userns -dir /tmp/foo pwd
  /tmp/foo
  ```

- Run a shell with isolated system identifiers.
  Requires CAP_SYS_ADMIN.
  See namespaces(7).

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
- Run a command in its own pid namespace.
  Requires CAP_SYS_ADMIN.
  See pid_namespaces(7).
  ```
  ; sudo isolate -pid sh
  # echo $$
  1
  ```

  A /proc filesystem shows (in the /proc/PID directories) only processes visible in the PID namespace of the process that performed the mount,
  even if the /proc filesystem is viewed from processes in other namespaces.

  After creating a new PID namespace, it is useful for the child to change its root directory and mount a new procfs instance at /proc so that tools such as ps(1) work correctly.

  To achieve this, isolate a shell with the `pid` and the `mount` options
  and mount a new proc filesystem on top of the old one as follows:
  ```
  $ mount -t proc proc /proc
  ```

- Run a command in its own mount namespace, see mount_namespaces(7).
  This isolates mount operations, as long as the propagation type
  of the mount point is not set to `MS_SHARED`
  (which is usually the case for '/').

  In this example we will bind mount `/proc` to `/tmp/proc` in a subshell,
  a fairly innocuous operation,
  without affecting the mount points of the parent shell.
  To make this example complete,
  the propagation type of '/' in the parent shell is `MS_SHARED`,
  so we will have to change it before the mounting operation in the subshell:
  ```
  ; mkdir /tmp/proc
  ; findmnt -o TARGET,PROPAGATION /
  TARGET PROPAGATION
  /      shared
  ; PS1='# ' sudo isolate -mount sh
  # 
  # 
  # mount --make-rprivate / # stop propagating mounts to the peer group
  # findmnt -o TARGET,PROPAGATION /
  TARGET PROPAGATION
  /      private
  # mount --bind /proc /tmp/proc
  # mount | grep 'proc '
  proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
  proc on /tmp/proc type proc (rw,nosuid,nodev,noexec,relatime)
  # exit
  ;
  ;
  ; findmnt -o TARGET,PROPAGATION /
  TARGET PROPAGATION
  /      shared
  ; mount | grep 'proc '
  proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
  ```
