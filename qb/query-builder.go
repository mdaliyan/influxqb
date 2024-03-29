package qb

import (
	"errors"
	"strings"

	"github.com/influxdata/influxdb/models"
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/mdaliyan/influxqb/inflx"
)

func NewHistogram(database string) *HistogramBuilder {
	h := HistogramBuilder{}
	h.database = database
	h.From(database)
	h.summaries = map[string]string{}
	h.dataSets = map[string]string{}
	return &h
}

type HistogramBuilder struct {
	database    string
	timeRange   string
	fill        string
	groupBy     string
	total       bool
	where       []string
	sum         models.Row
	dataSets    map[string]string
	summaries   map[string]string
	Response    *Indexer
	RawResponse *client.Response
}

func (h *HistogramBuilder) Total(enabled bool) (r map[string]interface{}) {
	h.total = enabled
	return
}

func (h *HistogramBuilder) Export() (r Response) {
	r = Response{
		Summary:  h.Response.GetSummary(),
		DataSets: DataSets{},
	}
	for key, _ := range h.dataSets {
		r.DataSets[key] = h.Response.GetTimeSeriesFor(key)
	}
	return
}

func (h *HistogramBuilder) Do(db string) (err error) {
	var r []client.Result
	h.RawResponse, r, err = inflx.Query(db, h.Query())
	if err == nil {
		if r[0].Err != "" {
			err = errors.New(r[0].Err)
			return
		}
		h.Response = NewHistogramData(r)
		return
	}
	return
}

func (h *HistogramBuilder) Query() string {
	var queries []string
	for key, field := range h.dataSets {
		queries = append(queries, h.buildQuery(map[string]string{key: field}, h.groupBy))
	}
	if len(h.summaries) > 0 {
		queries = append(queries, h.buildQuery(h.summaries, ""))
	}
	return strings.Join(queries, ";\n")
}

func (h *HistogramBuilder) buildQuery(set map[string]string, groupBy string) string {
	var selects []string
	for key, filed := range set {
		sel := filed
		if key != "" {
			sel += " as " + key
		}
		selects = append(selects, sel)
	}
	q := `select ` + strings.Join(selects, ", ") + ` from ` + h.database
	if h.where != nil {
		q += ` where ` + strings.Join(h.where, " and ")
	}
	if groupBy != "" {
		q += ` group by ` + groupBy
	}
	if h.fill != "" {
		q += ` fill(` + h.fill + ")"
	}
	return q
}

func (h *HistogramBuilder) From(database string) *HistogramBuilder {
	if h.total {
		h.database = "sanjagh.total_" + database
	} else {
		h.database = "sanjagh." + database
	}
	// todo : remove this bottom line
	h.database = "sanjagh." + database
	return h
}

func (h *HistogramBuilder) DataSet(as, field string) *HistogramBuilder {
	h.dataSets[as] = field
	return h
}

func (h *HistogramBuilder) Summary(as, field string) *HistogramBuilder {
	h.summaries[as] = field
	return h
}

func (h *HistogramBuilder) GroupBy(s string) *HistogramBuilder {
	h.groupBy = " " + s + " "
	return h
}

func (h *HistogramBuilder) GroupMinutely() *HistogramBuilder {
	h.From("minutely.statistics")
	h.groupBy = " time(60s) "
	return h
}

func (h *HistogramBuilder) GroupHourly() *HistogramBuilder {
	h.From("hourly.statistics")
	h.groupBy = " time(1h) "
	return h
}

func (h *HistogramBuilder) GroupDaily() *HistogramBuilder {
	h.From("daily.statistics")
	h.groupBy = " time(1d) "
	return h
}

func (h *HistogramBuilder) GroupMonthly() *HistogramBuilder {
	h.From("daily.statistics")
	h.groupBy = " time(30d) "
	return h
}

func (h *HistogramBuilder) GroupYearly() *HistogramBuilder {
	h.From("daily.statistics")
	h.groupBy = " time(365d) "
	return h
}

func (h *HistogramBuilder) Fill(s string) *HistogramBuilder {
	h.fill = s
	return h
}

func (h *HistogramBuilder) Where(s string) *HistogramBuilder {
	h.where = append(h.where, " "+s+" ")
	return h
}
