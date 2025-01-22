// +build freebsd linux netbsd openbsd solaris dragonfly

package lib

import "github.com/atotto/clipboard"

func setprimary() {
  clipboard.Primary = true
}
