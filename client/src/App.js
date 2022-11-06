import logo from './logo.svg';
import './App.css';

import {useState} from "react";

import {ethers} from "ethers";
import {
  L1_REGISTRY_CONTRACT_ADDRESS,
  L2_REGISTRY_CONTRACT_ADDRESS,
  CONTRACT_ABI
} from "./contract";

function App() {
  const [connectedAccount, setConnectedAccount] = useState("");

  const ensName = "paulio.eth";
  const txOptions = { ccipReadEnabled: true };

  const connectWallet = async () => {
    try {
      const {ethereum} = window;

      if (!ethereum) {
        alert("No wallet connect found")
        return
      }

      const accounts = await ethereum.request({method: "eth_requestAccounts"});
      setConnectedAccount(accounts[0]);
    } catch (err) {
      console.log(err);
    }
  }

  const createRegistryContract = (ethereum) => {
    return new ethers.Contract(
        L1_REGISTRY_CONTRACT_ADDRESS,
        CONTRACT_ABI,
        (new ethers.providers.Web3Provider(ethereum))
    );
  }

  const getOwner = async (ensName) => {
    try {
      const {ethereum} = window;

      if (ethereum) {
        // Initialize the registry contract connection
        const registry = createRegistryContract(ethereum);

        // Generate the nameHash for the corresponding ENS name
        const node = ethers.utils.namehash(ensName);

        // Query the L1 Registry contract & let ethers handle the CCIP reversion
        return await registry["owner(bytes32)"](node, txOptions)
      }
    } catch (err) {
      console.log(err);
    }
  }

  const getResolver = async () => {
    try {
      const {ethereum} = window;

      if (ethereum) {
        // Initialize the registry contract connection
        const registry = createRegistryContract(ethereum);

        // Generate the nameHash for the corresponding ENS name
        const node = ethers.utils.namehash(ensName);

        // Query the L1 Registry contract & let ethers handle the CCIP reversion
        return await registry["resolver(bytes32)"](node, txOptions)
      }
    } catch (err) {
      console.log(err);
    }
  }

  const getTTL = async () => {
    try {
      const {ethereum} = window;

      if (ethereum) {
        // Initialize the registry contract connection
        const registry = createRegistryContract(ethereum);

        // Generate the nameHash for the corresponding ENS name
        const node = ethers.utils.namehash(ensName);

        // Query the L1 Registry contract & let ethers handle the CCIP reversion
        return await registry["ttl(bytes32)"](node, txOptions)
      }
    } catch (err) {
      console.log(err);
    }
  }

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default App;
