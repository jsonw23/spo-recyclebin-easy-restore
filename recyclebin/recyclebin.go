package recyclebin

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/koltyakov/gosip"
	strategy "github.com/koltyakov/gosip-sandbox/strategies/ondemand"
	"github.com/koltyakov/gosip/api"
	"github.com/spf13/viper"
)

type Query struct {
	DeletedBefore time.Time
	DeletedAfter  time.Time
	DeletedBy     string
	ItemName      string
}

type Restore struct {
	Query *Query
}

func Authenticate() *api.SP {
	authCfg := &strategy.AuthCnfg{
		SiteURL: viper.GetString("siteUrl"),
	}

	client := &gosip.SPClient{AuthCnfg: authCfg}
	return api.NewSP(client)
}

func NewQuery(args []string) *Query {
	q := Query{
		DeletedBefore: viper.GetTime("before"),
		DeletedAfter:  viper.GetTime("after"),
		DeletedBy:     viper.GetString("by"),
	}
	if len(args) > 0 {
		q.ItemName = args[0]
	}
	return &q
}

func NewRestore(query *Query) *Restore {
	return &Restore{
		Query: query,
	}
}

func (q *Query) Results() *api.RecycleBinResp {
	sp := Authenticate()

	// assemble an OData filter statement that the SharePoint REST API will accept
	filters := []string{}
	if !q.DeletedBefore.IsZero() {
		filters = append(filters, fmt.Sprintf("DeletedDate le datetime'%s'", q.DeletedBefore.UTC().Format(time.RFC3339)))
	}
	if !q.DeletedAfter.IsZero() {
		filters = append(filters, fmt.Sprintf("DeletedDate ge datetime'%s'", q.DeletedAfter.UTC().Format(time.RFC3339)))
	}
	if q.DeletedBy != "" {
		filters = append(filters, fmt.Sprintf("DeletedByName eq '%s'", q.DeletedBy))
	}

	if q.ItemName != "" {
		filters = append(filters, fmt.Sprintf("substringof('%s', Title)", q.ItemName))
	}

	recycleBin := sp.Site().RecycleBin()
	if len(filters) > 0 {
		recycleBin = recycleBin.Filter(strings.Join(filters, " and "))
	}

	data, err := recycleBin.
		OrderBy("DeletedDate", false).
		Get()
	if err != nil {
		log.Fatal(err)
	}
	return &data
}

func (r *Restore) Run() {
	sp := Authenticate()
	for _, item := range r.Query.Results().Data() {
		d := item.Data()
		replacer := strings.NewReplacer("'", "''", "#", "%23", "%", "%25")
		odataReady := strings.Builder{}
		odataReady.WriteString("/")
		odataReady.WriteString(replacer.Replace(d.DirNamePath.DecodedURL))
		odataReady.WriteString("/")
		odataReady.WriteString(replacer.Replace(d.LeafNamePath.DecodedURL))
		file, err := sp.Web().GetFile(odataReady.String()).Get()
		if err != nil && !strings.Contains(err.Error(), "404 Not Found") {
			log.Println(err)
		}
		fileInfo := file.Data()
		if !fileInfo.Exists {
			log.Printf("restoring: /%s/%s", d.DirName, d.LeafName)
			if err := sp.Site().RecycleBin().GetByID(d.ID).Restore(); err != nil {
				log.Println(err)
			}
		}
	}
}
