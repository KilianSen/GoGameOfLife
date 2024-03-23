package Watchdogs

type callback func(updatedData *string) error
type Destructor func()
