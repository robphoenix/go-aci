package aci

import (
	"context"
	"fmt"
	"net/http"
)

var (
	listGeolocationPath = "api/node/class/geoSite.json?rsp-subtree=full"
)

// Site ...
type Site struct {
	Name        string
	Description string
	Buildings   []*Building
}

// Building ...
type Building struct {
	Name        string
	Description string
	Floors      []*Floor
}

// Floor ...
type Floor struct {
	Name        string
	Description string
	Rooms       []*Room
}

// Room ...
type Room struct {
	Name        string
	Description string
	Rows        []*Row
}

// Row ...
type Row struct {
	Name        string
	Description string
	Racks       []*Rack
}

// Rack ...
type Rack struct {
	Name        string
	Description string
}

// FabricInstanceContainer ...
type FabricInstanceContainer struct {
	FabricInstance `json:"fabricInst"`
}

// FabricInstance ...
type FabricInstance struct {
	GeoAttrs `json:"attributes"`
	GeoSites []GeoSiteContainer `json:"children"`
}

func newFabricInstanceContainer(site, action string) FabricInstanceContainer {
	c := []GeoSiteContainer{newGeoSiteContainer(site, "", action)}
	return FabricInstanceContainer{
		FabricInstance: FabricInstance{
			GeoAttrs: GeoAttrs{
				Dn:     "uni/fabric",
				Status: modify,
			},
			GeoSites: c,
		},
	}
}

// GeoSiteContainer ...
type GeoSiteContainer struct {
	GeoSite `json:"geoSite,omitempty"`
}

// GeoSite ...
type GeoSite struct {
	GeoAttrs     `json:"attributes,omitempty"`
	GeoBuildings []GeoBuildingContainer `json:"children,omitempty"`
}

func newGeoSiteContainer(site, building, action string) GeoSiteContainer {
	c := []GeoBuildingContainer{}
	// The building variable will be an empty string if
	// it is the site that is being added/deleted.
	// In this case we don't need to add any
	// GeoBuildingContainer's to c.
	// If it is the building that is being added/deleted
	// then the site just needs to be modified.
	if building != "" {
		c = append(c, newGeoBuildingContainer(site, building, "", action))
		action = modify
	}
	return GeoSiteContainer{
		GeoSite: GeoSite{
			GeoAttrs: GeoAttrs{
				Dn:     fmt.Sprintf("uni/fabric/site-%s", site),
				Status: action,
			},
			GeoBuildings: c,
		},
	}
}

// GeoBuildingContainer ...
type GeoBuildingContainer struct {
	GeoBuilding `json:"geoBuilding,omitempty"`
}

// GeoBuilding ...
type GeoBuilding struct {
	GeoAttrs  `json:"attributes,omitempty"`
	GeoFloors []GeoFloorContainer `json:"children,omitempty"`
}

func newGeoBuildingContainer(site, building, floor, action string) GeoBuildingContainer {
	c := []GeoFloorContainer{}
	// The floor variable will be an empty string if
	// it is the building that is being added/deleted.
	// In this case we don't need to add any
	// GeoFloorContainer's to c.
	// If it is the floor that is being added/deleted
	// then the building just needs to be modified.
	if floor != "" {
		c = append(c, newGeoFloorContainer(site, building, floor, "", action))
		action = modify
	}
	return GeoBuildingContainer{
		GeoBuilding: GeoBuilding{
			GeoAttrs: GeoAttrs{
				Dn:     fmt.Sprintf("uni/fabric/site-%s/building-%s", site, building),
				Status: action,
			},
			GeoFloors: c,
		},
	}
}

// GeoFloorContainer ...
type GeoFloorContainer struct {
	GeoFloor `json:"geoFloor,omitempty"`
}

// GeoFloor ...
type GeoFloor struct {
	GeoAttrs `json:"attributes,omitempty"`
	GeoRooms []GeoRoomContainer `json:"children,omitempty"`
}

