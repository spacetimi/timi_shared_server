package mongo_wrapper

type IDataItemDescriptor interface {
	GetDB() DBSpace
	GetCollectionName() string
	GetPrimaryKeys() []string
}

type IDataItem interface {
	GetDescriptor() IDataItemDescriptor
	IDirtyable
}

type IDataItemList interface {
	GetDataItems() []IDataItem
	AddDataItem(dataItem IDataItem)
}

// TODO: Rename this to something else?
type IDataItemCollection interface {
	GetDescriptor() IDataItemDescriptor
	GetDataItemFactory() IDataItemFactory
	IDirtyable
	IDataItemList
}

type IDirtyable interface {
	Dirty() bool
	SetDirty(dirty bool)
}

type IDataItemFactory interface {
	CreateDataItem() IDataItem
}
