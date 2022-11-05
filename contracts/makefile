include .env

testnet-init:
	@forge script --rpc-url $(LOCAL_TESTNET_RPC_URL__LIVE) \
	--private-key $(LOCAL_TESTNET_PRIVATE_KEY) \
	scripts/DeployLocalTestnet.s.sol :ContractScript \
	--broadcast -vvvv

goerli-init:
	@forge script --rpc-url $(RINKEBY_RPC_URL__LIVE) \
	--private-key $(PRIVATE_KEY) \
	script/Deploy.s.sol :ContractScript \
	--broadcast

install-deps:
	forge install transmissions11/solmate \
	OpenZeppelin/openzeppelin-contracts \
	OpenZeppelin/openzeppelin-contracts-upgradeable \
	ensdomains/ens \
	ensdomains/ens-contracts \
	--no-commit
