Relevant Commands:

dev-

dev-init.sh 
	initialises dev env
	previous env is destroyed
	-e
		launches hyperledger explorer

dev-validate.sh 
	validates the dev env
	installs/instantiates chaincode
dev-validate.sh skip
	validates the dev env
	does NOT install/instantiate chaincode
dev-stop.sh
	stops containers (non destructive)
	used for resetting containers
dev-start.sh
	starts stopped containers
dev-mode.sh
	check launch mode of env (dev or net)


Hyperledger Explorer
	view activities in the fabric network
	check status on env
	sandalone node based web app
	in the browser interface check
		nodes, blocks, transactions, chaincode, channels
	may be installed natively/containers
	uses a DB
		manage info abt network
	uses crypto folder

dev-init.sh -e
	launch env with explorer
exp-init.sh
	resets explorer runtime
	reset the container and relaunch the web app
	use if explorer doesnt work
exp-stop.sh
	stops running explorer container 
exp-start.sh
	starts explorer container


peer

source set-env.sh acme | budget
	changes env context
	sets up env variables
	sets up peer address for the appropriate org in env
	admin identity context for the org
show-env.sh
	prints current env context

peer channel list
	list channels peer has joined
	* must run set-env.sh first
peer channel getinfo -c [channel_name]
	get info about the channel blockchain
	* must run set-env.sh first
peer lifecycle chaincode command
	Flags
		--waitForEvent
			makes the command syncronous
			command submits transaction then waits for commit confirmation from peer before proceeding
			default is true
peer lifecycle chaincode queryinstalled
	get list of installed chaincode on a peer
peer lifecycle chaincode package [tar-filename] --flags
	package a chaincode

	Flags:
		-p path to chaincode
		--label label for the chaincode
				each org can have a unique labelling scheme
peer lifecycle chaincode install [tar-filename] --flags
	install a chaincode package on a peer

	Flags:
		--peerAddresses target peer; multiple peers
peer lifecycle chaincode approveformyorg
	approves an installed package
	Flags:
		-n -v name/version of CC
		-C channel ID
		-sequence sequence number for the CC on the channel
		--package-id ID of the installed package to be approved
		--init-required indicates if initialisation is required
		--peerAddresses
			address of endorsing peer
			all endorsing peers must have CC installed for it to be committed
		-o orderer address
peer lifecycle chaincode checkcommitreadiness
	show orgs that have approved a package
	Flags:
		-n -v name/version of CC
		-C channel ID
		--init-required indicates if initialisation is required
peer lifecycle chaincode commit --flags
	attempt to commit the chaincode
	if the --init-required is set to true and committing is successful, the CC is committed and it needs to initialise initial values before it can be invoked/queried
	Flags:
		-n -v name/version of CC
		-C channel ID
		-sequence sequence number for the CC on the channel
		--init-required indicates if initialisation is required
		--peerAddresses
			address of endorsing peer
			all endorsing peers must have CC installed for it to be committed
		-o orderer address
peer lifecycle chaincode querycommitted --flags
	check commit status of CC throughout the network
	flags:
		-n name of CC
		-C channel ID
peer chaincode invoke --flags
	instantiate a CC
	flags
		--isInit initialises the CC of a committed CC
peer chaincode query
	queries the peer
	only works after the CC has been initialised


Init/Query/Invoke

If the CC is committed with the --init-required
	Init must be executed before query or invoke
	Init means initialise the initial values (based on CC)
	Initial values must be set before any query or invoke happens

peer chaincode query
	query a peer

peer chaincode invoke
	invoke a peer

Update/Upgrade

Local CC Update within an organization
	local change to org specific implementation of CC is NOT an upgrade
	does not change specs
	if an org manages its own implementation of the committed  CC, they can change their code anytime w/o approval from other members (only org admin needs to approve)
	committing is not done for updates
	all peers within the same org must have the same updated chaincode
	local CC parameters MUST match with committed CC parameters on approval of the CC
		this includes the sequence

CC Upgrade
	network level upgrade
	may change: code, parameter
	requires approval and commit

	if only CC code is changed,
		requires the CC to be packaged and installed in org peers
		approval will be provided by the admins of the orgs using the package ID of the newly installed package and the parameters from the current committed CC
		Next Sequence num is used for the approval of the transaction

	if only CC parameter is changed,
		does NOT require re-installation of bytecode
		org admin approves CC upgrade by providing
			package ID of existing installed package
			a new set of parameters
			the next sequence number
		Ex of new parameter:
			Name stays the same
			Sequence WILL change
			everything else MAY change
			initialisation required
				may be wrong to always set this to true
				may accidentally reset every asset to initial data if always set to true
				upgrades may require it to be false





Chaincode Utility scripts

chain.sh
	execute peer chaincode commands
	simplifies execution of peer chaincode and peer lifecycle commands
	set the org context and chaincode env before using
		this means using source set-env.sh org_name
	Ex:
		$ peer lifecycle chaincode queryinstalled
		is equivalent to
		$ chain.sh queryinstalled

		$ peer lifecycle chaincode package [tar-filename]
		is equivalent to
		$ chain.sh package
		default label:
			CC_Name.CC_Version-internal_DEV_Counter.tar.gz
				internal_DEV_Counter
					- increments if the CC package with the current internal_DEV_Counter is already installed
					- different from version of CC
		tar file is created under $HOME/packages

		$ peer lifecycle chaincode install [tar-filename]
		is equivalent to
		$ chain.sh install
			looks for the package in ~/packages
		flags:
			-p creates the package then installs it

		$ peer lifecycle chaincode approveformyorg
		is equivalent to
		$ chain.sh approve

		$ peer lifecycle chaincode checkcommitreadiness
		is equivalent to
		$ chain.sh check

		$ peer lifecycle chaincode commit
		is equivalent to
		$ chain.sh commit

		$ peer lifecycle chaincode querycommitted
		is equivalent to
		$ chain.sh querycommitted

		peer chaincode invoke
		is equivalent to
		$ chain.sh invoke

		$ peer chaincode invoke --isInit
		is equivalent to
		$ chain.sh init



		approves the CC
		commits the CC
		checks if the approval was done with --init-required, if yes, it initialises the CC
		$ chain.sh instantiate

		list out all installed and committed CC
		$ chain.sh list

		packages, installs, approves, and commits
		retains the current version
		increments the sequence number
		$ chain.sh upgrade-auto

	syntax:
		chain.sh command flags
	Flags:
		-o 
			useful for learning commands
			return template command with default flags
set-chain-env.sh
	sets env variable for chaincode args in cc.env.sh
	flags:
		-h help
		-n name
		-v version
		-s new sequence number
		-I --init-required (true or false)
		-i invoke
		-q query
		-g chaincode endorsement policy
		-R private data collection
		-c initialization arg
show-chain-env.sh
	cat env vars in cc.env.sh
reset-chain-env.sh
	reset cc.env.sh to defaults
cc-build.sh
	builds the GoLang chaincode
	the same as
		go build path_to_file
		go build $CC_PATH
cc-run.sh
	runs the GoLang chaincode on the terminal
	the same as
		go run path_to_file
		go run $GOPATH/src/$CC_PATH/*.go