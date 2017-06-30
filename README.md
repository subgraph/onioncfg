# onioncfg
Onion Bridge Config GUI Utility in golang

Still very new.

## Building

Prereq: 
sudo apt-get install gtk+3.0 libgtk-3-dev

*Might* require:
go install -v -tags gtk_3_18 -gcflags "-N -l" onioncfg

## Using

Assumptions made by onioncfg:

* torrc on disk has DisableNetwork set to 1

* Pluggable transport has been defined in torrc already

* Some hardening on systemd Tor startup has been relaxed so that the obfs4proxy executable can be run (still figuring out what the ideal configuration is..)

