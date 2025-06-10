#!/bin/sh

nginx -g 'daemon off;' &
/app/datax-admin
