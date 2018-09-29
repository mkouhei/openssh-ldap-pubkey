package main

import (
	"testing"
)

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
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, "", "", ""}
	if _, err := l.connect(); err != nil {
		t.Fatal("Connect error")
	}
}

func TestConnectFail(t *testing.T) {
	l := &ldapEnv{"localhost", 9999, "dc=example,dc=org", defaultFilter, false, false, false, "", "", ""}
	if _, err := l.connect(); err == nil {
		t.Fatal("expecting fail to error.")
	}
}

func TestConnectTLS(t *testing.T) {
	l := &ldapEnv{"localhost", 636, "dc=example,dc=org", defaultFilter, true, true, false, "", "", ""}
	if _, err := l.connectTLS(); err != nil {
		t.Fatal("Connect error")
	}
}

func TestBind(t *testing.T) {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, "", "cn=admin,dc=example,dc=org", "password"}
	c, _ := l.connect()
	if err := simpleBind(c, l); err != nil {
		t.Fatal("Bind error")
	}
}

func TestBindFail(t *testing.T) {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, "", "cn=admin,dc=example,dc=org", ""}
	c, _ := l.connect()
	if err := simpleBind(c, l); err == nil {
		t.Fatal("Bind error")
	}
}

func Example_printPubkey() {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, "user0", "", ""}
	c, _ := l.connect()
	simpleBind(c, l)
	entries, _ := l.search(c)
	printPubkey(entries)
	// Output:
	// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8OldtAiW9lQ0/2VJcc9UpRW9nfcusGXEu2sS+p5kh05zTYWGd8xHgZD0vfoQfpTfSKuHsL6qlMyKQMfsULWQoMJmMhJZc2hU1LH4u9HXYwJxD7EFleGTfxgYw6F6+LWHPVTTyhq+oMgXp/qfE4lc5A0xd2En9Qc172naHD+cRHZxhfNNYEGhW7E6eYm02Gn4fBN8hSpuZzv3WlpRgFiAWGv9CqObdQUEFFnpYLnC2kmaHqz8lzkZ9c3jdJMn2zPYyDAqQ52GI8EKyX9SrbepGJUaa/DmGyEg8nIBu4U74Sigfcl6dsJmA2qlOqSxia21mnQEFiSARB74pakgiywFV user0@workstation
	// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCrMQOAP3o58yl96HjEsheDAO/qgQ/mLVJK7DW+VFbJ9dGJpJfB4CBXPoT9bfSn4y6dotqjBA1eDbpDyzrhLkIe1MWZrRjkFbzAtB54ydKSU48URsb+XtGnN6kKKpipolQRvr3CRV7Yu2ELJDq+9Oz1gILK4nc1W/iaORVO/tZRPA0vdQwP0qkUf//neUmXXbSxOSm+ekQvZI9KfJ2tWxe+mVSFt+PcC2P4A/bW9dCNplqZdFTMQxLYFpl5ZOz3fwWcy34Shcb5nSZbjpKZdNrpuUCLwq2FMxorupko8kf4RmvMYO3G6p6OqpoIt6raB8DDJ+v/f6jdgPA31HK0sejX user0@vm01
}

func Example_printPubkeyTLS() {
	l := &ldapEnv{"localhost", 636, "dc=example,dc=org", defaultFilter, true, true, false, "user0", "", ""}
	c, _ := l.connectTLS()
	simpleBind(c, l)
	entries, _ := l.search(c)
	printPubkey(entries)
	// Output:
	// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8OldtAiW9lQ0/2VJcc9UpRW9nfcusGXEu2sS+p5kh05zTYWGd8xHgZD0vfoQfpTfSKuHsL6qlMyKQMfsULWQoMJmMhJZc2hU1LH4u9HXYwJxD7EFleGTfxgYw6F6+LWHPVTTyhq+oMgXp/qfE4lc5A0xd2En9Qc172naHD+cRHZxhfNNYEGhW7E6eYm02Gn4fBN8hSpuZzv3WlpRgFiAWGv9CqObdQUEFFnpYLnC2kmaHqz8lzkZ9c3jdJMn2zPYyDAqQ52GI8EKyX9SrbepGJUaa/DmGyEg8nIBu4U74Sigfcl6dsJmA2qlOqSxia21mnQEFiSARB74pakgiywFV user0@workstation
	// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCrMQOAP3o58yl96HjEsheDAO/qgQ/mLVJK7DW+VFbJ9dGJpJfB4CBXPoT9bfSn4y6dotqjBA1eDbpDyzrhLkIe1MWZrRjkFbzAtB54ydKSU48URsb+XtGnN6kKKpipolQRvr3CRV7Yu2ELJDq+9Oz1gILK4nc1W/iaORVO/tZRPA0vdQwP0qkUf//neUmXXbSxOSm+ekQvZI9KfJ2tWxe+mVSFt+PcC2P4A/bW9dCNplqZdFTMQxLYFpl5ZOz3fwWcy34Shcb5nSZbjpKZdNrpuUCLwq2FMxorupko8kf4RmvMYO3G6p6OqpoIt6raB8DDJ+v/f6jdgPA31HK0sejX user0@vm01
}

