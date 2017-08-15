package aci

import (
	"context"
	"fmt"
	"net/http"
)

// FabricInstanceContainer ...
type FabricInstanceContainer struct {
	FabricInstance `json:"fabricInst"`
}

// FabricInstance ...
type FabricInstance struct {
	GeoAttrs `json:"attributes"`
	GeoSites []GeoSiteContainer `json:"children"`
}

// GeoSiteContainer ...
type GeoSiteContainer struct {
	GeoSite `json:"geoSite,omitempty"`
}

func newGeoSiteContainer(site *Site) GeoSiteContainer {
	var geoSite GeoSiteContainer
	geoSite.Name = site.Name()
	geoSite.Descr = site.Description()
	geoSite.DN = fmt.Sprintf("uni/fabric/site-%s", site.Name())
	geoSite.RN = fmt.Sprintf("site-%s", site.Name())
	geoSite.Status = site.Status()

	buildings := site.Buildings()
	if len(buildings) == 0 {
		return geoSite
	}

	var geoBuildings []GeoBuildingContainer
	for _, building := range buildings {
		geoBuilding := newGeoBuildingContainer(site.Name(), building)
		geoBuildings = append(geoBuildings, geoBuilding)
	}

	geoSite.GeoBuildings = geoBuildings

	return geoSite
}

// GeoSite ...
type GeoSite struct {
	GeoAttrs     `json:"attributes,omitempty"`
	GeoBuildings []GeoBuildingContainer `json:"children,omitempty"`
}

// GeoBuildingContainer ...
type GeoBuildingContainer struct {
	GeoBuilding `json:"geoBuilding,omitempty"`
}

func newGeoBuildingContainer(site string, building *Building) GeoBuildingContainer {
	var geoBuilding GeoBuildingContainer
	geoBuilding.Name = building.Name()
	geoBuilding.Descr = building.Description()
	geoBuilding.DN = fmt.Sprintf("uni/fabric/site-%s/building-%s", site, building.Name())
	geoBuilding.RN = fmt.Sprintf("building-%s", building.Name())
	geoBuilding.Status = building.Status()

	floors := building.Floors()
	if len(floors) == 0 {
		return geoBuilding
	}

	var geoFloors []GeoFloorContainer
	for _, floor := range floors {
		geoFloor := newGeoFloorContainer(site, building.Name(), floor)
		geoFloors = append(geoFloors, geoFloor)
	}

	geoBuilding.GeoFloors = geoFloors

	return geoBuilding
}

// GeoBuilding ...
type GeoBuilding struct {
	GeoAttrs  `json:"attributes,omitempty"`
	GeoFloors []GeoFloorContainer `json:"children,omitempty"`
}

// GeoFloorContainer ...
type GeoFloorContainer struct {
	GeoFloor `json:"geoFloor,omitempty"`
}

func newGeoFloorContainer(site, building string, floor *Floor) GeoFloorContainer {
	var geoFloor GeoFloorContainer
	geoFloor.Name = floor.Name()
	geoFloor.Descr = floor.Description()
	geoFloor.DN = fmt.Sprintf("uni/fabric/site-%s/building-%s/floor-%s", site, building, floor.Name())
	geoFloor.RN = fmt.Sprintf("floor-%s", floor.Name())
	geoFloor.Status = floor.Status()

	rooms := floor.Rooms()
	if len(rooms) == 0 {
		return geoFloor
	}

	var geoRooms []GeoRoomContainer
	for _, room := range rooms {
		geoRoom := newGeoRoomContainer(site, building, floor.Name(), room)
		geoRooms = append(geoRooms, geoRoom)
	}

	geoFloor.GeoRooms = geoRooms

	return geoFloor
}

// GeoFloor ...
type GeoFloor struct {
	GeoAttrs `json:"attributes,omitempty"`
	GeoRooms []GeoRoomContainer `json:"children,omitempty"`
}

// GeoRoomContainer ...
type GeoRoomContainer struct {
	GeoRoom `json:"geoRoom,omitempty"`
}

