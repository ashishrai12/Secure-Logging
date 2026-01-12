import requests
from flask import Flask, jsonify, request
from src.blockchain import Blockchain, Block
from src.identity import Identity
import time
import threading

class LoggingNode:
    def __init__(self, host='0.0.0.0', port=5000):
        self.blockchain = Blockchain()
        self.nodes = set()
        self.app = Flask(__name__)
        self.host = host
        self.port = port
        self._setup_routes()

    def _setup_routes(self):
        @self.app.route('/logs', methods=['POST'])
        def add_log():
            values = request.get_json()
            required = ['event', 'public_key', 'signature']
            if not all(k in values for k in required):
                return 'Missing values', 400

            # Verify signature
            is_valid = Identity.verify_event(
                values['public_key'],
                values['event'],
                values['signature']
            )
            
            if not is_valid:
                return 'Invalid signature', 401

            # Prepare data for block
            data = {
                "event": values['event'],
                "public_key": values['public_key'],
                "signature": values['signature']
            }

            # Mine block
            last_block = self.blockchain.last_block
            index = last_block.index + 1
            timestamp = time.time()
            proof = self.blockchain.proof_of_work(index, timestamp, data, last_block.hash)
            
            new_block = Block(index, timestamp, data, last_block.hash, proof)
            if self.blockchain.add_block(new_block):
                # Broadcast to other nodes
                self.broadcast_block(new_block)
                return jsonify({"message": "Log added successfully", "block": vars(new_block)}), 201
            
            return 'Could not add block', 500

        @self.app.route('/chain', methods=['GET'])
        def get_chain():
            return jsonify({
                "chain": self.blockchain.to_dict(),
                "length": len(self.blockchain.chain)
            }), 200

        @self.app.route('/nodes/register', methods=['POST'])
        def register_nodes():
            values = request.get_json()
            node_urls = values.get('nodes')
            if node_urls is None:
                return "Error: Please supply a valid list of nodes", 400

            for node in node_urls:
                self.nodes.add(node)

            return jsonify({"message": "New nodes added", "total_nodes": list(self.nodes)}), 201

        @self.app.route('/nodes/resolve', methods=['GET'])
        def consensus():
            replaced = self.resolve_conflicts()
            if replaced:
                return jsonify({"message": "Our chain was replaced", "new_chain": self.blockchain.to_dict()}), 200
            else:
                return jsonify({"message": "Our chain is authoritative", "chain": self.blockchain.to_dict()}), 200

        @self.app.route('/receive-block', methods=['POST'])
        def receive_block():
            block_data = request.get_json()
            block = Block(
                block_data['index'],
                block_data['timestamp'],
                block_data['data'],
                block_data['previous_hash'],
                block_data['proof']
            )
            if self.blockchain.add_block(block):
                return 'Block accepted', 201
            return 'Block rejected', 400

    def resolve_conflicts(self):
        """
        Consensus algorithm, resolves conflicts by replacing our chain with the longest one in the network.
        """
        neighbours = self.nodes
        new_chain = None
        max_length = len(self.blockchain.chain)

        for node in neighbours:
            try:
                response = requests.get(f'http://{node}/chain')
                if response.status_code == 200:
                    length = response.json()['length']
                    chain_data = response.json()['chain']

                    # Convert back to Block objects
                    temp_chain = []
                    for b in chain_data:
                        temp_chain.append(Block(b['index'], b['timestamp'], b['data'], b['previous_hash'], b['proof']))

                    if length > max_length and self.blockchain.is_chain_valid(temp_chain):
                        max_length = length
                        new_chain = temp_chain
            except:
                continue

        if new_chain:
            self.blockchain.chain = new_chain
            return True

        return False

    def broadcast_block(self, block: Block):
        for node in self.nodes:
            try:
                requests.post(f'http://{node}/receive-block', json=vars(block))
            except:
                continue

    def run(self):
        self.app.run(host=self.host, port=self.port)

if __name__ == '__main__':
    node = LoggingNode()
    node.run()
