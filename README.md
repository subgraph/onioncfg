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

* torrc on disk has DisableNetwork set to 1 so that the Tor client start making connections at system startup unless told by something (i.e. onioncfg) that it can

* Pluggable transport has been defined in torrc already, e.g.:

``` 
ClientTransportPlugin obfs3,obfs4 exec /usr/bin/obfs4proxy
```

* Some hardening on systemd Tor startup has been relaxed so that the obfs4proxy executable can be run (still figuring out what the ideal configuration is..)

* Tor control port is listening on TCP port 9111. This is because of a Tor quirk where if DisableNetwork is set to 1, the unix domain socket interface for
the Tor Control protocol isn't created even if specified (and it seems it should be). Therefore, for now, we require a TCP Tor Control port to be listening. We arbitrarily chose 9111.



