# Ignite Backend
​
Ignite Backend is a simple backend solution that allows indexing and retrieving data from a Cosmos blockchain through emitted events.

## Requirements
​
- [Go](https://go.dev/)
- [Protocol Buffer Compiler (protoc)](https://grpc.io/docs/protoc-installation/)
- [PostgreSQL](https://www.postgresql.org/)
- [GNU make utility](https://www.gnu.org/software/make/)
​
## Initial setup for development
​
The backend requires a PostgreSQL database server running.
​
First, create a "backend" database:
​
```bash
createdb --no-password backend
```
​
The required database tables will be created automatically by the `collector` the first time it is run.
​
Compile Ignite's backend by running `make` from the repository's root directory. The binary is generated
inside the `./bin` folder.
​
The next step is to start the `collector` service that will fetch all the transactions and events
starting from the first block until the current block height and populate the database:
​
```bash
bin/ignite-backend collector start --database-name backend --rpc-address IGNITE_CHAIN_ADDRESS -P sslmode=disable --log-level debug
```
​
Once the service is run it will keep collecting transactions as new blocks are generated.
​
Finally, run the `api` service to start the gRPC server:
​
```bash
bin/ignite-backend api start --database-name backend -P sslmode=disable --log-level debug
```