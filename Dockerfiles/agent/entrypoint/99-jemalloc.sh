#!/bin/bash

# Don't allow starting without an apikey set
if [[ ! -z "${JEMALLOC_CONF}" ]]; then
    export JEMALLOC_LD_PRELOAD="/opt/datadog-agent/embedded/lib/libjemalloc.so.2"
fi
