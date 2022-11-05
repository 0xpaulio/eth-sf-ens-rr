// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";

import "./utils/Utils.sol";

contract BaseTest is Test {
    Utils public utils;

    address public deployer;
    address public owner;
    address public dev;

    address public alice;
    address public bob;
    address public chris;

    function setUp() public virtual {
        utils = new Utils();

        deployer = utils.initializeAccount("deployer");
        owner = utils.initializeAccount("owner");
        dev = utils.initializeAccount("dev");

        alice = utils.initializeAccount("alice");
        bob   = utils.initializeAccount("bob");
        chris = utils.initializeAccount("chris");
    }
}
