#!/bin/sh

caddy run --config /etc/caddy/Caddyfile &
./datax-admin
