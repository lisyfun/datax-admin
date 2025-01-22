#!/bin/sh

./datax-admin &
caddy run --config /etc/caddy/Caddyfile
