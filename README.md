# Secure-Logging: Decentralized and Secure Blockchain Solution

A blockchain-based solution for tamper-proof logging events. This project ensures that log events are immutable once stored and can be verified using digital signatures.

## Features

- **Decentralized Ledger**: Uses a blockchain to store log events across multiple nodes.
- **Tamper-proof**: Each block is linked via hashes and secured by Proof of Work (PoW).
- **Secure Authentication**: Every log entry is signed using RSA digital signatures (Private/Public key pair).
- **P2P Consensus**: Longest-chain consensus algorithm to resolve conflicts between nodes.
- **REST API**: Simple interface to submit logs, view the chain, and manage nodes.
- **CI/CD**: Integrated GitHub Actions for automated testing.

## Getting Started

### Prerequisites

- Python 3.9+
- Docker & Docker Compose (optional for local network simulation)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/ashishrai12/Secure-Logging.git
   cd Secure-Logging
   ```

2. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

### Running Locally

To start a single node:
```bash
python src/node.py
```
The node will be available at `http://localhost:5000`.

### Running with Docker (Network of Nodes)

To simulate a decentralized network:
```bash
docker-compose up --build
```
This starts two nodes:
- Node 1: `http://localhost:5001`
- Node 2: `http://localhost:5002`

### Usage Example

You can use the provided client script to sign and submit a log:
```bash
# Make sure a node is running on 5001
export PYTHONPATH=$PYTHONPATH:.
python src/client.py
```

## API Documentation

- `POST /logs`: Submit a new log event.
  - Body: `{"event": "str", "public_key": "PEM_str", "signature": "Base64_str"}`
- `GET /chain`: Retrieve the full blockchain.
- `POST /nodes/register`: Register new nodes in the network.
  - Body: `{"nodes": ["host1:port", "host2:port"]}`
- `GET /nodes/resolve`: Trigger consensus algorithm to sync with the network.

## Testing

Run tests using pytest:
```bash
pytest tests/
```

## CI/CD Workflow

This project includes a GitHub Actions workflow in `.github/workflows/ci.yml` that automatically runs tests on every push to the `main` branch.
