package acir

import (
	"fmt"
	bbf "sunspot/go/acir/black_box_func"
	"sunspot/go/acir/brillig_call"
	"sunspot/go/acir/call"
	exp "sunspot/go/acir/expression"
	"sunspot/go/acir/memory_init"
	mem_op "sunspot/go/acir/memory_op"
	"sunspot/go/acir/msgpackutil"
	ops "sunspot/go/acir/opcodes"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
	"github.com/google/btree"
)

type Circuit[T shr.ACIRField, E constraint.Element] struct {
	CircuitName         string
	CurrentWitnessIndex uint32
	Opcodes             []ops.Opcode[E] `json:"opcodes"`            // Opcodes in the circuit
	PrivateParameters   btree.BTree     `json:"private_parameters"` // Witnesses
	PublicParameters    btree.BTree     `json:"public_parameters"`  // Witnesses
	ReturnValues        btree.BTree     `json:"return_values"`      // Witnesses
	MemoryBlocks        map[uint32]*logderivlookup.Table
}

func (c *Circuit[T, E]) UnmarshalReader(r *msgpackutil.Reader) error {
	c.PrivateParameters = *btree.New(2)
	c.PublicParameters = *btree.New(2)
	c.ReturnValues = *btree.New(2)

	// Reset the per-decode witness high-water mark so populateCurrentWitnessIndex
	// reads only what this circuit references (witness indices are local).
	r.ResetWitnessTracker()

	err := msgpackutil.ReadStruct(r, "Circuit", []msgpackutil.Field{
		{Name: "function_name", Decode: func(r *msgpackutil.Reader) error {
			s, err := r.ReadString()
			if err != nil {
				return err
			}
			c.CircuitName = s
			return nil
		}},
		{Name: "opcodes", Decode: c.readOpcodes},
		{Name: "private_parameters", Decode: func(r *msgpackutil.Reader) error { return readWitnessSet(r, &c.PrivateParameters) }},
		{Name: "public_parameters", Decode: func(r *msgpackutil.Reader) error { return readWitnessSet(r, &c.PublicParameters) }},
		{Name: "return_values", Decode: func(r *msgpackutil.Reader) error { return readWitnessSet(r, &c.ReturnValues) }},
		{Name: "assert_messages", Decode: msgpackutil.SkipField},
	})
	if err != nil {
		return err
	}

	if maxW, ok := r.MaxWitness(); ok {
		c.CurrentWitnessIndex = maxW
	}
	return nil
}

func (c *Circuit[T, E]) readOpcodes(r *msgpackutil.Reader) error {
	n, err := r.ReadArrayLen()
	if err != nil {
		return err
	}
	c.Opcodes = make([]ops.Opcode[E], n)
	for i := 0; i < n; i++ {
		var op ops.Opcode[E]
		if err := msgpackutil.ReadDispatchedEnum(r, "Opcode", []ops.Opcode[E]{
			&exp.Expression[T, E]{},
			&bbf.BlackBoxFuncCall[T, E]{},
			&mem_op.MemoryOp[T, E]{},
			&memory_init.MemoryInit[T, E]{},
			&brillig_call.BrilligCall[T, E]{},
			&call.Call[T, E]{},
		}, func(v ops.Opcode[E]) { op = v }); err != nil {
			return fmt.Errorf("opcode %d: %w", i, err)
		}
		c.Opcodes[i] = op
	}
	return nil
}

// readWitnessSet decodes a BTreeSet<Witness>, which serializes as a fixarray
// of witnesses. PublicInputs is a single-field tuple struct wrapping such a
// set; with EncodingStrategy::Array tuple structs are transparent on the
// wire, so the same decoder handles both.
func readWitnessSet(r *msgpackutil.Reader, dst *btree.BTree) error {
	n, err := r.ReadArrayLen()
	if err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		var w shr.Witness
		if err := w.UnmarshalReader(r); err != nil {
			return err
		}
		dst.ReplaceOrInsert(w)
	}
	return nil
}

