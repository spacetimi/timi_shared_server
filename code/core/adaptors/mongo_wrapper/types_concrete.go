package mongo_wrapper

/********************************************************************************/

/**
 * Implements IDataItemList
 */
type DataItemList struct {
	items []IDataItem
}

func (dil *DataItemList) GetDataItems() []IDataItem {
	return dil.items
}

func (dil *DataItemList) AddDataItem(item IDataItem) {
	dil.items = append(dil.items, item)
}

/********************************************************************************/

/**
 * Implements IDirtyable
 */
type Dirtyable struct {
	isDirty bool
}

func (dirtyable Dirtyable) Dirty() bool {
	return dirtyable.isDirty
}

func (dirtyable *Dirtyable) SetDirty(dirty bool) {
	dirtyable.isDirty = dirty
}

/********************************************************************************/

/**
 * Implements IDataItemDescriptor
 */
type DataItemDescriptor struct {
	DBType         DBSpace
	CollectionName string
	PrimaryKeys    []string
}

func (did *DataItemDescriptor) GetDB() DBSpace {
	return did.DBType
}

func (did *DataItemDescriptor) GetCollectionName() string {
	return did.CollectionName
}

func (did *DataItemDescriptor) GetPrimaryKeys() []string {
	return did.PrimaryKeys
}

/********************************************************************************/
