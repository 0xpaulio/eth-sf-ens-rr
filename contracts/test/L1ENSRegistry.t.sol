// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import "test/Base.t.sol";
import {L1ENSRegistry} from "src/l1/L1ENSRegistry.sol";

contract ENSRegistryTest is BaseTest {
    L1ENSRegistry l1Registry;

    uint64 constant private  _optimismChainId = 10;
    address constant private _l2RegistrarAddress = address(0x0);

    address constant private _goerliOVMAddressManager = address(0x0);    

    string[] _gatewayUrls;

    function setUp() public override virtual {
        BaseTest.setUp();

        _gatewayUrls = new string[](1);
        _gatewayUrls[0] = "http://localhost:3000/";

        // Switch to the deployer
        vm.startPrank(deployer);

        // Deploy the L1 Registry
        l1Registry = new L1ENSRegistry(
            _optimismChainId, _l2RegistrarAddress,
            _goerliOVMAddressManager,
            _gatewayUrls
        );
    }

    // function testSLO() public { }
}
