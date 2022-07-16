
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUBDIRS = rc

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = fdagentd fdagentcli 
fdagentd_SOURCES = fdasrv/fdaserv.go 
nodist_fdagentd_SOURCES = fdasrv/fdaconf.go
fdagentcli_SOURCES = fdacli/fdacli.go

EXTRA_fdagentcli_SOURCES = \
	fdaapi/servapi.go \
	fdaapi/userapi.go

EXTRA_fdagentd_SOURCES = \
	fdaapi/servapi.go \
	fdaapi/userapi.go \
	\
	fdasrv/fdaconf.go.in \
	\
	fdasrv/fdacont/contauth.go \
	fdasrv/fdacont/contcomm.go \
	fdasrv/fdacont/contserv.go \
	fdasrv/fdacont/contuser.go \
	\
	fdasrv/fdareg/regcomm.go \
	fdasrv/fdareg/reguser.go \
	\
	fdasrv/fdagent/storecomm.go \
	fdasrv/fdagent/storeuser.go


GOFLAGS = -ldflags="-s -w"

fdagentd$(EXEEXT): $(fdagentd_SOURCES) $(EXTRA_fdagentd_SOURCES)
	$(GO) build $(GOFLAGS) -o fdagentd$(EXEEXT) $(fdagentd_SOURCES) \
	   $(nodist_fdagentd_SOURCES)

fdagentcli$(EXEEXT): $(fdagentcli_SOURCES) $(EXTRA_fdagentcli_SOURCES)
	$(GO) build $(GOFLAGS) -o fdagentcli$(EXEEXT) $(fdagentcli_SOURCES)

EXTRA_DIST = \
	fdagent.conf

distclean-local:
	rm -rf autom4te.cache
	rm -rf tmp.log/ tmp.run/ tmp.data/
	rm -f fdasrv/fdasrv
	rm -f fdacli/fdacli

install-data-local:
	test -z $(DESTDIR)$(SRV_CONFDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_CONFDIR)
	test -z $(DESTDIR)$(SRV_LOGDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_LOGDIR)
	test -z $(DESTDIR)$(SRV_RUNDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_RUNDIR)
	test -z $(DESTDIR)$(SRV_DATADIR) || $(MKDIR_P) $(DESTDIR)$(SRV_DATADIR)

	test -z $(DESTDIR)$(SRV_LOGDIR) || $(CHOWN) $(SRV_USER):$(SRV_GROUP) $(DESTDIR)$(SRV_LOGDIR)
	test -z $(DESTDIR)$(SRV_RUNDIR) || $(CHOWN) $(SRV_USER):$(SRV_GROUP) $(DESTDIR)$(SRV_RUNDIR)
	test -z $(DESTDIR)$(SRV_DATADIR) || $(CHOWN) $(SRV_USER):$(SRV_GROUP) $(DESTDIR)$(SRV_DATADIR)

