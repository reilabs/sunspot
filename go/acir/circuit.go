package acir

import (
	"encoding/binary"
	"fmt"
	"io"
	ap "sunspot/go/acir/assertion_payload"
	bbf "sunspot/go/acir/black_box_func"
	"sunspot/go/acir/brillig"
	"sunspot/go/acir/call"
	exp "sunspot/go/acir/expression"
	"sunspot/go/acir/memory_init"
	mem_op "sunspot/go/acir/memory_op"
	ops "sunspot/go/acir/opcodes"
	shr "sunspot/go/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
	"github.com/google/btree"
)

type Circuit[T shr.ACIRField, E constraint.Element] struct {
	CircuitName         string
	CurrentWitnessIndex uint32                                           `json:"current_witness_index"`
	Opcodes             []ops.Opcode[E]                                  `json:"opcodes"`            // Opcodes in the circuit
	PrivateParameters   btree.BTree                                      `json:"private_parameters"` // Witnesses
	PublicParameters    btree.BTree                                      `json:"public_parameters"`  // Witnesses
	ReturnValues        btree.BTree                                      `json:"return_values"`      // Witnesses
	AssertMessages      map[ops.OpcodeLocation]ap.AssertionPayload[T, E] `json:"assert_messages"`    // Assert messages for the circuit
	MemoryBlocks        map[uint32]*logderivlookup.Table
}

