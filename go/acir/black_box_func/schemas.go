package blackboxfunc

import "sunspot/go/acir/msgpackutil"

// Per-variant serde field schemas. Field names and tag ordering follow
// noir's struct definitions at
// acvm-repo/acir/src/circuit/opcodes/black_box_function_call.rs.

var aes128EncryptSchema = msgpackutil.NewSchema(map[string]int{
	"inputs": 0, "iv": 1, "key": 2, "outputs": 3,
})

var andSchema = msgpackutil.NewSchema(map[string]int{
	"lhs": 0, "rhs": 1, "num_bits": 2, "output": 3,
})

var xorSchema = msgpackutil.NewSchema(map[string]int{
	"lhs": 0, "rhs": 1, "num_bits": 2, "output": 3,
})

var rangeSchema = msgpackutil.NewSchema(map[string]int{
	"input": 0, "num_bits": 1,
})

var blake2sSchema = msgpackutil.NewSchema(map[string]int{
	"inputs": 0, "outputs": 1,
})

var blake3Schema = msgpackutil.NewSchema(map[string]int{
	"inputs": 0, "outputs": 1,
})

var ecdsaSecp256k1Schema = msgpackutil.NewSchema(map[string]int{
	"public_key_x": 0, "public_key_y": 1, "signature": 2, "hashed_message": 3, "predicate": 4, "output": 5,
})

var ecdsaSecp256r1Schema = msgpackutil.NewSchema(map[string]int{
	"public_key_x": 0, "public_key_y": 1, "signature": 2, "hashed_message": 3, "predicate": 4, "output": 5,
})

var multiScalarMulSchema = msgpackutil.NewSchema(map[string]int{
	"points": 0, "scalars": 1, "predicate": 2, "outputs": 3,
})

var embeddedCurveAddSchema = msgpackutil.NewSchema(map[string]int{
	"input1": 0, "input2": 1, "predicate": 2, "outputs": 3,
})

var keccakf1600Schema = msgpackutil.NewSchema(map[string]int{
	"inputs": 0, "outputs": 1,
})

var recursiveAggregationSchema = msgpackutil.NewSchema(map[string]int{
	"verification_key": 0, "proof": 1, "public_inputs": 2, "key_hash": 3, "proof_type": 4, "predicate": 5,
})

var poseidon2PermutationSchema = msgpackutil.NewSchema(map[string]int{
	"inputs": 0, "outputs": 1,
})

var sha256CompressionSchema = msgpackutil.NewSchema(map[string]int{
	"inputs": 0, "hash_values": 1, "outputs": 2,
})

// schema() methods are the BlackBoxFunction interface hook that lets
// BlackBoxFuncCall.decodeBlackBoxFunction pass the right schema into
// msgpackutil.ReadStruct for each concrete variant.

func (*AES128Encrypt[T, E]) schema() msgpackutil.Schema { return aes128EncryptSchema }
func (*And[T, E]) schema() msgpackutil.Schema           { return andSchema }
func (*Xor[T, E]) schema() msgpackutil.Schema           { return xorSchema }
func (*Range[T, E]) schema() msgpackutil.Schema         { return rangeSchema }
func (*Blake2s[T, E]) schema() msgpackutil.Schema       { return blake2sSchema }
func (*Blake3[T, E]) schema() msgpackutil.Schema        { return blake3Schema }
func (*ECDSASECP256K1[T, E]) schema() msgpackutil.Schema {
	return ecdsaSecp256k1Schema
}
func (*ECDSASECP256R1[T, E]) schema() msgpackutil.Schema {
	return ecdsaSecp256r1Schema
}
func (*MultiScalarMul[T, E]) schema() msgpackutil.Schema   { return multiScalarMulSchema }
func (*EmbeddedCurveAdd[T, E]) schema() msgpackutil.Schema { return embeddedCurveAddSchema }
func (*Keccakf1600[T, E]) schema() msgpackutil.Schema      { return keccakf1600Schema }
func (*RecursiveAggregation[T, E]) schema() msgpackutil.Schema {
	return recursiveAggregationSchema
}
func (*Poseidon2Permutation[T, E]) schema() msgpackutil.Schema {
	return poseidon2PermutationSchema
}
func (*SHA256Compression[T, E]) schema() msgpackutil.Schema { return sha256CompressionSchema }
