package chf

// max value of uint32
const maxUint32 = ^uint32(0)

// default value if not found
var defaultValue Custom_value

// CustomsIndex indexes a company's custom column values for a single flow record at a time in a performant way,
// using a few slices (two []uint32 and one []Custom_value), no maps, no clean-up processing when finished with a flow.
// The only allocations made are three slices at initialization.
//
// Performance:
// - when indexing a flow, we scan its Customs list once, read once from a []uint32, write to two slices
// - when fetching a Custom by custom dimension ID, we read from 3 slices
type CustomsIndex struct {
	// company's max custom dimension ID
	maxCustomDimensionID uint32

	// incrementing value keeping track of which flow index we're currently indexing
	// - flow number runs from 1->maxUint32
	currentFlowNumber uint32

	// map a company's custom dimension ID to a sorted index of their IDs
	// - sized to to store up to the max custom dimension ID this company has access to
	columnIDToCompanyIndex []uint32

	// the indexed Custom for the company's custom dimension index
	// - sized to hold the number of custom dimensions this company has access to
	valueAtColumnIndex []Custom_value

	// the value of 'currentFlowNumber' when entry in 'valueAtColumnIndex' was set
	// - sized to hold the number of custom dimensions this company has access to
	flowNumberAtColumnIndex []uint32
}

// NewCustomsIndex builds a new CustomsIndex for a company, with all of the custom dimension IDs
// that will be considered.
func NewCustomsIndex(companyCIDs []uint32) *CustomsIndex {
	ret := CustomsIndex{
		currentFlowNumber:       0,
		valueAtColumnIndex:      make([]Custom_value, len(companyCIDs)),
		flowNumberAtColumnIndex: make([]uint32, len(companyCIDs)),
	}

	if len(companyCIDs) == 0 {
		ret.columnIDToCompanyIndex = make([]uint32, 0)
		return &ret
	}

	// populate the CID -> index
	Uint32Slice(companyCIDs).Sort()
	ret.maxCustomDimensionID = companyCIDs[len(companyCIDs)-1]
	ret.columnIDToCompanyIndex = make([]uint32, ret.maxCustomDimensionID+1)

	// build the columnID -> company column index
	for index, customDimensionID := range companyCIDs {
		// need to offset by 1, so a value of 0 means the column isn't in use by this company
		ret.columnIDToCompanyIndex[customDimensionID] = uint32(index + 1)
	}

	return &ret
}

// return the company's custom dimension index (0-based) for the input custom dimension ID
func (c *CustomsIndex) companyIndexForColumnID(customDimensionID uint32) (uint32, bool) {
	if customDimensionID <= c.maxCustomDimensionID {
		if companyColumnIndex := c.columnIDToCompanyIndex[customDimensionID]; companyColumnIndex > 0 {
			// we store the index in columnIDToCompanyIndex with increment, so 0 means not found - need to subtract 1 from it now
			return companyColumnIndex - 1, true
		}
	}
	return 0, false
}

// IndexFlow indexes the input flow
// - the flow's customs will only ever be scanned once
func (c *CustomsIndex) IndexFlow(flow CHF) {
	if c.currentFlowNumber == maxUint32 {
		// need to reset
		c.currentFlowNumber = 0
		count := len(c.flowNumberAtColumnIndex)
		for i := 0; i < count; i++ {
			c.flowNumberAtColumnIndex[i] = 0
		}
	}

	// this ensures that the flow number is always greater than 0
	c.currentFlowNumber++

	// put each custom in the index
	customs, err := flow.Custom()
	if err != nil {
		return
	}
	var custom Custom
	for i := 0; i < customs.Len(); i++ {
		custom = customs.At(i)
		customID := custom.Id()

		if companyColumnIndex, found := c.companyIndexForColumnID(customID); found {
			// store the custom in the index
			c.valueAtColumnIndex[companyColumnIndex] = custom.Value()
			c.flowNumberAtColumnIndex[companyColumnIndex] = c.currentFlowNumber
		}
	}
}

// CustomColumnWithID returns custom column for the input ID
// - make sure to check the second "found" return value. If false, you'll likely get dangerous garbage you shouldn't mess with.
func (c *CustomsIndex) CustomColumnWithID(customDimensionID uint32) (Custom_value, bool) {
	if companyColumnIndex, found := c.companyIndexForColumnID(customDimensionID); found {
		return c.valueAtColumnIndex[companyColumnIndex], c.flowNumberAtColumnIndex[companyColumnIndex] == c.currentFlowNumber
	}

	return defaultValue, false
}
