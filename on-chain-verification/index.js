const express = require('express');
const { ethers } = require('hardhat');
const cors = require('cors');


const app = express();
const port = 3003;
app.use(express.json());

app.use(cors());

const { Web3 } = require('web3');
const { poseidon } = require('@iden3/js-crypto');
const { SchemaHash } = require('@iden3/js-iden3-core');
const { prepareCircuitArrayValues } = require('@0xpolygonid/js-sdk');

const Operators = {
  NOOP : 0, // No operation, skip query verification in circuit
  EQ : 1, // equal
  LT : 2, // less than
  GT : 3, // greater than
  IN : 4, // in
  NIN : 5, // not in
  NE : 6   // not equal
}

app.post('/v1/contract/deploy', async (req, res) => {
  try {


    const builder= req.body.builder;
    const votesTheshhold= req.body.votesTheshhold;
    const tokenAddress= req.body.tokenAddress;

    const contractAddress = await deployContract(votesTheshhold, builder,tokenAddress);
    await setZKPRequest(contractAddress,true)
    await setZKPRequest(contractAddress,false)

    res.json({contractAddress: contractAddress });
  } catch (error) {
    console.error('Error deploying contract:', error);
    res.status(500).send('Internal Server Error');
  }
});

app.listen(port, () => {
  console.log(`Express app listening at http://localhost:${port}`);
});

async function deployContract(votesTheshhold,builder,tokenAddress) {
  const [deployer] = await ethers.getSigners();
  console.log('Deploying contract with address:', deployer.address);
  const verifierContract = "TokenTransferContract";

  const DAOVerifier = await ethers.getContractFactory(verifierContract);
  const daoVerifier = await DAOVerifier.deploy(
      votesTheshhold,builder,tokenAddress
  );

  await daoVerifier.deployed()


  console.log(" contract address:", daoVerifier.address);;
  return  daoVerifier.address;
}



function packValidatorParams(query, allowedIssuers,isBuilder = []) {
  let web3 = new Web3(Web3.givenProvider || 'wss://polygon-mumbai.g.alchemy.com/v2/W6qdHNAQ5hacjzn31F5_53PH5N2Rrn3a');
  return web3.eth.abi.encodeParameter(
      {
        CredentialAtomicQuery: {
          schema: 'uint256',
          claimPathKey: 'uint256',
          operator: 'uint256',
          slotIndex: 'uint256',
          value: isBuilder?'uint256[]':'uint256[]',
          queryHash: 'uint256',
          allowedIssuers: 'uint256[]',
          circuitIds: 'string[]',
          skipClaimRevocationCheck: 'bool',
          claimPathNotExists: 'uint256'
        }
      },
      {
        schema: query.schema,
        claimPathKey: query.claimPathKey,
        operator: query.operator,
        slotIndex: query.slotIndex,
        value: query.value,
        queryHash: query.queryHash,
        allowedIssuers: [],
        circuitIds: query.circuitIds,
        skipClaimRevocationCheck: query.skipClaimRevocationCheck,
        claimPathNotExists: query.claimPathNotExists
      }
  );
}

function coreSchemaFromStr(schemaIntString) {
  const schemaInt = BigInt(schemaIntString);
  return SchemaHash.newSchemaHashFromInt(schemaInt);
};

function calculateQueryHash(
    values,
    schema,
    slotIndex,
    operator,
    claimPathKey,
    claimPathNotExists
) {
  const expValue = prepareCircuitArrayValues(values, 64);
  console.log(expValue)
  const valueHash = poseidon.spongeHashX(expValue, 6);
  const schemaHash = coreSchemaFromStr(schema);
  const quaryHash = poseidon.hash([
    schemaHash.bigInt(),
    BigInt(slotIndex),
    BigInt(operator),
    BigInt(claimPathKey),
    BigInt(claimPathNotExists),
    valueHash
  ]);
  return quaryHash;
}


