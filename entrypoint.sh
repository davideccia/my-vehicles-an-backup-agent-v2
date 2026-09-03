#!/bin/sh
set -e

# Runs as root so it can fix ownership of bind-mounted volumes (e.g. the
# SQLite data directory), which Docker creates as root on the host when
# they don't already exist. su-exec then drops to the unprivileged
# appuser for everything else — migrate and the server itself never run
# as root.
chown -R appuser:appgroup /www/storage

exec su-exec appuser sh -c '/www/main artisan migrate && exec /www/main'