func newGeoFloorContainer(site, building, floor, room, action string) GeoFloorContainer {
	c := []GeoRoomContainer{}
	// The room variable will be an empty string if
	// it is the floor that is being added/deleted.
	// In this case we don't need to add any
	// GeoRowContainer's to c.
	// If it is the room that is being added/deleted
	// then the floor just needs to be modified.
	if room != "" {
		c = append(c, newGeoRoomContainer(site, building, floor, room, "", action))
		action = modify
	}
	return GeoFloorContainer{
		GeoFloor: GeoFloor{
			GeoAttrs: GeoAttrs{
				Dn: fmt.Sprintf(
					"uni/fabric/site-%s/building-%s/floor-%s",
					site,
					building,
					floor,
				),
				Status: action,
			},
			GeoRooms: c,
		},
	}
}

// GeoRoomContainer ...
type GeoRoomContainer struct {
	GeoRoom `json:"geoRoom,omitempty"`
}

// GeoRoom ...
type GeoRoom struct {
	GeoAttrs `json:"attributes,omitempty"`
	GeoRows  []GeoRowContainer `json:"children,omitempty"`
}

func newGeoRoomContainer(site, building, floor, room, row, action string) GeoRoomContainer {
	c := []GeoRowContainer{}
	// The row variable will be an empty string if
	// it is the room that is being added/deleted.
	// In this case we don't need to add any
	// GeoRowContainer's to c.
	// If it is the row that is being added/deleted
	// then the room just needs to be modified.
	if row != "" {
		c = append(c, newGeoRowContainer(site, building, floor, room, row, "", action))
		action = modify
	}
	return GeoRoomContainer{
		GeoRoom: GeoRoom{
			GeoAttrs: GeoAttrs{
				Dn: fmt.Sprintf(
					"uni/fabric/site-%s/building-%s/floor-%s/room-%s",
					site,
					building,
					floor,
					room,
				),
				Status: action,
			},
			GeoRows: c,
		},
	}
}

// GeoRowContainer ...
type GeoRowContainer struct {
	GeoRow `json:"geoRow,omitempty"`
}

// GeoRow ...
type GeoRow struct {
	GeoAttrs `json:"attributes,omitempty"`
	GeoRacks []GeoRackContainer `json:"children,omitempty"`
}

func newGeoRowContainer(site, building, floor, room, row, rack, action string) GeoRowContainer {
	c := []GeoRackContainer{}
	// The rack variable will be an empty string if
	// it is the row that is being added/deleted.
	// In this case we don't need to add any
	// GeoRackContainer's to c.
	// If it is the Rack that is being added/deleted
	// then the row just needs to be modified.
	if rack != "" {
		c = append(c, newGeoRackContainer(site, building, floor, room, row, rack, action))
		action = modify
	}
	return GeoRowContainer{
		GeoRow: GeoRow{
			GeoAttrs: GeoAttrs{
				Dn: fmt.Sprintf(
					"uni/fabric/site-%s/building-%s/floor-%s/room-%s/row-%s",
					site,
					building,
					floor,
					room,
					row,
				),
				Name:   row,
				Rn:     fmt.Sprintf("row-%s", row),
				Status: action,
			},
			GeoRacks: c,
		},
	}
}

// GeoRackContainer ...
type GeoRackContainer struct {
	GeoRack `json:"geoRack,omitempty"`
}

// GeoRack ...
type GeoRack struct {
	GeoAttrs `json:"attributes,omitempty"`
	GeoNodes []interface{} `json:"children,omitempty"`
}

func newGeoRackContainer(site, building, floor, room, row, rack, action string) GeoRackContainer {
	return GeoRackContainer{
		GeoRack: GeoRack{
			GeoAttrs: GeoAttrs{
				Dn: fmt.Sprintf(
					"uni/fabric/site-%s/building-%s/floor-%s/room-%s/row-%s/rack-%s",
					site,
					building,
					floor,
					room,
					row,
					rack,
				),
				Name:   rack,
				Rn:     fmt.Sprintf("rack-%s", rack),
				Status: action,
			},
			GeoNodes: nil,
		},
	}
}

