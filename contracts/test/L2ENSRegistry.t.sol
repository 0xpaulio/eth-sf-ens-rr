// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import "test/Base.t.sol";
import {L2ENSRegistry} from "src/l2/L2ENSRegistry.sol";

contract ENSRegistryTest is BaseTest {
    L2ENSRegistry l2Registry;

    event lbs(string, bytes32);

    function setUp() public override virtual {
        BaseTest.setUp();

        // Switch to the deployer
        vm.startPrank(deployer);

        // Deploy the L2 Registry
        l2Registry = new L2ENSRegistry();

        l2Registry.setSubnodeRecord(
            0x0,
            keccak256("eth"),
            deployer,
            address(0x1),
            64
        );

        // bytes32 _nullRoot = 0x0;
        // bytes32 _ethRoot = keccak256(abi.encodePacked(_nullRoot, keccak256("eth")));
        // emit lbs("_ethRoot", _ethRoot);
        // l2Registry.setSubnodeRecord(
        //     _ethRoot,
        //     keccak256("paulio"),
        //     deployer,
        //     address(0x1),
        //     64
        // );

        vm.stopPrank();
    }

    event las(address);
    event lus(uint64);

    function testSLO() public { 
        bytes32 _nullRoot;
        bytes32 _ethRoot = keccak256(abi.encodePacked(_nullRoot, keccak256("eth")));
        bytes32 _paulioRoot = keccak256(abi.encodePacked(_ethRoot, keccak256("paulio")));

        // keccak256(abi.encodePacked(node, uint256(1)));

        vm.startPrank(deployer);
        l2Registry.setSubnodeRecord(
            _ethRoot,
            keccak256("paulio"),
            deployer,
            address(0x1),
            64
        );

        bytes32 recordSLO_ = keccak256(
            abi.encodePacked(
                _paulioRoot, 
                uint256(0)
            )
        );
        
        bytes32 recordData_ = l2Registry.getSLO(recordSLO_);
        emit lbs("recordData_", recordData_);

        bytes32 shiftedData_ = l2Registry.getSLO(
            bytes32(
                uint256(recordSLO_) + uint256(0x1)
            )
        );
        emit lbs("shiftedData_", shiftedData_);

        bytes32 shiftedDatatwo_ = l2Registry.getSLO(
            bytes32(
                uint256(recordSLO_) + uint256(0x2)
            )
        );
        emit lbs("shiftedData_", shiftedDatatwo_);

        bytes32 actualSLO_ = l2Registry.getRecordSLO(_paulioRoot);
        
        bytes32 actualData_ = l2Registry.getSLO(
            bytes32(
                uint256(recordSLO_) + uint256(0x1)
            )
        );
        emit lbs("actualData_", actualData_);

        address res_ = address(uint160(uint256(actualData_)));
        uint64 ttl_;
        assembly {
            ttl_ := shr(0xa0, actualData_)
        }
        // uint64 ttl_ = shr(0xa0, actualData_);

        emit las(res_);
        emit lus(ttl_);

    }
}
