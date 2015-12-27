#!/bin/sh

## DATATRACK_DEBUG
## Activate debug messages.
## default: false
export DATATRACK_DEBUG=true

## DATATRACK_DATABASEPATH
## Determine the path of the database file.
## default: datatrack.db
# export DATATRACK_DATABASEPATH=datatrack.db

## Start datatrack server
bin/datatrack
