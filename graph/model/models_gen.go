// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Node interface {
	IsNode()
	GetID() string
}

type PaginatedList interface {
	IsPaginatedList()
	GetItems() []Node
	GetTotalItems() int
}

type Collection struct {
	ID               string      `json:"id"`
	Alias            string      `json:"alias"`
	Description      string      `json:"description"`
	Owner            string      `json:"owner"`
	OwnerMail        string      `json:"ownerMail"`
	Name             string      `json:"name"`
	Quality          int         `json:"quality"`
	TenantID         string      `json:"tenantId"`
	Tenant           *Tenant     `json:"tenant"`
	Objects          *ObjectList `json:"objects"`
	Files            *FileList   `json:"files"`
	TotalFileSize    int         `json:"totalFileSize"`
	TotalFileCount   int         `json:"totalFileCount"`
	TotalObjectCount int         `json:"totalObjectCount"`
}

func (Collection) IsNode()            {}
func (this Collection) GetID() string { return this.ID }

type CollectionList struct {
	Items      []*Collection `json:"items"`
	TotalItems int           `json:"totalItems"`
}

func (CollectionList) IsPaginatedList() {}
func (this CollectionList) GetItems() []Node {
	if this.Items == nil {
		return nil
	}
	interfaceSlice := make([]Node, 0, len(this.Items))
	for _, concrete := range this.Items {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}
func (this CollectionList) GetTotalItems() int { return this.TotalItems }

type CollectionListOptions struct {
	TenantID      *string            `json:"tenantId,omitempty"`
	Skip          *int               `json:"skip,omitempty"`
	Take          *int               `json:"take,omitempty"`
	SortDirection *SortDirection     `json:"sortDirection,omitempty"`
	SortKey       *CollectionSortKey `json:"sortKey,omitempty"`
	Search        *string            `json:"search,omitempty"`
}

type File struct {
	ID       string   `json:"id"`
	Checksum string   `json:"checksum"`
	Name     []string `json:"name"`
	MimeType string   `json:"mimeType"`
	Size     int      `json:"size"`
	Pronom   string   `json:"pronom"`
	Width    int      `json:"width"`
	Height   int      `json:"height"`
	Duration int      `json:"duration"`
	ObjectID string   `json:"objectId"`
	Object   *Object  `json:"object"`
}

func (File) IsNode()            {}
func (this File) GetID() string { return this.ID }

type FileList struct {
	Items      []*File `json:"items"`
	TotalItems int     `json:"totalItems"`
}

func (FileList) IsPaginatedList() {}
func (this FileList) GetItems() []Node {
	if this.Items == nil {
		return nil
	}
	interfaceSlice := make([]Node, 0, len(this.Items))
	for _, concrete := range this.Items {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}
func (this FileList) GetTotalItems() int { return this.TotalItems }

type FileListOptions struct {
	ObjectID      *string        `json:"objectId,omitempty"`
	CollectionID  *string        `json:"collectionId,omitempty"`
	Skip          *int           `json:"skip,omitempty"`
	Take          *int           `json:"take,omitempty"`
	SortDirection *SortDirection `json:"sortDirection,omitempty"`
	SortKey       *FileSortKey   `json:"sortKey,omitempty"`
	Search        *string        `json:"search,omitempty"`
}

type MimeType struct {
	ID        string `json:"id"`
	FileCount int    `json:"fileCount"`
}

func (MimeType) IsNode()            {}
func (this MimeType) GetID() string { return this.ID }

type MimeTypeList struct {
	Items      []*MimeType `json:"items"`
	TotalItems int         `json:"totalItems"`
}

func (MimeTypeList) IsPaginatedList() {}
func (this MimeTypeList) GetItems() []Node {
	if this.Items == nil {
		return nil
	}
	interfaceSlice := make([]Node, 0, len(this.Items))
	for _, concrete := range this.Items {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}
func (this MimeTypeList) GetTotalItems() int { return this.TotalItems }

type MimeTypeListOptions struct {
	CollectionID  *string          `json:"collectionId,omitempty"`
	Skip          *int             `json:"skip,omitempty"`
	Take          *int             `json:"take,omitempty"`
	SortDirection *SortDirection   `json:"sortDirection,omitempty"`
	SortKey       *MimeTypeSortKey `json:"sortKey,omitempty"`
}

type Object struct {
	ID                string              `json:"id"`
	Signature         string              `json:"signature"`
	Sets              []string            `json:"sets"`
	Identifiers       []string            `json:"identifiers"`
	Title             string              `json:"title"`
	AlternativeTitles []string            `json:"alternativeTitles"`
	Description       string              `json:"description"`
	Keywords          []string            `json:"keywords"`
	References        []string            `json:"references"`
	IngestWorkflow    string              `json:"ingestWorkflow"`
	User              string              `json:"user"`
	Address           string              `json:"address"`
	Created           string              `json:"created"`
	LastChanged       string              `json:"lastChanged"`
	Size              int                 `json:"size"`
	CollectionID      string              `json:"collectionId"`
	Collection        *Collection         `json:"collection"`
	Checksum          string              `json:"checksum"`
	ObjectInstances   *ObjectInstanceList `json:"objectInstances"`
	Files             *FileList           `json:"files"`
}

func (Object) IsNode()            {}
func (this Object) GetID() string { return this.ID }

type ObjectInstance struct {
	ID                   string                   `json:"id"`
	Path                 string                   `json:"path"`
	Created              string                   `json:"created"`
	Status               string                   `json:"status"`
	Size                 int                      `json:"size"`
	StoragePartitionID   string                   `json:"storagePartitionId"`
	StoragePartition     *StoragePartition        `json:"storagePartition"`
	ObjectID             string                   `json:"objectId"`
	Object               *Object                  `json:"object"`
	ObjectInstanceChecks *ObjectInstanceCheckList `json:"objectInstanceChecks"`
}

func (ObjectInstance) IsNode()            {}
func (this ObjectInstance) GetID() string { return this.ID }

type ObjectInstanceCheck struct {
	ID               string          `json:"id"`
	Checktime        string          `json:"checktime"`
	Error            bool            `json:"error"`
	Message          string          `json:"message"`
	ObjectInstanceID string          `json:"objectInstanceId"`
	ObjectInstance   *ObjectInstance `json:"objectInstance"`
}

func (ObjectInstanceCheck) IsNode()            {}
func (this ObjectInstanceCheck) GetID() string { return this.ID }

type ObjectInstanceCheckList struct {
	Items      []*ObjectInstanceCheck `json:"items"`
	TotalItems int                    `json:"totalItems"`
}

func (ObjectInstanceCheckList) IsPaginatedList() {}
func (this ObjectInstanceCheckList) GetItems() []Node {
	if this.Items == nil {
		return nil
	}
	interfaceSlice := make([]Node, 0, len(this.Items))
	for _, concrete := range this.Items {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}
func (this ObjectInstanceCheckList) GetTotalItems() int { return this.TotalItems }

type ObjectInstanceCheckListOptions struct {
	ObjectInstanceID *string                     `json:"objectInstanceId,omitempty"`
	Skip             *int                        `json:"skip,omitempty"`
	Take             *int                        `json:"take,omitempty"`
	SortDirection    *SortDirection              `json:"sortDirection,omitempty"`
	SortKey          *ObjectInstanceCheckSortKey `json:"sortKey,omitempty"`
	Search           *string                     `json:"search,omitempty"`
}

type ObjectInstanceList struct {
	Items      []*ObjectInstance `json:"items"`
	TotalItems int               `json:"totalItems"`
}

func (ObjectInstanceList) IsPaginatedList() {}
func (this ObjectInstanceList) GetItems() []Node {
	if this.Items == nil {
		return nil
	}
	interfaceSlice := make([]Node, 0, len(this.Items))
	for _, concrete := range this.Items {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}
func (this ObjectInstanceList) GetTotalItems() int { return this.TotalItems }

type ObjectInstanceListOptions struct {
	ObjectID      *string                `json:"ObjectId,omitempty"`
	Skip          *int                   `json:"skip,omitempty"`
	Take          *int                   `json:"take,omitempty"`
	SortDirection *SortDirection         `json:"sortDirection,omitempty"`
	SortKey       *ObjectInstanceSortKey `json:"sortKey,omitempty"`
	Search        *string                `json:"search,omitempty"`
}

type ObjectList struct {
	Items      []*Object `json:"items"`
	TotalItems int       `json:"totalItems"`
}

func (ObjectList) IsPaginatedList() {}
func (this ObjectList) GetItems() []Node {
	if this.Items == nil {
		return nil
	}
	interfaceSlice := make([]Node, 0, len(this.Items))
	for _, concrete := range this.Items {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}
func (this ObjectList) GetTotalItems() int { return this.TotalItems }

type ObjectListOptions struct {
	CollectionID  *string        `json:"collectionId,omitempty"`
	Skip          *int           `json:"skip,omitempty"`
	Take          *int           `json:"take,omitempty"`
	SortDirection *SortDirection `json:"sortDirection,omitempty"`
	SortKey       *ObjectSortKey `json:"sortKey,omitempty"`
	Search        *string        `json:"search,omitempty"`
}

type StorageLocation struct {
	ID                 string                `json:"id"`
	Alias              string                `json:"alias"`
	Type               string                `json:"type"`
	Vault              string                `json:"vault"`
	Connection         string                `json:"connection"`
	Quality            int                   `json:"quality"`
	Price              int                   `json:"price"`
	SecurityCompliency string                `json:"securityCompliency"`
	FillFirst          bool                  `json:"fillFirst"`
	OcflType           string                `json:"ocflType"`
	TenantID           string                `json:"tenantId"`
	Tenant             *Tenant               `json:"tenant"`
	NumberOfThreads    int                   `json:"numberOfThreads"`
	StoragePartitions  *StoragePartitionList `json:"storagePartitions"`
}

func (StorageLocation) IsNode()            {}
func (this StorageLocation) GetID() string { return this.ID }

type StorageLocationList struct {
	Items      []*StorageLocation `json:"items"`
	TotalItems int                `json:"totalItems"`
}

func (StorageLocationList) IsPaginatedList() {}
func (this StorageLocationList) GetItems() []Node {
	if this.Items == nil {
		return nil
	}
	interfaceSlice := make([]Node, 0, len(this.Items))
	for _, concrete := range this.Items {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}
func (this StorageLocationList) GetTotalItems() int { return this.TotalItems }

type StorageLocationListOptions struct {
	TenantID      *string                 `json:"tenantId,omitempty"`
	Skip          *int                    `json:"skip,omitempty"`
	Take          *int                    `json:"take,omitempty"`
	SortDirection *SortDirection          `json:"sortDirection,omitempty"`
	SortKey       *StorageLocationSortKey `json:"sortKey,omitempty"`
	Search        *string                 `json:"search,omitempty"`
}

type StoragePartition struct {
	ID                string              `json:"id"`
	Alias             string              `json:"alias"`
	Name              string              `json:"name"`
	MaxSize           int                 `json:"maxSize"`
	MaxObjects        int                 `json:"maxObjects"`
	CurrentSize       int                 `json:"currentSize"`
	CurrentObjects    int                 `json:"currentObjects"`
	StorageLocationID string              `json:"storageLocationId"`
	StorageLocation   *StorageLocation    `json:"storageLocation"`
	ObjectInstances   *ObjectInstanceList `json:"objectInstances"`
}

func (StoragePartition) IsNode()            {}
func (this StoragePartition) GetID() string { return this.ID }

type StoragePartitionList struct {
	Items      []*StoragePartition `json:"items"`
	TotalItems int                 `json:"totalItems"`
}

func (StoragePartitionList) IsPaginatedList() {}
func (this StoragePartitionList) GetItems() []Node {
	if this.Items == nil {
		return nil
	}
	interfaceSlice := make([]Node, 0, len(this.Items))
	for _, concrete := range this.Items {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}
func (this StoragePartitionList) GetTotalItems() int { return this.TotalItems }

type StoragePartitionListOptions struct {
	StorageLocationID *string                  `json:"storageLocationId,omitempty"`
	Skip              *int                     `json:"skip,omitempty"`
	Take              *int                     `json:"take,omitempty"`
	SortDirection     *SortDirection           `json:"sortDirection,omitempty"`
	SortKey           *StoragePartitionSortKey `json:"sortKey,omitempty"`
	Search            *string                  `json:"search,omitempty"`
}

type Tenant struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	Alias            string               `json:"alias"`
	Person           string               `json:"person"`
	Email            string               `json:"email"`
	Collections      *CollectionList      `json:"collections"`
	StorageLocations *StorageLocationList `json:"storageLocations"`
}

func (Tenant) IsNode()            {}
func (this Tenant) GetID() string { return this.ID }

type TenantList struct {
	Items      []*Tenant `json:"items"`
	TotalItems int       `json:"totalItems"`
}

func (TenantList) IsPaginatedList() {}
func (this TenantList) GetItems() []Node {
	if this.Items == nil {
		return nil
	}
	interfaceSlice := make([]Node, 0, len(this.Items))
	for _, concrete := range this.Items {
		interfaceSlice = append(interfaceSlice, concrete)
	}
	return interfaceSlice
}
func (this TenantList) GetTotalItems() int { return this.TotalItems }

type TenantListOptions struct {
	Skip          *int           `json:"skip,omitempty"`
	Take          *int           `json:"take,omitempty"`
	SortDirection *SortDirection `json:"sortDirection,omitempty"`
	SortKey       *TenantSortKey `json:"sortKey,omitempty"`
	Search        *string        `json:"search,omitempty"`
}

type CollectionSortKey string

const (
	CollectionSortKeyID               CollectionSortKey = "id"
	CollectionSortKeyName             CollectionSortKey = "name"
	CollectionSortKeyAlias            CollectionSortKey = "alias"
	CollectionSortKeyDescription      CollectionSortKey = "description"
	CollectionSortKeyOwner            CollectionSortKey = "owner"
	CollectionSortKeyOwnerMail        CollectionSortKey = "ownerMail"
	CollectionSortKeyTotalFileSize    CollectionSortKey = "totalFileSize"
	CollectionSortKeyTotalFileCount   CollectionSortKey = "totalFileCount"
	CollectionSortKeyTotalObjectCount CollectionSortKey = "totalObjectCount"
)

var AllCollectionSortKey = []CollectionSortKey{
	CollectionSortKeyID,
	CollectionSortKeyName,
	CollectionSortKeyAlias,
	CollectionSortKeyDescription,
	CollectionSortKeyOwner,
	CollectionSortKeyOwnerMail,
	CollectionSortKeyTotalFileSize,
	CollectionSortKeyTotalFileCount,
	CollectionSortKeyTotalObjectCount,
}

func (e CollectionSortKey) IsValid() bool {
	switch e {
	case CollectionSortKeyID, CollectionSortKeyName, CollectionSortKeyAlias, CollectionSortKeyDescription, CollectionSortKeyOwner, CollectionSortKeyOwnerMail, CollectionSortKeyTotalFileSize, CollectionSortKeyTotalFileCount, CollectionSortKeyTotalObjectCount:
		return true
	}
	return false
}

func (e CollectionSortKey) String() string {
	return string(e)
}

func (e *CollectionSortKey) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CollectionSortKey(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CollectionSortKey", str)
	}
	return nil
}

func (e CollectionSortKey) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type FileSortKey string

const (
	FileSortKeyID       FileSortKey = "id"
	FileSortKeyChecksum FileSortKey = "checksum"
	FileSortKeyMimeType FileSortKey = "mimeType"
	FileSortKeyPronom   FileSortKey = "pronom"
	FileSortKeyName     FileSortKey = "name"
)

var AllFileSortKey = []FileSortKey{
	FileSortKeyID,
	FileSortKeyChecksum,
	FileSortKeyMimeType,
	FileSortKeyPronom,
	FileSortKeyName,
}

func (e FileSortKey) IsValid() bool {
	switch e {
	case FileSortKeyID, FileSortKeyChecksum, FileSortKeyMimeType, FileSortKeyPronom, FileSortKeyName:
		return true
	}
	return false
}

func (e FileSortKey) String() string {
	return string(e)
}

func (e *FileSortKey) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FileSortKey(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FileSortKey", str)
	}
	return nil
}

func (e FileSortKey) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type MimeTypeSortKey string

const (
	MimeTypeSortKeyID        MimeTypeSortKey = "id"
	MimeTypeSortKeyFileCount MimeTypeSortKey = "fileCount"
)

var AllMimeTypeSortKey = []MimeTypeSortKey{
	MimeTypeSortKeyID,
	MimeTypeSortKeyFileCount,
}

func (e MimeTypeSortKey) IsValid() bool {
	switch e {
	case MimeTypeSortKeyID, MimeTypeSortKeyFileCount:
		return true
	}
	return false
}

func (e MimeTypeSortKey) String() string {
	return string(e)
}

func (e *MimeTypeSortKey) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MimeTypeSortKey(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MimeTypeSortKey", str)
	}
	return nil
}

func (e MimeTypeSortKey) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ObjectInstanceCheckSortKey string

const (
	ObjectInstanceCheckSortKeyID        ObjectInstanceCheckSortKey = "id"
	ObjectInstanceCheckSortKeyMessage   ObjectInstanceCheckSortKey = "message"
	ObjectInstanceCheckSortKeyChecktime ObjectInstanceCheckSortKey = "checktime"
)

var AllObjectInstanceCheckSortKey = []ObjectInstanceCheckSortKey{
	ObjectInstanceCheckSortKeyID,
	ObjectInstanceCheckSortKeyMessage,
	ObjectInstanceCheckSortKeyChecktime,
}

func (e ObjectInstanceCheckSortKey) IsValid() bool {
	switch e {
	case ObjectInstanceCheckSortKeyID, ObjectInstanceCheckSortKeyMessage, ObjectInstanceCheckSortKeyChecktime:
		return true
	}
	return false
}

func (e ObjectInstanceCheckSortKey) String() string {
	return string(e)
}

func (e *ObjectInstanceCheckSortKey) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ObjectInstanceCheckSortKey(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ObjectInstanceCheckSortKey", str)
	}
	return nil
}

func (e ObjectInstanceCheckSortKey) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ObjectInstanceSortKey string

const (
	ObjectInstanceSortKeyID     ObjectInstanceSortKey = "id"
	ObjectInstanceSortKeyPath   ObjectInstanceSortKey = "path"
	ObjectInstanceSortKeyStatus ObjectInstanceSortKey = "status"
)

var AllObjectInstanceSortKey = []ObjectInstanceSortKey{
	ObjectInstanceSortKeyID,
	ObjectInstanceSortKeyPath,
	ObjectInstanceSortKeyStatus,
}

func (e ObjectInstanceSortKey) IsValid() bool {
	switch e {
	case ObjectInstanceSortKeyID, ObjectInstanceSortKeyPath, ObjectInstanceSortKeyStatus:
		return true
	}
	return false
}

func (e ObjectInstanceSortKey) String() string {
	return string(e)
}

func (e *ObjectInstanceSortKey) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ObjectInstanceSortKey(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ObjectInstanceSortKey", str)
	}
	return nil
}

func (e ObjectInstanceSortKey) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ObjectSortKey string

const (
	ObjectSortKeyID                ObjectSortKey = "id"
	ObjectSortKeySignature         ObjectSortKey = "signature"
	ObjectSortKeyTitle             ObjectSortKey = "title"
	ObjectSortKeyDescription       ObjectSortKey = "description"
	ObjectSortKeyIngestWorkflow    ObjectSortKey = "ingestWorkflow"
	ObjectSortKeyUser              ObjectSortKey = "user"
	ObjectSortKeyAddress           ObjectSortKey = "address"
	ObjectSortKeyChecksum          ObjectSortKey = "checksum"
	ObjectSortKeyKeywords          ObjectSortKey = "keywords"
	ObjectSortKeyIdentifiers       ObjectSortKey = "identifiers"
	ObjectSortKeyAlternativeTitles ObjectSortKey = "alternativeTitles"
)

var AllObjectSortKey = []ObjectSortKey{
	ObjectSortKeyID,
	ObjectSortKeySignature,
	ObjectSortKeyTitle,
	ObjectSortKeyDescription,
	ObjectSortKeyIngestWorkflow,
	ObjectSortKeyUser,
	ObjectSortKeyAddress,
	ObjectSortKeyChecksum,
	ObjectSortKeyKeywords,
	ObjectSortKeyIdentifiers,
	ObjectSortKeyAlternativeTitles,
}

func (e ObjectSortKey) IsValid() bool {
	switch e {
	case ObjectSortKeyID, ObjectSortKeySignature, ObjectSortKeyTitle, ObjectSortKeyDescription, ObjectSortKeyIngestWorkflow, ObjectSortKeyUser, ObjectSortKeyAddress, ObjectSortKeyChecksum, ObjectSortKeyKeywords, ObjectSortKeyIdentifiers, ObjectSortKeyAlternativeTitles:
		return true
	}
	return false
}

func (e ObjectSortKey) String() string {
	return string(e)
}

func (e *ObjectSortKey) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ObjectSortKey(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ObjectSortKey", str)
	}
	return nil
}

func (e ObjectSortKey) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SortDirection string

const (
	SortDirectionAscending  SortDirection = "ASCENDING"
	SortDirectionDescending SortDirection = "DESCENDING"
)

var AllSortDirection = []SortDirection{
	SortDirectionAscending,
	SortDirectionDescending,
}

func (e SortDirection) IsValid() bool {
	switch e {
	case SortDirectionAscending, SortDirectionDescending:
		return true
	}
	return false
}

func (e SortDirection) String() string {
	return string(e)
}

func (e *SortDirection) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortDirection(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortDirection", str)
	}
	return nil
}

