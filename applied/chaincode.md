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
	GetStateByRange(startKey, endKey string) (StateQueryIteratorInterface, error)
	GetStateByRangeWithPagination(startKey, endKey string, pageSize int32, bookmark string) (StateQueryIteratorInterface, *peer.QueryResponseMetadata, error)

Private State Management - must be set with PDC
	PutPrivateData(PDCname string, key string, value []byte) error
	GetPrivateData(PDCname string, key string) ([]byte, error)
	DelPrivateData(PDCname string, key string) error
	GetStateByRange(PDCname string, startKey, endKey string) (StateQueryIteratorInterface, error)

Composite
	CreateCompositeKey(objectType string, atts []string) (string, error)
	SplitCompositeKey(compositeKey string) (string, []string, error)
	GetStateByPartialCompositeKey(objectType string, keys []string) (StateQueryIteratorInterface, error)
	GetStateByPartialCompositeKeyWithPagination(objectType string, keys []string, pageSize int32, bookmark string) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error)

	GetPrivateDataByPartialCompositeKey(PDCName, objectType string, keys []string) (StateQueryIteratorInterface, error)

Rich Queries (requires data to be modeled in JSON; peers must use CouchDB; Do NOT use for Update Transactions)
	GetQueryResult(query string) () (StateQueryIteratorInterface, error)
	GetQueryResultWithPagination(query string, pageSize int32, bookmark string) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error)
	GetPrivateDataQueryResult(PDCName, query string) (StateQueryIteratorInterface, error)

History (managed on a per peer basis; peers are configured to update/create history logs; peer manages on a per chaincode/asset basis; managed in GoLevelDB; NOT recommended for update transactions)
	GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error)

Access Control List
	GetID(stub ChaincodeStubInterface) (string, error)
	GetMSPID(stub ChaincodeStubInterface) (string, error)
	GetAttributeValue(stub ChaincodeStubInterface, attrName string) (value string, found bool, err error)
	AssertAttributeValue(stub ChaincodeStubInterface, attrName, attrValue string) error
	GetX509Certificate() (*x509.Certificate, error)