async function setZKPRequest(contractAddress,isBuilder) {

  // you can run https://go.dev/play/p/3id7HAhf-Wi to get schema hash and claimPathKey using YOUR schema
  // suggestion: Use your own go application with that code rather than using playground (it can give a timeout just because it’s restricted by the size of dependency package)
  const allowedIssuers = ['*'];
  //const allowedIssuers = ['did:polygonid:polygon:mumbai:2qNzSKEuYnHwN7NgdmVM8DMYgpWVnCtnup1esfeCJ1'];

  if(isBuilder){
    const schemaBigInt = "4174810409265315245331249378998401719"

    const type = 'ReFreshBuilder';
    const schemaUrl = 'ipfs://QmWamN8tGpJjKdr1NDEduedP16CH85W4Ua8oBi4NE3MAqA';
    // merklized path to field in the W3C credential according to JSONLD  schema e.g. birthday in the KYCAgeCredential under the url "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld"
    const schemaClaimPathKey = "16444678393588006201827608929727578052908209026381139130045458653892476365976"

    const requestId = 1702148496;

    const query = {
      requestId,
      schema: schemaBigInt,
      claimPathKey: schemaClaimPathKey,
      operator: Operators.EQ,
      slotIndex: 0,
      value: [1, ...new Array(63).fill(0)], // for operators 1-3 only first value matters
      circuitIds: ['credentialAtomicQuerySigV2OnChain'],
      skipClaimRevocationCheck: false,
      claimPathNotExists: 0
    };


    query.queryHash = calculateQueryHash(
        query.value,
        query.schema,
        query.slotIndex,
        query.operator,
        query.claimPathKey,
        query.claimPathNotExists
    ).toString();

    // add the address of the contract just deployed

    let daoVerifier = await ethers.getContractAt("TokenTransferContract", contractAddress)

    const validatorAddress = "0x1E4a22540E293C0e5E8c33DAfd6f523889cFd878"; // sig validator
    // const validatorAddress = "0x0682fbaA2E4C478aD5d24d992069dba409766121"; // mtp validator

    const invokeRequestMetadata = {
      id: '7f38a193-0918-4a48-9fac-36adfdb8b542',
      typ: 'application/iden3comm-plain-json',
      type: 'https://iden3-communication.io/proofs/1.0/contract-invoke-request',
      thid: '7f38a193-0918-4a48-9fac-36adfdb8b542',
      body: {
        reason: 'zk for Fund execution',
        transaction_data: {
          contract_address: contractAddress,
          method_id: 'b68967e2',
          chain_id: 80001,
          network: 'polygon-mumbai'
        },
        scope: [
          {
            id: query.requestId,
            circuitId: query.circuitIds[0],
            query: {
              allowedIssuers: allowedIssuers,
              context: schemaUrl,
              credentialSubject: {
                isBuilder: {
                  $eq: query.value[0]
                }
              },
              type
            }
          }
        ]
      }
    };

    try {
      const txId = await daoVerifier.setZKPRequest(
          requestId, {
            metadata: JSON.stringify(invokeRequestMetadata),
            validator: validatorAddress,
            data: packValidatorParams(query,allowedIssuers,isBuilder)
          });
      console.log("Request set: ", txId.hash);
    } catch (e) {
      console.log("error: ", e);
    }
  }else{
    const schemaBigInt = "278967279552142591952300574090397990154"

    const type = 'DAOVerificationToken';
    const schemaUrl = 'ipfs://QmQb3pfSfmFZNQapcXk3zdnDnmmpqmiZ6YWFcCwyq14ajM';
    // merklized path to field in the W3C credential according to JSONLD  schema e.g. birthday in the KYCAgeCredential under the url "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld"
    const schemaClaimPathKey = "18694576051158235627586407093896000546137977171154473590197366253690295065210"

    const requestId = 1701840378;

    const query = {
      requestId,
      schema: schemaBigInt,
      claimPathKey: schemaClaimPathKey,
      operator: Operators.GT,
      slotIndex: 0,
      value: [0, ...new Array(63).fill(0)], // for operators 1-3 only first value matters
      circuitIds: ['credentialAtomicQuerySigV2OnChain'],
      skipClaimRevocationCheck: false,
      claimPathNotExists: 0
    };

    console.log(query.value)


    query.queryHash = calculateQueryHash(
        query.value,
        query.schema,
        query.slotIndex,
        query.operator,
        query.claimPathKey,
        query.claimPathNotExists
    ).toString();

    // add the address of the contract just deployed

    let daoVerifier = await ethers.getContractAt("TokenTransferContract", contractAddress)


    const validatorAddress = "0x1E4a22540E293C0e5E8c33DAfd6f523889cFd878"; // sig validator
    // const validatorAddress = "0x0682fbaA2E4C478aD5d24d992069dba409766121"; // mtp validator

    const invokeRequestMetadata = {
      id: '7f38a193-0918-4a48-9fac-36adfdb8b542',
      typ: 'application/iden3comm-plain-json',
      type: 'https://iden3-communication.io/proofs/1.0/contract-invoke-request',
      thid: '7f38a193-0918-4a48-9fac-36adfdb8b542',
      body: {
        reason: 'vote for project',
        transaction_data: {
          contract_address: contractAddress,
          method_id: 'b68967e2',
          chain_id: 80001,
          network: 'polygon-mumbai'
        },
        scope: [
          {
            id: query.requestId,
            circuitId: query.circuitIds[0],
            query: {
              allowedIssuers: allowedIssuers,
              context: schemaUrl,
              credentialSubject: {
                token: {
                  $gt: query.value[0]
                }
              },
              type
            }
          }
        ]
      }
    };

    try {
      const txId = await daoVerifier.setZKPRequest(
          requestId, {
            metadata: JSON.stringify(invokeRequestMetadata),
            validator: validatorAddress,
            data: packValidatorParams(query,allowedIssuers,isBuilder)
          });
      console.log("Request set: ", txId.hash);
    } catch (e) {
      console.log("error: ", e);
    }
  }

}

  