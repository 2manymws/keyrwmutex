# keyrwmutex

`keyrwmutex` is a package that provides R/W locks on arbitrary strings using hash bucket.

This is created by combining [keymutex.hashedKeyMutex](https://pkg.go.dev/k8s.io/utils/keymutex#NewHashed) and [sync.RWMutex](https://pkg.go.dev/sync#RWMutex).

## References

- [k8s.io/utils/keymutex](github.com/kubernates/utils)
