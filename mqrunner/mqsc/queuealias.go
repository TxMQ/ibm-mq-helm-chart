package mqsc

/*
DEFINE QALIAS( q_name )
   [ CLUSNL( namelist_name ) ]             [ CLUSTER( cluster_name ) ]
   [ CLWLPRTY( integer ) ]                 [ CLWLRANK( integer ) ]
   [ CUSTOM( string ) ]
   [ DEFBIND( NOTFIXED | OPEN | GROUP ) ]  [ DEFPRESP( SYNC | ASYNC ) ]
   [ DEFPRTY( integer ) ]                  [ DEFPSIST( YES | NO ) ]
   [ DEFREADA( NO | YES | DISABLED ) ]     [ DESCR( string ) ]
   [ GET( ENABLED | DISABLED ) ]           [ LIKE( qalias_name ) ]
   [ PUT( ENABLED | DISABLED ) ]
   [ PROPCTL( COMPAT | NONE | ALL | FORCE | V6COMPAT ) ]
   [ REPLACE | NOREPLACE ]                 [ SCOPE( QMGR | CELL ) ]
   [ TARGET( string ) ]                    [ TARGTYPE( QUEUE | TOPIC ) ]
 */

type Queuealias struct {

	Name string
	Descr string
	Target string // target queue or topic
	Targtype string // queue/topic
	Put string // enabled/disabled
	Get string // enabled/disabled
	DefaultPriority int
	DefaultPersistence bool

	Authority []Authrec
	Alter []string
}

func (aq *Queuealias) Mqsc() string {
	return ""
}