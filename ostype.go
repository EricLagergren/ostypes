package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// globs are stolen from GNU coreutils' host-os.m4 file.
// see GNU-LICENSE.txt
var globs = map[string]string{
	"winnt*":        "Windows NT",
	"vos*":          "VOS",
	"sysv*":         "Unix System V",
	"superux*":      "SUPER-UX",
	"sunos*":        "SunOS",
	"stop*":         "STOP",
	"sco*":          "SCO Unix",
	"riscos*":       "RISC OS",
	"riscix*":       "RISCiX",
	"qnx*":          "QNX",
	"pw32*":         "PW32",
	"ptx*":          "ptx",
	"plan9*":        "Plan 9",
	"osf*":          "Tru64",
	"os2*":          "OS/2",
	"openbsd*":      "OpenBSD",
	"nsk*":          "NonStop Kernel",
	"nonstopux*":    "NonStop-UX",
	"netbsd*-gnu*":  "GNU/NetBSD",
	"netbsd*":       "NetBSD",
	"mirbsd*":       "MirBSD",
	"knetbsd*-gnu":  "GNU/kNetBSD",
	"kfreebsd*-gnu": "GNU/kFreeBSD",
	"msdosdjgpp*":   "DJGPP",
	"mpeix*":        "MPE/iX",
	"mint*":         "MiNT",
	"mingw*":        "MinGW",
	"lynxos*":       "LynxOS",
	"linux*":        "GNU/Linux",
	"hpux*":         "HP-UX",
	"hiux*":         "HI-UX",
	"gnu*":          "GNU",
	"freebsd*":      "FreeBSD",
	"dgux*":         "DG/UX",
	"bsdi*":         "BSD/OS",
	"bsd*":          "BSD",
	"beos*":         "BeOS",
	"aux*":          "A/UX",
	"atheos*":       "AtheOS",
	"amigaos*":      "Amiga OS",
	"aix*":          "AIX",
}

func match(pattern, name string) bool {
	ok, err := filepath.Match(pattern, name)
	return ok && err == nil
}

func main() {

	// Shell locals...
	cmd := exec.Command("sh", "-c", "echo $OSTYPE")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	if out[len(out)-1] == '\n' {
		out = out[:len(out)-1]
	}

	file, err := os.Create("host_os.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	rawHostOS := string(out)

	file.WriteString("package main\nconst HostOS = ")

	// This is what uname does.
	if rawHostOS == "" {
		file.WriteString("\"Unknown\"")
		return
	}

	var realHostOS string
	for glob, val := range globs {
		if match(glob, rawHostOS) {
			realHostOS = val
			break
		}
	}

	// No matches found, so apply the default heuristic.
	if realHostOS == "" {
		if match("[A-Za-z]*", rawHostOS) {
			upper := strings.ToUpper(string(rawHostOS[0]))

			if len(rawHostOS) > 1 {
				realHostOS = upper + rawHostOS[1:]
			} else {
				realHostOS = upper
			}

		} else {
			realHostOS = rawHostOS
		}
	}

	file.WriteString(fmt.Sprintf("%q", realHostOS))
}