func (e SortDirection) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type StorageLocationSortKey string

const (
	StorageLocationSortKeyID                 StorageLocationSortKey = "id"
	StorageLocationSortKeyAlias              StorageLocationSortKey = "alias"
	StorageLocationSortKeySecurityCompliency StorageLocationSortKey = "securityCompliency"
)

var AllStorageLocationSortKey = []StorageLocationSortKey{
	StorageLocationSortKeyID,
	StorageLocationSortKeyAlias,
	StorageLocationSortKeySecurityCompliency,
}

func (e StorageLocationSortKey) IsValid() bool {
	switch e {
	case StorageLocationSortKeyID, StorageLocationSortKeyAlias, StorageLocationSortKeySecurityCompliency:
		return true
	}
	return false
}

func (e StorageLocationSortKey) String() string {
	return string(e)
}

func (e *StorageLocationSortKey) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = StorageLocationSortKey(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid StorageLocationSortKey", str)
	}
	return nil
}

func (e StorageLocationSortKey) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type StoragePartitionSortKey string

const (
	StoragePartitionSortKeyID    StoragePartitionSortKey = "id"
	StoragePartitionSortKeyAlias StoragePartitionSortKey = "alias"
	StoragePartitionSortKeyName  StoragePartitionSortKey = "name"
)

