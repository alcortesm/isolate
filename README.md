# Isolate
A command that runs other commands in configurable isolated environments.

# Instalation

```
; go get github.com/alcortesm/isolate
```

# Examples

Run command without any isolation whatsoever:
```
; isolate echo no isolation at all
no isolation at all
```

Show exit code of the command:
```
; isolate -exitCode echo show me your exit code
show me your exit code
Exit Code 0
```

Isolate system identifiers:

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
