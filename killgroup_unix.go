//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || aix
// +build darwin dragonfly freebsd linux netbsd openbsd solaris aix

package main

import (
	"syscall"
)

func killGroup() {
	pgid, _ := syscall.Getpgid(0)
	syscall.Kill(pgid, syscall.SIGINT)
}
