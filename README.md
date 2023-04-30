# Ethereum Node Crawler

Crawls the Mainnet network and visualizes collected data. This repository includes backend, API and frontend for Ethereum network crawler. 

[Backend](./crawler) is based on [devp2p](https://github.com/ethereum/go-ethereum/tree/master/cmd/devp2p) tool. It tries to connect to discovered nodes, fetches info about them and creates a database. [API](./api) software reads raw node database, filters it, caches and serves as API. [Frontend](./frontend) is a web application which reads data from the API and visualizes them as a dashboard. 

Features:
- Advanced filtering, allows you to add filters for a customized dashboard
- Drilldown support, allows you to drill down the data to find interesting trends
- Network upgrade readiness overview
- Tabular view of all nodes crawled (along with IP Address and TCP/UDP Port)

### Theory

#### ENR (Ethereum Node Records)
An [ENR](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-778.md) is a data structure that represents a node in the Ethereum Consensus network. Nodes in the network have two roles - execution and consensus. Often a consensus layer is coupled with the task of validating transactions and proposing new blocks for the network. This is handled by the validating software of the stack.

#### Discovery Protocol
Discv4 and Discv5 are the two versions of the Node Discovery Protocol used in the application. One of the major contributions of this fork is to provide information about which node is a validator node and which isn't. It is not an evident process to determine this. Read more about discv5 protocol [here](https://github.com/ethereum/devp2p/blob/master/discv5/discv5-theory.md)

In the discovery implementation, a node queries another node using its ENR fields (such as IP and TCP) to set up a handshake. If the recipient node has no information regarding the calling node, a dial connection establishes an exchange of session keys to encrpyt the information flow. More details on types of packets are mentioned [here](https://github.com/ethereum/devp2p/blob/master/discv5/discv5-wire.md). Once the exchange of session keys is complete, write and read hello packets are transmitted to confirm the security of the session. The recipient node is then able to respond back to the. FindNode query that was initiatied earlier.

Of the many fields of the ENR, attnets record is used to determine whether a node is a validator node. If the attnets field of the Node record consists of non zero bytes, then with high confidence, it is a validator node. Read more on this [page](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/validator.md). The node discovery employs this to determine the validator status of the node, along with other information present on the ENR such as IP Address, Port, Record Public Key, Client Type, Language, Operating System.


## Usage
### Frontend 
#### Development
For local development with debugging, remoting, etc:
1. Copy `.env` into `.env.local` and replace the variables. 
1. And then `npm install` then `npm start`
1. Run tests to make sure the data processing is working good. `npm test`


### Crawler

#### Dependencies

- golang
- sqlite3

##### Country location

- `GeoLite2-Country.mmdb` file from [https://dev.maxmind.com/geoip/geolite2-free-geolocation-data?lang=en](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data?lang=en)
	- you will have to create an account to get access to this file

#### Development

```
cd crawler
go run .
```
Run crawler using `crawl` command. 
```
go run . crawl --table=nodetable
```
If you want to get the country that a Node is in you have to specify the location the geoIP database as well.

##### No GeoIP
```
crawler crawl --timeout 10m --table /path/to/database
```
##### With GeoIP

```
crawler crawl --timeout 10m --table /path/to/database --geoipdb GeoLite2-Country.mmdb
```

### Backend API

The API is using 2 databases. 1 of them is the raw data from the crawler and the other one is the API database.
Data will be moved from the crawler DB to the API DB regularly by this binary.
Make sure to start the crawler before the API if you intend to run them together during development.

#### Dependencies

- golang
- sqlite3

#### Development
```
go run ./ .
```

#### Sample Page (listing all nodes crawled)
![Nodes](https://github.com/VinayNR/node-crawler/blob/main/frontend/Nodes.png)
