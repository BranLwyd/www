package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/BranLwyd/acnh_flowers/flower"
	"github.com/BranLwyd/www/assets"
)

var (
	flowerTmpl = template.Must(template.New("flower-tmpl").Parse(string(assets.Asset["assets/flowers.html"])))
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

	var results []result
	var errStr string
	if speciesParam != "" && firstGenotypeParam != "" && secondGenotypeParam != "" {
		s, ok := species[speciesParam]
		if !ok {
			errStr = fmt.Sprintf("Unknown species %q", speciesParam)
			goto writeResult
		}

		gs, err := flower.NewGenotypeSerdeFromExample(firstGenotypeParam)
		if err != nil {
			errStr = fmt.Sprintf("Couldn't parse first genotype: %v", err)
			goto writeResult
		}
		if gs.GeneCount() != s.GeneCount() {
			errStr = fmt.Sprintf("First genotype has wrong number of genes (%d, wanted %d)", gs.GeneCount(), s.GeneCount())
			goto writeResult
		}

		ga, err := gs.ParseGenotype(firstGenotypeParam)
		if err != nil {
			errStr = fmt.Sprintf("Couldn't parse first genotype: %v", err)
			goto writeResult
		}
		gb, err := gs.ParseGenotype(secondGenotypeParam)
		if err != nil {
			errStr = fmt.Sprintf("Couldn't parse second genotype: %v", err)
			goto writeResult
		}

		gda := ga.ToGeneticDistribution()
		gdb := gb.ToGeneticDistribution()
		gdRslt := gda.Breed(gdb)

		var totalOdds uint64
		for _, p := range gdRslt {
			totalOdds += p
		}
		for g, p := range gdRslt {
			if p == 0 {
				continue
			}
			g := flower.Genotype(g)

			gcd := gcd(p, totalOdds)
			num, den := p/gcd, totalOdds/gcd

			results = append(results, result{
				Odds:      fmt.Sprintf("%.02f%% (%d/%d)", 100*float64(p)/float64(totalOdds), num, den),
				Genotype:  gs.RenderGenotype(g),
				Phenotype: s.Phenotype(g),
			})
		}
	}

writeResult:
	var buf bytes.Buffer
	if err := flowerTmpl.Execute(&buf, struct {
		Species        string
		FirstGenotype  string
		SecondGenotype string
		Results        []result
		Error          string
	}{
		Species:        speciesParam,
		FirstGenotype:  firstGenotypeParam,
		SecondGenotype: secondGenotypeParam,
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
)

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
