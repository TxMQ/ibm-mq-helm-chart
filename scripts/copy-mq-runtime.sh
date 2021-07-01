#!/bin/bash -xe

noidir=$1
installdir=$2

# select mq packages for installation

# inc32 "Does the runtime require 32-bit application support?"
export genmqpkg_inc32=0

# incnls "Does the runtime require support for languages other than English?"
export genmqpkg_incnls=0

# inccpp "Does the runtime require C++ libraries?"
export genmqpkg_inccpp=0

# inccbl "Does the runtime require COBOL libraries?"
export genmqpkg_inccbl=0

# incdnet "Does the runtime require .NET libraries?"
export genmqpkg_incdnet=0

# inctls "Does the runtime require SSL/TLS support?"
export genmqpkg_inctls=1

# incams "Does the runtime require AMS support?"
export genmqpkg_incams=0

# inccics "Does the runtime require CICS support [Y/N]?"
export genmqpkg_inccics=0

# incadm "Does the runtime require any administration tools?"
export genmqpkg_incadm=1

# incras "Does the runtime require any RAS tools?"
export genmqpkg_incras=1

# incsamp "Does the runtime require any sample applications"
export genmqpkg_incsamp=0

# incsdk "Does the runtime require the SDK to compile applications?"
export genmqpkg_incsdk=1

# incunthrd "Does the runtime require unthreaded application support?"
export genmqpkg_incunthrd=0

# incjre "Does the runtime require a Java Runtime Environment (JRE)"
export genmqpkg_incjre=1

# incamqp "Does the runtime require AMQP support? "
export genmqpkg_incamqp=0

# incman "Does the runtime require man pages?"
export genmqpkg_incman=0

# incmqbc "Does the runtime require the Bridge to blockchain?"
export genmqpkg_incmqbc=0

# incmqft "Does the runtime require Managed File Transfer?"
export genmqpkg_incmqft=0

# incmqsf "Does the runtime require the Bridge to Salesforce?"
export genmqpkg_incmqsf=0

# incmqxr "Does the runtime require Telemetry (MQXR) support?"
export genmqpkg_incmqxr=0

# incweb "Does the runtime require the MQ Console?"
export genmqpkg_incweb=1

#
#install mq
#

$noidir/bin/genmqpkg.sh -b $installdir
