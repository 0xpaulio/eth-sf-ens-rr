// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import {Lib_AddressResolver} from "optimism/contracts/libraries/resolver/Lib_AddressResolver.sol";
import {Lib_OVMCodec} from "optimism/contracts/libraries/codec/Lib_OVMCodec.sol";
import {Lib_SecureMerkleTrie} from "optimism/contracts/libraries/trie/Lib_SecureMerkleTrie.sol";
import {StateCommitmentChain} from "optimism/contracts/L1/rollup/StateCommitmentChain.sol";
import {Lib_RLPReader} from "optimism/contracts/libraries/rlp/Lib_RLPReader.sol";
import {Lib_BytesUtils} from "optimism/contracts/libraries/utils/Lib_BytesUtils.sol";

contract OptimismHelper is Lib_AddressResolver {
    /**
     * @dev Struct used to store an Optimsim L2 State Proof. The state proof is used to verify the inclusion of a specific record in an L2 account.
     * @param stateRoot .
     * @param stateRootBatchHeader .
     * @param stateRootProof .
     * @param stateTrieWitness .
     * @param storageTrieWitness .
     */
    struct L2StateProof {
        bytes32 stateRoot;
        Lib_OVMCodec.ChainBatchHeader stateRootBatchHeader;
        Lib_OVMCodec.ChainInclusionProof stateRootProof;
        bytes stateTrieWitness;
        bytes storageTrieWitness;
    }
    
    /// @dev Error to raise when the target contract ("account") does not exist in the state root.
    error AccountDNE(); 
    
    /// @dev Error to raise when the target storage slot does not exist in the corresponding account, under the latest state root.
    error StorageDNE(); 
    
    /// @dev Error to raise when the state root of the L2StateProof is invalid.
    error InvalidStateRoot(); 

    /**
     * @dev Modifier used to verify the provided state root proof.
     * @param _proof The L2 state proof for the corresponding CCIP lookup verification callback.
     */
    modifier validStateRootProof(L2StateProof memory _proof) {
        if (!verifyStateRootProof(_proof))
            revert InvalidStateRoot();
        _;
    }

    /*//////////////////////////////////////////////////////////////
                            CONSTRUCTOR
    //////////////////////////////////////////////////////////////*/

    /// @dev Constructs a Helper contract for verifying and interfacing with data stored on the optimism L2.
    constructor(address _ovmAddressManager) 
        Lib_AddressResolver(_ovmAddressManager) 
        { }

    /*//////////////////////////////////////////////////////////////
                      INTERNAL HELPER FUNCTIONS
    //////////////////////////////////////////////////////////////*/

    function verifyStateRootProof(L2StateProof memory proof)
        internal
        view
        returns (bool)
    {
        StateCommitmentChain ovmStateCommitmentChain = StateCommitmentChain(
            resolve("StateCommitmentChain")
        );
        return
            ovmStateCommitmentChain.verifyStateCommitment(
                proof.stateRoot,
                proof.stateRootBatchHeader,
                proof.stateRootProof
            );
    }

    function getStorageValue(
        address target,
        bytes32 slot,
        L2StateProof memory proof
    ) internal view returns (bytes32) {
        (
            bool exists,
            bytes memory encodedResolverAccount
        ) = Lib_SecureMerkleTrie.get(
                abi.encodePacked(target),
                proof.stateTrieWitness,
                proof.stateRoot
            );
        if (!exists)
            revert AccountDNE();
            
        Lib_OVMCodec.EVMAccount memory account = Lib_OVMCodec.decodeEVMAccount(
            encodedResolverAccount
        );
        (bool storageExists, bytes memory retrievedValue) = Lib_SecureMerkleTrie
            .get(
                abi.encodePacked(slot),
                proof.storageTrieWitness,
                account.storageRoot
            );
        if (!storageExists)
            revert StorageDNE();

        return toBytes32PadLeft(Lib_RLPReader.readBytes(retrievedValue));
    }

    function toBytes32PadLeft(bytes memory _bytes)
        internal
        pure
        returns (bytes32)
    {
        bytes32 ret;
        uint256 len = _bytes.length <= 32 ? _bytes.length : 32;
        assembly {
            ret := shr(mul(sub(32, len), 8), mload(add(_bytes, 32)))
        }
        return ret;
    }
}