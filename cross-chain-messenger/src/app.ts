import { ethers, BigNumber, Contract } from 'ethers';
import express from 'express';
import { JsonRpcProvider } from '@ethersproject/providers'
import { MerkleTree } from 'merkletreejs'
 import { getContractInterface } from '@eth-optimism/contracts'
import { getOEContract } from '@eth-optimism/sdk';
import * as dotenv from 'dotenv'
dotenv.config()

const optimismSDK = require("@eth-optimism/sdk")
const bodyParser = require('body-parser')

const app = express();
app.use(bodyParser.json())
const port = 41235;

const l1_provider = new ethers.providers.JsonRpcProvider(process.env.L1_RPC_URL)
const l2_provider = new ethers.providers.JsonRpcProvider(process.env.L2_RPC_URL)

const crossChainMessenger = new optimismSDK.CrossChainMessenger({
    l1ChainId: 31337,
    l2ChainId: 17,
    l1SignerOrProvider: l1_provider,
    l2SignerOrProvider: l2_provider
})

app.post('/storage_proof', async (req, res) => {
  console.log(req.body);
  try {
    let storageProof = await crossChainMessenger.getStorageProof(req.body["resolver_addr"], req.body["addr_slot"], {
      blockTag:'latest'
    });
    res.json(storageProof);
  } catch(e) {
    console.log("getStorageProof", e)
    res.status(400);
    res.send('getStorageProof failed');
  }
});

const ADDRESS_MANAGER_ADDRESS = '0xa6f73589243a6A7a9023b1Fa0651b1d89c177111';


const ZERO_ADDRESS = "0x" + "00".repeat(20);

export const loadContract = (
  name: string,
  address: string,
  provider: JsonRpcProvider
): Contract => {
  return new Contract(address, getContractInterface(name) as any, provider)
}

export const loadContractFromManager = async (
  name: string,
  Lib_AddressManager: Contract,
  provider: JsonRpcProvider
): Promise<Contract> => {
  const address = await Lib_AddressManager.getAddress(name)

  if (address === ZERO_ADDRESS) {
    throw new Error(
      `Lib_AddressManager does not have a record for a contract named: ${name}`
    )
  }

  return loadContract(name, address, provider)
}

// Instantiate the manager
const ovmAddressManager = loadContract('Lib_AddressManager', ADDRESS_MANAGER_ADDRESS, l1_provider);

interface StateRootBatchHeader {
    batchIndex: BigNumber
    batchRoot: string
    batchSize: BigNumber
    prevTotalElements: BigNumber
    extraData: string
}

async function getLatestStateBatchHeader(): Promise<{batch: StateRootBatchHeader, stateRoots: string[]}> {
    // Instantiate the state commitment chain
    const ovmStateCommitmentChain = await loadContractFromManager('StateCommitmentChain', ovmAddressManager, l1_provider);

    for(let endBlock = await l1_provider.getBlockNumber(); endBlock > 0; endBlock = Math.max(endBlock - 100, 0)) {
        const startBlock = Math.max(endBlock - 100, 1);
        const events: ethers.Event[] = await ovmStateCommitmentChain.queryFilter(
            ovmStateCommitmentChain.filters.StateBatchAppended(), startBlock, endBlock);
        if(events.length > 0) {
            const event = events[events.length - 1];
            const tx = await l1_provider.getTransaction(event.transactionHash);
            const [ stateRoots ] = ovmStateCommitmentChain.interface.decodeFunctionData('appendStateBatch', tx.data);
            return {
                batch: {
                    batchIndex: event.args?._batchIndex,
                    batchRoot: event.args?._batchRoot,
                    batchSize: event.args?._batchSize,
                    prevTotalElements: event.args?._prevTotalElements,
                    extraData: event.args?._extraData,
                },
                stateRoots,
            }
        }
    }
    throw Error("No state root batches found");
}

// app.post('/storage_proof', async (req, res) => {
//   const stateBatchHeader = await getLatestStateBatchHeader();
//   // The l2 block number we'll use is the last one in the state batch
//   const l2BlockNumber = stateBatchHeader.batch.prevTotalElements.add(stateBatchHeader.batch.batchSize);
//
//   // Construct a merkle proof for the state root we need
//   const elements = []
//   for (
//     let i = 0;
//     i < Math.pow(2, Math.ceil(Math.log2(stateBatchHeader.stateRoots.length)));
//     i++
//   ) {
//     if (i < stateBatchHeader.stateRoots.length) {
//       elements.push(stateBatchHeader.stateRoots[i])
//     } else {
//       elements.push(ethers.utils.keccak256('0x' + '00'.repeat(32)))
//     }
//   }
//   const hash = (el: Buffer | string): Buffer => {
//     return Buffer.from(ethers.utils.keccak256(el).slice(2), 'hex')
//   }
//   const leaves = elements.map((element) => {
//     return Buffer.from(element.slice(2), 'hex')
//   })
//   const index = elements.length - 1;
//   const tree = new MerkleTree(leaves, hash)
//   const treeProof = tree.getProof(leaves[index], index).map((element) => {
//     return element.data
//   });
//   const data = {
//         stateRoot: stateBatchHeader.stateRoots[index],
//         stateRootBatchHeader: stateBatchHeader.batch,
//         stateRootProof: {
//             index,
//             siblings: treeProof,
//         },
//   };
//   res.json(data);
// })

app.listen(port, () => {
  return console.log(`Express is listening at http://localhost:${port}`);
});

