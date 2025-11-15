// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Counter {
    uint256 private count;

    event CountChange(uint256 newValue, address indexed changer);

    constructor() {
        count = 0;
    }

    function incre() public  {
        count++;
    }

    function cc() public view returns(uint256) {
        return count;
    }

}