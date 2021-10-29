package jsonmodel

type ClientVersion struct {
	Platform string `json:"platform" binding:"required"`
	VersionString string `json:"version_string" binding:"required"`
	BuildIdentification string `json:"build_identification" binding:"required"`
}