#!/bin/bash
set -e

echo "Executing migration command"
/usr/dapper-labs/identity-server/bin/server migrate --config /usr/dapper-labs/identity-server/config/config.ini

echo "Starting identity server"
/usr/dapper-labs/identity-server/bin/server --config /usr/dapper-labs/identity-server/config/config.ini
