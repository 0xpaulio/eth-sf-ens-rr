include .env

testnet-init:
	@forge script --rpc-url $(LOCAL_TESTNET_RPC_URL__LIVE) \
	--private-key $(LOCAL_TESTNET_PRIVATE_KEY) \
	scripts/DeployLocalTestnet.s.sol :ContractScript \
	--broadcast -vvvv

goerli-l1-init:
	@forge script --rpc-url $(Goerli_ETH_RPC_URL__LIVE) \
	--private-key $(PRIVATE_KEY) \
	script/DeployL1Contracts.s.sol :ContractScript \
	--etherscan-api-key $(ETHERSCAN_API_KEY) --verify \
	--broadcast
	
goerli-l2-init:
	@forge script --rpc-url $(Goerli_OPT_RPC_URL__LIVE) \
	--private-key $(PRIVATE_KEY) \
	script/DeployL2Contracts.s.sol :ContractScript \
	--etherscan-api-key $(ETHERSCAN_API_KEY) --verify \
	--broadcast --legacy --slow

install-deps:
	forge install transmissions11/solmate \
	OpenZeppelin/openzeppelin-contracts \
	ensdomains/ens \
	ensdomains/ens-contracts \
	ethereum-optimism/optimism \
	--no-commit