func newGeoRoomContainer(site, building, floor string, room *Room) GeoRoomContainer {
	var geoRoom GeoRoomContainer
	geoRoom.Name = room.Name()
	geoRoom.Descr = room.Description()
	geoRoom.DN = fmt.Sprintf("uni/fabric/site-%s/building-%s/floor-%s/room-%s", site, building, floor, room.Name())
	geoRoom.RN = fmt.Sprintf("room-%s", room.Name())
	geoRoom.Status = room.Status()

	rows := room.Rows()
	if len(rows) == 0 {
		return geoRoom
	}

	var geoRows []GeoRowContainer
	for _, row := range rows {
		geoRow := newGeoRowContainer(site, building, floor, room.Name(), row)
		geoRows = append(geoRows, geoRow)
	}

	geoRoom.GeoRows = geoRows

	return geoRoom
}

// GeoRoom ...
type GeoRoom struct {
	GeoAttrs `json:"attributes,omitempty"`
	GeoRows  []GeoRowContainer `json:"children,omitempty"`
}

// GeoRowContainer ...
type GeoRowContainer struct {
	GeoRow `json:"geoRow,omitempty"`
}

func newGeoRowContainer(site, building, floor, room string, row *Row) GeoRowContainer {
	var geoRow GeoRowContainer
	geoRow.Name = row.Name()
	geoRow.Descr = row.Description()
	geoRow.DN = fmt.Sprintf("uni/fabric/site-%s/building-%s/floor-%s/room-%s/row-%s", site, building, floor, room, row.Name())
	geoRow.RN = fmt.Sprintf("row-%s", row.Name())
	geoRow.Status = row.Status()

	racks := row.Racks()
	if len(racks) == 0 {
		return geoRow
	}

	var geoRacks []GeoRackContainer
	for _, rack := range racks {
		geoRack := newGeoRackContainer(site, building, floor, room, row.Name(), rack)
		geoRacks = append(geoRacks, geoRack)
	}

	geoRow.GeoRacks = geoRacks

	return geoRow
}

// GeoRow ...
type GeoRow struct {
	GeoAttrs `json:"attributes,omitempty"`
	GeoRacks []GeoRackContainer `json:"children,omitempty"`
}

// GeoRackContainer ...
type GeoRackContainer struct {
	GeoRack `json:"geoRack,omitempty"`
}

func newGeoRackContainer(site, building, floor, room, row string, rack *Rack) GeoRackContainer {
	var geoRack GeoRackContainer
	geoRack.Name = rack.Name()
	geoRack.Descr = rack.Description()
	geoRack.DN = fmt.Sprintf("uni/fabric/site-%s/building-%s/floor-%s/room-%s/row-%s/rack-%s", site, building, floor, room, row, rack.Name())
	geoRack.RN = fmt.Sprintf("rack-%s", rack.Name())
	geoRack.Status = rack.Status()

	return geoRack
}

// GeoRack ...
type GeoRack struct {
	GeoAttrs `json:"attributes,omitempty"`
	GeoNodes []interface{} `json:"children,omitempty"`
}

// GeoAttrs ...
type GeoAttrs struct {
	Descr  string `json:"descr,omitempty"`
	DN     string `json:"dn,omitempty"`
	Name   string `json:"name,omitempty"`
	RN     string `json:"rn,omitempty"`
	Status string `json:"status,omitempty"`
}

// GeolocationResponse ...
type GeolocationResponse struct {
	TotalCount string        `json:"totalCount"`
	Imdata     []interface{} `json:"imdata"`
}

// GeolocationService handles communication with the geolocation related
// methods of the APIC API.
type GeolocationService service

// NewSite instantiates a Site.
func (s *GeolocationService) NewSite(name, description string) (*Site, error) {
	site := Site{}
	err := site.SetName(name)
	if err != nil {
		return &site, err
	}
	if description != "" {
		err := site.SetDescription(description)
		if err != nil {
			return &site, err
		}
	}
	return &site, nil
}

