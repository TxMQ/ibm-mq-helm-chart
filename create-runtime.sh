#!/bin/bash -xe

# select mq packages for installation

#askquestion inc32 "Does the runtime require 32-bit application support [Y/N]? " "" "" 32
export genmqpkg_inc32=0

#askquestion incnls "Does the runtime require support for languages other than English [Y/N]? " "" "Msg_cs Msg_de Msg_es Msg_fr Msg_hu Msg_it Msg_ja Msg_ko Msg_pl Msg_pt Msg_ru Msg_Zh_CN Msg_Zh_TW" nls
export genmqpkg_incnls=0

#askquestion inccpp "Does the runtime require C++ libraries [Y/N]? " "" "" cpp
export genmqpkg_inccpp=0

#askquestion inccbl "Does the runtime require COBOL libraries [Y/N]? " "" "" cobol
export genmqpkg_inccbl=0

#askquestion incdnet "Does the runtime require .NET libraries [Y/N]? " "" "" dotnet
export genmqpkg_incdnet=0

#askquestion inctls "Does the runtime require SSL/TLS support [Y/N]? "
export genmqpkg_inctls=1

#askquestion incams "Does the runtime require AMS support [Y/N]? " "" AMS
export genmqpkg_incams=0

#askquestion inccics "Does the runtime require CICS support [Y/N]? " "" "" cics
export genmqpkg_inccics=0

#askquestion incadm "Does the runtime require any administration tools [Y/N]? " "" "" adm
export genmqpkg_incadm=1

#askquestion incras "Does the runtime require any RAS tools [Y/N]? " "" "" ras
export genmqpkg_incras=1

#askquestion incsamp "Does the runtime require any sample applications [Y/N]? " "" Samples samp
export genmqpkg_incsamp=0

#askquestion incsdk "Does the runtime require the SDK to compile applications [Y/N]? " "" SDK
export genmqpkg_incsdk=1

#askquestion incunthrd "Does the runtime require unthreaded application support [Y/N]? " "" "" unthrd
export genmqpkg_incunthrd=0

#askquestion incjre "Does the runtime require a Java Runtime Environment (JRE) [Y/N]? " "$mqdir/java/jre64" JRE
export genmqpkg_incjre=1

#askquestion incamqp "Does the runtime require AMQP support [Y/N]? " "$mqdir/amqp" AMQP
export genmqpkg_incamqp=0

#askquestion incman "Does the runtime require man pages [Y/N]? " "$mqdir/man" Man
export genmqpkg_incman=0

#askquestion incmqbc "Does the runtime require the Bridge to blockchain [Y/N]? " "$mqdir/mqbc" BCBridge
export genmqpkg_incmqbc=0

#askquestion incmqft "Does the runtime require Managed File Transfer [Y/N]? " "$mqdir/mqft" "FTAgent FTBase FTLogger FTService FTTools" mft
export genmqpkg_incmqft=0

#askquestion incmqsf "Does the runtime require the Bridge to Salesforce [Y/N]? " "$mqdir/mqsf" SFBridge
export genmqpkg_incmqsf=0

#askquestion incmqxr "Does the runtime require Telemetry (MQXR) support [Y/N]? " "$mqdir/mqxr" XRService
export genmqpkg_incmqxr=1

#askquestion incweb "Does the runtime require the MQ Console [Y/N]? " "$mqdir/web" Web
export genmqpkg_incweb=1

#install mq
# wget https://public.dhe.ibm.com/ibmdl/export/pub/software/websphere/messaging/mqadv/9.2.2.0-IBM-MQ-Advanced-for-Developers-Non-Install-LinuxX64.tar.gz
/Users/simon/txmq/mq/non-install/bin/genmqpkg.sh -b /tmp/genmqpkg