// GeoAttrs ...
type GeoAttrs struct {
	Descr  string `json:"descr,omitempty"`
	Dn     string `json:"dn,omitempty"`
	Name   string `json:"name,omitempty"`
	Rn     string `json:"rn,omitempty"`
	Status string `json:"status,omitempty"`
}

// GeolocationService handles communication with the geolocation related
// methods of the APIC API.
type GeolocationService service

// GeolocationResponse ...
type GeolocationResponse struct {
	TotalCount string        `json:"totalCount"`
	Imdata     []interface{} `json:"imdata"`
}

// // AddFullSite ...
// func (s *GeolocationService) AddFullSite(ctx context.Context, site *Site) (GeolocationResponse, error) {
//
// }

// ListFullSites ...
func (s *GeolocationService) ListFullSites(ctx context.Context) ([]*Site, error) {
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

	var ss []*Site

	for _, site := range gs.Imdata {
		s := &Site{Name: site.Name}
		for _, building := range site.GeoBuildings {
			b := &Building{Name: building.Name}
			for _, floor := range building.GeoFloors {
				f := &Floor{Name: floor.Name}
				for _, room := range floor.GeoRooms {
					r1 := &Room{Name: room.Name}
					for _, row := range room.GeoRows {
						r2 := &Row{Name: row.Name}
						for _, rack := range row.GeoRacks {
							r3 := &Rack{Name: rack.Name}
							r2.Racks = append(r2.Racks, r3)
						}
						r1.Rows = append(r1.Rows, r2)
					}
					f.Rooms = append(f.Rooms, r1)
				}
				b.Floors = append(b.Floors, f)
			}
			s.Buildings = append(s.Buildings, b)
		}
		ss = append(ss, s)
	}

	return ss, nil
}

// AddSite ...
func (s *GeolocationService) AddSite(ctx context.Context, site string) (GeolocationResponse, error) {
	path := fmt.Sprintf("api/node/mo/uni/fabric/site-%s.json", site)

	var gr GeolocationResponse

	payload := newGeoSiteContainer(site, "", createModify)

	req, err := s.client.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return gr, err
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return gr, err
	}

	return gr, nil
}

// DeleteSite ...
func (s *GeolocationService) DeleteSite(ctx context.Context, site string) (GeolocationResponse, error) {
	path := "api/node/mo/uni/fabric.json"

	var gr GeolocationResponse

	payload := newFabricInstanceContainer(site, delete)

	req, err := s.client.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return gr, err
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return gr, err
	}

	return gr, nil
}

// ListSites ...
func (s *GeolocationService) ListSites(ctx context.Context) ([]*Site, error) {
	path := "api/node/class/geoSite.json"

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list sites: %v", err)
	}

	// structure of expected response
	var gr struct {
		Imdata []GeoSiteContainer `json:"imdata"`
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return nil, fmt.Errorf("list sites: %v", err)
	}

	var ss []*Site

	for _, s := range gr.Imdata {
		ss = append(ss, &Site{Name: s.Name})
	}

	return ss, nil
}

// AddBuilding ...
func (s *GeolocationService) AddBuilding(ctx context.Context, site, building string) (GeolocationResponse, error) {
	path := fmt.Sprintf("api/node/mo/uni/fabric/site-%s/building-%s.json", site, building)

	var gr GeolocationResponse

	payload := newGeoBuildingContainer(site, building, "", createModify)

	req, err := s.client.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return gr, err
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return gr, err
	}

	return gr, nil
}

// DeleteBuilding ...
func (s *GeolocationService) DeleteBuilding(ctx context.Context, site, building string) (GeolocationResponse, error) {
	path := fmt.Sprintf("api/node/mo/uni/fabric/site-%s.json", site)

	var gr GeolocationResponse

	payload := newGeoSiteContainer(site, building, delete)

	req, err := s.client.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return gr, err
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return gr, err
	}

	return gr, nil
}