func (c *Circuit[T, E]) UnmarshalReader(r io.Reader) error {
	var length uint64
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return err
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return err
	}

	c.CircuitName = string(data)
	if err := binary.Read(r, binary.LittleEndian, &c.CurrentWitnessIndex); err != nil {
		return err
	}

	var numOpcodes uint64
	if err := binary.Read(r, binary.LittleEndian, &numOpcodes); err != nil {
		return err
	}

	c.Opcodes = make([]ops.Opcode[E], numOpcodes)
	for i := uint64(0); i < numOpcodes; i++ {
		op, err := NewOpcode[T, E](r)
		if err != nil {
			return fmt.Errorf("failed to create opcode: %w", err)
		}
		if err := op.UnmarshalReader(r); err != nil {
			return fmt.Errorf("failed to unmarshal opcode at index %d: %w", i, err)
		}
		c.Opcodes[i] = op
	}

	var numPrivateParameters uint64
	if err := binary.Read(r, binary.LittleEndian, &numPrivateParameters); err != nil {
		return err
	}
	c.PrivateParameters = *btree.New(2)
	for i := uint64(0); i < numPrivateParameters; i++ {
		var witness shr.Witness
		if err := witness.UnmarshalReader(r); err != nil {
			return err
		}
		c.PrivateParameters.ReplaceOrInsert(witness)
	}

	var numPublicParameters uint64
	if err := binary.Read(r, binary.LittleEndian, &numPublicParameters); err != nil {
		return err
	}
	c.PublicParameters = *btree.New(2)
	for i := uint64(0); i < numPublicParameters; i++ {
		var witness shr.Witness
		if err := witness.UnmarshalReader(r); err != nil {
			return err
		}
		c.PublicParameters.ReplaceOrInsert(witness)
	}

	var numReturnValues uint64
	if err := binary.Read(r, binary.LittleEndian, &numReturnValues); err != nil {
		return err
	}
	c.ReturnValues = *btree.New(2)
	for i := uint64(0); i < numReturnValues; i++ {
		var witness shr.Witness
		if err := witness.UnmarshalReader(r); err != nil {
			return err
		}
		c.ReturnValues.ReplaceOrInsert(witness)
	}

	var numAssertMessages uint64
	if err := binary.Read(r, binary.LittleEndian, &numAssertMessages); err != nil {
		if err == io.EOF {
			c.AssertMessages = make(map[ops.OpcodeLocation]ap.AssertionPayload[T, E])
			return nil
		}
	}

	c.AssertMessages = make(map[ops.OpcodeLocation]ap.AssertionPayload[T, E], numAssertMessages)
	for i := uint64(0); i < numAssertMessages; i++ {
		var opcodeLocation ops.OpcodeLocation
		if err := opcodeLocation.UnmarshalReader(r); err != nil {
			return err
		}
		var payload ap.AssertionPayload[T, E]
		if err := payload.UnmarshalReader(r); err != nil {
			return err
		}
		c.AssertMessages[opcodeLocation] = payload
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
	currentWitnesses := c.collectCurrentWitnesses(witnesses, index)

	// 3. Add the constraints for the circuit
	if err := c.constrainCircuit(api, currentWitnesses); err != nil {
		return nil, nil, err
	}

	// 4. Connect call inputs/outputs
	c.constrainCircuitCalls(api, currentWitnesses, callConnections)

	// 5. Collect circuit inputs and outputs
	inputs := c.collectWitnesses(&c.PrivateParameters, currentWitnesses)
	outputs := c.collectWitnesses(&c.ReturnValues, currentWitnesses)

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
func (c *Circuit[T, E]) collectCurrentWitnesses(witnesses map[shr.Witness]frontend.Variable, index *uint32) map[shr.Witness]frontend.Variable {
	currentWitnesses := make(map[shr.Witness]frontend.Variable, c.CurrentWitnessIndex+1)

	for i := range c.CurrentWitnessIndex + 1 {
		v, ok := witnesses[shr.Witness(i+uint32(*index))]
		if !ok {
			// Sometimes circuits skip an index.
			// Insert dummy witness variable to make the number of witnesses
			// consistent with the `CurrentWitnessIndex` in the ACIR of the circuit
			currentWitnesses[shr.Witness(i)] = frontend.Variable(0)
			continue
		}
		currentWitnesses[shr.Witness(i)] = v
	}

	*index += c.CurrentWitnessIndex + 1
	return currentWitnesses
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

// Construct a list of input/ output variables of a circuit given a tree of witness indices and a index->variable mapping
func (c *Circuit[T, E]) collectWitnesses(tree *btree.BTree, currentWitnesses map[shr.Witness]frontend.Variable) []frontend.Variable {
	var vars []frontend.Variable
	tree.Ascend(func(it btree.Item) bool {
		witness, ok := it.(shr.Witness)
		if !ok {
			return false
		}
		vars = append(vars, currentWitnesses[witness])
		return true
	})
	return vars
}

// FillWitness adds the witnesses used by a circuit incremented by the starting index
func (c *Circuit[T, E]) FillWitnessTree(witnessTree *btree.BTree, resolve CircuitResolver[T, E], index uint32) (uint32, error) {
	if witnessTree == nil {
		return index, fmt.Errorf("no witness tree to fill")
	}

	for _, opcode := range c.Opcodes {
		if callOp, ok := opcode.(*call.Call[T, E]); ok {
			subCircuit, err := resolve(callOp.ID)
			if err != nil {
				return index, fmt.Errorf("failed to resolve circuit %d: %w", callOp.ID, err)
			}
			subCircuitWitnessTree := btree.New(2)
			subCircuit.FillWitnessTree(subCircuitWitnessTree, resolve, index)

			subCircuitWitnessTree.Ascend(func(it btree.Item) bool {
				witness, ok := it.(shr.Witness)
				if !ok {
					panic("Item in subwitness tree not of type witness")
				}
				witnessTree.ReplaceOrInsert(witness)
				index++
				return true
			})
		}
	}
	for _, opcode := range c.Opcodes {
		opcode.FillWitnessTree(witnessTree, index)
	}
	return index, nil
}

func NewOpcode[T shr.ACIRField, E constraint.Element](r io.Reader) (ops.Opcode[E], error) {
	var kind uint32
	if err := binary.Read(r, binary.LittleEndian, &kind); err != nil {
		return nil, err
	}
	switch kind {
	case 0:
		return &exp.Expression[T, E]{}, nil
	case 1:
		bbf, err := bbf.NewBlackBoxFunction[T, E](r)
		if err != nil {
			return nil, fmt.Errorf("unable to get opcode, error with black box:  %v", err)
		}
		return bbf, nil
	case 2:
		return &mem_op.MemoryOp[T, E]{}, nil
	case 3:
		return &memory_init.MemoryInit[T, E]{}, nil
	case 4:
		return &brillig.BrilligCall[T, E]{}, nil
	case 5:
		return &call.Call[T, E]{}, nil
	default:
		return nil, fmt.Errorf("unknown opcode kind: %d", kind)
	}
}
