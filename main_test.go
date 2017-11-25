package main

import (
	"fmt"
	"os"
	"testing"
)

func TestArgparse(t *testing.T) {
	f := "(&(objectClass=posixAccount)(uid=%s)(description=limited))"
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, ""}
	lc := &ldapEnv{"ldap.example.org", 9999, "ou=People,dc=example,dc=org", f, false, false, false, "user0"}
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
	l := &ldapEnv{"localhost", 636, "dc=example,dc=org", defaultFilter, false, false, false, ""}
	lc := &ldapEnv{"ldap.example.org", 9999, "ou=People,dc=example,dc=org", f, true, false, false, "user0"}
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
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, ""}
	lc := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, "user1"}
	os.Args = []string{"test_command"}
	os.Args = append(os.Args, "user1")
	l.argparse(os.Args, version)
	if *l != *lc {
		t.Fatalf("expecting:\n%v,\nbut:\n%v", lc, l)
	}
}

func TestArgparseNoArg(t *testing.T) {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, ""}
	os.Args = []string{"test_command"}
	if err := l.argparse(os.Args, version); err == nil {
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
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, ""}
	if _, err := l.connect(); err != nil {
		t.Fatal("Connect error")
	}
}

func TestConnectFail(t *testing.T) {
	l := &ldapEnv{"localhost", 9999, "dc=example,dc=org", defaultFilter, false, false, false, ""}
	if _, err := l.connect(); err == nil {
		t.Fatal("expecting fail to error.")
	}
}

func TestConnectTLS(t *testing.T) {
	l := &ldapEnv{"localhost", 636, "dc=example,dc=org", defaultFilter, true, true, false, ""}
	if _, err := l.connectTLS(); err != nil {
		t.Fatal("Connect error")
	}
}

func PrintPubkey() {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, "user0"}
	c, _ := l.connect()
	simpleBind(c)
	entries, _ := l.search(c)
	printPubkey(entries)
	// Output:
	// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8OldtAiW9lQ0/2VJcc9UpRW9nfcusGXEu2sS+p5kh05zTYWGd8xHgZD0vfoQfpTfSKuHsL6qlMyKQMfsULWQoMJmMhJZc2hU1LH4u9HXYwJxD7EFleGTfxgYw6F6+LWHPVTTyhq+oMgXp/qfE4lc5A0xd2En9Qc172naHD+cRHZxhfNNYEGhW7E6eYm02Gn4fBN8hSpuZzv3WlpRgFiAWGv9CqObdQUEFFnpYLnC2kmaHqz8lzkZ9c3jdJMn2zPYyDAqQ52GI8EKyX9SrbepGJUaa/DmGyEg8nIBu4U74Sigfcl6dsJmA2qlOqSxia21mnQEFiSARB74pakgiywFV user0@workstation
	// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCrMQOAP3o58yl96HjEsheDAO/qgQ/mLVJK7DW+VFbJ9dGJpJfB4CBXPoT9bfSn4y6dotqjBA1eDbpDyzrhLkIe1MWZrRjkFbzAtB54ydKSU48URsb+XtGnN6kKKpipolQRvr3CRV7Yu2ELJDq+9Oz1gILK4nc1W/iaORVO/tZRPA0vdQwP0qkUf//neUmXXbSxOSm+ekQvZI9KfJ2tWxe+mVSFt+PcC2P4A/bW9dCNplqZdFTMQxLYFpl5ZOz3fwWcy34Shcb5nSZbjpKZdNrpuUCLwq2FMxorupko8kf4RmvMYO3G6p6OqpoIt6raB8DDJ+v/f6jdgPA31HK0sejX user0@vm01
}

func PrintPubkeyTLS() {
	l := &ldapEnv{"localhost", 636, "dc=example,dc=org", defaultFilter, true, true, false, "user0"}
	c, _ := l.connectTLS()
	simpleBind(c)
	entries, _ := l.search(c)
	printPubkey(entries)
	// Output:
	// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8OldtAiW9lQ0/2VJcc9UpRW9nfcusGXEu2sS+p5kh05zTYWGd8xHgZD0vfoQfpTfSKuHsL6qlMyKQMfsULWQoMJmMhJZc2hU1LH4u9HXYwJxD7EFleGTfxgYw6F6+LWHPVTTyhq+oMgXp/qfE4lc5A0xd2En9Qc172naHD+cRHZxhfNNYEGhW7E6eYm02Gn4fBN8hSpuZzv3WlpRgFiAWGv9CqObdQUEFFnpYLnC2kmaHqz8lzkZ9c3jdJMn2zPYyDAqQ52GI8EKyX9SrbepGJUaa/DmGyEg8nIBu4U74Sigfcl6dsJmA2qlOqSxia21mnQEFiSARB74pakgiywFV user0@workstation
	// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCrMQOAP3o58yl96HjEsheDAO/qgQ/mLVJK7DW+VFbJ9dGJpJfB4CBXPoT9bfSn4y6dotqjBA1eDbpDyzrhLkIe1MWZrRjkFbzAtB54ydKSU48URsb+XtGnN6kKKpipolQRvr3CRV7Yu2ELJDq+9Oz1gILK4nc1W/iaORVO/tZRPA0vdQwP0qkUf//neUmXXbSxOSm+ekQvZI9KfJ2tWxe+mVSFt+PcC2P4A/bW9dCNplqZdFTMQxLYFpl5ZOz3fwWcy34Shcb5nSZbjpKZdNrpuUCLwq2FMxorupko8kf4RmvMYO3G6p6OqpoIt6raB8DDJ+v/f6jdgPA31HK0sejX user0@vm01
}

func PrintPubkeyDoesNotUseSSHPublicKey() {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, "user2"}
	c, _ := l.connect()
	simpleBind(c)
	entries, _ := l.search(c)
	printPubkey(entries)
	// Output:
	//
}

func PrintPubkeyDoesNotExistUser() {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, "user5"}
	c, _ := l.connect()
	simpleBind(c)
	entries, _ := l.search(c)
	printPubkey(entries)
	// Output:
	//
}

func ShowVersion() {
	l := &ldapEnv{}
	os.Args = append(os.Args, "-version")
	ver = "X.X.X"
	l.argparse(os.Args[4:], ver)
	// Output:
	// openssh-ldap-pubkey X.X.X
	//
	// Copyright (C) 2015 Kouhei Maeda
	// License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>.
	// This is free software, and you are welcome to redistribute it.
	// There is NO WARRANTY, to the extent permitted by law.
}