// ListBuildings ...
func (s *GeolocationService) ListBuildings(ctx context.Context, site string) ([]*Building, error) {
	path := fmt.Sprintf("api/node/mo/uni/fabric/site-%s.json?query-target=children&target-subtree-class=geoRow", site)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list buildings: %v", err)
	}

	// structure of expected response
	var gr struct {
		Imdata []GeoBuildingContainer `json:"imdata"`
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return nil, fmt.Errorf("list buildings: %v", err)
	}

	var bs []*Building

	for _, b := range gr.Imdata {
		bs = append(bs, &Building{Name: b.Name})
	}

	return bs, nil
}

// AddFloor ...
func (s *GeolocationService) AddFloor(ctx context.Context, site, building, floor string) (GeolocationResponse, error) {
	path := fmt.Sprintf(
		"api/node/mo/uni/fabric/site-%s/building-%s/floor-%s.json",
		site,
		building,
		floor,
	)

	var gr GeolocationResponse

	payload := newGeoFloorContainer(site, building, floor, "", createModify)

	req, err := s.client.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return gr, err
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return gr, err
	}

	return gr, nil
}

// DeleteFloor ...
func (s *GeolocationService) DeleteFloor(ctx context.Context, site, building, floor string) (GeolocationResponse, error) {
	path := fmt.Sprintf(
		"api/node/mo/uni/fabric/site-%s/building-%s.json",
		site,
		building,
	)

	var gr GeolocationResponse

	payload := newGeoBuildingContainer(site, building, floor, delete)

	req, err := s.client.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return gr, err
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return gr, err
	}

	return gr, nil
}

// ListFloors ...
func (s *GeolocationService) ListFloors(ctx context.Context, site, building string) ([]*Floor, error) {
	path := fmt.Sprintf(
		"api/node/mo/uni/fabric/site-%s/building-%s.json?query-target=children&target-subtree-class=geoRow",
		site,
		building,
	)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list floors: %v", err)
	}

	// structure of expected response
	var gr struct {
		Imdata []GeoFloorContainer `json:"imdata"`
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return nil, fmt.Errorf("list floors: %v", err)
	}

	var fs []*Floor

	for _, f := range gr.Imdata {
		fs = append(fs, &Floor{Name: f.Name})
	}

	return fs, nil
}

// AddRoom ...
func (s *GeolocationService) AddRoom(ctx context.Context, site, building, floor, room string) (GeolocationResponse, error) {
	path := fmt.Sprintf(
		"api/node/mo/uni/fabric/site-%s/building-%s/floor-%s/room-%s.json",
		site,
		building,
		floor,
		room,
	)

	var gr GeolocationResponse

	payload := newGeoRoomContainer(site, building, floor, room, "", createModify)

	req, err := s.client.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return gr, err
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return gr, err
	}

	return gr, nil
}

// DeleteRoom ...
func (s *GeolocationService) DeleteRoom(ctx context.Context, site, building, floor, room string) (GeolocationResponse, error) {
	path := fmt.Sprintf(
		"api/node/mo/uni/fabric/site-%s/building-%s/floor-%s.json",
		site,
		building,
		floor,
	)

	var gr GeolocationResponse

	payload := newGeoFloorContainer(site, building, floor, room, delete)

	req, err := s.client.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return gr, err
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return gr, err
	}

	return gr, nil
}

// ListRooms ...
func (s *GeolocationService) ListRooms(ctx context.Context, site, building, floor string) ([]*Room, error) {
	path := fmt.Sprintf(
		"api/node/mo/uni/fabric/site-%s/building-%s/floor-%s.json?query-target=children&target-subtree-class=geoRow",
		site,
		building,
		floor,
	)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list rooms: %v", err)
	}

	// structure of expected response
	var gr struct {
		Imdata []GeoRoomContainer `json:"imdata"`
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return nil, fmt.Errorf("list rooms: %v", err)
	}

	var rs []*Room

	for _, r := range gr.Imdata {
		rs = append(rs, &Room{Name: r.Name})
	}

	return rs, nil
}

