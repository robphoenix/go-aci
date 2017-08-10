package aci

import (
	"fmt"
	"regexp"
)

// location
type location struct {
	name        string
	description string
	locations   []*location
	status      string
}

func (l *location) String() string {
	if l == nil {
		return ""
	}
	s := ""
	s += l.name
	s += "\n"
	for _, location := range l.locations {
		s += "\t"
		s += location.String()
		s += "\n"
	}
	return s
}

// SetName validates and sets the name of a location
//
// A site name can be up to 64 characters and must begin with an
// alphanumeric character. It can only contain alphanumeric
// characters and the following symbols: -.:_
func (l *location) SetName(name string) error {
	valid := regexp.MustCompile(`(^[a-zA-Z0-9]{1,2}$|^[a-zA-Z0-9][a-zA-Z0-9-.:_]{0,63}$)`)
	if !valid.MatchString(name) {
		return fmt.Errorf("invalid name: %s", name)
	}
	l.name = name
	return nil
}

// Name returns the name of a location.
func (l *location) Name() string {
	return l.name
}

// SetDescription validates and sets the description of a location
//
// A site description can be up to 128 characters. It can only
// contain alphanumeric characters and the following symbols:
// !#$%()*,-.:;@_{|}~?&+
func (l *location) SetDescription(description string) error {
	valid := regexp.MustCompile(`^[a-zA-Z0-9!#$%()*,-.:;@_{|}~?&+\s]{0,128}$`)
	if !valid.MatchString(description) {
		return fmt.Errorf("invalid description: %s", description)
	}
	l.description = description
	return nil
}

// Description resturns a location description.
func (l *location) Description() string {
	return l.description
}

// SetCreated sets the status of the location to "created,modified".
func (l *location) SetCreated() string {
	l.status = createdModified
	return l.status
}

// SetDeleted sets the status of the location to "deleted".
func (l *location) SetDeleted() string {
	l.status = deleted
	return l.status
}

// Status returns the status of the location.
func (l *location) Status() string {
	return l.status
}

func (l *location) addLocation(loc location) {
	l.locations = append(l.locations, &loc)
}

func (l *location) deleteLocation(loc location) {
	var index int
	for i, ll := range l.locations {
		if ll == &loc {
			index = i
			break
		}
	}
	l.locations = append(l.locations[:index], l.locations[index+1:]...)
}

// Site ...
type Site struct {
	location
}

// Buildings ...
func (site *Site) Buildings() []*Building {
	var buildings []*Building
	for _, location := range site.locations {
		buildings = append(buildings, &Building{*location})
	}
	return buildings
}

// AddBuilding ...
func (site *Site) AddBuilding(building *Building) {
	site.addLocation(building.location)
}

// DeleteBuilding ...
func (site *Site) DeleteBuilding(building *Building) {
	site.deleteLocation(building.location)
}

// Building ...
type Building struct {
	location
}

// Floors ...
func (building *Building) Floors() []*Floor {
	var floors []*Floor
	for _, location := range building.locations {
		floors = append(floors, &Floor{*location})
	}
	return floors
}

// AddFloor ...
func (building *Building) AddFloor(floor *Floor) {
	building.addLocation(floor.location)
}

// DeleteFloor ...
func (building *Building) DeleteFloor(floor *Floor) {
	building.deleteLocation(floor.location)
}

// Floor ...
type Floor struct {
	location
}

// Rooms ...
func (floor *Floor) Rooms() []*Room {
	var rooms []*Room
	for _, location := range floor.locations {
		rooms = append(rooms, &Room{*location})
	}
	return rooms
}

// AddRoom ...
func (floor *Floor) AddRoom(room *Room) {
	floor.addLocation(room.location)
}

// DeleteRoom ...
func (floor *Floor) DeleteRoom(room *Room) {
	floor.deleteLocation(room.location)
}

// Room ...
type Room struct {
	location
}

// Rows ...
func (room *Room) Rows() []*Row {
	var rows []*Row
	for _, location := range room.locations {
		rows = append(rows, &Row{*location})
	}
	return rows
}

// AddRow ...
func (room *Room) AddRow(row *Row) {
	room.addLocation(row.location)
}

// DeleteRow ...
func (room *Room) DeleteRow(row *Row) {
	room.deleteLocation(row.location)
}

// Row ...
type Row struct {
	location
}

// Racks ...
func (row *Row) Racks() []*Rack {
	var racks []*Rack
	for _, location := range row.locations {
		racks = append(racks, &Rack{*location})
	}
	return racks
}

// AddRack ...
func (row *Row) AddRack(rack *Rack) {
	row.addLocation(rack.location)
}

// DeleteRack ...
func (row *Row) DeleteRack(rack *Rack) {
	row.deleteLocation(rack.location)
}

// Rack ...
type Rack struct {
	location
}
