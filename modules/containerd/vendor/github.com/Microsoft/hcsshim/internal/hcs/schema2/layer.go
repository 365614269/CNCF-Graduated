/*
 * HCS API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 2.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package hcsschema

type FileSystemFilterType string

const (
	UnionFS FileSystemFilterType = "UnionFS"
	WCIFS   FileSystemFilterType = "WCIFS"
)

type Layer struct {
	Id string `json:"Id,omitempty"`

	Path string `json:"Path,omitempty"`

	PathType string `json:"PathType,omitempty"`

	//  Unspecified defaults to Enabled
	Cache string `json:"Cache,omitempty"`
}
