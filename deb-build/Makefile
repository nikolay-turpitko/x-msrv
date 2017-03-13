# source files should be mapped into $GOPATH/src/...
# outputs will be in /deb

SHELL=/bin/bash

APPDIR=$(GOPATH)/src/github.com/nikolay-turpitko/x-msrv

PKG_NAME="x-msrv"
PKG_VERSION="0.0.1"
DEPLOY_DIR=/deploy/$(PKG_NAME)/$(PKG_NAME)-$(PKG_VERSION)

export DEBFULLNAME="Nikolay Turpitko"
export DEBEMAIL="xxx@yyy.com"

all:
	# flags for linker tells it to link for linux and omit the symbol table and debug information
	@cd $(APPDIR) && \
		glide -q --no-color install && \
		go test -compiler gc -ldflags '-H linux -s' && \
		go clean -i -r && \
		go install -compiler gc -ldflags '-H linux -s'
	@rm -rf /deploy
	@mkdir -p $(DEPLOY_DIR)/{bin,man,etc,systemd}
	@cp -p $(GOPATH)/bin/x-msrv $(DEPLOY_DIR)/bin
	@cp -r $(APPDIR)/systemd $(DEPLOY_DIR)
	@cp -r $(APPDIR)/etc $(DEPLOY_DIR)
	@pandoc $(APPDIR)/x-msrv.md -s -t man > $(DEPLOY_DIR)/man/x-msrv.1
	@cd $(DEPLOY_DIR) && \
		dh_make \
		--yes \
		--native \
		--copyright mit \
		--email $(DEBEMAIL) \
		--single \
		--packagename $(PKG_NAME)
	@cd $(DEPLOY_DIR)/debian && \
		rm -f ./*.ex ./*.EX ./*.docs README.* README
	@cp -r /debian $(DEPLOY_DIR)
	@cd $(DEPLOY_DIR) && debuild -i -us -uc -b --lintian-opts --profile ubuntu
	@cp /deploy/$(PKG_NAME)/*.deb /deb
