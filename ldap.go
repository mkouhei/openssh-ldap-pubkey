package main

import (
	"errors"
	"fmt"
	"net"
	"strings"

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
	host, err := l.getHost()

	if err != nil {
		logging(err)
		return nil, err
	}
	return ldap.Dial("tcp", fmt.Sprintf("%s:%d", host, l.port))
}

func (l *ldapEnv) connectTLS() (*ldap.Conn, error) {
	host, err := l.getHost()
	if err != nil {
		logging(err)
		return nil, err
	}

	certs := *x509.NewCertPool()
	tlsConfig := &tls.Config{
		RootCAs: &certs,
	}

	if !isAddr(host) {
		tlsConfig.ServerName = host
	}
	if isAddr(host) || l.skip {
		tlsConfig.InsecureSkipVerify = true
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
	ip := net.ParseIP(host)
	return isIPv4(ip) || isIPv6(ip)
}

func (l *ldapEnv) validateIPv6LinkLocal() (string, error) {
	addr := strings.Split(l.host, "%")

	ip := net.ParseIP(addr[0])

	if isIPv6(ip) && ip.IsLinkLocalUnicast() {
		return fmt.Sprintf("[%s]", l.host), nil
	}
	return "", errors.New("Invalid IPv6 link-local address")
}

func isIPv6(ip net.IP) bool {
	return ip.To4() == nil && ip.To16() != nil
}

func isIPv4(ip net.IP) bool {
	return ip.To4() != nil && ip.To16() != nil
}

func (l *ldapEnv) validateIPv6() (string, error) {
	if strings.Contains(l.host, "%") {
		return l.validateIPv6LinkLocal()
	}
	ip := net.ParseIP(l.host)

	// global scope
	if isIPv6(ip) {
		return fmt.Sprintf("[%s]", l.host), nil
	}
	return "", errors.New("Invalid IPv6 address")
}

func (l *ldapEnv) validateHost() (string, error) {
	var host string
	var err error
	addrs, err := net.LookupHost(l.host)
	if err == nil && len(addrs) > 0 {
		host = l.host
	} else {
		host = ""
		err = errors.New("Invalid hostname / FQDN")
	}
	return host, err
}

func (l *ldapEnv) getHost() (string, error) {
	var host string
	var err error

	if strings.HasPrefix(l.host, "[") || strings.HasSuffix(l.host, "]") {
		err = errors.New("Invalid host")
	} else {
		// IPv6
		if strings.Contains(l.host, ":") {
			host, err = l.validateIPv6()
		} else if isIPv4(net.ParseIP(l.host)) {
			// ipv4
			host = l.host
		} else {
			// fqdn
			host, err = l.validateHost()
		}
	}

	return host, err
}
