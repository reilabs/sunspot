package acir

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	bbf "nr-groth16/acir/black_box_func"
	"nr-groth16/acir/brillig"
	"nr-groth16/acir/call"
	exp "nr-groth16/acir/expression"
	"nr-groth16/acir/memory_init"
	mem_op "nr-groth16/acir/memory_op"
	ops "nr-groth16/acir/opcodes"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
	"github.com/google/btree"
)

type Circuit[T shr.ACIRField, E constraint.Element] struct {
	CurrentWitnessIndex uint32                                        `json:"current_witness_index"`
	Opcodes             []ops.Opcode[E]                               `json:"opcodes"`            // Opcodes in the circuit
	ExpressionWidth     exp.ExpressionWidth                           `json:"expression_width"`   // Width of the expressions in the circuit
	PrivateParameters   btree.BTree                                   `json:"private_parameters"` // Witnesses
	PublicParameters    btree.BTree                                   `json:"public_parameters"`  // Witnesses
	ReturnValues        btree.BTree                                   `json:"return_values"`      // Witnesses
	AssertMessages      map[ops.OpcodeLocation]AssertionPayload[T, E] `json:"assert_messages"`    // Assert messages for the circuit
	MemoryBlocks        map[uint32]*logderivlookup.Table
}

func (c *Circuit[T, E]) UnmarshalReader(r io.Reader) error {
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

	if err := c.ExpressionWidth.UnmarshalReader(r); err != nil {
		return err
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
			c.AssertMessages = make(map[ops.OpcodeLocation]AssertionPayload[T, E])
			return nil
		}
	}

	/*c.AssertMessages = make(map[ops.OpcodeLocation]AssertionPayload[T], numAssertMessages)
	for i := uint64(0); i < numAssertMessages; i++ {
		log.Trace().Msg("Unmarshalling AssertMessage at index: " + fmt.Sprint(i))
		var opcodeLocation ops.OpcodeLocation
		if err := opcodeLocation.UnmarshalReader(r); err != nil {
			return err
		}
		var payload AssertionPayload[T]
		if err := payload.UnmarshalReader(r); err != nil {
			return err
		}
		c.AssertMessages[opcodeLocation] = payload
	}*/

	return nil
}

func (c *Circuit[T, E]) Define(api frontend.Builder[E], witnesses map[shr.Witness]frontend.Variable, resolve CircuitResolver[T, E]) error {
	c.MemoryBlocks = make(map[uint32]*logderivlookup.Table)
	for _, opcode := range c.Opcodes {

		mem_init, ok := opcode.(*memory_init.MemoryInit[T, E])
		if ok {
			table := logderivlookup.New(api)
			mem_init.Table = &table
			c.MemoryBlocks[mem_init.BlockID] = &table
		}
		mem_op, ok := opcode.(*mem_op.MemoryOp[T, E])
		if ok {
			mem_op.Memory = c.MemoryBlocks
		}

		// --- Nested circuit calls ---
		if callOp, ok := opcode.(*call.Call[T, E]); ok {
			subCircuit, err := resolve(callOp.ID)
			if err != nil {
				return fmt.Errorf("failed to resolve circuit %d: %w", callOp.ID, err)
			}

			// Prepare subcircuit input witnesses
			subWitnesses := make(map[shr.Witness]frontend.Variable, len(callOp.Inputs))
			for i := range callOp.Inputs {
				inputWitness := callOp.Inputs[i]
				v, ok := witnesses[inputWitness]
				if !ok {
					return fmt.Errorf("call: missing input witness %d", inputWitness)
				}
				subWitnesses[shr.Witness(i)] = v
			}

			// Run subcircuit definition
			if err := subCircuit.Define(api, subWitnesses, resolve); err != nil {
				return fmt.Errorf("failed to define subcircuit %d: %w", callOp.ID, err)
			}

			// Collect subcircuit outputs (these are stored in subCircuit.ReturnValues)
			idx := 0
			subCircuit.ReturnValues.Ascend(func(i btree.Item) bool {
				if idx >= len(callOp.Outputs) {
					return false
				}
				w := i.(shr.Witness)
				v := subWitnesses[w]

				// TODO: Replace commented code with something that legit works
				// Apply predicate if applicable
				// if callOp.Predicate != nil {
				// 	v = api.Select(predicate, v, 0)
				// }

				witnesses[callOp.Outputs[idx]] = v
				idx++
				return true
			})

			// Continue with next opcode
			continue
		}

		if err := opcode.Define(api, witnesses); err != nil {
			return err
		}
	}
	return nil
}

func (c *Circuit[T, E]) FillWitnessTree(witnessTree *btree.BTree) {
	if witnessTree == nil {
		return
	}

	for _, opcode := range c.Opcodes {
		opcode.FillWitnessTree(witnessTree)
	}
}

func (c *Circuit[T, E]) CollectConstantsAsWitnesses(start uint32, witnessTree *btree.BTree) {
	if witnessTree == nil {
		return
	}

	for _, opcode := range c.Opcodes {
		opcode.CollectConstantsAsWitnesses(start, witnessTree)
	}
}

func (c *Circuit[T, E]) FeedConstantsAsWitnesses() []*big.Int {
	values := make([]*big.Int, 0)

	for _, opcode := range c.Opcodes {
		values = append(values, opcode.FeedConstantsAsWitnesses()...)
	}
	return values
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
