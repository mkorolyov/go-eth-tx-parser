# Golang HTTP Service for Ethereum Transactions

This is a simple HTTP service built in Go that interacts with an Ethereum JSON RPC client to observe transactions for subscribed addresses. The service provides APIs to subscribe to an address, fetch transactions for a subscribed address, and retrieve the latest processed Ethereum block.

## Features

- Subscribe to an Ethereum address to monitor transactions.
- Fetch all transactions associated with a subscribed address.
- Get the current Ethereum block.
- In-memory storage for demonstration purposes.

## Usage

Run the service:

```bash
make run
```

The server will start on `localhost:8080` by default.

## API Endpoints

### 1. Subscribe to an Address

**Endpoint:** `/address/{address}/subscribe`

**Method:** `POST`

**Path Parameters:**
- `address`: Ethereum address to subscribe to (e.g., `0x1234...`)

**Response:**
- `200 OK` on success
- `500 Internal Server Error` if subscription fails

**Example:**

```bash
curl -X POST "http://localhost:8080/subscribe/0x1234567890abcdef1234567890abcdef12345678"
```

---

### 2. Get Transactions for an Address

**Endpoint:** `/transactions`

**Method:** `GET`

**Query Parameters:**
- `address`: Ethereum address to retrieve transactions for (e.g., `0x1234...`)

**Response:**
- `200 OK` with a JSON array of transactions
- `204 No Content` if no transactions are found for the address
- `500 Internal Server Error` if fetching transactions fails

**Example:**

```bash
curl -X GET "http://localhost:8080/transactions?address=0x1234567890abcdef1234567890abcdef12345678"
```

**Sample Response:**

```json
[
  {
    "hash": "0xabc123",
    "from": "0x1234567890abcdef1234567890abcdef12345678",
    "to": "0xabcdef1234567890abcdef1234567890abcdef12",
    "value": "1000000000000000000"
  }
]
```

---

### 3. Get Current Ethereum Block

**Endpoint:** `/current_block`

**Method:** `GET`

**Response:**
- `200 OK` with the current block number in JSON format
- `500 Internal Server Error` if fetching the block fails

**Example:**

```bash
curl -X GET "http://localhost:8080/current_block"
```

**Sample Response:**

```json
{
  "current_block": 12345678
}
```

## Notes

- This service uses in-memory storage (`NewInMemoryStorage`) for simplicity and demonstration purposes. For production, replace it with a persistent storage solution (e.g., a database).
- This services uses most naive http server implementation. For production, consider using a more complex approach, with proper middlewares, logging, and error handling.

## License

This project is licensed under the MIT License. Feel free to use and modify it as needed.

