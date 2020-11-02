fabric-network package
	Gateway class - connection point for accessing the fabric network
		new
			create gateway instance
		<async>connect(connection-profile, connection-options)
			connects network to gateway and connects identity to gateway
		<async>getNetwork(network-name) 
			get network instance
		<async>getClient() 
			get client instance
		<async>disconnect()
	Network class - represents set of peers belonging to a network/channel
		getContract(chaincode-id) returns object representing chaincode
		getChannel() returns Channel class
	Contract class - API for invoke and query
		<async>submitTransaction(name [,args])
			invoke
			waits for commit event (default)
				TimeoutError - if commit event is not received withint timeout duration
				Custom:
					remove event listener
					create custom event listener
		<async>evaluateTransaction(name [,args])
			query
	Transaction class - API for invoke and query (provide finer control)
		<async>createTransaction(name)
			a new transaction must be created for each separate transaction
		<async>submit([args])
			invoke
		<async>evaluate([args])
			query
		<async>setTransient(transientMap)
		getTransactionId()
		getName()
		getNetwork()


	Wallet - manages identity of a user in different channels

	X509WalletMixin class - creates identities for wallets

	Wallet Interface
		label - identities in a walled are identified using this
		<async>import(label,identity) - import an identity into a wallet
		<async>export(label) - export an identity from a wallet
		<async>list() - list the identities
		<async>exists(label) - check if the identity with label exists
		<async>delete(label) - delete identity in wallet
	InMemoryWallet class
		- creates wallet at runtime
		- does not persist
		- needs custom identities store
	FileSystemWallet class
		- profiles are managed in user's filesystem
		- requires folder location
	CouchDBWallet class
		- profiles managed in centralized DB
		- requires URL to CouchDB server

fabric-client package
	Client class - API for interacting with peers/orderers (install/instantiate chaincode, query peers, creating channels)
		Create an instance of the class
			a) <async, static>loadFromConfig(object | path-to-file)
			b) new() --> <async>loadFromConfig(object | path-to-file)
		Connecting to User Context (User class)
			a) ernroll user with fabric-ca-server to generate crypto
			b) create user context using cert/key of existing user
		Initialise credentials store
			<async>initCredentialStores()
			<async>createUser(opts)
			<async>loadUserFromStateStore(name)
			<async>setUserContext(userContext | userNamePasswordObject, skipPersistent)
	User class - User context/credentials
		can be persisted
