import requests
from src.identity import Identity
import json

def run_demo():
    # 1. Initialize Identity
    print("Initializing Identity...")
    identity = Identity()
    pub_key = identity.get_public_key_string()
    
    # 2. Define Log Event
    event = "SYSTEM_STARTUP: All services initialized successfully."
    print(f"Signing event: {event}")
    
    # 3. Sign Event
    signature = identity.sign_event(event)
    
    # 4. Send to Node (assuming node is running on 5001)
    url = "http://localhost:5001/logs"
    payload = {
        "event": event,
        "public_key": pub_key,
        "signature": signature
    }
    
    print("Sending log to blockchain...")
    try:
        response = requests.post(url, json=payload)
        if response.status_code == 201:
            print("Successfully added log to blockchain!")
            print(json.dumps(response.json(), indent=2))
        else:
            print(f"Failed to add log: {response.text}")
    except Exception as e:
        print(f"Error: {e}. Is the node running on port 5001?")

if __name__ == "__main__":
    run_demo()
