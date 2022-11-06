// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import "forge-std/Script.sol";

import {L1ENSRegistry} from "src/l1/L1ENSRegistry.sol";

contract ContractScript is Script {
    L1ENSRegistry l1Registry;

    // uint64 constant private  _optimismChainId = 10;
    uint64  _optimismChainId = 420;
    address _l2RegistrarAddress = 0x5833395C4BDf4Bb50ac7F4868Fe922AE091D1129;

    address _goerliOVMAddressManager = 0xa6f73589243a6A7a9023b1Fa0651b1d89c177111;

    string[] _gatewayUrls;

    function setUp() public {
        _gatewayUrls = new string[](1);
        _gatewayUrls[0] = "http://localhost:3000/gateway/{sender}/{data}.json";
    }

    function run() public {
        vm.startBroadcast();

        // Deploy the L1 Registry
        l1Registry = new L1ENSRegistry(
            _optimismChainId, _l2RegistrarAddress,
            _goerliOVMAddressManager,
            _gatewayUrls
        );
    }
}
