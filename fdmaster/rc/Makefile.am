
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

install-data-hook:
if FREEBSD_OS
	chmod a+x $(DESTDIR)/${etcdir}/rc.d/fdmasterd
endif
	$(INSTALL) -d -m 750 -o $(SRV_USER) -g $(SRV_GROUP) $(DESTDIR)$(SRV_LOGDIR)
	$(INSTALL) -d -m 750 -o $(SRV_USER) -g $(SRV_GROUP) $(DESTDIR)$(SRV_RUNDIR)
	for data in $(nobase_conf_DATA);do \
	  chmod 0644 $(DESTDIR)$(SRV_CONFDIR)/$$data; \
	done

if FREEBSD_OS
etcdir = @prefix@/etc
nobase_etc_SCRIPTS = rc.d/fdmasterd
endif

if LINUX_OS
if SYSTEMD
systemddir = /lib/systemd/system
nobase_systemd_DATA = fdmasterd.service
endif
etcdir = @sysconfdir@
nobase_etc_SCRIPTS = init.d/fdmasterd
endif

EXTRA_DIST = \
	fdmasterd.service.in \
	init.d/fdmasterd.in \
	rc.d/fdmasterd.in

