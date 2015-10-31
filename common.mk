LDFLAGS = -X cli.version=$(VERSION) -B 0x$(shell head -c20 /dev/urandom|od -An -tx1|tr -d ' \n')
NAME = go-irc-bot

all: $(NAME)
$(NAME): $(shell find . -type f -name '*.go')
	go build -a -ldflags "$(LDFLAGS)" \
	-v -x

release:
	mkdir -p $(NAME)-"$(VERSION)"/src/$(NAME)

	rsync -avzr --delete \
		--filter='- $(NAME)-*' \
		--filter='- /$(NAME)' \
		--filter='- .*' \
		. $(NAME)-"$(VERSION)"/src/$(NAME)

	tar czf $(NAME)-"$(VERSION)".tgz $(NAME)-"$(VERSION)"
