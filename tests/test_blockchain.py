import pytest
import time
from src.blockchain import Blockchain, Block
from src.identity import Identity

def test_blockchain_genesis():
    bc = Blockchain()
    assert len(bc.chain) == 1
    assert bc.chain[0].data == "Genesis Block"

def test_add_block():
    bc = Blockchain(difficulty=1)
    last_block = bc.last_block
    index = last_block.index + 1
    data = "Test Log"
    timestamp = time.time()
    proof = bc.proof_of_work(index, timestamp, data, last_block.hash)
    
    new_block = Block(index, timestamp, data, last_block.hash, proof)
    assert bc.add_block(new_block) is True
    assert len(bc.chain) == 2

def test_invalid_block():
    bc = Blockchain(difficulty=1)
    new_block = Block(1, time.time(), "Fake Block", "wrong_hash", 0)
    assert bc.add_block(new_block) is False

def test_identity_signing():
    identity = Identity()
    data = "Important System Event"
    signature = identity.sign_event(data)
    pub_key = identity.get_public_key_string()
    
    assert Identity.verify_event(pub_key, data, signature) is True
    assert Identity.verify_event(pub_key, "Tampered Data", signature) is False

def test_chain_validity():
    bc = Blockchain(difficulty=1)
    for i in range(3):
        last_block = bc.last_block
        index = last_block.index + 1
        data = f"Log {i}"
        ts = time.time()
        proof = bc.proof_of_work(index, ts, data, last_block.hash)
        bc.add_block(Block(index, ts, data, last_block.hash, proof))
    
    assert bc.is_chain_valid(bc.chain) is True
    
    # Tamper with the chain
    bc.chain[1].data = "Hacked"
    assert bc.is_chain_valid(bc.chain) is False
