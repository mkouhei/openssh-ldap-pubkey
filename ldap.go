package main

import (
	"errors"
	"fmt"
	"net"

	"crypto/tls"
	"crypto/x509"

	"gopkg.in/ldap.v2"
)

type ldapEnv struct {
	host   string
	port   int
	base   string
	filter string
	tls    bool
	skip   bool
	debug  bool
	uid    string
	binddn string
	bindpw string
}

func (l *ldapEnv) connect() (*ldap.Conn, error) {
	host := l.host
	if !isAddr(host) {
		host = l.getAddr()
	}
	return ldap.Dial("tcp", fmt.Sprintf("%s:%d", host, l.port))
}

func (l *ldapEnv) connectTLS() (*ldap.Conn, error) {
	certs := *x509.NewCertPool()
	tlsConfig := &tls.Config{
		RootCAs: &certs,
	}

	if isAddr(l.host) || l.skip {
		tlsConfig.InsecureSkipVerify = true
	}
	host := l.host
	if !isAddr(l.host) {
		tlsConfig.ServerName = l.host
		host = l.getAddr()
	}
	return ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", host, l.port), tlsConfig)
}

func simpleBind(c *ldap.Conn, l *ldapEnv) error {
	bindRequest := ldap.NewSimpleBindRequest(l.binddn, l.bindpw, nil)
	_, err := c.SimpleBind(bindRequest)
	return err
}

func (l *ldapEnv) search(c *ldap.Conn) ([]*ldap.Entry, error) {
	searchRequest := ldap.NewSearchRequest(
		l.base, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, 0, false,
		fmt.Sprintf(l.filter, l.uid), []string{sshPublicKeyName}, nil)
	sr, err := c.Search(searchRequest)
	return sr.Entries, err
}

func printPubkey(entries []*ldap.Entry) error {
	if len(entries) != 1 {
		return errors.New("User does not exist or too many entries returned")
	}

	if len(entries[0].GetAttributeValues("sshPublicKey")) == 0 {
		return errors.New("User does not use ldapPublicKey")
	}
	for _, pubkey := range entries[0].GetAttributeValues("sshPublicKey") {
		fmt.Println(pubkey)
	}
	return nil
}

func isAddr(host string) bool {
	return !(net.ParseIP(host).To4() == nil && net.ParseIP(host).To16() == nil)
}

func (l *ldapEnv) getAddr() string {
	addrs, err := net.LookupHost(l.host)
	if err != nil {
		return l.host
	}
	if net.ParseIP(addrs[0]).To16() != nil {
		return fmt.Sprintf("[%s]", addrs[0])
	}
	return addrs[0]
}
