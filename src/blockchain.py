import hashlib
import json
import time
from typing import List, Dict, Any

class Block:
    def __init__(self, index: int, timestamp: float, data: Any, previous_hash: str, proof: int = 0):
        self.index = index
        self.timestamp = timestamp
        self.data = data
        self.previous_hash = previous_hash
        self.proof = proof
        self.hash = self.calculate_hash()

    def calculate_hash(self) -> str:
        block_string = json.dumps({
            "index": self.index,
            "timestamp": self.timestamp,
            "data": self.data,
            "previous_hash": self.previous_hash,
            "proof": self.proof
        }, sort_keys=True).encode()
        return hashlib.sha256(block_string).hexdigest()

    def __repr__(self):
        return f"Block(index={self.index}, hash={self.hash[:8]}, prev={self.previous_hash[:8]})"

class Blockchain:
    def __init__(self, difficulty: int = 4):
        self.chain: List[Block] = []
        self.difficulty = difficulty
        self.create_genesis_block()

    def create_genesis_block(self):
        genesis_block = Block(0, time.time(), "Genesis Block", "0")
        self.chain.append(genesis_block)

    @property
    def last_block(self) -> Block:
        return self.chain[-1]

    def add_block(self, block: Block) -> bool:
        if self.is_valid_new_block(block, self.last_block):
            self.chain.append(block)
            return True
        return False

    def proof_of_work(self, index: int, timestamp: float, data: Any, previous_hash: str) -> int:
        proof = 0
        while True:
            hash_val = self._calculate_hash_with_proof(index, timestamp, data, previous_hash, proof)
            if hash_val.startswith('0' * self.difficulty):
                return proof
            proof += 1

    def _calculate_hash_with_proof(self, index: int, timestamp: float, data: Any, previous_hash: str, proof: int) -> str:
        block_string = json.dumps({
            "index": index,
            "timestamp": timestamp,
            "data": data,
            "previous_hash": previous_hash,
            "proof": proof
        }, sort_keys=True).encode()
        return hashlib.sha256(block_string).hexdigest()

    def is_valid_new_block(self, new_block: Block, previous_block: Block) -> bool:
        if previous_block.index + 1 != new_block.index:
            return False
        if previous_block.hash != new_block.previous_hash:
            return False
        if new_block.hash != new_block.calculate_hash():
            return False
        if not new_block.hash.startswith('0' * self.difficulty):
            return False
        return True

    def is_chain_valid(self, chain: List[Block]) -> bool:
        for i in range(1, len(chain)):
            if not self.is_valid_new_block(chain[i], chain[i-1]):
                return False
        return True

    def to_dict(self):
        return [vars(b) for b in self.chain]
