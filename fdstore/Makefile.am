
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUBDIRS = rc

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = fdstored fdstorecli 
fdstored_SOURCES = fdssrv/fdsserv.go 
nodist_fdstored_SOURCES = fdssrv/fdsconf.go
fdstorecli_SOURCES = fdscli/fdscli.go

EXTRA_fdstorecli_SOURCES = \
	fdsapi/servapi.go \
	fdsapi/userapi.go

EXTRA_fdstored_SOURCES = \
	fdsapi/servapi.go \
	fdsapi/userapi.go \
	\
	fdssrv/fdsconf.go.in \
	\
	fdssrv/fdscont/contauth.go \
	fdssrv/fdscont/contcomm.go \
	fdssrv/fdscont/contserv.go \
	fdssrv/fdscont/contuser.go \
	\
	fdssrv/fdsreg/regcomm.go \
	fdssrv/fdsreg/reguser.go \
	\
	fdssrv/fdstore/storecomm.go \
	fdssrv/fdstore/storeuser.go


GOFLAGS = -ldflags="-s -w"

fdstored$(EXEEXT): $(fdstored_SOURCES) $(EXTRA_fdstored_SOURCES)
	$(GO) build $(GOFLAGS) -o fdstored$(EXEEXT) $(fdstored_SOURCES) \
	   $(nodist_fdstored_SOURCES)

fdstorecli$(EXEEXT): $(fdstorecli_SOURCES) $(EXTRA_fdstorecli_SOURCES)
	$(GO) build $(GOFLAGS) -o fdstorecli$(EXEEXT) $(fdstorecli_SOURCES)

EXTRA_DIST = \
	fdstore.conf

distclean-local:
	rm -rf autom4te.cache
	rm -rf tmp.log/ tmp.run/ tmp.data/
	rm -f fdssrv/fdssrv
	rm -f fdscli/fdscli

install-data-local:
	test -z $(DESTDIR)$(SRV_CONFDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_CONFDIR)
	test -z $(DESTDIR)$(SRV_LOGDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_LOGDIR)
	test -z $(DESTDIR)$(SRV_RUNDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_RUNDIR)
	test -z $(DESTDIR)$(SRV_DATADIR) || $(MKDIR_P) $(DESTDIR)$(SRV_DATADIR)

	test -z $(DESTDIR)$(SRV_LOGDIR) || $(CHOWN) $(SRV_USER):$(SRV_GROUP) $(DESTDIR)$(SRV_LOGDIR)
	test -z $(DESTDIR)$(SRV_RUNDIR) || $(CHOWN) $(SRV_USER):$(SRV_GROUP) $(DESTDIR)$(SRV_RUNDIR)
	test -z $(DESTDIR)$(SRV_DATADIR) || $(CHOWN) $(SRV_USER):$(SRV_GROUP) $(DESTDIR)$(SRV_DATADIR)