var AllStoragePartitionSortKey = []StoragePartitionSortKey{
	StoragePartitionSortKeyID,
	StoragePartitionSortKeyAlias,
	StoragePartitionSortKeyName,
}

func (e StoragePartitionSortKey) IsValid() bool {
	switch e {
	case StoragePartitionSortKeyID, StoragePartitionSortKeyAlias, StoragePartitionSortKeyName:
		return true
	}
	return false
}

func (e StoragePartitionSortKey) String() string {
	return string(e)
}

func (e *StoragePartitionSortKey) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = StoragePartitionSortKey(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid StoragePartitionSortKey", str)
	}
	return nil
}

func (e StoragePartitionSortKey) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TenantSortKey string

const (
	TenantSortKeyID     TenantSortKey = "id"
	TenantSortKeyName   TenantSortKey = "name"
	TenantSortKeyAlias  TenantSortKey = "alias"
	TenantSortKeyPerson TenantSortKey = "person"
	TenantSortKeyEmail  TenantSortKey = "email"
)

var AllTenantSortKey = []TenantSortKey{
	TenantSortKeyID,
	TenantSortKeyName,
	TenantSortKeyAlias,
	TenantSortKeyPerson,
	TenantSortKeyEmail,
}

func (e TenantSortKey) IsValid() bool {
	switch e {
	case TenantSortKeyID, TenantSortKeyName, TenantSortKeyAlias, TenantSortKeyPerson, TenantSortKeyEmail:
		return true
	}
	return false
}

func (e TenantSortKey) String() string {
	return string(e)
}

func (e *TenantSortKey) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TenantSortKey(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TenantSortKey", str)
	}
	return nil
}

func (e TenantSortKey) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
