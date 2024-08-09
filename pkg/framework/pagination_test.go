// Merge data list from multiple pages

package framework

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"
)

type IMockPaginator interface {
	GetConf() def.ConfPaginator
	GetPaginator() IPaginator
}

type mockPaginator struct {
	conf def.ConfPaginator
	// Indicate if length of last page is the same as limit
	fullLastPage bool
	// Indicate if TotalCount in NextCondition is set to negative (means not given)
	noTotalCount bool
	// If not null, return the result of the callback function
	fnCallback cbMockPaginator
}

func (p *mockPaginator) GetConf() def.ConfPaginator {
	return p.conf
}

func (p *mockPaginator) GetPaginator() IPaginator {
	return p
}

type cbMockPaginator func() ([]*json.RawMessage, NextCondition, error)

func (p *mockPaginator) setFullLastPage(flag bool) *mockPaginator {
	p.fullLastPage = flag
	return p
}

func (p *mockPaginator) setNoTotalCount(flag bool) *mockPaginator {
	p.noTotalCount = flag
	return p
}

func (p *mockPaginator) setFnCallback(fn cbMockPaginator) *mockPaginator {
	p.fnCallback = fn
	return p
}

func (p *mockPaginator) GetOnePage(paginationParam map[string]any) ([]*json.RawMessage, NextCondition, error) {
	if p.fnCallback != nil {
		return p.fnCallback()
	}

	switch p.conf.PaginationType {
	case def.PAGE_OFFSET_LIMIT:
		offset, ok := paginationParam[p.conf.OffsetName]
		if !ok {
			return nil, NextCondition{}, errors.New("invalid test suite")
		}
		iOffset := parseToInt(offset, p.conf.OffsetType)
		if iOffset < 0 {
			return nil, NextCondition{}, errors.New("invalid test suite")
		}
		limit, ok := paginationParam[p.conf.LimitName]
		if !ok {
			return nil, NextCondition{}, errors.New("invalid test suite")
		}
		iLimit := parseToInt(limit, p.conf.LimitType)
		if iLimit <= 0 {
			return nil, NextCondition{}, errors.New("invalid test suite")
		}

		iTotal := iLimit * 2
		if !p.fullLastPage && iLimit > 1 {
			iTotal -= 1
		}
		iItems := iTotal - iOffset
		if iItems < 0 {
			iItems = 0
		}
		if iItems > iLimit {
			iItems = iLimit
		}

		if p.noTotalCount {
			iTotal = -1 // iTotal not output in data
		}

		rm, _ := internal.JsonMarshal(
			map[string]any{
				p.conf.RespTotalName: iTotal,
				"Items":              make([]int, iItems),
			})

		return ResultDataParse(rm, p.conf, "$.Items")
	case def.PAGE_CURPAGE_SIZE:
		offset, ok := paginationParam[p.conf.OffsetName]
		if !ok {
			return nil, NextCondition{}, errors.New("invalid test suite")
		}
		iOffset := parseToInt(offset, p.conf.OffsetType)
		// iOffset is curpage starts with 1
		if iOffset <= 0 {
			return nil, NextCondition{}, errors.New("invalid test suite")
		}
		limit, ok := paginationParam[p.conf.LimitName]
		if !ok {
			return nil, NextCondition{}, errors.New("invalid test suite")
		}
		iLimit := parseToInt(limit, p.conf.LimitType)
		if iLimit <= 0 {
			return nil, NextCondition{}, errors.New("invalid test suite")
		}

		iTotal := iLimit * 2
		if !p.fullLastPage && iLimit > 1 {
			iTotal -= 1
		}
		iItems := iLimit
		if iOffset == 2 {
			iItems = iTotal - iLimit
		} else if iOffset > 2 {
			iItems = 0
		} // else iOffset==1, iItems = iLimit

		if p.noTotalCount {
			iTotal = -1 // iTotal not output in data
		}

		rm, _ := internal.JsonMarshal(
			map[string]any{
				p.conf.RespTotalName: iTotal,
				"Items":              make([]int, iItems),
			})

		return ResultDataParse(rm, p.conf, "$.Items")
	case def.PAGE_NOPAGEINATION:
		rm, _ := internal.JsonMarshal(
			map[string]any{
				"Items": make([]int, 8),
			})
		return ResultDataParse(rm, p.conf, "$.Items")
	case def.PAGE_MARKER:
		limit, ok := paginationParam[p.conf.LimitName]
		if !ok {
			return nil, NextCondition{}, errors.New("invalid test suite")
		}
		iLimit := parseToInt(limit, p.conf.LimitType)
		if iLimit <= 0 {
			return nil, NextCondition{}, errors.New("invalid test suite")
		}
		marker, ok := paginationParam[p.conf.MarkerName].(string)
		if !ok {
			return nil, NextCondition{}, errors.New("invalid test suite")
		}

		iItems := iLimit
		if len(marker) > 0 && !p.fullLastPage && iLimit > 1 {
			iItems -= 1
		}

		r := make(map[string]any)
		r["Items"] = make([]int, iItems)
		if len(marker) == 0 {
			r[p.conf.NextMarkerName] = "marker"
		}
		if len(p.conf.TruncatedName) > 0 {
			r[p.conf.TruncatedName] = len(marker) == 0
		}

		rm, _ := internal.JsonMarshal(r)
		return ResultDataParse(rm, p.conf, "$.Items")
	default:
		return nil, NextCondition{}, errors.New("invalid PaginationType")
	}
}

