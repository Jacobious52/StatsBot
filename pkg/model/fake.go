package model

// FakeSeries defines testing data
var FakeSeries = Timeseries{
	1: map[string]uint64{
		"User A": 4,
		"User B": 10,
		"User C": 5,
	},
	2: map[string]uint64{
		"User A": 1,
		"User B": 9,
		"User C": 4,
	},
	3: map[string]uint64{
		"User A": 6,
		"User B": 12,
		"User C": 5,
	},
	4: map[string]uint64{
		"User A": 1,
		"User B": 2,
		"User C": 9,
	},
	5: map[string]uint64{
		"User A": 4,
		"User B": 4,
		"User C": 3,
	},
}
