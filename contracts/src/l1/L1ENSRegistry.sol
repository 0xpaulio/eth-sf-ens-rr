// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import {ENS} from "src/l1/interfaces/ENS.sol";

import {OptimismHelper} from "src/l1/types/OptimismHelper.sol";

/**
 * The ENS registry contract.
 */
contract L1ENSRegistry is 
    ENS, 
    
    OptimismHelper {

    /*//////////////////////////////////////////////////////////////
                                ERRORS
    //////////////////////////////////////////////////////////////*/

    /**
     * @dev EIP-5559 - Error to raise when mutations are being deferred to an L2.
     * @param chainId Chain ID to perform the deferred mutation to.
     * @param contractAddress Contract Address at which the deferred mutation should transact with.
     */
    error StorageHandledByL2(
        uint256 chainId, 
        address contractAddress
    );

    /**
     * @dev EIP-3668 - Error to raise when a CCIP Request is required to complete the query.
     * @param sender The address of the smartcontract that triggers the CCIP request.
     * @param urls A list of gateway urls that can be used to resolve the CCIP request.
     * @param callData The calldata to pass to the CCIP gateway.
     * @param callbackFunction The function selector of the function to verify & decode the response from the CCIP gateway.
     * @param extraData Extra data to be passed to the calback function.
     */
    error OffchainLookup(
        address sender,
        string[] urls,
        bytes callData,
        bytes4 callbackFunction,
        bytes extraData
    );

    /*//////////////////////////////////////////////////////////////
                              CONSTANTS
    //////////////////////////////////////////////////////////////*/

    /**
     * @dev the storage slot of the records mapping in the L2 registry.
     * mapping(bytes32 => Record) records;
     */
    uint256 constant public SLO__L2_REGISTRY__RECORDS = 0;
    
    /**
     * @dev the storage slot of the operators mapping in the L2 registry.
     * mapping(address => mapping(address => bool)) operators;
     */
    uint256 constant public SLO__L2_REGISTRY__OPERATORS = 1;

    /*//////////////////////////////////////////////////////////////
                              VARIABLES
    //////////////////////////////////////////////////////////////*/

    /// @dev the chain id of the L2 Registry
    uint64 public L2_REGISTRY_CHAIN_ID;

    /// @dev the contract address of the L2 Registry
    address public L2_REGISTRY_CONTRACT_ADDRESS;

    /// @dev the list of gateway urls that can handle the  contract address of the L2 Registrar 
    string[] public gatewayUrls;

    /*//////////////////////////////////////////////////////////////
                             CONSTRUCTOR
    //////////////////////////////////////////////////////////////*/

    /** 
     * @dev Constructs a new "entry gateway" ENS registry, where all data is stored in and resolved from the L2. 
     * This contract uses: 
     *      - EIP-3668 for cross-chain reads 
     *      - EIP-5559 for async cross-chain mutations.
     */
    constructor(
        uint64 _chainId, address _l2Registrar, 
        address _ovmAddressManager, 
        string[] memory _gatewayUrls
    )
        OptimismHelper(_ovmAddressManager) 
    {
        L2_REGISTRY_CHAIN_ID = _chainId;
        L2_REGISTRY_CONTRACT_ADDRESS = _l2Registrar;

        gatewayUrls = _gatewayUrls;
    }

    /*//////////////////////////////////////////////////////////////
                       PUBLIC MUTATOR FUNCTIONS
    //////////////////////////////////////////////////////////////*/

    /**
     * @dev Sets the record for a node.
     * @param node The node to update.
     * @param owner The address of the new owner.
     * @param resolver The address of the resolver.
     * @param ttl The TTL in seconds.
     */
    function setRecord(
        bytes32 node,
        address owner,
        address resolver,
        uint64 ttl
    ) external virtual override {
        _writeDeferral();
    }

    /**
     * @dev Sets the record for a subnode.
     * @param node The parent node.
     * @param label The hash of the label specifying the subnode.
     * @param owner The address of the new owner.
     * @param resolver The address of the resolver.
     * @param ttl The TTL in seconds.
     */
    function setSubnodeRecord(
        bytes32 node,
        bytes32 label,
        address owner,
        address resolver,
        uint64 ttl
    ) external virtual override {
        _writeDeferral();
    }

    /**
     * @dev Transfers ownership of a node to a new address. May only be called by the current owner of the node.
     * @param node The node to transfer ownership of.
     * @param owner The address of the new owner.
     */
    function setOwner(bytes32 node, address owner)
        public
        virtual
        override
    {
        _writeDeferral();
    }

    /**
     * @dev Transfers ownership of a subnode keccak256(node, label) to a new address. May only be called by the owner of the parent node.
     * @param node The parent node.
     * @param label The hash of the label specifying the subnode.
     * @param owner The address of the new owner.
     */
    function setSubnodeOwner(
        bytes32 node,
        bytes32 label,
        address owner
    ) public virtual override returns (bytes32) {
        _writeDeferral();
    }

    /**
     * @dev Sets the resolver address for the specified node.
     * @param node The node to update.
     * @param resolver The address of the resolver.
     */
    function setResolver(bytes32 node, address resolver)
        public
        virtual
        override
    {
        _writeDeferral();
    }

    /**
     * @dev Sets the TTL for the specified node.
     * @param node The node to update.
     * @param ttl The TTL in seconds.
     */
    function setTTL(bytes32 node, uint64 ttl)
        public
        virtual
        override
    {
        _writeDeferral();
    }

    /**
     * @dev Enable or disable approval for a third party ("operator") to manage
     *  all of `msg.sender`'s ENS records. Emits the ApprovalForAll event.
     * @param operator Address to add to the set of authorized operators.
     * @param approved True if the operator is approved, false to revoke approval.
     */
    function setApprovalForAll(address operator, bool approved)
        external
        virtual
        override
    {
        _writeDeferral();
    }

    /*//////////////////////////////////////////////////////////////
                        PUBLIC ACCESSOR FUNCTIONS
    //////////////////////////////////////////////////////////////*/

    /**
     * @dev Returns the address that owns the specified node.
     * @param _node The specified node.
     * @return address of the owner.
     */
    function owner(bytes32 _node)
        public
        view
        virtual
        override
        returns (address)
    {
        revert OffchainLookup(
            address(this),
            gatewayUrls,
            msg.data,
            this.ownerWithProof.selector,
            abi.encode(_node)
        );
    }
    
    function ownerWithProof(bytes32 _node, L2StateProof memory _proof)
        public
        view
        virtual
        validStateRootProof(_proof)
        returns (address)
    {
        // Calculate the SLO of the corresponding Perm
        bytes32 slot_ = keccak256(
            abi.encodePacked(
                _node, 
                SLO__L2_REGISTRY__RECORDS
            )
        );

        // Grab the Perm flag with the calculated SLO
        bytes32 value = getStorageValue(
            L2_REGISTRY_CONTRACT_ADDRESS, 
            slot_, 
            _proof
        );

        address addr = address(uint160(uint256(value)));
        return addr == address(this) ? address(0x0) : addr;
    }

    /**
     * @dev Returns the address of the resolver for the specified node.
     * @param _node The specified node.
     * @return address of the resolver.
     */
    function resolver(bytes32 _node)
        public
        view
        virtual
        override
        returns (address)
    {
        revert OffchainLookup(
            address(this),
            gatewayUrls,
            msg.data,
            this.resolverWithProof.selector,
            abi.encode(_node)
        );
    }
    
    function resolverWithProof(bytes32 _node, L2StateProof memory _proof)
        public
        view
        virtual
        validStateRootProof(_proof)
        returns (address)
    {
        // Calculate the SLO of the corresponding Perm
        bytes32 slot_ = keccak256(
            abi.encodePacked(
                _node, 
                SLO__L2_REGISTRY__RECORDS
            )
        );

        // Shift the slot_ pos to get resolver
        slot_ = bytes32(
            uint256(slot_) + uint256(0x1)
        );

        // Grab the Perm flag with the calculated SLO
        bytes32 value = getStorageValue(
            L2_REGISTRY_CONTRACT_ADDRESS, 
            slot_, 
            _proof
        );
        return address(uint160(uint256(value)));
    }

    /**
     * @dev Returns the TTL of a node, and any records associated with it.
     * @param _node The specified node.
     * @return ttl of the node.
     */
    function ttl(bytes32 _node) public view virtual override returns (uint64) {   
        revert OffchainLookup(
            address(this),
            gatewayUrls,
            msg.data,
            this.ttlWithProof.selector,
            abi.encode(_node)
        );
    }
    
    function ttlWithProof(bytes32 _node, L2StateProof memory _proof) 
        public 
        view 
        virtual 
        validStateRootProof(_proof)
        returns (uint64) 
    {
        // Calculate the SLO of the corresponding Perm
        bytes32 slot_ = keccak256(
            abi.encodePacked(
                _node, 
                SLO__L2_REGISTRY__RECORDS
            )
        );

        // Shift the slot_ pos to get resolver
        slot_ = bytes32(
            uint256(slot_) + uint256(0x1)
        );

        // Grab the Perm flag with the calculated SLO
        bytes32 value = getStorageValue(
            L2_REGISTRY_CONTRACT_ADDRESS, 
            slot_, 
            _proof
        );

        // shift the data to grab the ttl
        uint64 ttl_;
        assembly {
            ttl_ := shr(0xa0, value)
        }
        return ttl_;
    }

    /**
     * @dev Returns whether a record has been imported to the registry.
     * @param _node The specified node.
     * @return Bool if record exists
     */
    function recordExists(bytes32 _node)
        public
        view
        virtual
        override
        returns (bool)
    {        
        revert OffchainLookup(
            address(this),
            gatewayUrls,
            msg.data,
            this.recordExistsWithProof.selector,
            abi.encode(_node)
        );
    }
     
    function recordExistsWithProof(bytes32 _node, L2StateProof memory _proof)
        public
        view
        virtual
        validStateRootProof(_proof)
        returns (bool)
    {
        return ownerWithProof(_node, _proof) != address(0);
    }

    /**
     * @dev Query if an address is an authorized operator for another address.
     * @param _owner The address that owns the records.
     * @param _operator The address that acts on behalf of the owner.
     * @return True if `operator` is an approved operator for `owner`, false otherwise.
     */
    function isApprovedForAll(address _owner, address _operator)
        external
        view
        virtual
        override
        returns (bool)
    {
        revert OffchainLookup(
            address(this),
            gatewayUrls,
            msg.data,
            this.isApprovedForAllWithProof.selector,
            abi.encode(_owner, _operator)
        );
    }
    
    function isApprovedForAllWithProof(address _owner, address _operator, L2StateProof memory _proof)
        external
        view
        virtual
        validStateRootProof(_proof)
        returns (bool approved_)
    {
        // Calculate the SLO of the corresponding Perm
        bytes32 slot_ = keccak256(
            abi.encodePacked(
                _operator,
                keccak256(
                    abi.encodePacked(
                        _owner, 
                        SLO__L2_REGISTRY__OPERATORS
                    )
                )
            )
        );

        // Grab the Perm flag with the calculated SLO
        bytes32 value = getStorageValue(
            L2_REGISTRY_CONTRACT_ADDRESS, 
            slot_, 
            _proof
        );

        assembly {
            approved_ := value
        }
        return approved_;
    }

    /*//////////////////////////////////////////////////////////////
                            ENS CCIP LOGIC
    //////////////////////////////////////////////////////////////*/

    /**
     * @notice Builds an OffchainLookup error.
     * @param callData The calldata for the corresponding lookup.
     * @return Always reverts with an OffchainLookup error.
     */
    function _offChainLookup(bytes4 callbackFunction, bytes calldata callData) private view returns(bytes memory) {
        // // TODO: implement this
        // bytes memory extraData_ = new bytes(callData.length);
        // assembly {
        //     extraData_ := add(callData, 0x04)
        // }

        revert OffchainLookup(
            address(this),
            gatewayUrls,
            callData,
            callbackFunction,
            abi.encode(callData[4:])
        );
    }

    /*//////////////////////////////////////////////////////////////
                           ENS CCWDP LOGIC
    //////////////////////////////////////////////////////////////*/

    /**
     * @dev Builds write deferral StorageHandledByL2 reversion.
     */
    function _writeDeferral() internal view {
        revert StorageHandledByL2(
            L2_REGISTRY_CHAIN_ID, 
            L2_REGISTRY_CONTRACT_ADDRESS
        );
    }
}
