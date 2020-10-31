All CC must implement chaincode interface
	init
		initialisation logic
		initialisation of CC
		called when app submits transaction with init flag
		can be invoked using $ peer chaincode invoke --is-init
	invoke
		business logic
		can be invoked from query/invoke wo init api 
		can be invoked using
			peer chaincode invoke
			peer chaincode query
	main
		register the CC instance with fabric runtime

fabric runtime maps CC name to CC instance
client applications invoke CC by passing CC name to fabric runtime

Logging
	NET Mode - container file system
	DEV Mode - console

	Minimum Logging Level Heirarchy
		in net mode, set in core.yml file of the peer that instantiates the CC container. (can be overridden by env vars in chaincode/logging/)
		in dev mode, set in env variables
		Value:
			DEBUG
			INFO
			NOTICE
			WARNING
			ERROR
			CRITICAL

Supplementary Logger
gocc/src/acflogger/acflogger.go
	uses env variables to set logging level

cc-logs.sh
	read log messages of CC
	flags
		-o org name
		-p peer name

token/v2/token.go
	uses custom acflogger
	when invoke is executed, it emits logs at all levels

Get Arguments in stub
	GetArgs() [][]byte
	GetStringArgs() []string
	GetFunctionAndParameters() (string, []string)
	GetArgsSlice() ([]byte, error) 

State Management
	PutState(key string, value []byte) error
	GetState(key string) ([]byte, error)
	DelState(key string) error

Private State Management - must be set with PDC
	PutPrivateData(PDCname string, key string, value []byte) error
	GetPrivateData(PDCname string, key string) ([]byte, error)
	DelPrivateData(PDCname string, key string) error