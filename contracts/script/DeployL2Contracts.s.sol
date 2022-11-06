// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import "forge-std/Script.sol";

import {L2ENSRegistry} from "src/l2/L2ENSRegistry.sol";

contract ContractScript is Script {
    L2ENSRegistry l2Registry;

    address _owner = 0x84C970BFcD59a0e98eC6f13Cbdf24AA1a741f033;
    address _goerliENSResolver = 0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e;    

    string[] _gatewayUrls;

    function setUp() public {
        _gatewayUrls = new string[](1);
        _gatewayUrls[0] = "http://localhost:3000/";
    }

    function run() public {
        vm.startBroadcast();

        // Deploy the L2 Registry
        l2Registry = new L2ENSRegistry();

        // Register the eth domain
        l2Registry.setSubnodeRecord(
            0x0,
            keccak256("eth"),
            _owner,
            _goerliENSResolver,
            64
        );

        bytes32 _nullRoot;
        bytes32 _ethRoot = keccak256(
            abi.encodePacked(
                _nullRoot, 
                keccak256("eth")
            )
        );

        // Register the paulio.eth subdomain
        l2Registry.setSubnodeRecord(
            _ethRoot,
            keccak256("paulio"),
            _owner,
            _goerliENSResolver,
            64
        );

        bytes32 _paulioRoot = keccak256(
            abi.encodePacked(
                _ethRoot, 
                keccak256("paulio")
            )
        );

        // Register the lmeow.paulio.eth subdomain
        l2Registry.setSubnodeRecord(
            _paulioRoot,
            keccak256("lmeow"),
            _owner,
            _goerliENSResolver,
            64
        );

        l2Registry.owner(_paulioRoot);
        l2Registry.resolver(_paulioRoot);
        l2Registry.ttl(_paulioRoot);
    }
}
