=====================
 openssh-ldap-pubkey
=====================


Status
======

.. image:: https://travis-ci.org/mkouhei/openssh-ldap-pubkey.svg?branch=master
   :target: https://travis-ci.org/mkouhei/openssh-ldap-pubkey
.. image:: https://coveralls.io/repos/mkouhei/openssh-ldap-pubkey/badge.svg?branch=master&service=github
   :target: https://coveralls.io/github/mkouhei/openssh-ldap-pubkey?branch=master
.. image:: https://readthedocs.org/projects/openssh-ldap-pubkey/badge/?version=latest
   :target: https://openssh-ldap-pubkey.readthedocs.org/en/latest/?badge=latest
   :alt: Documentation Status


Requirements
============

LDAP server
-----------

* Add `openssh-lpk schema <https://storage.googleapis.com/google-code-archive-downloads/v2/code.google.com/openssh-lpk/openssh-lpk_openldap.schema>`_.
* Add an objectClass ldapPublicKey to user entry.
* Add one or more sshPublicKey attribute to user entry.

OpenSSH server
--------------

* OpenSSH over 6.2.
* Installing this utility.
* Setup ``AuthorozedKeysCommand`` and ``AuthorizedKeysCommandUser`` in ``sshd_config``.

See also
========

* `OpenSSH 6.2 release <https://www.openssh.com/txt/release-6.2>`_
* `openssh-lpk <https://code.google.com/p/openssh-lpk/wiki/Main>`_

