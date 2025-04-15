/*
Copyright © 2025 Nicolò Piovan <nicopiovan@gmail.com>
*/

package couchdb

type Database struct {
	IstanceStartTime  string  `json:"instance_start_time"`
	DbName            string  `json:"db_name"`
	PurgeSeq          string  `json:"purge_seq"`
	UpdateSeq         string  `json:"update_seq"`
	Sizes             Sizes   `json:"sizes"`
	Props             Props   `json:"props"`
	DocDelCount       int     `json:"doc_del_count"`
	DocCount          int     `json:"doc_count"`
	DiskFormatVersion int     `json:"disk_format_version"`
	CompactRunning    bool    `json:"compact_running"`
	Cluster           Cluster `json:"cluster"`
}

type Sizes struct {
	File     int `json:"file"`
	External int `json:"external"`
	Active   int `json:"active"`
}

type Props struct {
	Partitioned bool `json:"partitioned"`
}

type Cluster struct {
	N int `json:"n"`
	Q int `json:"q"`
	R int `json:"r"`
	W int `json:"w"`
}
