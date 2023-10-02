#!/bin/sh

# simple commandline unit test script

t="../tablizer"
fail=0

ex() {
    # execute a test, report+exit on error, stay silent otherwise
    log="/tmp/test-tablizer.$$.log"
    name=$1
    shift

    echo -n "TEST $name "

    $* > $log 2>&1

    if test $? -ne 0; then
        echo "failed, see $log"
        fail=1
    else
        echo "ok"
        rm -f $log
    fi
}

# only use files in test dir
cd $(dirname $0)

echo "Executing commandline tests ..."

# io pattern tests
ex io-pattern-and-file $t bk7 testtable
cat testtable | ex io-pattern-and-stdin $t bk7
cat testtable | ex io-pattern-and-stdin-dash $t bk7 -

# same w/o pattern
ex io-just-file $t testtable
cat testtable | ex io-just-stdin $t
cat testtable | ex io-just-stdin-dash $t -

if test $fail -ne 0; then
    echo "!!! Some tests failed !!!"
    exit 1
fi
