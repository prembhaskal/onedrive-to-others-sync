package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DirectoryObjectCollectionResponse 
type DirectoryObjectCollectionResponse struct {
    BaseCollectionPaginationCountResponse
}
// NewDirectoryObjectCollectionResponse instantiates a new DirectoryObjectCollectionResponse and sets the default values.
func NewDirectoryObjectCollectionResponse()(*DirectoryObjectCollectionResponse) {
    m := &DirectoryObjectCollectionResponse{
        BaseCollectionPaginationCountResponse: *NewBaseCollectionPaginationCountResponse(),
    }
    return m
}
// CreateDirectoryObjectCollectionResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDirectoryObjectCollectionResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDirectoryObjectCollectionResponse(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DirectoryObjectCollectionResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseCollectionPaginationCountResponse.GetFieldDeserializers()
    res["value"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                res[i] = v.(DirectoryObjectable)
            }
            m.SetValue(res)
        }
        return nil
    }
    return res
}
// GetValue gets the value property value. The value property
func (m *DirectoryObjectCollectionResponse) GetValue()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("value")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DirectoryObjectCollectionResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseCollectionPaginationCountResponse.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetValue() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetValue()))
        for i, v := range m.GetValue() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("value", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetValue sets the value property value. The value property
func (m *DirectoryObjectCollectionResponse) SetValue(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("value", value)
    if err != nil {
        panic(err)
    }
}
// DirectoryObjectCollectionResponseable 
type DirectoryObjectCollectionResponseable interface {
    BaseCollectionPaginationCountResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetValue()([]DirectoryObjectable)
    SetValue(value []DirectoryObjectable)()
}
