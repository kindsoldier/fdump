
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUBDIRS = rc

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = fdmasterd fdmastercli 
fdmasterd_SOURCES = fdmsrv/fdmserv.go 
nodist_fdmasterd_SOURCES = fdmsrv/fdmconf.go
fdmastercli_SOURCES = fdmcli/fdmcli.go

EXTRA_fdmastercli_SOURCES = \
	fdmapi/servapi.go \
	fdmapi/userapi.go

EXTRA_fdmasterd_SOURCES = \
	fdmapi/servapi.go \
	fdmapi/userapi.go \
	\
	fdmsrv/fdmconf.go.in \
	\
	fdmsrv/fdmcont/contauth.go \
	fdmsrv/fdmcont/contcomm.go \
	fdmsrv/fdmcont/contserv.go \
	fdmsrv/fdmcont/contuser.go \
	\
	fdmsrv/fdmreg/regcomm.go \
	fdmsrv/fdmreg/reguser.go \
	\
	fdmsrv/fdmaster/storecomm.go \
	fdmsrv/fdmaster/storeuser.go


GOFLAGS = -ldflags="-s -w"

fdmasterd$(EXEEXT): $(fdmasterd_SOURCES) $(EXTRA_fdmasterd_SOURCES)
	$(GO) build $(GOFLAGS) -o fdmasterd$(EXEEXT) $(fdmasterd_SOURCES) \
	   $(nodist_fdmasterd_SOURCES)

fdmastercli$(EXEEXT): $(fdmastercli_SOURCES) $(EXTRA_fdmastercli_SOURCES)
	$(GO) build $(GOFLAGS) -o fdmastercli$(EXEEXT) $(fdmastercli_SOURCES)

EXTRA_DIST = \
	fdmaster.conf

distclean-local:
	rm -rf autom4te.cache
	rm -rf tmp.log/ tmp.run/ tmp.data/
	rm -f fdmsrv/fdmsrv
	rm -f fdmcli/fdmcli

install-data-local:
	test -z $(DESTDIR)$(SRV_CONFDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_CONFDIR)
	test -z $(DESTDIR)$(SRV_LOGDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_LOGDIR)
	test -z $(DESTDIR)$(SRV_RUNDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_RUNDIR)
	test -z $(DESTDIR)$(SRV_DATADIR) || $(MKDIR_P) $(DESTDIR)$(SRV_DATADIR)

	test -z $(DESTDIR)$(SRV_LOGDIR) || $(CHOWN) $(SRV_USER):$(SRV_GROUP) $(DESTDIR)$(SRV_LOGDIR)
	test -z $(DESTDIR)$(SRV_RUNDIR) || $(CHOWN) $(SRV_USER):$(SRV_GROUP) $(DESTDIR)$(SRV_RUNDIR)
	test -z $(DESTDIR)$(SRV_DATADIR) || $(CHOWN) $(SRV_USER):$(SRV_GROUP) $(DESTDIR)$(SRV_DATADIR)

