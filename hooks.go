package fuel

// HookSave allows for validations to be performed on both PUT and POST
// REST API calls (i.e before insert and updation to DB)
type HookSave interface {
	PreSave(a Aide) error
}

// HookInsert allows for validations to be done on POST call,
// i.e before inserting to DB
type HookInsert interface {
	PreInsert(a Aide) error
}

// HookUpdate allows for validations to be done on PUT call,
// i.e before updating a record to DB
type HookUpdate interface {
	PreUpdate(a Aide) error
}
