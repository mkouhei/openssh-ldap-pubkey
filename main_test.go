package main

import (
	"fmt"
	"os"
	"testing"
)

func TestArgparse(t *testing.T) {
	f := "(&(objectClass=posixAccount)(uid=%s)(description=limited))"
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, "", "", ""}
	lc := &ldapEnv{"ldap.example.org", 9999, "ou=People,dc=example,dc=org", f, false, false, false, "user0", "", ""}
	os.Args = []string{"test_command"}
	os.Args = append(os.Args, "-host=ldap.example.org")
	os.Args = append(os.Args, "-port=9999")
	os.Args = append(os.Args, "-base=ou=People,dc=example,dc=org")
	os.Args = append(os.Args, fmt.Sprintf("-filter=%s", f))
	os.Args = append(os.Args, "-tls=false")
	os.Args = append(os.Args, "user0")
	l.argparse(os.Args, version)
	if *l != *lc {
		t.Fatalf("expecting:\n%v,\nbut:\n%v", lc, l)
	}
}

func TestArgparseTLS(t *testing.T) {
	f := "(&(objectClass=posixAccount)(uid=%s)(description=limited))"
	l := &ldapEnv{"localhost", 636, "dc=example,dc=org", defaultFilter, false, false, false, "", "", ""}
	lc := &ldapEnv{"ldap.example.org", 9999, "ou=People,dc=example,dc=org", f, true, false, false, "user0", "", ""}
	os.Args = []string{"test_command"}
	os.Args = append(os.Args, "-host=ldap.example.org")
	os.Args = append(os.Args, "-port=9999")
	os.Args = append(os.Args, "-base=ou=People,dc=example,dc=org")
	os.Args = append(os.Args, fmt.Sprintf("-filter=%s", f))
	os.Args = append(os.Args, "-tls=true")
	os.Args = append(os.Args, "-skip=false")
	os.Args = append(os.Args, "user0")
	l.argparse(os.Args, version)
	if *l != *lc {
		t.Fatalf("expecting:\n%v,\nbut:\n%v", lc, l)
	}
}

func TestArgparseNoOptions(t *testing.T) {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, "", "", ""}
	lc := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, "user1", "", ""}
	os.Args = []string{"test_command"}
	os.Args = append(os.Args, "user1")
	l.argparse(os.Args, version)
	if *l != *lc {
		t.Fatalf("expecting:\n%v,\nbut:\n%v", lc, l)
	}
}

func TestArgparseNoArg(t *testing.T) {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, "", "", ""}
	os.Args = []string{"test_command"}
	if err := l.argparse(os.Args, version); err == nil {
		t.Fatal("expecting: error without user argument.")
	}
}

func ShowVersion() {
	l := &ldapEnv{}
	os.Args = append(os.Args, "-version")
	ver = "X.X.X"
	l.argparse(os.Args[4:], ver)
	// Output:
	// openssh-ldap-pubkey X.X.X
	//
	// Copyright (C) 2015-2018 Kouhei Maeda
	// License GPLv3+: GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>.
	// This is free software, and you are welcome to redistribute it.
	// There is NO WARRANTY, to the extent permitted by law.
}
