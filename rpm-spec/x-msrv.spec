Name: x-msrv
Version: 0.0.1
Release: 1%{?dist}
Summary: sample microservice
Group: Applications/Internet
License: BSD
URL: https://github.com/nikolay-turpitko/x-msrv

Requires: systemd,rsyslog
Requires(post): systemd
Requires(preun): systemd
Requires(postun): systemd

%description
Sample microservice.

%build
appdir=$GOPATH/src/github.com/nikolay-turpitko/x-msrv
cd $appdir
glide -q --no-color install
# flags for linker tells it to link for linux and omit the symbol table and debug information
go test -compiler gc -ldflags '-H linux -s'
go clean -i -r
go install -compiler gc -ldflags '-H linux -s'
rm -rf %{_builddir}
mkdir -p %{_builddir}%{_bindir}
mkdir -p %{_builddir}%{_sysconfdir}/x-msrv/
mkdir -p %{_builddir}%{_unitdir}
mkdir -p %{_builddir}%{_mandir}/man1
mkdir -p %{_builddir}%{_sysconfdir}/x-msrv
cp -p $appdir/LICENSE %{_builddir}
cp -p $GOPATH/bin/x-msrv %{_builddir}%{_bindir}
cp -p $appdir/systemd/x-msrv.service %{_builddir}%{_unitdir}
cp -p $appdir/systemd/x-msrv.timer %{_builddir}%{_unitdir}
cp -p $appdir/etc/* %{_builddir}%{_sysconfdir}/x-msrv/
pandoc $appdir/x-msrv.md -s -t man > %{_builddir}%{_mandir}/man1/x-msrv.1

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}%{_bindir}
mkdir -p %{buildroot}%{_sysconfdir}/x-msrv/
mkdir -p %{buildroot}%{_unitdir}
mkdir -p %{buildroot}%{_mandir}/man1
mkdir -p %{buildroot}%{_sysconfdir}/x-msrv/
cp -p %{_builddir}%{_bindir}/x-msrv %{buildroot}%{_bindir}
cp -p %{_builddir}%{_unitdir}/x-msrv.service %{buildroot}%{_unitdir}
cp -p %{_builddir}%{_unitdir}/x-msrv.timer %{buildroot}%{_unitdir}
cp -p %{_builddir}%{_mandir}/man1/x-msrv.1 %{buildroot}%{_mandir}/man1
cp -p %{_builddir}%{_sysconfdir}/x-msrv/* %{buildroot}%{_sysconfdir}/x-msrv/

%clean
rm -rf %{_builddir}
rm -rf %{buildroot}

%post
%systemd_post x-msrv.service
%systemd_post x-msrv.timer

if [ $1 -eq 1 ] ; then
	# Initial installation
	
	# Create user and group for application
	getent group x-msrv >/dev/null 2>&1 || groupadd -r x-msrv || :
	getent passwd x-msrv >/dev/null 2>&1 || \
		useradd -r -g x-msrv -d /var/lib/x-msrv -s /sbin/nologin \
		-c "Sample microservice" x-msrv >/dev/null 2>&1 || :
fi

%preun
%systemd_preun x-msrv.service
%systemd_preun x-msrv.timer

%postun
%systemd_postun x-msrv.service
%systemd_postun x-msrv.timer

%files
%doc LICENSE
%{_mandir}/man1/*
%{_bindir}/x-msrv
%config %{_unitdir}/x-msrv.service
%config %{_unitdir}/x-msrv.timer
%config %{_sysconfdir}/x-msrv/*


%changelog
* Sun Mar 12 2017 Nikolay Turpitko - 0.0.1-1
- initial build version