func parseToInt(v any, p def.ParamType) int {
	r := -1
	ok := false
	switch p {
	case def.PARAM_INT:
		r, ok = v.(int)
	case def.PARAM_STRING:
		var s string
		if s, ok = v.(string); ok {
			if r2, err := strconv.Atoi(s); err != nil {
				ok = false
			} else {
				r = r2
			}
		}
	case def.PARAM_STRING_LIST:
		var sl []string
		if sl, ok = v.([]string); ok && len(sl) > 0 {
			s := sl[0]
			if r2, err := strconv.Atoi(s); err != nil {
				ok = false
			} else {
				r = r2
			}
		} else {
			ok = false
		}
	}

	if !ok {
		return -1
	} else {
		return r
	}
}

type mockEmptyNopaginationPaginator struct {
	mockPaginator
}

func (p *mockEmptyNopaginationPaginator) GetPaginator() IPaginator {
	return p
}

func (p *mockEmptyNopaginationPaginator) GetOnePage(paginationParam map[string]any) ([]*json.RawMessage, NextCondition, error) {
	return nil, NextCondition{}, auth.ProfileNotDefinedError{}
}

func TestGetEntireList(t *testing.T) {
	SetPageSize(5)
	menp := mockEmptyNopaginationPaginator{
		mockPaginator: mockPaginator{
			conf: def.ConfPaginator{
				PaginationType: def.PAGE_NOPAGEINATION,
			},
		}}

	type args struct {
		p IMockPaginator
	}
	tests := []struct {
		name     string
		args     args
		wantSize int
		wantErr  bool
	}{
		{
			"Valid PAGE_OFFSET_LIMIT with fullLastPage==true",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_OFFSET_LIMIT,
						OffsetName:     "offset",
						OffsetType:     def.PARAM_INT,
						LimitName:      "limit",
						LimitType:      def.PARAM_INT,
						RespTotalName:  "total",
					},
				}).
					setFullLastPage(true),
			},
			10,
			false,
		},
		{
			"Valid PAGE_OFFSET_LIMIT with fullLastPage==false",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_OFFSET_LIMIT,
						OffsetName:     "offset",
						OffsetType:     def.PARAM_INT,
						LimitName:      "limit",
						LimitType:      def.PARAM_INT,
						RespTotalName:  "total",
					},
				}),
			},
			9,
			false,
		},
		{
			"Valid PAGE_OFFSET_LIMIT with PARAM_STRING",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_OFFSET_LIMIT,
						OffsetName:     "offset",
						OffsetType:     def.PARAM_STRING,
						LimitName:      "limit",
						LimitType:      def.PARAM_STRING,
						RespTotalName:  "total",
					},
				}).
					setFullLastPage(true),
			},
			10,
			false,
		},
		{
			"Valid PAGE_OFFSET_LIMIT with PARAM_STRING_LIST",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_OFFSET_LIMIT,
						OffsetName:     "offset",
						OffsetType:     def.PARAM_STRING_LIST,
						LimitName:      "limit",
						LimitType:      def.PARAM_STRING_LIST,
						RespTotalName:  "total",
					},
				}).
					setFullLastPage(true),
			},
			10,
			false,
		},
		{
			"Valid PAGE_OFFSET_LIMIT with fullLastPage==true and noTotalcount==true",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_OFFSET_LIMIT,
						OffsetName:     "offset",
						OffsetType:     def.PARAM_INT,
						LimitName:      "limit",
						LimitType:      def.PARAM_INT,
						RespTotalName:  "total",
					},
				}).
					setFullLastPage(true).
					setNoTotalCount(true),
			},
			10,
			false,
		},
		{
			"Valid PAGE_OFFSET_LIMIT with fullLastPage==false and noTotalcount==true",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_OFFSET_LIMIT,
						OffsetName:     "offset",
						OffsetType:     def.PARAM_INT,
						LimitName:      "limit",
						LimitType:      def.PARAM_INT,
						RespTotalName:  "total",
					},
				}).
					setNoTotalCount(true),
			},
			9,
			false,
		},
		{
			"Valid PAGE_CURPAGE_SIZE with fullLastPage==true",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_CURPAGE_SIZE,
						OffsetName:     "curpage",
						OffsetType:     def.PARAM_INT,
						LimitName:      "size",
						LimitType:      def.PARAM_INT,
						RespTotalName:  "total",
					},
				}).
					setFullLastPage(true),
			},
			10,
			false,
		},
		{
			"Valid PAGE_CURPAGE_SIZE with fullLastPage==false",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_CURPAGE_SIZE,
						OffsetName:     "curpage",
						OffsetType:     def.PARAM_INT,
						LimitName:      "size",
						LimitType:      def.PARAM_INT,
						RespTotalName:  "total",
					},
				}),
			},
			9,
			false,
		},
		{
			"Valid PAGE_CURPAGE_SIZE with fullLastPage==true and noTotalcount==true",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_CURPAGE_SIZE,
						OffsetName:     "curpage",
						OffsetType:     def.PARAM_INT,
						LimitName:      "size",
						LimitType:      def.PARAM_INT,
						RespTotalName:  "total",
					},
				}).
					setFullLastPage(true).
					setNoTotalCount(true),
			},
			10,
			false,
		},
		{
			"Valid PAGE_CURPAGE_SIZE with fullLastPage==false and noTotalcount==true",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_CURPAGE_SIZE,
						OffsetName:     "curpage",
						OffsetType:     def.PARAM_INT,
						LimitName:      "size",
						LimitType:      def.PARAM_INT,
						RespTotalName:  "total",
					},
				}).
					setNoTotalCount(true),
			},
			9,
			false,
		},
		{
			"Valid PAGE_NOPAGEINATION",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_NOPAGEINATION,
					},
				}),
			},
			8,
			false,
		},
		{
			"Valid PAGE_NOPAGEINATION with empty data",
			args{
				&menp,
			},
			0,
			false,
		},
		{
			"Valid PAGE_MARKER with fullLastPage==true",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_MARKER,
						LimitName:      "max",
						LimitType:      def.PARAM_INT,
						MarkerName:     "token",
						NextMarkerName: "nexttoken",
						RespTotalName:  "total",
					},
				}).
					setFullLastPage(true),
			},
			10,
			false,
		},
		{
			"Valid PAGE_MARKER with fullLastPage==false",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_MARKER,
						LimitName:      "max",
						LimitType:      def.PARAM_INT,
						MarkerName:     "token",
						NextMarkerName: "nexttoken",
						RespTotalName:  "total",
					},
				}),
			},
			9,
			false,
		},
		{
			"Valid PAGE_MARKER with fullLastPage==true and TruncatedName",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_MARKER,
						LimitName:      "max",
						LimitType:      def.PARAM_INT,
						MarkerName:     "token",
						NextMarkerName: "nexttoken",
						TruncatedName:  "isTrun",
						RespTotalName:  "total",
					},
				}).
					setFullLastPage(true),
			},
			10,
			false,
		},
		{
			"Valid PAGE_MARKER with fullLastPage==false and TruncatedName",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_MARKER,
						LimitName:      "max",
						LimitType:      def.PARAM_INT,
						MarkerName:     "token",
						NextMarkerName: "nexttoken",
						TruncatedName:  "isTrun",
						RespTotalName:  "total",
					},
				}),
			},
			9,
			false,
		},
		{
			"pagination type not set",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGEINATION_DEFAULT,
					},
				}),
			},
			0,
			true,
		},
		{
			"invalid PaginationType",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: -1,
					},
				}),
			},
			0,
			true,
		},
		////////////////
		{
			"valid with ProfileNotDefinedError",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_OFFSET_LIMIT,
						OffsetName:     "offset",
						OffsetType:     def.PARAM_INT,
						LimitName:      "limit",
						LimitType:      def.PARAM_INT,
						RespTotalName:  "total",
					},
				}).
					setFullLastPage(true).
					setFnCallback(func() ([]*json.RawMessage, NextCondition, error) {
						return nil, NextCondition{},
							fmt.Errorf("%w", auth.ProfileNotDefinedError{})
					}),
			},
			0,
			false,
		},
		{
			"valid with actual data less than total count",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_OFFSET_LIMIT,
						OffsetName:     "offset",
						OffsetType:     def.PARAM_INT,
						LimitName:      "limit",
						LimitType:      def.PARAM_INT,
						RespTotalName:  "total",
					},
				}).
					setFullLastPage(true).
					setFnCallback(func() ([]*json.RawMessage, NextCondition, error) {
						return make([]*json.RawMessage, 0),
							NextCondition{TotalCount: 100},
							nil
					}),
			},
			0,
			false,
		},
		{
			"valid with actual data more than total count",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_OFFSET_LIMIT,
						OffsetName:     "offset",
						OffsetType:     def.PARAM_INT,
						LimitName:      "limit",
						LimitType:      def.PARAM_INT,
						RespTotalName:  "total",
					},
				}).
					setFullLastPage(true).
					setFnCallback(func() ([]*json.RawMessage, NextCondition, error) {
						return make([]*json.RawMessage, 2),
							NextCondition{TotalCount: 1},
							nil
					}),
			},
			2,
			false,
		},
		{
			"failed to get data",
			args{
				(&mockPaginator{
					conf: def.ConfPaginator{
						PaginationType: def.PAGE_OFFSET_LIMIT,
						OffsetName:     "offset",
						OffsetType:     def.PARAM_INT,
						LimitName:      "limit",
						LimitType:      def.PARAM_INT,
						RespTotalName:  "total",
					},
				}).
					setFullLastPage(true).
					setFnCallback(func() ([]*json.RawMessage, NextCondition, error) {
						return nil, NextCondition{}, errors.New("mock error")
					}),
			},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEntireList(tt.args.p.GetPaginator(), tt.args.p.GetConf())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEntireList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !(len(got) == tt.wantSize) {
				t.Errorf("GetEntireList() size = %v, want %v", len(got), tt.wantSize)
			}
		})
	}
}

