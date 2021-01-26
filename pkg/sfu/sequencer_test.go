package sfu

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_sequencer(t *testing.T) {
	seq := newSequencer()
	off := uint16(15)

	for i := uint16(1); i < 100; i++ {
		seq.push(i, i+off, 123, true)
	}

	req := []uint16{17, 18, 22, 33}
	res := seq.getSeqNoPairs(req)
	assert.Equal(t, len(req), len(res))
	for i, val := range res {
		assert.Equal(t, val.getTargetSeqNo(), req[i])
		assert.Equal(t, val.getSourceSeqNo(), req[i]-off)
	}
	res = seq.getSeqNoPairs(req)
	assert.Equal(t, 0, len(res))
	time.Sleep(60 * time.Millisecond)
	res = seq.getSeqNoPairs(req)
	assert.Equal(t, len(req), len(res))
	for i, val := range res {
		assert.Equal(t, val.getSourceSeqNo(), req[i])
		assert.Equal(t, val.getSourceSeqNo(), req[i]-off)
	}
}

func Test_sequencer_getNACKSeqNo(t *testing.T) {
	type args struct {
		seqNo []uint16
	}
	type fields struct {
		input  []uint16
		offset uint16
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []uint16
	}{
		{
			name: "Should get correct seq numbers",
			fields: fields{
				input:  []uint16{2, 3, 4, 7, 8},
				offset: 5,
			},
			args: args{
				seqNo: []uint16{4 + 5, 5 + 5, 8 + 5},
			},
			want: []uint16{4, 8},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			n := newSequencer()

			for _, i := range tt.fields.input {
				n.push(i, i+tt.fields.offset, 123, true)
			}

			g := n.getSeqNoPairs(tt.args.seqNo)
			var got []uint16
			for _, sn := range g {
				got = append(got, sn.getSourceSeqNo())
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSeqNoPairs() = %v, want %v", got, tt.want)
			}
		})
	}
}
