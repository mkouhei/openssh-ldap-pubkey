package main

import (
	"os"
	"testing"

	"path/filepath"
)

func TestGetNslcdConfPath(t *testing.T) {
	if conf := getNslcdConfPath(); conf != nslcdConf {
		t.Fatalf("expecting: %s, but: %s", nslcdConf, conf)
	}
	os.Setenv("NSLCD_CONF", "/path/to/nslcd.conf")
	if conf := getNslcdConfPath(); conf != "/path/to/nslcd.conf" {
		t.Fatalf("expecting: /path/to/nslcd.conf, but: %s", conf)
	}
}

func TestLoadNslcdConf(t *testing.T) {
	conf, err := filepath.Abs("testdata/nslcd.conf")
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("NSLCD_CONF", conf)
	l := &ldapEnv{}
	l.loadNslcdConf()
	lc := &ldapEnv{"ldap.example.org", 389, "dc=example,dc=org", defaultFilter, false, false, false, ""}
	if *lc != *l {
		t.Fatal("Failed to load testdata/nslcd.conf via NSLCD_CONF env var.")
	}
}

func TestFailLoadNslcdConf(t *testing.T) {
	conf, err := filepath.Abs("testdata/nslcd-noexist.conf")
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("NSLCD_CONF", conf)
	l := &ldapEnv{}
	l.loadNslcdConf()
	lc := &ldapEnv{"localhost", 389, "dc=example,dc=org", defaultFilter, false, false, false, ""}
	if *lc != *l {
		t.Fatal("Failed to load default configuration.")
	}
}

func TestLoadNslcdConfWithTLS(t *testing.T) {
	conf, err := filepath.Abs("testdata/nslcd-tls.conf")
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("NSLCD_CONF", conf)
	l := &ldapEnv{}
	l.loadNslcdConf()
	lc := &ldapEnv{"ldap.example.org", 686, "ou=People,dc=example,dc=org", "(&(objectClass=posixAccount)(uid=%s)(description=limited))", true, false, false, ""}
	if *lc != *l {
		t.Fatal("Failed to load default configuration.")
	}
}

func TestLoadNslcdConfWithTLSAllow(t *testing.T) {
	conf, err := filepath.Abs("testdata/nslcd-tls-allow.conf")
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("NSLCD_CONF", conf)
	l := &ldapEnv{}
	l.loadNslcdConf()
	lc := &ldapEnv{"ldap.example.org", 686, "ou=People,dc=example,dc=org", "(&(objectClass=posixAccount)(uid=%s)(description=limited))", true, true, false, ""}
	if *lc != *l {
		t.Fatal("Failed to load default configuration.")
	}
}

func TestLoadNslcdConfWithTLSNever(t *testing.T) {
	conf, err := filepath.Abs("testdata/nslcd-tls-never.conf")
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("NSLCD_CONF", conf)
	l := &ldapEnv{}
	l.loadNslcdConf()
	lc := &ldapEnv{"ldap.example.org", 686, "ou=People,dc=example,dc=org", "(&(objectClass=posixAccount)(uid=%s)(description=limited))", true, true, false, ""}
	if *lc != *l {
		t.Fatal("Failed to load default configuration.")
	}
}

func TestLoadNslcdConfWithTLSSkip(t *testing.T) {
	conf, err := filepath.Abs("testdata/nslcd-tls-skip.conf")
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("NSLCD_CONF", conf)
	l := &ldapEnv{}
	l.loadNslcdConf()
	lc := &ldapEnv{"192.0.2.100", 686, "ou=People,dc=example,dc=org", "(&(objectClass=posixAccount)(uid=%s)(description=limited))", true, true, false, ""}
	if *lc != *l {
		t.Fatal("Failed to load default configuration.")
	}
}

func TestLoadNslcdConfWithPort(t *testing.T) {
	conf, err := filepath.Abs("testdata/nslcd-tls-port.conf")
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("NSLCD_CONF", conf)
	l := &ldapEnv{}
	l.loadNslcdConf()
	lc := &ldapEnv{"example.org", 686, "ou=People,dc=example,dc=org", "(&(objectClass=posixAccount)(uid=%s)(description=limited))", true, false, false, ""}
	if *lc != *l {
		t.Fatal("Failed to load default configuration.")
	}
}

func TestLoadNslcdConfNoFilter(t *testing.T) {
	conf, err := filepath.Abs("testdata/nslcd-no-filter.conf")
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("NSLCD_CONF", conf)
	l := &ldapEnv{}
	l.loadNslcdConf()
	lc := &ldapEnv{"ldap.example.org", 389, "dc=example,dc=org", defaultFilter, false, false, false, ""}
	if *lc != *l {
		t.Fatal("Failed to load default configuration.")
	}
}

func TestLoadNslcdConfNoUsername(t *testing.T) {
	conf, err := filepath.Abs("testdata/nslcd-no-username.conf")
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("NSLCD_CONF", conf)
	l := &ldapEnv{}
	l.loadNslcdConf()
	lc := &ldapEnv{"ldap.example.org", 389, "dc=example,dc=org", "(objectClass=posixAccount)", false, false, false, ""}
	if *lc != *l {
		t.Fatal("Failed to load default configuration.")
	}
}

func TestLoadNslcdConfInvalidURL(t *testing.T) {
	conf, err := filepath.Abs("testdata/nslcd-invalid-url.conf")
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("NSLCD_CONF", conf)
	l := &ldapEnv{}
	l.loadNslcdConf()
	lc := &ldapEnv{"", 389, "dc=example,dc=org", defaultFilter, false, false, false, ""}
	if *lc == *l {
		t.Fatal("Failed to parse url.")
	}
}

func TestLoadNslcdConfInvalidPort(t *testing.T) {
	conf, err := filepath.Abs("testdata/nslcd-invalid-port.conf")
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("NSLCD_CONF", conf)
	l := &ldapEnv{}
	l.loadNslcdConf()
	lc := &ldapEnv{"", 389, "dc=example,dc=org", defaultFilter, false, false, false, ""}
	if *lc == *l {
		t.Fatal("Failed to validate port syntax.")
	}
}
