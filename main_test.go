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

func TestIsAddrWithV4(t *testing.T) {
	addr := "192.0.2.100"
	if !isAddr(addr) {
		t.Fatalf("expecting: %s is true, but false", addr)
	}
}

func TestIsAddrWithV6(t *testing.T) {
	addr := "2001:db8::100"
	if !isAddr(addr) {
		t.Fatalf("expecting: %s is true, but false", addr)
	}
}

func TestIsAddrWithFQDN(t *testing.T) {
	host := "ldap.example.org"
	if isAddr(host) {
		t.Fatalf("expecting: %s is false, but true", host)
	}
}

func TestConnect(t *testing.T) {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, ""}
	c := l.connect()
	if c == nil {
		t.Fatal("Connect error")
	}
}

func Example_PrintPubkey() {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, ""}
	c := l.connect()
	simpleBind(c)
	printPubkey(l.search(c))
	// Output:
	// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8OldtAiW9lQ0/2VJcc9UpRW9nfcusGXEu2sS+p5kh05zTYWGd8xHgZD0vfoQfpTfSKuHsL6qlMyKQMfsULWQoMJmMhJZc2hU1LH4u9HXYwJxD7EFleGTfxgYw6F6+LWHPVTTyhq+oMgXp/qfE4lc5A0xd2En9Qc172naHD+cRHZxhfNNYEGhW7E6eYm02Gn4fBN8hSpuZzv3WlpRgFiAWGv9CqObdQUEFFnpYLnC2kmaHqz8lzkZ9c3jdJMn2zPYyDAqQ52GI8EKyX9SrbepGJUaa/DmGyEg8nIBu4U74Sigfcl6dsJmA2qlOqSxia21mnQEFiSARB74pakgiywFV user0@workstation
	// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCrMQOAP3o58yl96HjEsheDAO/qgQ/mLVJK7DW+VFbJ9dGJpJfB4CBXPoT9bfSn4y6dotqjBA1eDbpDyzrhLkIe1MWZrRjkFbzAtB54ydKSU48URsb+XtGnN6kKKpipolQRvr3CRV7Yu2ELJDq+9Oz1gILK4nc1W/iaORVO/tZRPA0vdQwP0qkUf//neUmXXbSxOSm+ekQvZI9KfJ2tWxe+mVSFt+PcC2P4A/bW9dCNplqZdFTMQxLYFpl5ZOz3fwWcy34Shcb5nSZbjpKZdNrpuUCLwq2FMxorupko8kf4RmvMYO3G6p6OqpoIt6raB8DDJ+v/f6jdgPA31HK0sejX user0@vm01
}
