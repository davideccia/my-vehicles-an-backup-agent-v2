#!/bin/sh
set -e

/www/main artisan migrate

exec /www/main
