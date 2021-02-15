package scraper

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func Test_rangeBlockCache_Set(t *testing.T) {
	type fields struct {
		c map[rangeInfo]types.Header
	}
	type args struct {
		r        rangeInfo
		h        types.Header
		expected map[rangeInfo]types.Header
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "simple add",
			fields: fields{c: make(map[rangeInfo]types.Header)},
			args: args{
				r: rangeInfo{from: 1, to: 20},
				h: types.Header{},
				expected: map[rangeInfo]types.Header{
					{from: 1, to: 20}: types.Header{},
				},
			},
		},
		{
			name:   "simple add - with header",
			fields: fields{c: make(map[rangeInfo]types.Header)},
			args: args{
				r: rangeInfo{from: 1, to: 20},
				h: types.Header{Number: big.NewInt(12), Time: 12345},
				expected: map[rangeInfo]types.Header{
					{from: 12, to: 20}: types.Header{
						Number: big.NewInt(12),
						Time:   12345,
					},
				},
			},
		},
		{
			name: "join left - empty",
			fields: fields{c: map[rangeInfo]types.Header{
				{from: 21, to: 40}: types.Header{},
			},
			},
			args: args{
				r: rangeInfo{from: 1, to: 20},
				h: types.Header{},
				expected: map[rangeInfo]types.Header{
					{from: 1, to: 40}: types.Header{},
				},
			},
		},
		{
			name: "join left - header",
			fields: fields{c: map[rangeInfo]types.Header{
				{from: 21, to: 40}: types.Header{},
			},
			},
			args: args{
				r: rangeInfo{from: 1, to: 20},
				h: types.Header{Number: big.NewInt(12), Time: 12345},
				expected: map[rangeInfo]types.Header{
					{from: 12, to: 40}: types.Header{Number: big.NewInt(12), Time: 12345},
				},
			},
		},
		{
			name: "join right - empty",
			fields: fields{c: map[rangeInfo]types.Header{
				{from: 1, to: 20}: types.Header{},
			},
			},
			args: args{
				r: rangeInfo{from: 21, to: 40},
				h: types.Header{},
				expected: map[rangeInfo]types.Header{
					{from: 1, to: 40}: types.Header{},
				},
			},
		},
		{
			name: "join right - header",
			fields: fields{c: map[rangeInfo]types.Header{
				{from: 1, to: 20}: types.Header{},
			},
			},
			args: args{
				r: rangeInfo{from: 21, to: 40},
				h: types.Header{Number: big.NewInt(25), Time: 12345},
				expected: map[rangeInfo]types.Header{
					{from: 25, to: 40}: types.Header{Number: big.NewInt(25), Time: 12345},
				},
			},
		},
		{
			name: "Intersection",
			fields: fields{c: map[rangeInfo]types.Header{
				{from: 1, to: 40}: types.Header{},
			},
			},
			args: args{
				r: rangeInfo{from: 10, to: 30},
				h: types.Header{},
				expected: map[rangeInfo]types.Header{
					{from: 1, to: 40}: types.Header{},
				},
			},
		},
		{
			name: "Intersection - header",
			fields: fields{c: map[rangeInfo]types.Header{
				{from: 1, to: 40}: types.Header{},
			},
			},
			args: args{
				r: rangeInfo{from: 10, to: 30},
				h: types.Header{Number: big.NewInt(25), Time: 12345},
				expected: map[rangeInfo]types.Header{
					{from: 25, to: 40}: types.Header{Number: big.NewInt(25), Time: 12345},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rbc := &rangeBlockCache{
				c: tt.fields.c,
			}
			rbc.Set(tt.args.r, tt.args.h)

			require.Equal(t, tt.args.expected, rbc.c)
		})
	}
}
