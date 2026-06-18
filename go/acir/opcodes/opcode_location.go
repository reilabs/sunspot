package opcodes

import (
	"encoding/json"
	"fmt"
	"sunspot/go/acir/msgpackutil"
)

type OpcodeLocation struct {
	ACIRAddress  *uint64
	ACIRIndex    *uint64
	BrilligIndex *uint64
}

func (o OpcodeLocation) MarshalJSON() ([]byte, error) {
	fieldsMap := make(map[string]interface{})

	if o.ACIRAddress != nil {
		fieldsMap["ACIRAddress"] = *o.ACIRAddress
	} else if o.ACIRIndex != nil {
		fieldsMap["ACIRIndex"] = *o.ACIRIndex
	} else if o.BrilligIndex != nil {
		fieldsMap["BrilligIndex"] = *o.BrilligIndex
	} else {
		return nil, fmt.Errorf("unknown OpcodeLocation Kind")
	}
	return json.Marshal(fieldsMap)
}

// On the wire OpcodeLocation is an int-keyed single-entry fixmap whose tag
// selects the variant (0=Acir(usize), 1=Brillig{acir_index, brillig_index}).
// The Brillig payload itself is a tagged struct (or positional 2-array under
// EncodingStrategy::Array, which is the active strategy for non-Program
// types) — see acvm-repo/acir/src/circuit/mod.rs in noir.
func (o *OpcodeLocation) UnmarshalReader(r *msgpackutil.Reader) error {
	return msgpackutil.ReadEnum(r, opcodeLocationSchema, o.decode)
}

func (o *OpcodeLocation) decode(f msgpackutil.Field, r *msgpackutil.Reader) error {
	switch f.Tag {
	case 0:

		v, err := r.ReadUint()
		if err != nil {
			return err
		}
		o.ACIRAddress = new(uint64)
		*o.ACIRAddress = v
		return nil
	case 1:
		o.ACIRIndex = new(uint64)
		o.BrilligIndex = new(uint64)
		return msgpackutil.ReadStruct(r, opcodeLocationBrilligSchema, func(fld msgpackutil.Field, r *msgpackutil.Reader) error {
			switch fld.Tag {
			case 0:
				v, err := r.ReadUint()
				if err != nil {
					return err
				}
				*o.ACIRIndex = v
				return nil
			case 1:
				v, err := r.ReadUint()
				if err != nil {
					return err
				}
				*o.BrilligIndex = v
				return nil
			default:
				return fmt.Errorf("OpcodeLocation.Brillig: unknown field %s", fld)
			}
		})
	default:
		return fmt.Errorf("OpcodeLocation: unknown variant %v", f)
	}
}

var opcodeLocationSchema = msgpackutil.NewSchema(map[string]int{
	"Acir": 0, "Brillig": 1,
})

// Brillig variant inner struct fields.
var opcodeLocationBrilligSchema = msgpackutil.NewSchema(map[string]int{
	"acir_index": 0, "brillig_index": 1,
})

func (o *OpcodeLocation) Equals(other *OpcodeLocation) bool {
	return uint64PtrEquals(o.ACIRAddress, other.ACIRAddress) &&
		uint64PtrEquals(o.ACIRIndex, other.ACIRIndex) &&
		uint64PtrEquals(o.BrilligIndex, other.BrilligIndex)
}

func uint64PtrEquals(a, b *uint64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