func Example_printPubkeyDoesNotUseSSHPublicKey() {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, "user2", "", ""}
	c, _ := l.connect()
	simpleBind(c, l)
	entries, _ := l.search(c)
	printPubkey(entries)
	// Output:
	//
}

func Example_printPubkeyDoesNotExistUser() {
	l := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, "user5", "", ""}
	c, _ := l.connect()
	simpleBind(c, l)
	entries, _ := l.search(c)
	printPubkey(entries)
	// Output:
	//
}

func TestGetHost(t *testing.T) {
	requestHost := "example.net"
	l := &ldapEnv{requestHost, 389, "dc=example,dc=org", defaultFilter, false, false, false, "user5", "", ""}
	if host, err := l.getHost(); host != requestHost && err != nil {
		t.Fatalf("expected: %s, but got: %s", requestHost, host)
	}

	requestHost = "localhost"
	l = &ldapEnv{requestHost, 389, "dc=example,dc=org", defaultFilter, false, false, false, "user1", "", ""}
	if host, err := l.getHost(); host != requestHost && err != nil {
		t.Fatalf("expected: %s, but got: %s", requestHost, host)
	}

	requestHost = "2001:db8::192:0:2:100"
	l = &ldapEnv{requestHost, 389, "dc=example,dc=org", defaultFilter, false, false, false, "user1", "", ""}
	if host, err := l.getHost(); host != "[%s]" && err != nil {
		t.Fatalf("expected: [%s], but got: %s", requestHost, host)
	}

	requestHost = "2001:db8::192:0:2:100%enp0s0"
	l = &ldapEnv{requestHost, 389, "dc=example,dc=org", defaultFilter, false, false, false, "user1", "", ""}
	if host, err := l.getHost(); host != "" && err == nil {
		t.Fatalf("expected error: 'invalid host / IP address', but got: %s", host)
	}

	requestHost = "fe80::192:0:2:100%enp0s0"
	l = &ldapEnv{requestHost, 389, "dc=example,dc=org", defaultFilter, false, false, false, "user1", "", ""}
	if host, err := l.getHost(); host != "[%s]" && err != nil {
		t.Fatalf("expected: [%s], but got: %s", requestHost, host)
	}

	requestHost = "[2001:db8::192:0:2:100]"
	l = &ldapEnv{requestHost, 389, "dc=example,dc=org", defaultFilter, false, false, false, "user1", "", ""}
	if host, err := l.getHost(); host != "" && err == nil {
		t.Fatalf("expected error: 'invalid host / IP address', but got: %s", host)
	}

	requestHost = "192.0.2.100"
	l = &ldapEnv{requestHost, 389, "dc=example,dc=org", defaultFilter, false, false, false, "user1", "", ""}
	if host, err := l.getHost(); host != requestHost || err != nil {
		t.Fatalf("expected: %s, but got: %s", requestHost, host)
	}

	requestHost = "192.168.1"
	l = &ldapEnv{requestHost, 389, "dc=example,dc=org", defaultFilter, false, false, false, "user1", "", ""}
	if host, err := l.getHost(); host != "" && err == nil {
		t.Fatalf("expected error: 'invalid host / IP address', but got: %s", host)
	}

	requestHost = "invalid:host"
	l = &ldapEnv{requestHost, 389, "dc=example,dc=org", defaultFilter, false, false, false, "user1", "", ""}
	if host, err := l.getHost(); host != "" && err == nil {
		t.Fatalf("expected error: 'invalid host / IP address', but got: %s", host)
	}

	requestHost = "invalid_host"
	l = &ldapEnv{requestHost, 389, "dc=example,dc=org", defaultFilter, false, false, false, "user1", "", ""}
	if host, err := l.getHost(); host != "" && err == nil {
		t.Fatalf("expected error: 'invalid host / IP address', but got: %s", host)
	}

	requestHost = "invalid host"
	l = &ldapEnv{requestHost, 389, "dc=example,dc=org", defaultFilter, false, false, false, "user1", "", ""}
	if host, err := l.getHost(); host != "" && err == nil {
		t.Fatalf("expected error: 'invalid host / IP address', but got: %s", host)
	}
}
