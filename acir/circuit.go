package acir

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	exp "nr-groth16/acir/expression"
	ops "nr-groth16/acir/opcodes"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
	"github.com/rs/zerolog/log"
)

type Circuit[T shr.ACIRField] struct {
	CurrentWitnessIndex uint32                                     `json:"current_witness_index"`
	Opcodes             []ops.Opcode[T]                            `json:"opcodes"`            // Opcodes in the circuit
	ExpressionWidth     exp.ExpressionWidth                        `json:"expression_width"`   // Width of the expressions in the circuit
	PrivateParameters   btree.BTree                                `json:"private_parameters"` // Witnesses
	PublicParameters    btree.BTree                                `json:"public_parameters"`  // Witnesses
	ReturnValues        btree.BTree                                `json:"return_values"`      // Witnesses
	AssertMessages      map[ops.OpcodeLocation]AssertionPayload[T] `json:"assert_messages"`    // Assert messages for the circuit
	Recursive           bool                                       `json:"recursive"`          // Whether the circuit is recursive
}

func (c *Circuit[T]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &c.CurrentWitnessIndex); err != nil {
		return err
	}
	log.Trace().Msg("Unmarshalling Circuit with current witness index: " + fmt.Sprint(c.CurrentWitnessIndex))

	var numOpcodes uint64
	if err := binary.Read(r, binary.LittleEndian, &numOpcodes); err != nil {
		return err
	}
	log.Trace().Msg("Unmarshalling Circuit with number of opcodes: " + fmt.Sprint(numOpcodes))

	c.Opcodes = make([]ops.Opcode[T], numOpcodes)
	for i := uint64(0); i < numOpcodes; i++ {
		log.Trace().Msg("Unmarshalling Opcode at index: " + fmt.Sprint(i))
		if err := c.Opcodes[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	log.Trace().Msg("Unmarshalling expression width")
	if err := c.ExpressionWidth.UnmarshalReader(r); err != nil {
		return err
	}

	var numPrivateParameters uint64
	if err := binary.Read(r, binary.LittleEndian, &numPrivateParameters); err != nil {
		return err
	}
	log.Trace().Msg("Unmarshalling Circuit with number of private parameters: " + fmt.Sprint(numPrivateParameters))
	c.PrivateParameters = *btree.New(2)
	for i := uint64(0); i < numPrivateParameters; i++ {
		log.Trace().Msg("Unmarshalling PrivateParameter at index: " + fmt.Sprint(i))
		var witness shr.Witness
		if err := witness.UnmarshalReader(r); err != nil {
			return err
		}
		log.Trace().Msg("Unmarshalling PrivateParameter: " + fmt.Sprint(witness))
		c.PrivateParameters.ReplaceOrInsert(witness)
	}

	var numPublicParameters uint64
	if err := binary.Read(r, binary.LittleEndian, &numPublicParameters); err != nil {
		return err
	}
	log.Trace().Msg("Unmarshalling Circuit with number of public parameters: " + fmt.Sprint(numPublicParameters))
	c.PublicParameters = *btree.New(2)
	for i := uint64(0); i < numPublicParameters; i++ {
		log.Trace().Msg("Unmarshalling PublicParameter at index: " + fmt.Sprint(i))
		var witness shr.Witness
		if err := witness.UnmarshalReader(r); err != nil {
			return err
		}
		log.Trace().Msg("Unmarshalling PublicParameter: " + fmt.Sprintf("%x", witness))
		c.PublicParameters.ReplaceOrInsert(witness)
	}

	var numReturnValues uint32
	if err := binary.Read(r, binary.LittleEndian, &numReturnValues); err != nil {
		return err
	}
	log.Trace().Msg("Unmarshalling Circuit with number of return values: " + fmt.Sprint(numReturnValues))
	c.ReturnValues = *btree.New(2)
	for i := uint32(0); i < numReturnValues; i++ {
		log.Trace().Msg("Unmarshalling ReturnValue at index: " + fmt.Sprint(i))
		var witness shr.Witness
		if err := witness.UnmarshalReader(r); err != nil {
			return err
		}
		log.Trace().Msg("Unmarshalling ReturnValue: " + fmt.Sprintf("%x", witness))
		c.ReturnValues.ReplaceOrInsert(witness)
	}

	log.Trace().Msg("Unmarshalling Circuit with assert messages")

	var numAssertMessages uint64
	if err := binary.Read(r, binary.LittleEndian, &numAssertMessages); err != nil {
		if err == io.EOF {
			log.Trace().Msg("No assert messages found, continuing without them")
			c.AssertMessages = make(map[ops.OpcodeLocation]AssertionPayload[T])
			return nil
		}
		log.Trace().Msg("Error reading number of assert messages: " + err.Error())
	}
	log.Trace().Msg("Unmarshalling Circuit with number of assert messages: " + fmt.Sprintf("%x %d", numAssertMessages, numAssertMessages))

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

	var recursiveFlag uint8
	if err := binary.Read(r, binary.LittleEndian, &recursiveFlag); err != nil {
		if err == io.EOF {
			log.Trace().Msg("No recursive flag found, continuing without it")
			c.Recursive = false
			return nil
		}
		return err
	}
	c.Recursive = recursiveFlag != 0
	log.Trace().Msg("Unmarshalling Circuit with recursive flag: " + fmt.Sprint(c.Recursive))

	return nil
}

func (c *Circuit[T]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {
	for _, opcode := range c.Opcodes {
		//if index != 1 || index == 4 || index == 5 || index == 6 || index == 8 {
		//	continue
		//}
		if err := opcode.Define(api, witnesses); err != nil {
			return err
		}
	}
	return nil
}

func (c *Circuit[T]) FillWitnessTree(witnessTree *btree.BTree) {
	if witnessTree == nil {
		return
	}

	for _, opcode := range c.Opcodes {
		opcode.FillWitnessTree(witnessTree)
	}
}

func (c *Circuit[T]) CollectConstantsAsWitnesses(start uint32, witnessTree *btree.BTree) {
	if witnessTree == nil {
		return
	}

	for _, opcode := range c.Opcodes {
		opcode.CollectConstantsAsWitnesses(start, witnessTree)
	}
}

func (c *Circuit[T]) FeedConstantsAsWitnesses() []*big.Int {
	values := make([]*big.Int, 0)

	for _, opcode := range c.Opcodes {
		values = append(values, opcode.FeedConstantsAsWitnesses()...)
	}
	return values
}
