package opcodes

import (
	"encoding/json"
	"fmt"
	"sunspot/go/acir/msgpackutil"
)

type OpcodeLocationKind uint8

const (
	OpcodeLocationAcir OpcodeLocationKind = iota
	OpcodeLocationBrillig
)

type OpcodeLocation struct {
	Kind         OpcodeLocationKind
	ACIRAddress  uint64 // Acir variant
	ACIRIndex    uint64 // Brillig variant
	BrilligIndex uint64 // Brillig variant
}

func (o OpcodeLocation) MarshalJSON() ([]byte, error) {
	fieldsMap := make(map[string]interface{})
	switch o.Kind {
	case OpcodeLocationAcir:
		fieldsMap["ACIRAddress"] = o.ACIRAddress
	case OpcodeLocationBrillig:
		fieldsMap["ACIRIndex"] = o.ACIRIndex
		fieldsMap["BrilligIndex"] = o.BrilligIndex
	default:
		return nil, fmt.Errorf("unknown OpcodeLocation Kind: %d", o.Kind)
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
		o.Kind = OpcodeLocationAcir
		o.ACIRAddress = v
		return nil
	case 1:
		o.Kind = OpcodeLocationBrillig
		return msgpackutil.ReadStruct(r, opcodeLocationBrilligSchema, func(fld msgpackutil.Field, r *msgpackutil.Reader) error {
			switch fld.Tag {
			case 0:
				v, err := r.ReadUint()
				if err != nil {
					return err
				}
				o.ACIRIndex = v
				return nil
			case 1:
				v, err := r.ReadUint()
				if err != nil {
					return err
				}
				o.BrilligIndex = v
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
