# wstunnel
Tunneling SOCKS5 proxy over Websockets.

This project contains two binaries: client and server.
Running them in tandem creates a tunnel over Websocket protocol that proxies any information.

## Example
Alice wants to connect to Bob's computer via SSH (port 22), but Alice is connected to Eve's Wifi,
and Eve has a firewall in place that blocks port 22. In fact, Eve's firewall only lets web traffic
on ports 80 and 443 go through.

Alice can ask Faythe, who has no firewall, to set up the server (from this package) on her computer:
    bazel run :server

And run the client locally:
    bazel run :client -- -host=faythe.com

Now, any SOCKS5 message sent to localhost:8080 will be tunneled via websockets to Faythe's computer.
If Alice wants to SSH Bob now, she can simply do:
    ssh -o "ProxyCommand=nc -X 5 -x localhost:8080 %h %p" bob.com

If, Except for the Firewall, Alice's connection to the internet must go through a SOCKS5 proxy, she
can run the client locally using:
    all_proxy=socks5://outbound.alice.com:12345/ bazel run :client -- -host=faythe.com