// AddRow ...
func (s *GeolocationService) AddRow(ctx context.Context, site, building, floor, room, row string) (GeolocationResponse, error) {
	path := fmt.Sprintf(
		"api/node/mo/uni/fabric/site-%s/building-%s/floor-%s/room-%s/row-%s.json",
		site,
		building,
		floor,
		room,
		row,
	)

	var gr GeolocationResponse

	payload := newGeoRowContainer(site, building, floor, room, row, "", createModify)

	req, err := s.client.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return gr, err
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return gr, err
	}

	return gr, nil
}

// DeleteRow ...
func (s *GeolocationService) DeleteRow(ctx context.Context, site, building, floor, room, row string) (GeolocationResponse, error) {
	path := fmt.Sprintf(
		"api/node/mo/uni/fabric/site-%s/building-%s/floor-%s/room-%s.json",
		site,
		building,
		floor,
		room,
	)

	var gr GeolocationResponse

	payload := newGeoRoomContainer(site, building, floor, room, row, delete)

	req, err := s.client.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return gr, err
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return gr, err
	}

	return gr, nil
}

// ListRows ...
func (s *GeolocationService) ListRows(ctx context.Context, site, building, floor, room string) ([]*Row, error) {
	path := fmt.Sprintf(
		"api/node/mo/uni/fabric/site-%s/building-%s/floor-%s/room-%s.json?query-target=children&target-subtree-class=geoRow",
		site,
		building,
		floor,
		room,
	)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list rows: %v", err)
	}

	// structure of expected response
	var gr struct {
		Imdata []GeoRowContainer `json:"imdata"`
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return nil, fmt.Errorf("list rows: %v", err)
	}

	var rs []*Row

	for _, r := range gr.Imdata {
		rs = append(rs, &Row{Name: r.Name})
	}

	return rs, nil
}

// AddRack ...
func (s *GeolocationService) AddRack(ctx context.Context, site, building, floor, room, row, rack string) (GeolocationResponse, error) {
	path := fmt.Sprintf(
		"api/node/mo/uni/fabric/site-%s/building-%s/floor-%s/room-%s/row-%s/rack-%s.json",
		site,
		building,
		floor,
		room,
		row,
		rack,
	)

	var gr GeolocationResponse

	payload := newGeoRackContainer(site, building, floor, room, row, rack, createModify)

	req, err := s.client.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return gr, err
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return gr, err
	}

	return gr, nil
}

// DeleteRack ...
func (s *GeolocationService) DeleteRack(ctx context.Context, site, building, floor, room, row, rack string) (GeolocationResponse, error) {
	path := fmt.Sprintf(
		"api/node/mo/uni/fabric/site-%s/building-%s/floor-%s/room-%s/row-%s.json",
		site,
		building,
		floor,
		room,
		row,
	)

	var gr GeolocationResponse

	payload := newGeoRowContainer(site, building, floor, room, row, rack, delete)

	req, err := s.client.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return gr, err
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return gr, err
	}

	return gr, nil
}

// ListRacks ...
func (s *GeolocationService) ListRacks(ctx context.Context, site, building, floor, room, row string) ([]*Rack, error) {
	path := fmt.Sprintf(
		"api/node/mo/uni/fabric/site-%s/building-%s/floor-%s/room-%s/row-%s.json?query-target=children&target-subtree-class=geoRack",
		site,
		building,
		floor,
		room,
		row,
	)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list racks: %v", err)
	}

	// structure of expected response
	var gr struct {
		Imdata []GeoRackContainer `json:"imdata"`
	}

	_, err = s.client.Do(ctx, req, &gr)
	if err != nil {
		return nil, fmt.Errorf("list racks: %v", err)
	}

	var rs []*Rack

	for _, r := range gr.Imdata {
		rs = append(rs, &Rack{Name: r.Name})
	}

	return rs, nil
}
