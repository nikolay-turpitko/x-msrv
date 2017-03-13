% x-msrv(1)
% Nikolay Turpitko
% March 2017

# NAME

x-msrv - sample blueprint microservice.

# DESCRIPTION

x-msrv is a sample microservice which is managed by systemd.

It is periodically executed as systemd timer service, connects to NSQ to obtain
bunch of JSON messages, validates and parses them and stores into Aerospike.

When service configuration changed, it should be restarted using commands below:

        sudo systemctl daemon-reload
        sudo systemctl restart x-msrv

# FILES

**/etc/x-msrv.yml**
:   Configuration file containing queue and DB connection details for this service.

**/lib/systemd/system/x-msrv.service**
:   Configuration file of systemd unit for this service.

**/lib/systemd/system/x-msrv.timer**
:   Configuration file of systemd timer for this service.

# SEE ALSO

systemctl(1), journalctl(1)
