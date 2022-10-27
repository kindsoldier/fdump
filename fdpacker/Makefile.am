
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = fdpacker
fdpacker_SOURCES = fdpacker.go

GOFLAGS = -ldflags="-s -w"

fdpacker$(EXEEXT): $(fdpacker_SOURCES)
	$(GO) build $(GOFLAGS) -o fdpacker$(EXEEXT) $(fdpacker_SOURCES)

distclean-local:
	rm -rf autom4te.cache
	rm -f $(sbin_PROGRAMS)
