import {
      ComputeBudgetProgram,
  Connection,
  Keypair,
  PublicKey,
  SendTransactionError,
  Transaction,
  TransactionInstruction,
  sendAndConfirmTransaction,
} from "@solana/web3.js";
import {
  TOKEN_PROGRAM_ID,
} from "@solana/spl-token";
import fs from "fs";

const PROGRAM_ID = new PublicKey("9R6bqjw54kNcVgxdTU6i6fXLpnmB29Badep2b3S1XwM6");
const MINT = new PublicKey("47VHscQ5d95Axasw38Axi5ttmcxBKLqim3NnLZ5ddTvN");

// === CONNECTION & WALLET ===
const connection = new Connection("https://api.devnet.solana.com", "confirmed");
const payer = Keypair.fromSecretKey(
  Uint8Array.from(JSON.parse(fs.readFileSync(`${process.env.HOME}/.config/solana/id.json`, "utf-8")))
);

// === DERIVE PDA ===
const [mintAuthorityPda] = PublicKey.findProgramAddressSync(
  [Buffer.from("ofac_check_mint_auth")],
  PROGRAM_ID
);

const computeIx = ComputeBudgetProgram.setComputeUnitLimit({
  units: 520_000,
});


// === READ COMMAND LINE ARGUMENTS ===
// Usage: npx ts-node invokeProgram.ts <destination> <proof.bin> <witness.bin>
const [destinationArg, proofPath, witnessPath] = process.argv.slice(2);

if (!destinationArg || !proofPath || !witnessPath) {
  console.error(
    "❌ Usage: npx ts-node invokeProgram.ts <destination> <proof.bin> <witness.bin>"
  );
  process.exit(1);
}

if (!fs.existsSync(proofPath) || !fs.existsSync(witnessPath)) {
  console.error(
    `❌ File not found: ${!fs.existsSync(proofPath) ? proofPath : witnessPath}`
  );
  process.exit(1);
}

const proofBytes = fs.readFileSync(proofPath);
const witnessBytes = fs.readFileSync(witnessPath);
console.log(`Loaded ${proofBytes.length} bytes from ${proofPath}`);
console.log(`Loaded ${witnessBytes.length} bytes from ${witnessPath}`);

// === CONCATENATE BYTES ===
const instructionData = Buffer.concat([proofBytes, witnessBytes]);
console.log(`Total instruction data length: ${instructionData.length} bytes`);

// === BUILD INSTRUCTION ===
const DESTINATION = new PublicKey(destinationArg);

const ix = new TransactionInstruction({
  programId: PROGRAM_ID, keys: [
    { pubkey: MINT, isSigner: false, isWritable: true },
    { pubkey: DESTINATION, isSigner: false, isWritable: true },
    { pubkey: TOKEN_PROGRAM_ID, isSigner: false, isWritable: false },
    { pubkey: mintAuthorityPda, isSigner: false, isWritable: false },
  ],
  data: instructionData,
});

(async () => {
  try {
    const sig = await sendAndConfirmTransaction(
      connection,
      new Transaction().add(ix).add(computeIx),
      [payer]
    );
    console.log("✅ Transaction successful!", sig);
  } catch (err) {
    console.error(
      "❌ Transaction failed!",
      err instanceof SendTransactionError ? err.logs : err
    );
  }
})();
