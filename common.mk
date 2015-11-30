LDFLAGS = \
	-X $(NAME)/$(SDIR)/cli.version=$(VERSION) \
	-X $(NAME)/$(SDIR)/bot.version=$(VERSION) \
	-B 0x$(shell head -c20 /dev/urandom|od -An -tx1|tr -d ' \n')

all: $(NAME)
$(NAME): *.go $(SDIR)/*/*.go
	go build -a -ldflags "$(LDFLAGS)" -v

release:
	# /src/ in this section are not $$(SDIR) !
	mkdir -p $(NAME)-"$(VERSION)"/src/$(NAME)

	rsync -avzr --delete \
		--filter='- $(NAME)-*' \
		--filter='- /$(NAME)' \
		--filter='- *.db' \
		--filter='+ /.git/' \
		--filter='- *~' \
		--filter='- *.org' \
		--filter='- .*' \
		. $(NAME)-"$(VERSION)"/src/$(NAME)

	tar czf $(NAME)-"$(VERSION)".tgz $(NAME)-"$(VERSION)"
