package graph

import "api-project/graphql-api/gql/graph/cablemodems"

// these stubs do nothing.
// gqlgen needs its subresolver functions like .CableModems(), .Transponders()
// to return non-nil pointers, but those pointers point to empty structs.
// so we allocate exactly once at program start, here.
// they can't be next to where they're used because gqlgen generate will move them out of the resolver files.
var (
	stubModemsResolver = new(cablemodems.CableModems)
)