// NewBuilding instantiates a Building.
func (s *GeolocationService) NewBuilding(name, description string) (*Building, error) {
	building := Building{}
	err := building.SetName(name)
	if err != nil {
		return &building, err
	}
	if description != "" {
		err := building.SetDescription(description)
		if err != nil {
			return &building, err
		}
	}
	return &building, nil
}

// NewFloor instantiates a Floor.
func (s *GeolocationService) NewFloor(name, description string) (*Floor, error) {
	floor := Floor{}
	err := floor.SetName(name)
	if err != nil {
		return &floor, err
	}
	if description != "" {
		err := floor.SetDescription(description)
		if err != nil {
			return &floor, err
		}
	}
	return &floor, nil
}

// NewRoom instantiates a Room.
func (s *GeolocationService) NewRoom(name, description string) (*Room, error) {
	room := Room{}
	err := room.SetName(name)
	if err != nil {
		return &room, err
	}
	if description != "" {
		err := room.SetDescription(description)
		if err != nil {
			return &room, err
		}
	}
	return &room, nil
}

// NewRow instantiates a Row.
func (s *GeolocationService) NewRow(name, description string) (*Row, error) {
	row := Row{}
	err := row.SetName(name)
	if err != nil {
		return &row, err
	}
	if description != "" {
		err := row.SetDescription(description)
		if err != nil {
			return &row, err
		}
	}
	return &row, nil
}

// NewRack instantiates a Rack.
func (s *GeolocationService) NewRack(name, description string) (*Rack, error) {
	rack := Rack{}
	err := rack.SetName(name)
	if err != nil {
		return &rack, err
	}
	if description != "" {
		err := rack.SetDescription(description)
		if err != nil {
			return &rack, err
		}
	}
	return &rack, nil
}

// UpdateSite ...
func (s *GeolocationService) UpdateSite(ctx context.Context, site *Site) (GeolocationResponse, error) {
	path := fmt.Sprintf("api/node/mo/uni/fabric/site-%s.json", site.Name())
	payload := newGeoSiteContainer(site)

	var gr GeolocationResponse

	req, err := s.client.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return gr, err
	}

	_, err = s.client.Do(ctx, req, &gr)
	return gr, err
}

// ListSites ...
func (s *GeolocationService) ListSites(ctx context.Context) ([]*Site, error) {
	path := "api/node/class/geoSite.json?rsp-subtree=full"

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list all: %v", err)
	}

	// structure of expected response
	var gs struct {
		Imdata []GeoSiteContainer `json:"imdata"`
	}

	_, err = s.client.Do(ctx, req, &gs)
	if err != nil {
		return nil, fmt.Errorf("list all: %v", err)
	}

	var sites []*Site

	for _, siteLocation := range gs.Imdata {
		site, err := s.NewSite(siteLocation.Name, siteLocation.Descr)
		if err != nil {
			return sites, err
		}
		for _, buildingLocation := range siteLocation.GeoBuildings {
			building, err := s.NewBuilding(buildingLocation.Name, buildingLocation.Descr)
			if err != nil {
				return sites, err
			}
			for _, floorLocation := range buildingLocation.GeoFloors {
				floor, err := s.NewFloor(floorLocation.Name, floorLocation.Descr)
				if err != nil {
					return sites, err
				}
				for _, roomLocation := range floorLocation.GeoRooms {
					room, err := s.NewRoom(roomLocation.Name, roomLocation.Descr)
					if err != nil {
						return sites, err
					}
					for _, rowLocation := range roomLocation.GeoRows {
						row, err := s.NewRow(rowLocation.Name, rowLocation.Descr)
						if err != nil {
							return sites, err
						}
						for _, rackLocation := range rowLocation.GeoRacks {
							rack, err := s.NewRack(rackLocation.Name, rackLocation.Descr)
							if err != nil {
								return sites, err
							}
							row.AddRack(rack)
						}
						room.AddRow(row)
					}
					floor.AddRoom(room)
				}
				building.AddFloor(floor)
			}
			site.AddBuilding(building)
		}
		sites = append(sites, site)
	}
	return sites, nil
}
