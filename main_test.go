package main

import (
	"fmt"
	"os"
	"testing"
)

func TestArgparse(t *testing.T) {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, ""}
	f := "(&(objectClass=posixAccount)(uid=%s)(description=limited))"
	lc := &ldapEnv{"ldap.example.org", 9999, "ou=People,dc=example,dc=org", f, true, false, "user0"}
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

func TestArgparseTLS(t *testing.T) {
	l := &ldapEnv{"localhost", 636, "dc=example,dc=org", defaultFilter, true, false, ""}
	f := "(&(objectClass=posixAccount)(uid=%s)(description=limited))"
	lc := &ldapEnv{"ldap.example.org", 9999, "ou=People,dc=example,dc=org", f, true, true, "user0"}
	os.Args = append(os.Args, "-host=ldap.example.org")
	os.Args = append(os.Args, "-port=9999")
	os.Args = append(os.Args, "-base=ou=People,dc=example,dc=org")
	os.Args = append(os.Args, fmt.Sprintf("-filter=%s", f))
	os.Args = append(os.Args, "-tls=true")
	os.Args = append(os.Args, "-skip=true")
	os.Args = append(os.Args, "user0")
	l.argparse(os.Args[4:])
	if *l != *lc {
		t.Fatalf("expecting: %v,but %v", lc, l)
	}
	os.Args = os.Args[:5]
}

func TestArgparseNoOptions(t *testing.T) {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, ""}
	lc := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, "user1"}
	os.Args = append(os.Args, "user1")
	l.argparse(os.Args[4:])
	if *l != *lc {
		t.Fatalf("expecting: %v,but %v", lc, l)
	}
	os.Args = os.Args[:5]
}

func TestArgparseNoArg(t *testing.T) {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, ""}
	if err := l.argparse(os.Args[4:]); err == nil {
		t.Fatal("expecting: error without user argument.")
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
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, ""}
	if _, err := l.connect(); err != nil {
		t.Fatal("Connect error")
	}
}

func TestConnectFail(t *testing.T) {
	l := &ldapEnv{"localhost", 9999, "dc=example,dc=org", defaultFilter, false, false, ""}
	if _, err := l.connect(); err == nil {
		t.Fatal("expecting fail to error.")
	}
}

func TestConnectTLS(t *testing.T) {
	l := &ldapEnv{"localhost", 636, "dc=example,dc=org", defaultFilter, true, true, ""}
	if _, err := l.connectTLS(); err != nil {
		t.Fatal("Connect error")
	}
}

func Example_PrintPubkey() {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, "user0"}
	c, _ := l.connect()
	simpleBind(c)
	entries, _ := l.search(c)
	printPubkey(entries)
	// Output:
	// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8OldtAiW9lQ0/2VJcc9UpRW9nfcusGXEu2sS+p5kh05zTYWGd8xHgZD0vfoQfpTfSKuHsL6qlMyKQMfsULWQoMJmMhJZc2hU1LH4u9HXYwJxD7EFleGTfxgYw6F6+LWHPVTTyhq+oMgXp/qfE4lc5A0xd2En9Qc172naHD+cRHZxhfNNYEGhW7E6eYm02Gn4fBN8hSpuZzv3WlpRgFiAWGv9CqObdQUEFFnpYLnC2kmaHqz8lzkZ9c3jdJMn2zPYyDAqQ52GI8EKyX9SrbepGJUaa/DmGyEg8nIBu4U74Sigfcl6dsJmA2qlOqSxia21mnQEFiSARB74pakgiywFV user0@workstation
	// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCrMQOAP3o58yl96HjEsheDAO/qgQ/mLVJK7DW+VFbJ9dGJpJfB4CBXPoT9bfSn4y6dotqjBA1eDbpDyzrhLkIe1MWZrRjkFbzAtB54ydKSU48URsb+XtGnN6kKKpipolQRvr3CRV7Yu2ELJDq+9Oz1gILK4nc1W/iaORVO/tZRPA0vdQwP0qkUf//neUmXXbSxOSm+ekQvZI9KfJ2tWxe+mVSFt+PcC2P4A/bW9dCNplqZdFTMQxLYFpl5ZOz3fwWcy34Shcb5nSZbjpKZdNrpuUCLwq2FMxorupko8kf4RmvMYO3G6p6OqpoIt6raB8DDJ+v/f6jdgPA31HK0sejX user0@vm01
}

func Example_PrintPubkeyTLS() {
	l := &ldapEnv{"localhost", 636, "dc=example,dc=org", defaultFilter, true, true, "user0"}
	c, _ := l.connectTLS()
	simpleBind(c)
	entries, _ := l.search(c)
	printPubkey(entries)
	// Output:
	// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8OldtAiW9lQ0/2VJcc9UpRW9nfcusGXEu2sS+p5kh05zTYWGd8xHgZD0vfoQfpTfSKuHsL6qlMyKQMfsULWQoMJmMhJZc2hU1LH4u9HXYwJxD7EFleGTfxgYw6F6+LWHPVTTyhq+oMgXp/qfE4lc5A0xd2En9Qc172naHD+cRHZxhfNNYEGhW7E6eYm02Gn4fBN8hSpuZzv3WlpRgFiAWGv9CqObdQUEFFnpYLnC2kmaHqz8lzkZ9c3jdJMn2zPYyDAqQ52GI8EKyX9SrbepGJUaa/DmGyEg8nIBu4U74Sigfcl6dsJmA2qlOqSxia21mnQEFiSARB74pakgiywFV user0@workstation
	// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCrMQOAP3o58yl96HjEsheDAO/qgQ/mLVJK7DW+VFbJ9dGJpJfB4CBXPoT9bfSn4y6dotqjBA1eDbpDyzrhLkIe1MWZrRjkFbzAtB54ydKSU48URsb+XtGnN6kKKpipolQRvr3CRV7Yu2ELJDq+9Oz1gILK4nc1W/iaORVO/tZRPA0vdQwP0qkUf//neUmXXbSxOSm+ekQvZI9KfJ2tWxe+mVSFt+PcC2P4A/bW9dCNplqZdFTMQxLYFpl5ZOz3fwWcy34Shcb5nSZbjpKZdNrpuUCLwq2FMxorupko8kf4RmvMYO3G6p6OqpoIt6raB8DDJ+v/f6jdgPA31HK0sejX user0@vm01
}

func Example_PrintPubkeyDoesNotUseSSHPublicKey() {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, "user2"}
	c, _ := l.connect()
	simpleBind(c)
	entries, _ := l.search(c)
	printPubkey(entries)
	// Output:
	//
}

func Example_PrintPubkeyDoesNotExistUser() {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, "user5"}
	c, _ := l.connect()
	simpleBind(c)
	entries, _ := l.search(c)
	printPubkey(entries)
	// Output:
	//
}
