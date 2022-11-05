// SPDX-License-Identifier: Unlicense
pragma solidity >=0.8.0;

import {DSTest} from "ds-test/test.sol";
import {Vm} from "forge-std/Vm.sol";

contract Utils is DSTest {
    Vm internal immutable vm = Vm(HEVM_ADDRESS);
    bytes32 internal nextUser =
        keccak256(abi.encodePacked("generateAddress()"));

    function getNextUserAddress() external returns (address payable) {
        address payable user = payable(address(uint160(uint256(nextUser))));
        nextUser = keccak256(abi.encodePacked(nextUser));
        return user;
    }

    function createUsers(uint256 userNum)
        external
        returns (address payable[] memory)
    {
        address payable[] memory users = new address payable[](userNum);
        for (uint256 i = 0; i < userNum; i++) {
            address payable user = this.getNextUserAddress();
            vm.deal(user, 100 ether);
            users[i] = user;
        }
        return users;
    }

    function mineBlocks(uint256 numBlocks) external {
        uint256 targetBlock = block.number + numBlocks;
        vm.roll(targetBlock);
    }

    function initializeAccounts(string[] memory names)
        external
        returns (address payable[] memory users)
    {
        users = new address payable[](names.length);
        for (uint256 index; index < names.length; index += 1) {
            users[index] = this.getNextUserAddress();
            vm.label(users[index], names[index]);
        }
    }

    function initializeAccount(string memory name)
        external
        returns (address payable user)
    {
        user = this.getNextUserAddress();
        vm.label(user, name);
    }
}