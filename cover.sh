COVER_PROFILE=cover.out
go test fretcon/fretcon -covermode=count -coverprofile=$COVER_PROFILE
go tool cover -html=$COVER_PROFILE
