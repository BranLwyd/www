package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/BranLwyd/acnh_flowers/flower"
	"github.com/BranLwyd/www/assets"
)

var (
	colorMap = map[string]string{
		"White":         "W",
		"White (seed)":  "W",
		"Pink":          "P",
		"Red":           "R",
		"Red (seed)":    "R",
		"Orange":        "O",
		"Yellow":        "Y",
		"Yellow (seed)": "Y",
		"Green":         "G",
		"Blue":          "B",
		"Purple":        "U",
		"Black":         "K",
	}

	flowerTmpl = template.Must(template.New("flower-tmpl").Funcs(template.FuncMap{
		"pic": func(species, phenotype string) string {
			s := strings.ToUpper(species[0:1])
			c := colorMap[phenotype]
			return fmt.Sprintf("/img/%s%s.png", s, c)
		},
	}).Parse(string(assets.Asset["assets/flowers.html"])))
)

type flowerHandler struct{}

func (flowerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	speciesParam := r.Form.Get("species")
	firstGenotypeParam := r.Form.Get("firstGenotype")
	secondGenotypeParam := r.Form.Get("secondGenotype")

	type result struct {
		Odds      string
		Genotype  string
		Phenotype string
	}

	var resultString string
	var results []result
	var errStr string
	if speciesParam != "" && firstGenotypeParam != "" && secondGenotypeParam != "" {
		s, ok := species[speciesParam]
		if !ok {
			errStr = fmt.Sprintf("Unknown species %q", speciesParam)
			goto writeResult
		}

		gs, err := flower.NewGenotypeSerdeFromExampleDistribution(firstGenotypeParam)
		if err != nil {
			errStr = fmt.Sprintf("Couldn't parse first genotype: %v", err)
			goto writeResult
		}
		if gs.GeneCount() != s.GeneCount() {
			errStr = fmt.Sprintf("First genotype has wrong number of genes (%d, wanted %d)", gs.GeneCount(), s.GeneCount())
			goto writeResult
		}

		gda, err := gs.ParseGeneticDistribution(firstGenotypeParam)
		if err != nil {
			errStr = fmt.Sprintf("Couldn't parse first genotype: %v", err)
			goto writeResult
		}
		gdb, err := gs.ParseGeneticDistribution(secondGenotypeParam)
		if err != nil {
			errStr = fmt.Sprintf("Couldn't parse second genotype: %v", err)
			goto writeResult
		}

		gdRslt := gda.Breed(gdb)
		resultString = gs.RenderGeneticDistribution(gdRslt)

		var totalOdds uint64
		gdRslt.Visit(func(g flower.Genotype, p uint64) { totalOdds += p })
		gdRslt.Visit(func(g flower.Genotype, p uint64) {
			gcd := gcd(p, totalOdds)
			num, den := p/gcd, totalOdds/gcd

			results = append(results, result{
				Odds:      fmt.Sprintf("%.02f%% (%d/%d)", 100*float64(p)/float64(totalOdds), num, den),
				Genotype:  gs.RenderGenotype(g),
				Phenotype: renderPhenotype(s, g),
			})
		})
	}

writeResult:
	var buf bytes.Buffer
	if err := flowerTmpl.Execute(&buf, struct {
		Species        string
		FirstGenotype  string
		SecondGenotype string
		ResultString   string
		Results        []result
		Error          string
	}{
		Species:        speciesParam,
		FirstGenotype:  firstGenotypeParam,
		SecondGenotype: secondGenotypeParam,
		ResultString:   resultString,
		Results:        results,
		Error:          errStr,
	}); err != nil {
		log.Printf("Could not server flower-breeding page: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(buf.Bytes()))
}

func renderPhenotype(s flower.Species, g flower.Genotype) string {
	if _, ok := seedGenotypes[s.Name()][g]; ok {
		return fmt.Sprintf("%s (seed)", s.Phenotype(g))
	}
	return s.Phenotype(g)
}

var (
	species = map[string]flower.Species{
		"cosmos":      flower.Cosmos(),
		"hyacinths":   flower.Hyacinths(),
		"lilies":      flower.Lilies(),
		"mums":        flower.Mums(),
		"pansies":     flower.Pansies(),
		"roses":       flower.Roses(),
		"tulips":      flower.Tulips(),
		"windflowers": flower.Windflowers(),
	}

	seedGenotypes map[string]map[flower.Genotype]struct{}
)

func init() {
	seedGenotypes = map[string]map[flower.Genotype]struct{}{}
	for _, x := range []struct {
		species       flower.Species
		seedGenotypes []string
	}{
		{flower.Cosmos(), []string{"rryySs", "rrYYSs", "RRyyss"}},
	} {
		m := map[flower.Genotype]struct{}{}
		for _, sg := range x.seedGenotypes {
			g, err := x.species.ParseGenotype(sg)
			if err != nil {
				panic(fmt.Sprintf("Couldn't parse seed genotype %q for species %s", sg, x.species.Name()))
			}
			m[g] = struct{}{}
		}
		seedGenotypes[x.species.Name()] = m
	}
}

// Based on https://en.wikipedia.org/wiki/Binary_GCD_algorithm#Iterative_version_in_C.
func gcd(u, v uint64) uint64 {
	// Base cases.
	if u == 0 {
		return v
	}
	if v == 0 {
		return u
	}

	// Remove largest factor of 2.
	shift := 0
	for (u|v)&1 == 0 {
		shift++
		u >>= 1
		v >>= 1
	}

	// Remove additional, non-common factors of 2 from u.
	for u&1 == 0 {
		u >>= 1
	}

	// Loop invariant: u is odd.
	for v != 0 {
		for v&1 == 0 {
			v >>= 1
		}
		if u > v {
			u, v = v, u
		}
		v -= u
	}
	return u << shift
}
