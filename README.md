# A Toy Distributed Hash Table

Based on a [B-Tree](https://github.com/baziotis/golang-btree). This is really **NOT** for production
right now, and it doesn't ever intend to be.

There are multiple servers and a single client. The client receives keys from the user
and distributes them to the servers based on their modulo. Eventually, the goal
is to implement consistent hashing and Chord.

The DHT is persistent. If you restart it, you will be able to retrieve the values based on keys.
You have to delete the `.db` files created by the servers to clear it

## Usage

Commands:
- `GET`: Retrieve a value from the DHT based on a key. Example: `GET:10`
- `INSERT`: Insert a key-value pair to the DHT. Example: `INSERT:10:20`
- `DEL`: Delete a key (and its value). Example: `DEL:10`
- `EXIT`: Exit the client and terminate the servers.