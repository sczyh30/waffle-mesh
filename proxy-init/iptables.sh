#!/usr/bin/env bash

PROXY_INBOUND_PORT=9081
PROXY_OUTBOUND_PORT=9080
PROXY_UID=2186

# Remove old
iptables -t nat -D OUTPUT -p tcp -j WAFFLE_OUTPUT 2>/dev/null

iptables -t nat -F WAFFLE_OUTPUT 2>/dev/null
iptables -t nat -X WAFFLE_OUTPUT 2>/dev/null
iptables -t nat -F WAFFLE_INBOUND 2>/dev/null
iptables -t nat -X WAFFLE_INBOUND 2>/dev/null
iptables -t nat -F WAFFLE_OUTBOUND_REDIRECT 2>/dev/null
iptables -t nat -X WAFFLE_OUTBOUND_REDIRECT 2>/dev/null
iptables -t nat -F WAFFLE_INBOUND_REDIRECT 2>/dev/null
iptables -t nat -X WAFFLE_INBOUND_REDIRECT 2>/dev/null

# Inbound redirect
iptables -t nat -N WAFFLE_INBOUND_REDIRECT
iptables -t nat -A WAFFLE_INBOUND_REDIRECT -p tcp -j REDIRECT --to-port ${PROXY_INBOUND_PORT}

# Outbound redirect
iptables -t nat -N WAFFLE_OUTBOUND_REDIRECT
iptables -t nat -A WAFFLE_OUTBOUND_REDIRECT -p tcp -j REDIRECT --to-port ${PROXY_OUTBOUND_PORT}

# Input chain
iptables -t nat -N WAFFLE_INBOUND
iptables -t nat -A PREROUTING -p tcp -j WAFFLE_INBOUND

iptables -t nat -A WAFFLE_INBOUND -p tcp --dport 22 -j RETURN
iptables -t nat -A WAFFLE_INBOUND -p tcp -j WAFFLE_INBOUND_REDIRECT

# Output chain
iptables -t nat -N WAFFLE_OUTPUT

# Jump to the WAFFLE_OUTPUT chain from OUTPUT chain for all tcp traffic.
iptables -t nat -A OUTPUT -p tcp -j WAFFLE_OUTPUT

iptables -t nat -A WAFFLE_OUTPUT -m owner --uid-owner ${PROXY_UID} -j RETURN
iptables -t nat -A WAFFLE_OUTPUT -m owner --gid-owner ${PROXY_UID} -j RETURN

# Ignore loopback
iptables -t nat -A WAFFLE_OUTPUT -o lo -j RETURN

# Redirect remaining outbound traffic to Proxy
iptables -t nat -A WAFFLE_OUTPUT -j WAFFLE_OUTBOUND_REDIRECT

# Output result
iptables -t nat -vnL
echo "iptables OK"





