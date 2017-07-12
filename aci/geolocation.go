package aci

var (
	addGeoSitePath = "api/node/mo/uni/fabric/site-%s.json" // geo-site name
)

// Site ...
type Site struct {
	Name      string
	Buildings []Building
}

// Building ...
type Building struct {
	Name   string
	Floors []Floor
}

// Floor ...
type Floor struct {
	Name  string
	Rooms []Room
}

// Room ...
type Room struct {
	Name string
	Rows []Row
}

// Row ...
type Row struct {
	Name  string
	Racks []Rack
}

// Rack ...
type Rack struct {
	Name string
	// Nodes
}

// GeoSiteContainer ...
type GeoSiteContainer struct {
	GeoSite `json:"geoSite"`
}

// GeoSite ...
type GeoSite struct {
	GeoAttrs `json:"attributes"`
	Children []GeoBuildingContainer `json:"children"`
}

// GeoBuildingContainer ...
type GeoBuildingContainer struct {
	GeoBuilding `json:"geoBuilding"`
}

// GeoBuilding ...
type GeoBuilding struct {
	GeoAttrs `json:"attributes"`
	Children []GeoFloorContainer `json:"children"`
}

// GeoFloorContainer ...
type GeoFloorContainer struct {
	GeoFloor `json:"geoFloor"`
}

// GeoFloor ...
type GeoFloor struct {
	GeoAttrs `json:"attributes"`
	Children []GeoRoomContainer `json:"children"`
}

// GeoRoomContainer ...
type GeoRoomContainer struct {
	GeoRoom `json:"geoRoom"`
}

// GeoRoom ...
type GeoRoom struct {
	GeoAttrs `json:"attributes"`
	Children []GeoRowContainer `json:"children"`
}

// GeoRowContainer ...
type GeoRowContainer struct {
	GeoRow `json:"geoRow"`
}

// GeoRow ...
type GeoRow struct {
	GeoAttrs `json:"attributes"`
	Children []GeoRackContainer `json:"children"`
}

// GeoRackContainer ...
type GeoRackContainer struct {
	GeoRack `json:"geoRack"`
}

// GeoRack ...
type GeoRack struct {
	GeoAttrs `json:"attributes"`
	// TODO: specify node struct ??
	Children []interface{} `json:"children"`
}

// GeoAttrs ...
type GeoAttrs struct {
	Dn     string `json:"dn"`
	Name   string `json:"name"`
	Rn     string `json:"rn"`
	Status string `json:"status"`
}

// api/node/mo/uni/fabric/site-%s.json // geo-site name
// {
//   "geoSite": {
//     "attributes": {
//       "dn": "uni/fabric/site-ABC-Liverpool",
//       "name": "ABC-Liverpool",
//       "rn": "site-ABC-Liverpool",
//       "status": "created"
//     },
//     "children": [
//       {
//         "geoBuilding": {
//           "attributes": {
//             "dn": "uni/fabric/site-ABC-Liverpool/building-Node4-Liverpool",
//             "name": "Node4-Liverpool",
//             "rn": "building-Node4-Liverpool",
//             "status": "created"
//           },
//           "children": [
//             {
//               "geoFloor": {
//                 "attributes": {
//                   "dn": "uni/fabric/site-ABC-Liverpool/building-Node4-Liverpool/floor-Floor-000",
//                   "name": "Floor-000",
//                   "rn": "floor-Floor-000",
//                   "status": "created"
//                 },
//                 "children": [
//                   {
//                     "geoRoom": {
//                       "attributes": {
//                         "dn": "uni/fabric/site-ABC-Liverpool/building-Node4-Liverpool/floor-Floor-000/room-Hall-001",
//                         "name": "Hall-001",
//                         "rn": "room-Hall-001",
//                         "status": "created"
//                       },
//                       "children": [
//                         {
//                           "geoRow": {
//                             "attributes": {
//                               "dn": "uni/fabric/site-ABC-Liverpool/building-Node4-Liverpool/floor-Floor-000/room-Hall-001/row-Row-000",
//                               "name": "Row-000",
//                               "rn": "row-Row-000",
//                               "status": "created"
//                             },
//                             "children": [
//                               {
//                                 "geoRack": {
//                                   "attributes": {
//                                     "dn": "uni/fabric/site-ABC-Liverpool/building-Node4-Liverpool/floor-Floor-000/room-Hall-001/row-Row-000/rack-Rack-1234",
//                                     "name": "Rack-1234",
//                                     "rn": "rack-Rack-1234",
//                                     "status": "created"
//                                   },
//                                   "children": []
//                                 }
//                               }
//                             ]
//                           }
//                         }
//                       ]
//                     }
//                   }
//                 ]
//               }
//             }
//           ]
//         }
//       }
//     ]
//   }
// }
// response: {"totalCount":"0","imdata":[]}
