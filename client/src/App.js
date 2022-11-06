import './App.css';

import {useState} from "react";

import {ethers} from "ethers";
import {CONTRACT_ABI, L1_REGISTRY_CONTRACT_ADDRESS} from "./contract";

function App() {
    const [connectedAccount, setConnectedAccount] = useState("");

    const nullModal = {open: false, functionSelectorName: "", response: ""}
    const [modal, setModal] = useState(nullModal);
    const [loading, setLoading] = useState("");

    const ensName = "paulio.eth";
    const txOptions = {ccipReadEnabled: true};

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
            handleError(err)
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
                setLoading("Owner")

                // Initialize the registry contract connection
                const registry = createRegistryContract(ethereum);

                // Generate the nameHash for the corresponding ENS name
                const node = ethers.utils.namehash(ensName);

                // Query the L1 Registry contract & let ethers handle the CCIP reversion
                const owner = await registry.owner(node, txOptions)
                setModal(
                    {
                        open: true,
                        functionSelectorName: "Owner()",
                        response: owner
                    }
                )
                setLoading("")
                return owner
            }
        } catch (err) {
            handleError(err)
        }
    }

    const getResolver = async () => {
        try {
            const {ethereum} = window;

            if (ethereum) {
                setLoading("Resolver")

                // Initialize the registry contract connection
                const registry = createRegistryContract(ethereum);

                // Generate the nameHash for the corresponding ENS name
                const node = ethers.utils.namehash(ensName);

                // Query the L1 Registry contract & let ethers handle the CCIP reversion
                const resolver = await registry.resolver(node, txOptions)
                setModal(
                    {
                        open: true,
                        functionSelectorName: "Resolver()",
                        response: resolver
                    }
                )
                setLoading("")
                return resolver
            }
        } catch (err) {
            handleError(err)
        }
    }

    const getTTL = async () => {
        try {
            const {ethereum} = window;

            if (ethereum) {
                setLoading("TTL")

                // Initialize the registry contract connection
                const registry = createRegistryContract(ethereum);

                // Generate the nameHash for the corresponding ENS name
                const node = ethers.utils.namehash(ensName);

                // Query the L1 Registry contract & let ethers handle the CCIP reversion
                const ttl = await registry.ttl(node, txOptions)
                setModal(
                    {
                        open: true,
                        functionSelectorName: "TTL()",
                        response: ttl
                    }
                )
                setLoading("")
                return ttl
            }
        } catch (err) {
            handleError(err)
        }
    }
    const handleError = (err) => {
        console.log("Handling error")
        console.log(err);
        console.log("opening modal")
        setModal(
            {
                open: true,
                functionSelectorName: "<Error/>",
                response: `${err}`
            }
        )
        setLoading("")
    }

    const loadingButton = (name, handler) => {
        if (loading === name) {
            return (
                <div className="button">
                    <div className="spin"/>
                </div>
            )
        }
        return (
            <div
                className="button"
                onClick={() => handler(ensName)}
            >
                Get {name}
            </div>
        )
    }

    return (
        <div className="App">
            {
                modal.open &&
                <div className="modal-box">
                    <div className="modal-container">
                        <h1>{modal.functionSelectorName}</h1>
                        <div className="left-align">
                            <h2 className="row">ENS Name: <p className="data">{ensName}</p></h2>
                            <h2 className="row">Response: <p className="data">{modal.response}</p></h2>
                        </div>
                        <div
                            className="button center"
                            onClick={() => setModal(nullModal)}
                        >
                            Close Modal
                        </div>
                    </div>
                </div>
            }

            {/* Header */}
            <div className="Header">
                <h1>Trustless L2 ENS Registry Rollup Solution</h1>

                {
                    connectedAccount === "" ?
                        <div
                            className="button connect-wallet"
                            onClick={connectWallet}
                        >
                            Connect Wallet
                        </div>
                        :
                        <div
                            className="button connect-wallet disabled"
                        >
                            {`${connectedAccount.substr(0, 6)}...${connectedAccount.substr(connectedAccount.length - 4)}`}
                        </div>
                }
            </div>

            {/* Body */}
            <div className="Body">
                <h2>Cross Chain Accessor Functions:</h2>
                {loadingButton("Owner", getOwner)}
                {loadingButton("Resolver", getResolver)}
                {loadingButton("TTL", getTTL)}
            </div>

            {/* Footer */}
            <div className="Footer">
                <h1>Made with ❤</h1>
            </div>
        </div>
    );
}

export default App;
