# *.spec files should be mapped into the /spec dir
# source files should be mapped into $GOPATH/src/...
# spec file should use go build tools with $GOPATH var (or defaults)
# outputs will be in /rpms

SHELL = /bin/bash

all:
	@rm -rf /root/rpmbuild
	@mkdir -p /root/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
	@cp /spec/*.spec /root/rpmbuild/SPECS
	@rpmbuild -bb /root/rpmbuild/SPECS/*.spec
	@cp -r --remove-destination /root/rpmbuild/RPMS/* /rpms
