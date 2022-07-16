
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

install-data-hook:
if FREEBSD_OS
	chmod a+x $(DESTDIR)/${etcdir}/rc.d/fdagentd
endif
	$(INSTALL) -d -m 750 -o $(SRV_USER) -g $(SRV_GROUP) $(DESTDIR)$(SRV_LOGDIR)
	$(INSTALL) -d -m 750 -o $(SRV_USER) -g $(SRV_GROUP) $(DESTDIR)$(SRV_RUNDIR)
	for data in $(nobase_conf_DATA);do \
	  chmod 0644 $(DESTDIR)$(SRV_CONFDIR)/$$data; \
	done

if FREEBSD_OS
etcdir = @prefix@/etc
nobase_etc_SCRIPTS = rc.d/fdagentd
endif

if LINUX_OS
if SYSTEMD
systemddir = /lib/systemd/system
nobase_systemd_DATA = fdagentd.service
endif
etcdir = @sysconfdir@
nobase_etc_SCRIPTS = init.d/fdagentd
endif

EXTRA_DIST = \
	fdagentd.service.in \
	init.d/fdagentd.in \
	rc.d/fdagentd.in

