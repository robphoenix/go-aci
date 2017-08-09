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

func newFabricInstanceContainer(st Site, action string) FabricInstanceContainer {
	c := []GeoSiteContainer{newGeoSiteContainer(st, Building{}, action)}
	return FabricInstanceContainer{
		FabricInstance: FabricInstance{
			GeoAttrs: GeoAttrs{
				Dn:     "uni/fabric",
				Status: modified,
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

func newGeoSiteContainer(st Site, b Building, action string) GeoSiteContainer {
	c := []GeoBuildingContainer{}
	// The building variable will be an empty string if
	// it is the site that is being added/deleted.
	// In this case we don't need to add any
	// GeoBuildingContainer's to c.
	// If it is the building that is being added/deleted
	// then the site just needs to be modified.
	if b.Name != "" {
		c = append(c, newGeoBuildingContainer(st, b, Floor{}, action))
		action = modified
	}
	return GeoSiteContainer{
		GeoSite: GeoSite{
			GeoAttrs: GeoAttrs{
				Dn:     fmt.Sprintf("uni/fabric/site-%s", st.Name),
				Descr:  st.Description,
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

func newGeoBuildingContainer(st Site, b Building, f Floor, action string) GeoBuildingContainer {
	c := []GeoFloorContainer{}
	// The floor variable will be an empty string if
	// it is the building that is being added/deleted.
	// In this case we don't need to add any
	// GeoFloorContainer's to c.
	// If it is the floor that is being added/deleted
	// then the building just needs to be modified.
	if f.Name != "" {
		c = append(c, newGeoFloorContainer(st, b, f, Room{}, action))
		action = modified
	}
	return GeoBuildingContainer{
		GeoBuilding: GeoBuilding{
			GeoAttrs: GeoAttrs{
				Dn:     fmt.Sprintf("uni/fabric/site-%s/building-%s", st.Name, b.Name),
				Descr:  b.Description,
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

func newGeoFloorContainer(st Site, b Building, f Floor, rm Room, action string) GeoFloorContainer {
	c := []GeoRoomContainer{}
	// The room variable will be an empty string if
	// it is the floor that is being added/deleted.
	// In this case we don't need to add any
	// GeoRowContainer's to c.
	// If it is the room that is being added/deleted
	// then the floor just needs to be modified.
	if rm.Name != "" {
		c = append(c, newGeoRoomContainer(st, b, f, rm, Row{}, action))
		action = modified
	}
	return GeoFloorContainer{
		GeoFloor: GeoFloor{
			GeoAttrs: GeoAttrs{
				Dn: fmt.Sprintf(
					"uni/fabric/site-%s/building-%s/floor-%s",
					st.Name, b.Name, f.Name,
				),
				Descr:  f.Description,
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

func newGeoRoomContainer(st Site, b Building, f Floor, rm Room, rw Row, action string) GeoRoomContainer {
	c := []GeoRowContainer{}
	// The row variable will be an empty string if
	// it is the room that is being added/deleted.
	// In this case we don't need to add any
	// GeoRowContainer's to c.
	// If it is the row that is being added/deleted
	// then the room just needs to be modified.
	if rw.Name != "" {
		c = append(c, newGeoRowContainer(st, b, f, rm, rw, Rack{}, action))
		action = modified
	}
	return GeoRoomContainer{
		GeoRoom: GeoRoom{
			GeoAttrs: GeoAttrs{
				Dn: fmt.Sprintf(
					"uni/fabric/site-%s/building-%s/floor-%s/room-%s",
					st.Name, b.Name, f.Name, rm.Name,
				),
				Descr:  rm.Description,
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

func newGeoRowContainer(st Site, b Building, f Floor, rm Room, rw Row, rk Rack, action string) GeoRowContainer {
	c := []GeoRackContainer{}
	// The rack variable will be an empty string if
	// it is the row that is being added/deleted.
	// In this case we don't need to add any
	// GeoRackContainer's to c.
	// If it is the Rack that is being added/deleted
	// then the row just needs to be modified.
	if rk.Name != "" {
		c = append(c, newGeoRackContainer(st, b, f, rm, rw, rk, action))
		action = modified
	}
	return GeoRowContainer{
		GeoRow: GeoRow{
			GeoAttrs: GeoAttrs{
				Dn: fmt.Sprintf(
					"uni/fabric/site-%s/building-%s/floor-%s/room-%s/row-%s",
					st.Name, b.Name, f.Name, rm.Name, rw.Name,
				),
				Name:   rw.Name,
				Descr:  rw.Description,
				Rn:     fmt.Sprintf("row-%s", rw.Name),
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

func newGeoRackContainer(st Site, b Building, f Floor, rm Room, rw Row, rk Rack, action string) GeoRackContainer {
	return GeoRackContainer{
		GeoRack: GeoRack{
			GeoAttrs: GeoAttrs{
				Dn: fmt.Sprintf(
					"uni/fabric/site-%s/building-%s/floor-%s/room-%s/row-%s/rack-%s",
					st.Name, b.Name, f.Name, rm.Name, rw.Name, rk.Name,
				),
				Name:   rk.Name,
				Descr:  rk.Description,
				Rn:     fmt.Sprintf("rack-%s", rk.Name),
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

// // UpdateSite ...
// func (s *GeolocationService) UpdateSite(ctx context.Context, site *Site) (GeolocationResponse, error) {
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

	var sites []*Site

	for _, site := range gs.Imdata {
		st := &Site{Name: site.Name}
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
			st.Buildings = append(st.Buildings, b)
		}
		sites = append(sites, st)
	}

	return sites, nil
}

//
// // AddSite ...
// func (s *GeolocationService) AddSite(ctx context.Context, st Site) (GeolocationResponse, error) {
//         path := fmt.Sprintf("api/node/mo/uni/fabric/site-%s.json", st.Name)
//
//         var gr GeolocationResponse
//
//         payload := newGeoSiteContainer(st, Building{}, createModify)
//
//         req, err := s.client.NewRequest(http.MethodPost, path, payload)
//         if err != nil {
//                 return gr, err
//         }
//
//         _, err = s.client.Do(ctx, req, &gr)
//         if err != nil {
//                 return gr, err
//         }
//
//         return gr, nil
// }
