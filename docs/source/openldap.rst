==========================================
 How to setup LDAP server for openssh-lpk
==========================================

Precondition
============

This article restricts OpenLDAP with ``slapd_config`` on Debian systems only.

Requirements
============

* Debian Wheezy later or Ubuntu Precise later.
* OpenLDAP(slapd) 2.4.28 over.
* debconf-utils
* ldap-utils
* ldapvi
* `openssh-lpk schema <https://openssh-lpk.googlecode.com/svn/trunk/schemas/openssh-lpk_openldap.schema>`_

Install
=======

1. Prepare debconf configuration for slad. Replace each parameters for your envirionment.

   .. code-block:: shell

      $ cat << EOF > debconf.txt
      slapd	slapd/password1	password	
      slapd	slapd/internal/adminpw	password	
      slapd	slapd/internal/generated_adminpw	password	
      slapd	slapd/password2	password	
      slapd	slapd/unsafe_selfwrite_acl	note	
      slapd	slapd/allow_ldap_v2	boolean	false
      slapd	shared/organization	string	example.org
      slapd	slapd/move_old_database	boolean	true
      slapd	slapd/password_mismatch	note	
      slapd	slapd/dump_database	select	when needed
      slapd	slapd/dump_database_destdir	string	/var/backups/slapd-VERSION
      slapd	slapd/invalid_config	boolean	true
      slapd	slapd/domain	string	example.org
      slapd	slapd/backend	select	HDB
      slapd	slapd/purge_database	boolean	true
      slapd	slapd/no_configuration	boolean	false
      EOF

   .. note::
      debconf separator is ``tab``.
   
   See :download:`sample debconf configuration <./_static/debconf.txt>`.   

2. Install packages except of slapd.

   .. code-block:: shell

      $ sudo apt-get install debconf-utils ldap-utils ldapvi

3. Download openssh-lpk schema and convert to LDIF.

   .. code-block:: shell
                   
      $ curl https://openssh-lpk.googlecode.com/svn/trunk/schemas/openssh-lpk_openldap.schema | sed "
      1i\dn: cn=openssh-lpk,cn=schema,cn=config\nobjectClass: olcSchemaConfig\ncn: openssh-lpk
      /^#/d
      /^$/d
      :a
      / $/N
      / $/b a
      s/\n//g
      s/\t//g
      /octetStringMatch$/N
      s/\n/ /
      /AUXILIARY$/N
      s/\n/ /
      /objectclass'$/N
      s/\n//
      s/^attributetype (/olcAttributeTypes: {0}(/
      s/^objectclass (/olcObjectClasses: {0}(/
      :b
      / $/N
      / $/b b
      s/\n//g
      s/\t//g
      " > openssh-lpk.ldif

   See :download:`the convert script <./_static/conv.sh.txt>`, :download:`openssh-lpk schema ldif <./_static/openssh-lpk.ldif.txt>`.

4. Prepare the LDIF for changing for rootdn password.

   .. code-block:: shell

      $ cat << EOF > rootdnpw.ldif
      dn: olcDatabase={1}hdb,cn=config
      changetype: modify
      replace: olcRootPW
      olcRootPW: {SSHA}BADfSMMJo53/L/gaFiG0xqKnOsmds4fW
      EOF

   Replace the ``olcRootPW`` value by generated with ``slappasswd`` command. [#]_

   See :download:`the change rootdn password LDIF <./_static/rootdnpw.ldif.txt>`.
                   
5. Prepare the LDIF of ``organizationalUnit`` entry.

   .. code-block:: shell

      $ cat <<EOF > ou.ldif
      dn: ou=People,dc=example,dc=org
      objectClass: organizationalUnit
      ou: People
      EOF

   Replace the ``dn`` and ``ou`` value.

   See :download:`the adding ou LDIF <./_static/ou.ldif.txt>`.

6. Prepare the LDIF of user entry.

   .. code-block:: shell
                
      $ cat << EOF > users.ldfi
      dn: uid=user0,ou=People,dc=example,dc=org
      cn: user0
      objectClass: inetOrgPerson
      objectClass: posixAccount
      objectClass: shadowAccount
      objectClass: ldapPublicKey
      loginShell: /bin/bash
      uidNumber: 1000
      gidNumber: 1000
      sn: user0
      homeDirectory: /home/user0
      mail: user0@example.org
      uid: user0
      sshPublicKey: ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8OldtAiW9lQ0/2VJcc9UpRW9nfcusGXEu2sS+p5kh05zTYWGd8xHgZD0vfoQfpTfSKuHsL6qlMyKQMfsULWQoMJmMhJZc2hU1LH4u9HXYwJxD7EFleGTfxgYw6F6+LWHPVTTyhq+oMgXp/qfE4lc5A0xd2En9Qc172naHD+cRHZxhfNNYEGhW7E6eYm02Gn4fBN8hSpuZzv3WlpRgFiAWGv9CqObdQUEFFnpYLnC2kmaHqz8lzkZ9c3jdJMn2zPYyDAqQ52GI8EKyX9SrbepGJUaa/DmGyEg8nIBu4U74Sigfcl6dsJmA2qlOqSxia21mnQEFiSARB74pakgiywFV user0@workstation
      sshPublicKey: ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCrMQOAP3o58yl96HjEsheDAO/qgQ/mLVJK7DW+VFbJ9dGJpJfB4CBXPoT9bfSn4y6dotqjBA1eDbpDyzrhLkIe1MWZrRjkFbzAtB54ydKSU48URsb+XtGnN6kKKpipolQRvr3CRV7Yu2ELJDq+9Oz1gILK4nc1W/iaORVO/tZRPA0vdQwP0qkUf//neUmXXbSxOSm+ekQvZI9KfJ2tWxe+mVSFt+PcC2P4A/bW9dCNplqZdFTMQxLYFpl5ZOz3fwWcy34Shcb5nSZbjpKZdNrpuUCLwq2FMxorupko8kf4RmvMYO3G6p6OqpoIt6raB8DDJ+v/f6jdgPA31HK0sejX user0@vm01
      userPassword:{SSHA}eKfVPm3raZmYPx5Os+KGKVUPVb6P+766
     
      dn: uid=user1,ou=People,dc=example,dc=org
      (snip)
      EOF

   Replace the values of ``dn``, ``cn``, ``loginShell``, ``uidNumber``, ``gidNumber``, ``sn``, ``homeDirectory``, ``mail``, ``uid``, ``sshPublicKey``.

   See :download:`the adding users LDIF <./_static/users.ldif.txt>`.

7. Change slapd configuration.

   .. code-block:: shell

      $ sudo ldapadd -H ldapi:/// -Y EXTERNAL -f openssh-lpk.ldif
      $ sudo ldapmodify -H ldapi:/// -Y EXTERNAL -f rootdnpw.ldif
      $ sudo ldapadd -x -h localhost -D cn=admin,dc=example,dc=org -W -f ou.ldif
      $ sudo ldapadd -x -h localhost -D cn=admin,dc=example,dc=org -W -f users.ldif
     
.. rubric:: footnote

.. [#] ``slappasswd`` command is contained in ``slapd`` package. Use ``slappasswd`` command in other system.
