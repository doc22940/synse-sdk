// Synse SDK
// Copyright (c) 2019 Vapor IO
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package output

import (
	"fmt"

	synse "github.com/vapor-ware/synse-server-grpc/go"
)

// Reading describes a single device reading at a given time. The timestamp
// for a reading is represented using the RFC3339 layout.
type Reading struct {
	// Timestamp is the RFC3339-formatted time at which the reading was taken.
	Timestamp string

	// Type is the type of the reading, as defined by the Reading's output.
	Type string

	// Info provides additional information about a reading.
	Info string

	// Unit describes the unit of measure for the reading.
	Unit *Unit

	// Value is the reading value itself.
	Value interface{}

	// output is the Output used to render and format the reading.
	output *Output
}

// GetOutput gets the associated output for a Reading.
func (reading *Reading) GetOutput() *Output {
	return reading.output
}

// Encode translates the Reading to its corresponding gRPC message.
func (reading *Reading) Encode() *synse.V3Reading {
	var unit = &Unit{}
	if reading.Unit != nil {
		unit = reading.Unit
	}

	r := synse.V3Reading{
		Timestamp: reading.Timestamp,
		Type:      reading.Type,
		Context:   map[string]string{}, // todo: adding context to reading
		Unit:      unit.Encode(),
	}

	if reading.Info != "" {
		r.Context["info"] = reading.Info
	}

	switch t := reading.Value.(type) {
	case string:
		r.Value = &synse.V3Reading_StringValue{StringValue: t}
	case bool:
		r.Value = &synse.V3Reading_BoolValue{BoolValue: t}
	case float64:
		r.Value = &synse.V3Reading_Float64Value{Float64Value: t}
	case float32:
		r.Value = &synse.V3Reading_Float32Value{Float32Value: t}
	case int64:
		r.Value = &synse.V3Reading_Int64Value{Int64Value: t}
	case int32:
		r.Value = &synse.V3Reading_Int32Value{Int32Value: t}
	case int16:
		r.Value = &synse.V3Reading_Int32Value{Int32Value: int32(t)}
	case int8:
		r.Value = &synse.V3Reading_Int32Value{Int32Value: int32(t)}
	case int:
		r.Value = &synse.V3Reading_Int64Value{Int64Value: int64(t)}
	case []byte:
		r.Value = &synse.V3Reading_BytesValue{BytesValue: t}
	case uint64:
		r.Value = &synse.V3Reading_Uint64Value{Uint64Value: t}
	case uint32:
		r.Value = &synse.V3Reading_Uint32Value{Uint32Value: t}
	case uint16:
		r.Value = &synse.V3Reading_Uint32Value{Uint32Value: uint32(t)}
	case uint8:
		r.Value = &synse.V3Reading_Uint32Value{Uint32Value: uint32(t)}
	case uint:
		r.Value = &synse.V3Reading_Uint64Value{Uint64Value: uint64(t)}
	case nil:
		r.Value = nil
	default:
		// If the reading type isn't one of the above, panic. The plugin should
		// terminate. This is indicative of the plugin is providing bad data.
		panic(fmt.Sprintf("unsupported reading value type: %s", t))
	}
	return &r
}
