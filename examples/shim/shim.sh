#!/bin/bash

SHIM_ROOT=~/shims
mkdir -p ${SHIM_ROOT}

case :$PATH:
  in *:${SHIM_ROOT}:*) ;;
     *) echo "note: ${SHIM_ROOT} is not in your PATH ($PATH)" >&2;;
esac

if [ -z $1 ]; then
    ls ${SHIM_ROOT}
    exit 0
fi

if [ -z $2 ]; then
    NAME=$1
else
    NAME=$2
fi

DESTINATION="${SHIM_ROOT}/${NAME}"

if [ -e ${DESTINATION} ]; then
    echo "${DESTINATION} already exists." >&2
    exit 1
fi

ln -sv $(realpath $1) $(realpath ${DESTINATION})
chmod +x ${DESTINATION}