func TestResultDataParse(t *testing.T) {
	type args struct {
		v                any
		conf             def.ConfPaginator
		dataListJsonPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// Valid cases are under test of TestGetEntireList
		{
			"failed to unmarshal as map in PAGE_OFFSET_LIMIT",
			args{
				[]int{},
				def.ConfPaginator{PaginationType: def.PAGE_OFFSET_LIMIT},
				"",
			},
			true,
		},
		{
			"invalid response, missing key",
			args{
				make(map[string]any),
				def.ConfPaginator{
					PaginationType: def.PAGE_OFFSET_LIMIT,
					RespTotalName:  "total"},
				"",
			},
			true,
		},
		{
			"failed to convert to number",
			args{
				map[string]any{
					"total": "invalid",
				},
				def.ConfPaginator{
					PaginationType: def.PAGE_OFFSET_LIMIT,
					RespTotalName:  "total"},
				"",
			},
			true,
		},
		{
			"failed to convert to int",
			args{
				map[string]any{
					"total": 2.5,
				},
				def.ConfPaginator{
					PaginationType: def.PAGE_OFFSET_LIMIT,
					RespTotalName:  "total"},
				"",
			},
			true,
		},
		{
			"failed to convert to list in PAGE_OFFSET_LIMIT",
			args{
				map[string]any{
					"total": 1,
				},
				def.ConfPaginator{
					PaginationType: def.PAGE_OFFSET_LIMIT,
					RespTotalName:  "total"},
				"$.total",
			},
			true,
		},
		{
			"failed to convert to list in PAGE_NOPAGEINATION",
			args{
				map[string]any{
					"total": 1,
				},
				def.ConfPaginator{PaginationType: def.PAGE_NOPAGEINATION},
				"$.total",
			},
			true,
		},
		{
			"failed to unmarshal as map in PAGE_MARKER",
			args{
				[]int{},
				def.ConfPaginator{PaginationType: def.PAGE_MARKER},
				"",
			},
			true,
		},
		{
			"invalid response, missing key",
			args{
				make(map[string]any),
				def.ConfPaginator{
					PaginationType: def.PAGE_MARKER,
					TruncatedName:  "isTrun"},
				"",
			},
			true,
		},
		{
			"failed to convert to bool",
			args{
				map[string]any{
					"isTrun": "invalid",
				},
				def.ConfPaginator{
					PaginationType: def.PAGE_MARKER,
					TruncatedName:  "isTrun"},
				"",
			},
			true,
		},
		{
			"config of NextMarkerName is empty",
			args{
				make(map[string]any),
				def.ConfPaginator{PaginationType: def.PAGE_MARKER},
				"",
			},
			true,
		},
		{
			"failed to convert to string",
			args{
				map[string]any{
					"nexttoken": 1,
				},
				def.ConfPaginator{
					PaginationType: def.PAGE_MARKER,
					NextMarkerName: "nexttoken"},
				"",
			},
			true,
		},
		{
			"failed to convert to list in PAGE_MARKER",
			args{
				map[string]any{
					"total": 1,
				},
				def.ConfPaginator{
					PaginationType: def.PAGE_MARKER,
					NextMarkerName: "nexttoken"},
				"$.total",
			},
			true,
		},
		{
			"failed to deal with PaginationType",
			args{
				make(map[string]any),
				def.ConfPaginator{PaginationType: def.PAGEINATION_DEFAULT},
				"",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, _ := internal.JsonMarshal(tt.args.v)
			_, _, err := ResultDataParse(rm, tt.args.conf, tt.args.dataListJsonPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResultDataParse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
