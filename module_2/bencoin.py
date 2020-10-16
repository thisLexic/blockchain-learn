import datetime, hashlib, json, requests
from flask import Flask, jsonify, request
from uuid import uuid4
from urllib.parse import urlparse


class Blockchain:

    def __init__(self):
        # include length of the chain
        # transactions is the mempool
        self.chain = []
        self.transactions = []
        self.create_block(proof=1, previous_hash="0")
        self.nodes = set()

    def create_block(self, proof, previous_hash):
        """
        the proof key is the nonce
        """
        block = {
            # include current hash
            'index': len(self.chain) + 1,
            # time_stamp will initially be different for other nodes for demo purposes
            # it will be corrected after the consensus protocol
            'time_stamp': str(datetime.datetime.now()),
            'proof': proof,
            'previous_hash': previous_hash,
            'transactions': self.transactions,
        }
        self.transactions = []
        self.chain.append(block)
        return block

    def get_previous_block(self):
        """
        Get the last block in the chain
        it is called get_previous_block  because when you add
        a new block to a blockchain, you will append it to the last block in the chain
        """
        return self.chain[-1]

    def proof_of_work(self, previous_proof):
        """ 
        Finds the nonce (proof) for a block
        Ensures that a blocks hash is valid and is below the target

        for demo purposes, hash_operation is implemented based only on the previous proof
        in production, this can be modified to hash the whole block instead
        """
        new_proof = 1
        check_proof = False
        while not check_proof:
            hash_operation = hashlib.sha256(str(new_proof**2 - previous_proof**2).encode()).hexdigest()
            if hash_operation[:4] == '0000':
                check_proof = True
            else:
                new_proof += 1
        return new_proof

    def hash(self, block):
        """
        returns hash of block
        """
        encoded_block = json.dumps(block, sort_keys=True).encode()
        return hashlib.sha256(encoded_block).hexdigest()

    def is_chain_valid(self, chain):
        """chain is part of the class so you don't need it as a parameter """
        previous_block = chain[0]
        block_index = 1
        while block_index < len(chain):
            block = chain[block_index]
            if block['previous_hash'] != self.hash(previous_block):
                # checks if the current block points to the previous block
                return False
            
            previous_proof = previous_block['proof']
            proof = block['proof']
            hash_operation = hashlib.sha256(str(proof**2 - previous_proof**2).encode()).hexdigest()
            if hash_operation[:4] != "0000":
                # checks if the hash of the current block is below the target
                return False

            previous_block = block
            block_index += 1
        return True

    def add_transaction(self, sender, receiver, amount):
        self.transactions.append({
                'sender': sender,
                'receiver': receiver,
                'amount': amount,
            })
        # returns the index of the new block to be added
        previous_block = self.get_previous_block()
        return previous_block['index'] + 1

    def add_node(self, address):
        """
        address = 'http://127.0.0.1:5000/'
        parsed_url = parced address
        parsed_url.netloc = '127.0.0.1:5000'
        """
        parsed_url = urlparse(address)
        self.nodes.add(parsed_url.netloc)

    def replace_chain(self):
        # find the chain in the network that is the longest and also valid
        network = self.nodes
        longest_chain = None
        max_length = len(self.chain)
        for node in network:
            response = requests.get(f'http://{node}/get_chain')
            if response.status_code == 200:
                length = response.json()['len']
                chain = response.json()['chain']
                if length > max_length and self.is_chain_valid(chain):
                    max_length = length
                    longest_chain = chain
        if longest_chain:
            self.chain = longest_chain
            return True
        return False


# Creates the Web app
app = Flask(__name__)
app.config['JSONIFY_PRETTYPRINT_REGULAR'] = False

# Create address for node in port 5000 (where the reward for mining comes from)
node_address = str(uuid4()).replace('-', '')

# Creates the blockchain
blockchain = Blockchain()

# add block to blockchain
@app.route('/mine_block', methods=['GET'])
def mine_block():
    previous_block = blockchain.get_previous_block()
    previous_proof = previous_block['proof']
    proof = blockchain.proof_of_work(previous_proof)
    previous_hash = blockchain.hash(previous_block)
    # Ben mined the block for {node_address} with 1 Ben coin as a reward
    blockchain.add_transaction(sender=node_address, receiver='Ben', amount=1)
    block = blockchain.create_block(proof, previous_hash)
    response = {
        'message': 'Congratulations, you just mined a block!',
        'index': block['index'],
        'time_stamp': block['time_stamp'],
        'proof': block['proof'],
        'previous_hash': block['previous_hash'],
        'transactions': block['transactions'],
    }
    return jsonify(response), 200

# get full blockchain
@app.route('/get_chain', methods=['GET'])
def get_chain():
    response = {
        'chain': blockchain.chain,
        'len': len(blockchain.chain),
    }
    return jsonify(response), 200

# check if blockchain is valid
@app.route('/is_valid', methods=['GET'])
def is_valid():
    valid = blockchain.is_chain_valid(blockchain.chain)
    response ={
        'valid': valid,
    }
    return jsonify(response), 200

# add a new trasaction to the block
@app.route('/add_transaction', methods=['POST'])
def add_transaction():
    json = request.get_json()
    transaction_keys = ['sender', 'receiver', 'amount']
    if not all(key in json for key in transaction_keys):
        return 'Some keys in the json are missing', 400
    # For demo purposes theres no signature but in application there is
    index = blockchain.add_transaction(json['sender'], json['receiver'], json['amount'])
    response = {
        'message': f'This transaction will be added to Block {index}',
    }
    return jsonify(response), 201

# add new node/s to the network
# this is done per node in the demo
@app.route('/connect_node', methods=['POST'])
def connect_node():
    json = request.get_json()
    nodes = json.get('nodes')
    if nodes is None:
        return "No node provided", 400
    for node in nodes:
        blockchain.add_node(node)
    response = {
        'message': 'Nodes have been successfully added to the network',
        'total_nodes': list(blockchain.nodes),
    }
    return jsonify(response), 201

# replacing the chain by the longest chain (if needed)
@app.route('/replace_chain', methods=['GET'])
def replace_chain():
    is_chain_replaced = blockchain.replace_chain()
    if is_chain_replaced:
        response = {
            'message': 'Chain replaced by the longest chain in the network',
            'chain': blockchain.chain,
        }
    else:
        response = {
            'message': 'Chain is the longest chain in the network',
            'chain': blockchain.chain,
        }
    return jsonify(response), 200


app.run(host='0.0.0.0', port=5000)