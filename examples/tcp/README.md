# TCP Muxy Example

Start up three shells, one to listen on a TCP port:

_NOTE_: Ensure that [netcat](http://nc110.sourceforge.net/) is installed first.
```
nc -k -l 8000
```

...and the other to run `muxy`:

```
muxy proxy --config tcp.yml
```

...and the other to send some messages to:

```
telnet localhost 8001`
```

In this final shell, put in a message such as "hello". In the other shell,
you should see the message being overridden with "wow, new request!".

In that shell, type in "goodbye!". Back in the original shell, you should see
this message overridden with "wow, new response!".
