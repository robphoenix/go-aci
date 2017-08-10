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

// GeoSite ...
type GeoSite struct {
	GeoAttrs     `json:"attributes,omitempty"`
	GeoBuildings []GeoBuildingContainer `json:"children,omitempty"`
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

// GeoFloorContainer ...
type GeoFloorContainer struct {
	GeoFloor `json:"geoFloor,omitempty"`
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

// GeoRoom ...
type GeoRoom struct {
	GeoAttrs `json:"attributes,omitempty"`
	GeoRows  []GeoRowContainer `json:"children,omitempty"`
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

// GeoRackContainer ...
type GeoRackContainer struct {
	GeoRack `json:"geoRack,omitempty"`
}

// GeoRack ...
type GeoRack struct {
	GeoAttrs `json:"attributes,omitempty"`
	GeoNodes []interface{} `json:"children,omitempty"`
}

// GeoAttrs ...
type GeoAttrs struct {
	Descr  string `json:"descr,omitempty"`
	Dn     string `json:"dn,omitempty"`
	Name   string `json:"name,omitempty"`
	Rn     string `json:"rn,omitempty"`
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

// // UpdateSite ...
// func (s *GeolocationService) UpdateSite(ctx context.Context, location *Site) (GeolocationResponse, error) {
//
// }

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
		site, err := NewSite(siteLocation.Name, siteLocation.Descr)
		if err != nil {
			return sites, err
		}
		for _, buildingLocation := range siteLocation.GeoBuildings {
			building, err := NewBuilding(buildingLocation.Name, buildingLocation.Descr)
			if err != nil {
				return sites, err
			}
			for _, floorLocation := range buildingLocation.GeoFloors {
				floor, err := NewFloor(floorLocation.Name, floorLocation.Descr)
				if err != nil {
					return sites, err
				}
				for _, roomLocation := range floorLocation.GeoRooms {
					room, err := NewRoom(roomLocation.Name, roomLocation.Descr)
					if err != nil {
						return sites, err
					}
					for _, rowLocation := range roomLocation.GeoRows {
						row, err := NewRow(rowLocation.Name, rowLocation.Descr)
						if err != nil {
							return sites, err
						}
						for _, rackLocation := range rowLocation.GeoRacks {
							rack, err := NewRack(rackLocation.Name, rackLocation.Descr)
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
