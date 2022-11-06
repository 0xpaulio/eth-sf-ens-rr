// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import "forge-std/Script.sol";

import {L1ENSRegistry} from "src/l1/L1ENSRegistry.sol";

contract CounterScript is Script {
    L1ENSRegistry l1Registry;

    uint64 constant private  _optimismChainId = 10;
    address constant private _l2RegistrarAddress = address(0x0);

    address constant private _goerliOVMAddressManager = address(0xa6f73589243a6A7a9023b1Fa0651b1d89c177111);

    string[] _gatewayUrls;

    function setUp() public {
        _gatewayUrls = new string[](1);
        _gatewayUrls[0] = "http://localhost:3000/";
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
