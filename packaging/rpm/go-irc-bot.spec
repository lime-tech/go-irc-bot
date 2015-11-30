Summary: Extendable IRC bot
Name: go-irc-bot
Version: 2015122205.f0aae3ef
Release: 1
License: GPLv2+
Source0: %{name}-%{version}.tgz
Source1: %{name}.service
Source2: %{name}.tmpfiles.d.conf
Packager: Bob
BuildRequires: golang

%description
Just small IRC workflow bot, see README.

%prep
%setup

%build
export GOPATH="$(pwd)"
export PATH="$PATH:$(pwd)/bin"
export root='src/'%{name}
source "$root"/bootstrap
make %{name}

%post
%tmpfiles_create %{name}.conf
%systemd_post %{name}.service
chown -R %{name} %{_sysconfdir}/%{name}
chmod 700 %{_sysconfdir}/%{name}
chmod 600 %{_sysconfdir}/%{name}/*.toml

%pre
# Add the "%{name}" user
getent group %{name} >/dev/null || groupadd -r %{name}
mkdir -p /var/lib/%{name} >/dev/null || true
getent passwd %{name} >/dev/null || \
    useradd -r -g %{name} -s /sbin/nologin -d "/var/lib/"%{name} -c "%{summary}" %{name}
chown %{name}:%{name} /var/lib/%{name}
exit 0

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun_with_restart %{name}.service

%clean
rm -rf $RPM_BUILD_ROOT

%install
install -m0755 -D            src/%{name}/%{name}       $RPM_BUILD_ROOT%{_bindir}/%{name}
install -m0600 -D            src/%{name}/config.toml   $RPM_BUILD_ROOT%{_sysconfdir}/%{name}/config.toml
install        -D -p -m 0644 %{SOURCE1}                $RPM_BUILD_ROOT%{_unitdir}/%{name}.service
install        -D -p -m 0644 %{SOURCE2}                $RPM_BUILD_ROOT%{_tmpfilesdir}/%{name}.conf

%files
%defattr(-,root,root,0755)
%doc src/%{name}/README.md
%config(noreplace) %{_sysconfdir}/%{name}/config.toml
%{_bindir}/%{name}
%{_tmpfilesdir}/%{name}.conf
%{_unitdir}/%{name}.service
