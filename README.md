# Isolate

A Linux command to run commands with configurable levels of isolation.
Inspired by some popular containerization tools,
like Docker.

# Instalation

```
; go get github.com/alcortesm/isolate
```

# Examples

Run a command without any isolation whatsoever:
```
; isolate echo no isolation at all
no isolation at all
```

Run a command and show it exit code:
```
; isolate -exitCode echo show me your exit code
show me your exit code
Exit Code 0
```

Runs a shell with isolated system identifiers:

```
; hostname
cherry
; isolate -uts bash
$ hostname foo
$ hostname
foo
$ exit
; hostname
cherry
```
