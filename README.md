# Leonard DAO Vinci

## Project Description
Excited by the prospect of computationally generated art, we saw the opportunity for a project which combines this with the decentralisation potential which blockchain technology brings. Running our project as a DAO offers opportunities to reward our users and community. This is done with profit-sharing in the sales of the artwork on [OpenSea](https://opensea.io/), without middlemen being involved. We wanted to build a project that would be engaging for new users to the blockchain space, and would allow them to enter the world of crypto without having to go through an exchange.

## Tech Stack
- **UI**: ReactJS, Rimble, web3.js
- **Middleware**: Go, NodeJs
- **Art Generators**: NodeJs, Python
- **Ethereum**: Zeppelin's ZepKit, Go Ethereum, Solidity, Truffle, OpenZeppelin, Infura, OpenSea

# Running the project
## User Interface
From the root directory run the following, long-running, command:s
```
cd client
npm install
npm run start
```

## Server for communicating with the OpenSea API
```
cd openseaIntegration
npm install
. .env; node index.js
```

Note: make sure that you are running Node 8.11.4

The server handles requests on the URL `http://localhost:3001/auction/token/<tokenId>`
