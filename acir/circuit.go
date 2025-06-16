package acir

import (
	"encoding/binary"
	"io"
	exp "nr-groth16/acir/expression"
	ops "nr-groth16/acir/opcodes"
	shr "nr-groth16/acir/shared"

	"github.com/consensys/gnark/frontend"
	"github.com/google/btree"
)

type Circuit[T shr.ACIRField] struct {
	CurrentWitnessIndex uint32
	Opcodes             []ops.Opcode[T]
	ExpressionWidth     exp.ExpressionWidth
	PrivateParameters   btree.BTree // Witnesses
	PublicParameters    btree.BTree // Witnesses
	ReturnValues        btree.BTree // Witnesses
	AssertMessages      map[ops.OpcodeLocation]AssertionPayload[T]
}

func (c *Circuit[T]) UnmarshalReader(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &c.CurrentWitnessIndex); err != nil {
		return err
	}

	var numOpcodes uint64
	if err := binary.Read(r, binary.LittleEndian, &numOpcodes); err != nil {
		return err
	}
	c.Opcodes = make([]ops.Opcode[T], numOpcodes)
	for i := uint64(0); i < numOpcodes; i++ {
		if err := c.Opcodes[i].UnmarshalReader(r); err != nil {
			return err
		}
	}

	var numPrivateParameters uint32
	if err := binary.Read(r, binary.LittleEndian, &numPrivateParameters); err != nil {
		return err
	}
	c.PrivateParameters = *btree.New(2)
	for i := uint32(0); i < numPrivateParameters; i++ {
		var witness shr.Witness
		if err := witness.UnmarshalReader(r); err != nil {
			return err
		}
		c.PrivateParameters.ReplaceOrInsert(witness)
	}

	var numPublicParameters uint32
	if err := binary.Read(r, binary.LittleEndian, &numPublicParameters); err != nil {
		return err
	}
	c.PublicParameters = *btree.New(2)
	for i := uint32(0); i < numPublicParameters; i++ {
		var witness shr.Witness
		if err := witness.UnmarshalReader(r); err != nil {
			return err
		}
		c.PublicParameters.ReplaceOrInsert(witness)
	}

	var numReturnValues uint32
	if err := binary.Read(r, binary.LittleEndian, &numReturnValues); err != nil {
		return err
	}
	c.ReturnValues = *btree.New(2)
	for i := uint32(0); i < numReturnValues; i++ {
		var witness shr.Witness
		if err := witness.UnmarshalReader(r); err != nil {
			return err
		}
		c.ReturnValues.ReplaceOrInsert(witness)
	}

	var numAssertMessages uint32
	if err := binary.Read(r, binary.LittleEndian, &numAssertMessages); err != nil {
		return err
	}
	c.AssertMessages = make(map[ops.OpcodeLocation]AssertionPayload[T], numAssertMessages)
	for i := uint32(0); i < numAssertMessages; i++ {
		var opcodeLocation ops.OpcodeLocation
		if err := opcodeLocation.UnmarshalReader(r); err != nil {
			return err
		}
		var payload AssertionPayload[T]
		if err := payload.UnmarshalReader(r); err != nil {
			return err
		}
		c.AssertMessages[opcodeLocation] = payload
	}
	return nil
}

func (c *Circuit[T]) Define(api frontend.API, witnesses map[shr.Witness]frontend.Variable) error {
	for _, opcode := range c.Opcodes {
		if err := opcode.Define(api, witnesses); err != nil {
			return err
		}
	}
	return nil
}
