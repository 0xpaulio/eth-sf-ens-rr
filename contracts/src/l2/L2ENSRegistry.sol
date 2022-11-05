// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import "src/l1/interfaces/ENS.sol";

/**
 * The ENS registry contract.
 */
contract L2ENSRegistry is ENS {
    struct Record {
        address owner;
        address resolver;
        uint64 ttl;
    }

    mapping(bytes32 => Record) records;
    mapping(address => mapping(address => bool)) operators;

    // Permits modifications only by the owner of the specified node.
    modifier authorised(bytes32 node) {
        address owner_ = records[node].owner;
        require(owner_ == msg.sender || operators[owner_][msg.sender]);
        _;
    }

    /**
     * @dev Constructs a new ENS registry.
     */
    constructor() {
        records[0x0].owner = msg.sender;
    }

    // DEBUGGING STUFF
    function getSLO(bytes32 _pos) 
        external 
        view
        returns (bytes32 data_) 
    {
        assembly {
            data_ := sload(_pos)
        }
    } 
    
    event lbs_r(bytes32);

    function getRecordSLO(bytes32 _node) 
        external 
        returns (bytes32) 
    {
        Record storage r_ = records[_node];

        bytes32 pos_;
        
        assembly {
            pos_ := r_.slot
        }

        emit lbs_r(pos_);
        return pos_;
    } 

    /**
     * @dev Sets the record for a node.
     * @param _node The node to update.
     * @param _owner The address of the new owner.
     * @param _resolver The address of the resolver.
     * @param _ttl The TTL in seconds.
     */
    function setRecord(
        bytes32 _node,
        address _owner,
        address _resolver,
        uint64  _ttl
    ) external virtual override {
        setOwner(_node, _owner);
        _setResolverAndTTL(_node, _resolver, _ttl);
    }

    /**
     * @dev Sets the record for a subnode.
     * @param _node The parent node.
     * @param _label The hash of the label specifying the subnode.
     * @param _owner The address of the new owner.
     * @param _resolver The address of the resolver.
     * @param _ttl The TTL in seconds.
     */
    function setSubnodeRecord(
        bytes32 _node,
        bytes32 _label,
        address _owner,
        address _resolver,
        uint64 _ttl
    ) external virtual override {
        bytes32 subnode_ = setSubnodeOwner(_node, _label, _owner);
        _setResolverAndTTL(subnode_, _resolver, _ttl);
    }

    /**
     * @dev Transfers ownership of a node to a new address. May only be called by the current owner of the node.
     * @param _node The node to transfer ownership of.
     * @param _owner The address of the new owner.
     */
    function setOwner(bytes32 _node, address _owner)
        public
        virtual
        override
        authorised(_node)
    {
        _setOwner(_node, _owner);
        emit Transfer(_node, _owner);
    }

    /**
     * @dev Transfers ownership of a subnode keccak256(node, label) to a new address. May only be called by the owner of the parent node.
     * @param _node The parent node.
     * @param _label The hash of the label specifying the subnode.
     * @param _owner The address of the new owner.
     */
    function setSubnodeOwner(
        bytes32 _node,
        bytes32 _label,
        address _owner
    ) public virtual override authorised(_node) returns (bytes32) {
        bytes32 subnode_ = keccak256(abi.encodePacked(_node, _label));
        _setOwner(subnode_, _owner);
        emit NewOwner(_node, _label, _owner);
        return subnode_;
    }

    /**
     * @dev Sets the resolver address for the specified node.
     * @param _node The node to update.
     * @param _resolver The address of the resolver.
     */
    function setResolver(bytes32 _node, address _resolver)
        public
        virtual
        override
        authorised(_node)
    {
        emit NewResolver(_node, _resolver);
        records[_node].resolver = _resolver;
    }

    /**
     * @dev Sets the TTL for the specified node.
     * @param _node The node to update.
     * @param _ttl The TTL in seconds.
     */
    function setTTL(bytes32 _node, uint64 _ttl)
        public
        virtual
        override
        authorised(_node)
    {
        emit NewTTL(_node, _ttl);
        records[_node].ttl = _ttl;
    }

    /**
     * @dev Enable or disable approval for a third party ("operator") to manage
     *  all of `msg.sender`'s ENS records. Emits the ApprovalForAll event.
     * @param _operator Address to add to the set of authorized operators.
     * @param _approved True if the operator is approved, false to revoke approval.
     */
    function setApprovalForAll(address _operator, bool _approved)
        external
        virtual
        override
    {
        operators[msg.sender][_operator] = _approved;
        emit ApprovalForAll(msg.sender, _operator, _approved);
    }

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
        address addr = records[_node].owner;
        if (addr == address(this)) {
            return address(0x0);
        }

        return addr;
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
        return records[_node].resolver;
    }

    /**
     * @dev Returns the TTL of a node, and any records associated with it.
     * @param _node The specified node.
     * @return ttl of the node.
     */
    function ttl(bytes32 _node) public view virtual override returns (uint64) {
        return records[_node].ttl;
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
        return records[_node].owner != address(0x0);
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
        return operators[_owner][_operator];
    }

    function _setOwner(bytes32 _node, address _owner) internal virtual {
        records[_node].owner = _owner;
    }

    function _setResolverAndTTL(
        bytes32 _node,
        address _resolver,
        uint64 _ttl
    ) internal {
        if (_resolver != records[_node].resolver) {
            records[_node].resolver = _resolver;
            emit NewResolver(_node, _resolver);
        }

        if (_ttl != records[_node].ttl) {
            records[_node].ttl = _ttl;
            emit NewTTL(_node, _ttl);
        }
    }
}
