# vim:ts=8

PREFIX	?= /usr/local
BINDIR	?= $(DESTDIR)$(PREFIX)/bin
MANDIR	?= $(DESTDIR)$(PREFIX)/man/man1

INSTALL_PROGRAM ?= install -s
INSTALL_DATA ?= install
SUDO ?= sudo

PROG	= podfetcher

all: install

install: all
	$(SUDO) $(INSTALL_PROGRAM) $(PROG) $(BINDIR)
	@mkdir -p ~/.podfetcher
	$(INSTALL_DATA) -m 644 config.toml ~/.podfetcher/config.toml
	$(INSTALL_DATA) -m 644 feeds ~/.podfetcher/feeds

.PHONY: all install
