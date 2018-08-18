=============================
 How to setup OpenSSH server
=============================

Precondition
============

This article restricts OpenSSH 6.2 over on Debian systems only.

.. note::
   You can use `openssh-ldap package <https://apps.fedoraproject.org/packages/openssh-ldap>`_ instead of this utility in the distribution based RHEL.

Requirements
============

* Debian Jessie later or Ubuntu Trusty later.
* OpenSSH 6.2 over
* openssh-ldap-pubkey
* Go 1.2 over

Optional
--------

* nslcd

Install with nslcd (recommend)
==============================

When the following precondition is sufficient,
``openssh-ldap-pubkey`` can loads parameters from ``/etc/nslcd.conf``.

* ``nslcd`` package is installed.
* There is ``/etc/nslcd.conf``.
* Set ``root`` to ``AuthorizedKeysCommandUser`` of ``/etc/ssh/sshd_config``.

The parameters are follows.

.. list-table:: nslcd.conf keys compare openssh-ldap-pubkey options.
   :header-rows: 1

   * - nslcd.conf
     - openssh-ldap-pubkey
   * - | ``uri``
       | ldap://example.org
       | ldaps://example.org
     - | ``host``, ``port``, ``tls``
       | example.org, 389, false
       | example.org, 636, true
   * - | ``base``
       | dc=example,dc=org
     - | ``base``
       | dc=example,dc=org
   * - | ``pam_authz_search``
       | (&(objectClass=posixAccount)(uid=$username))
     - | ``filter``
       | (&(objectClass=posixAccount)(uid=%s))
   * - | ``tls_reqcert``
       | never, allow
       | try, demand, hard
     - | ``skip``
       | true
       | false
   * - | ``binddn`` (option for bind)
       | cn=admin,dc=example,dc=org
     - | n/a
   * - | ``bindpw`` (option for bind)
       | examplepassword
     - | n/a

1. Download binary.

   .. code-block:: shell

      $ export GOPATH=/path/to/gocode
      $ go get github.com/mkouhei/openssh-ldap-pubkey
      $ chmod 0755 /path/to/gocode/bin/openssh-ldap-pubkey
      $ sudo chown root: /path/to/gocode/bin/openssh-ldap-pubkey

2. Setup sshd_config.

   Appends ``AuthorizedKeysCommand`` and ``AuthorizedKeysCommandUser``.

   .. code-block:: ini

      AuthorizedKeysCommand /path/to/openssh-ldap-pubkey
      AuthorizedKeysCommandUser root

3. Restart sshd.

   .. code-block:: shell

      $ sudo service ssh restart


Install without nslcd
=====================

If ``nslcd`` is not installed and there is not ``/etc/nslcd.conf``,
you should prepare wrapper script of ``openssh-ldap-pubkey``.

1. Download binary.

   .. code-block:: shell

      $ export GOPATH=/path/to/gocode
      $ go get github.com/mkouhei/openssh-ldap-pubkey
      $ chmod 0755 /path/to/gocode/bin/openssh-ldap-pubkey
      $ sudo chown root: /path/to/gocode/bin/openssh-ldap-pubkey


2. Prepare wrapper script.

   without TLS,

   .. code-block:: shell

      $ sudo bash -c "cat << EOF > /etc/ssh/openssh-ldap-pubkey.sh
      #!/bin/sh -e
      /path/to/openssh-ldap-pubkey -host=ldap.example.org -base=dc=example,dc=org $1
      EOF
      $ sudo chmod +x /etc/ssh/openssh-ldap-pubkey.sh

   with TLS.

   .. code-block:: shell

      $ sudo bash -c "cat << EOF > /etc/ssh/openssh-ldap-pubkey.sh
      #!/bin/sh -e
      /path/to/openssh-ldap-pubkey -host=ldap.example.org -port 636 -base=dc=example,dc=org -tls=true $1
      EOF
      $ sudo chmod +x /etc/ssh/openssh-ldap-pubkey.sh

3. Setup sshd_config.

   Appends ``AuthorizedKeysCommand`` and ``AuthorizedKeysCommandUser``.

   .. code-block:: ini

      AuthorizedKeysCommand /etc/ssh/openssh-ldap-pubkey.sh
      AuthorizedKeysCommandUser root

4. Restart sshd.

   .. code-block:: shell

      $ sudo service ssh restart

