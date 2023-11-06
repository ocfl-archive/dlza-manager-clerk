## Example of queries :  


```query tenants {
  tenants {
    items{
    	id
      name
    }
    totalItems
  }
  }
  
  
  query GetAllTenants($optionsTenant: TenantListOptions){
  tenants(options: $optionsTenant) {
    items {
      id
      name
            person
      storageLocations(options: {
        take: 3
      }) {
        items {
          id,
          alias
        }
        totalItems
      }
      collections(options: {
        take: 4
      }) {
        items {
          id,
          alias
        }
        totalItems
      }
    }
    totalItems
  }
}

query GetObjectInstancesByPartition($optionsObjectInstance: ObjectInstanceListOptions){
  objectInstances(options: $optionsObjectInstance){
    items {
      id
      path
      objectInstanceChecks(options: {
        take: 2
        skip:0
      }) {
        items{
          id
          error
        }
        totalItems
      }
    }
    totalItems
  }
}
query ObjectInstanceCheck($optionsObjectInstanceCheck: ObjectInstanceCheckListOptions){
  objectInstanceChecks(options: $optionsObjectInstanceCheck){
    items{
      id
      error
      objectInstance{
        id
        path
      }
    }
    totalItems
  }
}
query StoragePartition($id: ID!){
  storagePartition(id: $id){
    id
    maxSize
    objectInstances(options:{
      take: 2
    }){
      items{
        id
        path
      }
      totalItems
    }
  }
}
```
## Variables 
```
{
  "optionsTenant": {
    "skip": 0,
    "take": 19,
    "sortKey": "PERSON",
    "sortDirection": "DESCENDING"
  },
    "optionsObjectInstance": {
    "ObjectId":  "317a2db6-3048-4f67-a467-031c9f3996e6",
    "skip": 0,
    "take": 20,
    "sortKey": "ID",
    "sortDirection": "DESCENDING"
    },
  "optionsObjectInstanceCheck": {
    "ObjectInstanceId": "b87a8e5c-46db-44e2-92ab-f050f50f43de",
    "skip": 0,
    "take": 20,
    "sortKey": "ID",
    "sortDirection": "ASCENDING"
  },
  "id": "5170a927-3ca1-45ab-873f-d478020b51e1"
}
```
