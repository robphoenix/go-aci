package aci

import "time"

var (
	addGeoSitePath      = "api/node/mo/uni/fabric/site-%s.json" // geo-site name
	listGeolocationPath = "api/node/class/geoSite.json?rsp-subtree=full"
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

// GeolocationService handles communication with the geolocation related
// methods of the APIC API.
type GeolocationService service

type ListResponse struct {
	Imdata []struct {
		GeoSite struct {
			Attributes struct {
				ChildAction string    `json:"childAction"`
				Descr       string    `json:"descr"`
				Dn          string    `json:"dn"`
				LcOwn       string    `json:"lcOwn"`
				ModTs       time.Time `json:"modTs"`
				Name        string    `json:"name"`
				NameAlias   string    `json:"nameAlias"`
				OwnerKey    string    `json:"ownerKey"`
				OwnerTag    string    `json:"ownerTag"`
				Status      string    `json:"status"`
				UID         string    `json:"uid"`
			} `json:"attributes"`
			Children []struct {
				GeoBuilding struct {
					Attributes struct {
						ChildAction string    `json:"childAction"`
						Descr       string    `json:"descr"`
						LcOwn       string    `json:"lcOwn"`
						ModTs       time.Time `json:"modTs"`
						Name        string    `json:"name"`
						NameAlias   string    `json:"nameAlias"`
						Rn          string    `json:"rn"`
						Status      string    `json:"status"`
						UID         string    `json:"uid"`
					} `json:"attributes"`
					Children []struct {
						GeoFloor struct {
							Attributes struct {
								ChildAction string    `json:"childAction"`
								Descr       string    `json:"descr"`
								LcOwn       string    `json:"lcOwn"`
								ModTs       time.Time `json:"modTs"`
								Name        string    `json:"name"`
								NameAlias   string    `json:"nameAlias"`
								Rn          string    `json:"rn"`
								Status      string    `json:"status"`
								UID         string    `json:"uid"`
							} `json:"attributes"`
							Children []struct {
								GeoRoom struct {
									Attributes struct {
										ChildAction string    `json:"childAction"`
										Descr       string    `json:"descr"`
										LcOwn       string    `json:"lcOwn"`
										ModTs       time.Time `json:"modTs"`
										Name        string    `json:"name"`
										NameAlias   string    `json:"nameAlias"`
										Rn          string    `json:"rn"`
										Status      string    `json:"status"`
										UID         string    `json:"uid"`
									} `json:"attributes"`
									Children []struct {
										GeoRack struct {
											Attributes struct {
												ChildAction string    `json:"childAction"`
												Descr       string    `json:"descr"`
												LcOwn       string    `json:"lcOwn"`
												ModTs       time.Time `json:"modTs"`
												Name        string    `json:"name"`
												NameAlias   string    `json:"nameAlias"`
												Rn          string    `json:"rn"`
												Status      string    `json:"status"`
												UID         string    `json:"uid"`
											} `json:"attributes"`
										} `json:"geoRack"`
									} `json:"children"`
								} `json:"geoRoom"`
							} `json:"children"`
						} `json:"geoFloor"`
					} `json:"children"`
				} `json:"geoBuilding"`
			} `json:"children"`
		} `json:"geoSite"`
	} `json:"imdata"`
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

// // ADD SITE
// "api/node/mo/uni/fabric/site-%s.json" site name
// {
//    "geoSite":{
//       "attributes":{
//          "dn":"uni/fabric/site-Site01",
//          "name":"Site01",
//          "rn":"site-Site01",
//          "status":"created"
//       },
//       "children":[]
//    }
// }
//
// // DELETE SITE
// api/node/mo/uni/fabric.json
// {
//   "fabricInst": {
//     "attributes": {
//       "dn": "uni/fabric",
//       "status": "modified"
//     },
//     "children": [
//       {
//         "geoSite": {
//           "attributes": {
//             "dn": "uni/fabric/site-Site01",
//             "status": "deleted"
//           },
//           "children": []
//         }
//       }
//     ]
//   }
// }
//
// // LIST SITES
// api/node/class/geoSite.json
// response:
// {
//   "totalCount": "2",
//   "imdata": [
//     {
//       "geoSite": {
//         "attributes": {
//           "childAction": "",
//           "descr": "",
//           "dn": "uni/fabric/site-default",
//           "lcOwn": "local",
//           "modTs": "2017-07-17T08:22:15.217+00:00",
//           "monPolDn": "uni/fabric/monfab-default",
//           "name": "default",
//           "ownerKey": "",
//           "ownerTag": "",
//           "status": "",
//           "uid": "0"
//         }
//       }
//     },
//     {
//       "geoSite": {
//         "attributes": {
//           "childAction": "",
//           "descr": "",
//           "dn": "uni/fabric/site-Site01",
//           "lcOwn": "local",
//           "modTs": "2017-07-17T13:54:14.796+00:00",
//           "monPolDn": "uni/fabric/monfab-default",
//           "name": "Site01",
//           "ownerKey": "",
//           "ownerTag": "",
//           "status": "",
//           "uid": "15374"
//         }
//       }
//     }
//   ]
// }
//
// // ADD BUILDING
// "api/node/mo/uni/fabric/site-Site01/building-Building02.json"
// {
//   "geoBuilding": {
//     "attributes": {
//       "dn": "uni/fabric/site-Site01/building-Building02",
//       "name": "Building02",
//       "rn": "building-Building02",
//       "status": "created"
//     },
//     "children": []
//   }
// }
//
// // DELETE BUILDING
// "api/node/mo/uni/fabric/site-Site01.json"
// {
//   "geoSite": {
//     "attributes": {
//       "dn": "uni/fabric/site-Site01",
//       "status": "modified"
//     },
//     "children": [
//       {
//         "geoBuilding": {
//           "attributes": {
//             "dn": "uni/fabric/site-Site01/building-Building01",
//             "status": "deleted"
//           },
//           "children": []
//         }
//       }
//     ]
//   }
// }
//
// // LIST BUILDINGS
// "api/node/mo/uni/fabric/site-Site01.json?query-target=children&target-subtree-class=geoBuilding"
// response:
// {
//   "totalCount": "2",
//   "imdata": [
//     {
//       "geoBuilding": {
//         "attributes": {
//           "childAction": "",
//           "descr": "",
//           "dn": "uni/fabric/site-Site01/building-default",
//           "lcOwn": "local",
//           "modTs": "2017-07-17T13:54:14.796+00:00",
//           "monPolDn": "uni/fabric/monfab-default",
//           "name": "default",
//           "status": "",
//           "uid": "0"
//         }
//       }
//     },
//     {
//       "geoBuilding": {
//         "attributes": {
//           "childAction": "",
//           "descr": "",
//           "dn": "uni/fabric/site-Site01/building-Building01",
//           "lcOwn": "local",
//           "modTs": "2017-07-17T13:54:14.796+00:00",
//           "monPolDn": "uni/fabric/monfab-default",
//           "name": "Building01",
//           "status": "",
//           "uid": "15374"
//         }
//       }
//     }
//   ]
// }
//
// // ADD FLOOR
// "api/node/mo/uni/fabric/site-Site01/building-Building01/floor-Floor01.json"
// {
//   "geoFloor": {
//     "attributes": {
//       "dn": "uni/fabric/site-Site01/building-Building01/floor-Floor01",
//       "name": "Floor01",
//       "rn": "floor-Floor01",
//       "status": "created"
//     },
//     "children": []
//   }
// }
//
// // DELETE FLOOR
// "api/node/mo/uni/fabric/site-Site01/building-Building01.json"
// {
//   "geoBuilding": {
//     "attributes": {
//       "dn": "uni/fabric/site-Site01/building-Building01",
//       "status": "modified"
//     },
//     "children": [
//       {
//         "geoFloor": {
//           "attributes": {
//             "dn": "uni/fabric/site-Site01/building-Building01/floor-Floor01",
//             "status": "deleted"
//           },
//           "children": []
//         }
//       }
//     ]
//   }
// }
//
// // LIST FLOORS
// "api/node/mo/uni/fabric/site-Site01/building-Building01.json?query-target=children&target-subtree-class=geoFloor"
// response:
// {
//   "totalCount": "1",
//   "imdata": [
//     {
//       "geoFloor": {
//         "attributes": {
//           "childAction": "",
//           "descr": "",
//           "dn": "uni/fabric/site-Site01/building-Building01/floor-Floor01",
//           "lcOwn": "local",
//           "modTs": "2017-07-17T13:54:14.796+00:00",
//           "monPolDn": "uni/fabric/monfab-default",
//           "name": "Floor01",
//           "status": "",
//           "uid": "15374"
//         }
//       }
//     }
//   ]
// }
//
// // ADD ROOM
// "api/node/mo/uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01.json"
// {
//   "geoRoom": {
//     "attributes": {
//       "dn": "uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01",
//       "name": "Room01",
//       "rn": "room-Room01",
//       "status": "created"
//     },
//     "children": []
//   }
// }
//
// // DELETE ROOM
// "api/node/mo/uni/fabric/site-Site01/building-Building01/floor-Floor01.json"
// {
//   "geoFloor": {
//     "attributes": {
//       "dn": "uni/fabric/site-Site01/building-Building01/floor-Floor01",
//       "status": "modified"
//     },
//     "children": [
//       {
//         "geoRoom": {
//           "attributes": {
//             "dn": "uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01",
//             "status": "deleted"
//           },
//           "children": []
//         }
//       }
//     ]
//   }
// }
//
// LIST ROOMS
// "api/node/mo/uni/fabric/site-Site01/building-Building01/floor-Floor01.json?query-target=children&target-subtree-class=geoRoom"
// response:
// {
//   "totalCount": "1",
//   "imdata": [
//     {
//       "geoRoom": {
//         "attributes": {
//           "childAction": "",
//           "descr": "",
//           "dn": "uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01",
//           "lcOwn": "local",
//           "modTs": "2017-07-17T13:54:14.796+00:00",
//           "monPolDn": "uni/fabric/monfab-default",
//           "name": "Room01",
//           "status": "",
//           "uid": "15374"
//         }
//       }
//     }
//   ]
// }
//
// // ADD ROW
// "api/node/mo/uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01/row-Row01.json"
// {
//   "geoRow": {
//     "attributes": {
//       "dn": "uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01/row-Row01",
//       "name": "Row01",
//       "rn": "row-Row01",
//       "status": "created"
//     },
//     "children": []
//   }
// }
//
// // DELETE ROW
// "api/node/mo/uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01.json"
// {
//   "geoRoom": {
//     "attributes": {
//       "dn": "uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01",
//       "status": "modified"
//     },
//     "children": [
//       {
//         "geoRow": {
//           "attributes": {
//             "dn": "uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01/row-Row01",
//             "status": "deleted"
//           },
//           "children": []
//         }
//       }
//     ]
//   }
// }
//
// LIST ROWS
// "api/node/mo/uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01.json?query-target=children&target-subtree-class=geoRow"
// response:
// {
//   "totalCount": "1",
//   "imdata": [
//     {
//       "geoRow": {
//         "attributes": {
//           "childAction": "",
//           "descr": "",
//           "dn": "uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01/row-Row01",
//           "lcOwn": "local",
//           "modTs": "2017-07-17T13:54:14.796+00:00",
//           "monPolDn": "uni/fabric/monfab-default",
//           "name": "Row01",
//           "status": "",
//           "uid": "15374"
//         }
//       }
//     }
//   ]
// }
//
// // ADD RACK
// "api/node/mo/uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01/row-Row01/rack-Rack01.json"
// {
//   "geoRack": {
//     "attributes": {
//       "dn": "uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01/row-Row01/rack-Rack01",
//       "name": "Rack01",
//       "rn": "rack-Rack01",
//       "status": "created"
//     },
//     "children": []
//   }
// }
//
// // DELETE RACK
// "api/node/mo/uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01/row-Row01.json"
// {
//   "geoRow": {
//     "attributes": {
//       "dn": "uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01/row-Row01",
//       "status": "modified"
//     },
//     "children": [
//       {
//         "geoRack": {
//           "attributes": {
//             "dn": "uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01/row-Row01/rack-Rack01",
//             "status": "deleted"
//           },
//           "children": []
//         }
//       }
//     ]
//   }
// }
//
// LIST RACKS
// "api/node/mo/uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01/row-Row01.json?query-target=children&target-subtree-class=geoRack"
// response:
// {
//   "totalCount": "1",
//   "imdata": [
//     {
//       "geoRack": {
//         "attributes": {
//           "childAction": "",
//           "descr": "",
//           "dn": "uni/fabric/site-Site01/building-Building01/floor-Floor01/room-Room01/row-Row01/rack-Rack01",
//           "lcOwn": "local",
//           "modTs": "2017-07-17T13:54:14.796+00:00",
//           "monPolDn": "uni/fabric/monfab-default",
//           "name": "Rack01",
//           "status": "",
//           "uid": "15374"
//         }
//       }
//     }
//   ]
// }