// Define the constraints for a circuit
// This returns the input and output variables of the circuit,
// so that circuits that call the circuit can check that the values they called the
// circuit with are consistent with the true value.
func (c *Circuit[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable, resolve CircuitResolver[T, E], index *uint32) ([]frontend.Variable, []frontend.Variable, error) {
	c.MemoryBlocks = make(map[uint32]*logderivlookup.Table)

	// 1. Resolve and define all subcircuits
	callConnections, err := c.defineSubcircuits(api, witnesses, resolve, index)
	if err != nil {
		return nil, nil, err
	}

	// 2. Collect witnesses for current circuit
	currentWitnesses, err := c.collectCurrentWitnesses(witnesses, index)
	if err != nil {
		return nil, nil, err
	}

	// 3. Add the constraints for the circuit
	if err := c.constrainCircuit(api, currentWitnesses); err != nil {
		return nil, nil, err
	}

	// 4. Connect call inputs/outputs
	c.constrainCircuitCalls(api, currentWitnesses, callConnections)

	// 5. Collect circuit inputs and outputs
	inputs := c.collectWitnesses(currentWitnesses, &c.PrivateParameters, &c.PublicParameters)
	outputs := c.collectWitnesses(currentWitnesses, &c.ReturnValues)

	return inputs, outputs, nil
}

// Run the definition function for the circuits called by the circuit
func (c *Circuit[T, E]) defineSubcircuits(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable, resolve CircuitResolver[T, E], index *uint32) (map[int]struct {
	Inputs  []frontend.Variable
	Outputs []frontend.Variable
}, error) {
	callConnections := make(map[int]struct {
		Inputs  []frontend.Variable
		Outputs []frontend.Variable
	})

	for i, opcode := range c.Opcodes {
		callOp, ok := opcode.(*call.Call[T, E])
		if !ok {
			continue
		}

		subCircuit, err := resolve(callOp.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve circuit %d: %w", callOp.ID, err)
		}

		// Run subcircuit definition
		in, out, err := subCircuit.Define(api, witnesses, resolve, index)
		if err != nil {
			return nil, fmt.Errorf("failed to define subcircuit %d: %w", callOp.ID, err)
		}

		if len(in) > len(callOp.Inputs) {
			return nil, fmt.Errorf("input count mismatch: subcircuit %d requires more inputs than are given by the outer circuit", callOp.ID)
		}
		if len(out) > len(callOp.Outputs) {
			return nil, fmt.Errorf("output count mismatch: subcircuit %d provides more outputs than the outer circuit is expecting", callOp.ID)
		}

		callConnections[i] = struct {
			Inputs  []frontend.Variable
			Outputs []frontend.Variable
		}{Inputs: in, Outputs: out}
	}

	return callConnections, nil
}

// Get the partial witness for a particular circuit call
// The partial witness for a whole programme consists of a concatenation of a postorder traversal of the programme tree.
// We perform a postorder traversal and 'pop' the witness that we need by incrementing the global index
func (c *Circuit[T, E]) collectCurrentWitnesses(witnesses map[shr.Witness]frontend.Variable, index *uint32) (map[shr.Witness]frontend.Variable, error) {
	currentWitnesses := make(map[shr.Witness]frontend.Variable, c.CurrentWitnessIndex+1)

	for i := range c.CurrentWitnessIndex + 1 {
		global := shr.Witness(i + uint32(*index))
		v, ok := witnesses[global]
		if !ok {
			// Compile allocates a variable for every slot in [0, totalSlots),
			// so this lookup must succeed if the global index is advanced
			// consistently with the layout computed by Program.WitnessLayout.
			return nil, fmt.Errorf(
				"circuit %q: missing witness for slot %d (global %d); "+
					"witness layout did not allocate this index",
				c.CircuitName, i, global,
			)
		}
		currentWitnesses[shr.Witness(i)] = v
	}

	*index += c.CurrentWitnessIndex + 1
	return currentWitnesses, nil
}

// Add constraints for a specific circuit call within a programme
func (c *Circuit[T, E]) constrainCircuit(api frontend.Builder[E], currentWitnesses map[shr.Witness]frontend.Variable) error {
	for _, opcode := range c.Opcodes {
		memInit, ok := opcode.(*memory_init.MemoryInit[T, E])
		if ok {
			table := logderivlookup.New(api)
			memInit.Table = &table
			c.MemoryBlocks[memInit.BlockID] = &table
		}

		memOp, ok := opcode.(*mem_op.MemoryOp[T, E])
		if ok {
			memOp.Memory = c.MemoryBlocks
		}

		if err := opcode.Define(api, currentWitnesses); err != nil {
			return err
		}
	}
	return nil
}

// Ensure that the input and return values of a circuit call are consistent with the values
// that are in the partial witness for the outer circuit
func (c *Circuit[T, E]) constrainCircuitCalls(api frontend.Builder[E], currentWitnesses map[shr.Witness]frontend.Variable, callConnections map[int]struct {
	Inputs  []frontend.Variable
	Outputs []frontend.Variable
}) {
	for i, opcode := range c.Opcodes {
		callOp, ok := opcode.(*call.Call[T, E])
		if !ok {
			continue
		}
		connection := callConnections[i]
		for j, inputWitness := range callOp.Inputs {
			api.AssertIsEqual(currentWitnesses[inputWitness], connection.Inputs[j])
		}
		for j, outputWitness := range callOp.Outputs {
			api.AssertIsEqual(currentWitnesses[outputWitness], connection.Outputs[j])
		}
	}
}

// Construct a list of input/ output variables of a circuit given trees of witness
// indices and a index->variable mapping. Entries from all trees are visited in
// global witness-index order (Witness.Less sorts by uint32).
func (c *Circuit[T, E]) collectWitnesses(currentWitnesses map[shr.Witness]frontend.Variable, trees ...*btree.BTree) []frontend.Variable {
	merged := btree.New(2)
	for _, tree := range trees {
		tree.Ascend(func(it btree.Item) bool {
			merged.ReplaceOrInsert(it)
			return true
		})
	}
	var vars []frontend.Variable
	merged.Ascend(func(it btree.Item) bool {
		witness, ok := it.(shr.Witness)
		if !ok {
			return false
		}
		vars = append(vars, currentWitnesses[witness])
		return true
	})
	return vars
}

// countSubcircuitSlots walks the call tree in opcode order and returns the
// global slot at which this circuit's own witnesses begin: the starting
// `index` plus the total slot count of every transitively-called subcircuit.
func (c *Circuit[T, E]) countSubcircuitSlots(resolve CircuitResolver[T, E], index uint32) (uint32, error) {
	for _, opcode := range c.Opcodes {
		callOp, ok := opcode.(*call.Call[T, E])
		if !ok {
			continue
		}
		subCircuit, err := resolve(callOp.ID)
		if err != nil {
			return index, fmt.Errorf("failed to resolve circuit %d: %w", callOp.ID, err)
		}
		subOwnStart, err := subCircuit.countSubcircuitSlots(resolve, index)
		if err != nil {
			return index, err
		}
		index = subOwnStart + subCircuit.CurrentWitnessIndex + 1
	}
	return index, nil
}
