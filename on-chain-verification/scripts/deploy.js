async function deployContract() {
  const verifierContract = "TokenTransferContract";
  const votesTheshhold = 1;
  const builder="0xa763ebb58Fc66220F208e697E585a4197A941c84";
  const tokenAddress="0x844c811c0dc060808ac024b6e300499cbbd574b7";
  
  const ERC20Verifier = await ethers.getContractFactory(verifierContract);
  const erc20Verifier = await ERC20Verifier.deploy(
    votesTheshhold,tokenAddress
  );
  
  await erc20Verifier.deployed();
  console.log(votesTheshhold, " contract address:", erc20Verifier.address);
}
deployContract()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });