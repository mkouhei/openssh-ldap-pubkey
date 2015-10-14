package main

import (
	"fmt"
	"os"
	"testing"
)

func TestArgparse(t *testing.T) {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, ""}
	f := "(&(objectClass=posixAccount)(uid=%s)(description=limited))"
	lc := &ldapEnv{"ldap.example.org", 9999, "ou=People,dc=example,dc=org", f, true, "user0"}
	os.Args = append(os.Args, "-host=ldap.example.org")
	os.Args = append(os.Args, "-port=9999")
	os.Args = append(os.Args, "-base=ou=People,dc=example,dc=org")
	os.Args = append(os.Args, fmt.Sprintf("-filter=%s", f))
	os.Args = append(os.Args, "-tls=true")
	os.Args = append(os.Args, "user0")
	l.argparse(os.Args[4:])
	if *l != *lc {
		t.Fatalf("expecting: %v,but %v", lc, l)
	}
	os.Args = os.Args[:5]
}

func TestArgparseNoOptions(t *testing.T) {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, ""}
	lc := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, "user1"}
	os.Args = append(os.Args, "user1")
	l.argparse(os.Args[4:])
	if *l != *lc {
		t.Fatalf("expecting: %v,but %v", lc, l)
	}
}
