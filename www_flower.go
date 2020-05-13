package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

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

	var results []string
	var errStr string
	if speciesParam != "" && firstGenotypeParam != "" && secondGenotypeParam != "" {
		s, ok := species[speciesParam]
		if !ok {
			errStr = fmt.Sprintf("Unknown species %q", speciesParam)
			goto writeResult
		}

		ga, err := s.ParseGenotype(firstGenotypeParam)
		if err != nil {
			errStr = fmt.Sprintf("Couldn't parse first genotype: %v", err)
			goto writeResult
		}
		gb, err := s.ParseGenotype(secondGenotypeParam)
		if err != nil {
			errStr = fmt.Sprintf("Couldn't parse second genotype: %v", err)
			goto writeResult
		}

		gda := s.ToGeneticDistribution(ga)
		gdb := s.ToGeneticDistribution(gb)
		gdRslt := gda.Breed(gdb)

		for g, p := range gdRslt.dist {
			if p == 0 {
				continue
			}
			g := Genotype(g)
			results = append(results, fmt.Sprintf("%d: %s (%s)", p, s.RenderGenotype(g), s.phenotypes[g]))
		}
	}

writeResult:
	var buf bytes.Buffer
	if err := flowerTmpl.Execute(&buf, struct {
		Species        string
		FirstGenotype  string
		SecondGenotype string
		Results        []string
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
	species = map[string]Species{
		"cosmos": MustSpecies("Mums", 3, map[string]string{
			"rryyww": "White",
			"rryyWw": "White",
			"rryyWW": "White",
			"rrYyww": "Yellow",
			"rrYyWw": "Yellow",
			"rrYyWW": "White",
			"rrYYww": "Yellow",
			"rrYYWw": "Yellow",
			"rrYYWW": "Yellow",
			"Rryyww": "Pink",
			"RryyWw": "Pink",
			"RryyWW": "Pink",
			"RrYyww": "Orange",
			"RrYyWw": "Orange",
			"RrYyWW": "Pink",
			"RrYYww": "Orange",
			"RrYYWw": "Orange",
			"RrYYWW": "Orange",
			"RRyyww": "Red",
			"RRyyWw": "Red",
			"RRyyWW": "Red",
			"RRYyww": "Orange",
			"RRYyWw": "Orange",
			"RRYyWW": "Red",
			"RRYYww": "Black",
			"RRYYWw": "Black",
			"RRYYWW": "Red",
		}),

		"hyacinths": MustSpecies("Mums", 3, map[string]string{
			"rryyww": "White",
			"rryyWw": "White",
			"rryyWW": "Blue",
			"rrYyww": "Yellow",
			"rrYyWw": "Yellow",
			"rrYyWW": "White",
			"rrYYww": "Yellow",
			"rrYYWw": "Yellow",
			"rrYYWW": "Yellow",
			"Rryyww": "Red",
			"RryyWw": "Pink",
			"RryyWW": "White",
			"RrYyww": "Orange",
			"RrYyWw": "Yellow",
			"RrYyWW": "Yellow",
			"RrYYww": "Orange",
			"RrYYWw": "Yellow",
			"RrYYWW": "Yellow",
			"RRyyww": "Red",
			"RRyyWw": "Red",
			"RRyyWW": "Red",
			"RRYyww": "Blue",
			"RRYyWw": "Red",
			"RRYyWW": "Red",
			"RRYYww": "Purple",
			"RRYYWw": "Purple",
			"RRYYWW": "Purple",
		}),

		"lilies": MustSpecies("Mums", 3, map[string]string{
			"rryyww": "White",
			"rryyWw": "White",
			"rryyWW": "White",
			"rrYyww": "Yellow",
			"rrYyWw": "White",
			"rrYyWW": "White",
			"rrYYww": "Yellow",
			"rrYYWw": "Yellow",
			"rrYYWW": "White",
			"Rryyww": "Red",
			"RryyWw": "Pink",
			"RryyWW": "White",
			"RrYyww": "Orange",
			"RrYyWw": "Yellow",
			"RrYyWW": "Yellow",
			"RrYYww": "Orange",
			"RrYYWw": "Yellow",
			"RrYYWW": "Yellow",
			"RRyyww": "Black",
			"RRyyWw": "Red",
			"RRyyWW": "Pink",
			"RRYyww": "Black",
			"RRYyWw": "Red",
			"RRYyWW": "Pink",
			"RRYYww": "Orange",
			"RRYYWw": "Orange",
			"RRYYWW": "White",
		}),

		"mums": MustSpecies("Mums", 3, map[string]string{
			"rryyww": "White",
			"rryyWw": "White",
			"rryyWW": "Purple",
			"rrYyww": "Yellow",
			"rrYyWw": "Yellow",
			"rrYyWW": "White",
			"rrYYww": "Yellow",
			"rrYYWw": "Yellow",
			"rrYYWW": "Yellow",
			"Rryyww": "Pink",
			"RryyWw": "Pink",
			"RryyWW": "Pink",
			"RrYyww": "Yellow",
			"RrYyWw": "Red",
			"RrYyWW": "Pink",
			"RrYYww": "Purple",
			"RrYYWw": "Purple",
			"RrYYWW": "Purple",
			"RRyyww": "Red",
			"RRyyWw": "Red",
			"RRyyWW": "Red",
			"RRYyww": "Purple",
			"RRYyWw": "Purple",
			"RRYyWW": "Red",
			"RRYYww": "Green",
			"RRYYWw": "Green",
			"RRYYWW": "Red",
		}),

		"pansies": MustSpecies("Mums", 3, map[string]string{
			"rryyww": "White",
			"rryyWw": "White",
			"rryyWW": "Blue",
			"rrYyww": "Yellow",
			"rrYyWw": "Yellow",
			"rrYyWW": "Blue",
			"rrYYww": "Yellow",
			"rrYYWw": "Yellow",
			"rrYYWW": "Yellow",
			"Rryyww": "Red",
			"RryyWw": "Red",
			"RryyWW": "Blue",
			"RrYyww": "Orange",
			"RrYyWw": "Orange",
			"RrYyWW": "Orange",
			"RrYYww": "Yellow",
			"RrYYWw": "Yellow",
			"RrYYWW": "Yellow",
			"RRyyww": "Red",
			"RRyyWw": "Red",
			"RRyyWW": "Purple",
			"RRYyww": "Red",
			"RRYyWw": "Red",
			"RRYyWW": "Purple",
			"RRYYww": "Orange",
			"RRYYWw": "Orange",
			"RRYYWW": "Purple",
		}),

		"roses": MustSpecies("Roses", 4, map[string]string{
			"rryywwbb": "White",
			"rryywwBb": "White",
			"rryywwBB": "White",
			"rryyWwbb": "White",
			"rryyWwBb": "White",
			"rryyWwBB": "White",
			"rryyWWbb": "Purple",
			"rryyWWBb": "Purple",
			"rryyWWBB": "Purple",
			"rrYywwbb": "Yellow",
			"rrYywwBb": "Yellow",
			"rrYywwBB": "Yellow",
			"rrYyWwbb": "White",
			"rrYyWwBb": "White",
			"rrYyWwBB": "White",
			"rrYyWWbb": "Purple",
			"rrYyWWBb": "Purple",
			"rrYyWWBB": "Purple",
			"rrYYwwbb": "Yellow",
			"rrYYwwBb": "Yellow",
			"rrYYwwBB": "Yellow",
			"rrYYWwbb": "Yellow",
			"rrYYWwBb": "Yellow",
			"rrYYWwBB": "Yellow",
			"rrYYWWbb": "White",
			"rrYYWWBb": "White",
			"rrYYWWBB": "White",
			"Rryywwbb": "Red",
			"RryywwBb": "Pink",
			"RryywwBB": "White",
			"RryyWwbb": "Red",
			"RryyWwBb": "Pink",
			"RryyWwBB": "White",
			"RryyWWbb": "Red",
			"RryyWWBb": "Pink",
			"RryyWWBB": "Purple",
			"RrYywwbb": "Orange",
			"RrYywwBb": "Yellow",
			"RrYywwBB": "Yellow",
			"RrYyWwbb": "Red",
			"RrYyWwBb": "Pink",
			"RrYyWwBB": "White",
			"RrYyWWbb": "Red",
			"RrYyWWBb": "Pink",
			"RrYyWWBB": "Purple",
			"RrYYwwbb": "Orange",
			"RrYYwwBb": "Yellow",
			"RrYYwwBB": "Yellow",
			"RrYYWwbb": "Orange",
			"RrYYWwBb": "Yellow",
			"RrYYWwBB": "Yellow",
			"RrYYWWbb": "Red",
			"RrYYWWBb": "Pink",
			"RrYYWWBB": "White",
			"RRyywwbb": "Black",
			"RRyywwBb": "Red",
			"RRyywwBB": "Pink",
			"RRyyWwbb": "Black",
			"RRyyWwBb": "Red",
			"RRyyWwBB": "Pink",
			"RRyyWWbb": "Black",
			"RRyyWWBb": "Red",
			"RRyyWWBB": "Pink",
			"RRYywwbb": "Orange",
			"RRYywwBb": "Orange",
			"RRYywwBB": "Yellow",
			"RRYyWwbb": "Red",
			"RRYyWwBb": "Red",
			"RRYyWwBB": "White",
			"RRYyWWbb": "Black",
			"RRYyWWBb": "Red",
			"RRYyWWBB": "Purple",
			"RRYYwwbb": "Orange",
			"RRYYwwBb": "Orange",
			"RRYYwwBB": "Yellow",
			"RRYYWwbb": "Orange",
			"RRYYWwBb": "Orange",
			"RRYYWwBB": "Yellow",
			"RRYYWWbb": "Blue",
			"RRYYWWBb": "Red",
			"RRYYWWBB": "White",
		}),

		"tulips": MustSpecies("Tulips", 3, map[string]string{
			"rryyww": "White",
			"rryyWw": "White",
			"rryyWW": "White",
			"rrYyww": "Yellow",
			"rrYyWw": "Yellow",
			"rrYyWW": "White",
			"rrYYww": "Yellow",
			"rrYYWw": "Yellow",
			"rrYYWW": "Yellow",
			"Rryyww": "Red",
			"RryyWw": "Pink",
			"RryyWW": "White",
			"RrYyww": "Orange",
			"RrYyWw": "Yellow",
			"RrYyWW": "Yellow",
			"RrYYww": "Orange",
			"RrYYWw": "Yellow",
			"RrYYWW": "Yellow",
			"RRyyww": "Black",
			"RRyyWw": "Red",
			"RRyyWW": "Red",
			"RRYyww": "Black",
			"RRYyWw": "Red",
			"RRYyWW": "Red",
			"RRYYww": "Purple",
			"RRYYWw": "Purple",
			"RRYYWW": "Purple",
		}),

		"windflowers": MustSpecies("Windflowers", 3, map[string]string{
			"rryyww": "White",
			"rryyWw": "White",
			"rryyWW": "Blue",
			"rrYyww": "Orange",
			"rrYyWw": "Orange",
			"rrYyWW": "Blue",
			"rrYYww": "Orange",
			"rrYYWw": "Orange",
			"rrYYWW": "Orange",
			"Rryyww": "Red",
			"RryyWw": "Red",
			"RryyWW": "Blue",
			"RrYyww": "Pink",
			"RrYyWw": "Pink",
			"RrYyWW": "Pink",
			"RrYYww": "Orange",
			"RrYYWw": "Orange",
			"RrYYWW": "Orange",
			"RRyyww": "Red",
			"RRyyWw": "Red",
			"RRyyWW": "Purple",
			"RRYyww": "Red",
			"RRYyWw": "Red",
			"RRYyWW": "Purple",
			"RRYYww": "Pink",
			"RRYYWw": "Pink",
			"RRYYWW": "Purple",
		}),
	}

	gene0 = [3]string{"rr", "Rr", "RR"}
	gene1 = [3]string{"yy", "Yy", "YY"}
	gene2 = [3]string{"ww", "Ww", "WW"}
	gene3 = [3]string{"bb", "Bb", "BB"}
)

// Species represents a specific species of flower, such as Windflower or Mum.
type Species struct {
	name       string              // a human-readable name for this species, e.g. "Windflower".
	geneCount  int                 // how many genes this flower has; assumed to be 3 or 4.
	phenotypes map[Genotype]string // phenotypes by genotype
}

func NewSpecies(name string, geneCount int, phenotypes map[string]string) (Species, error) {
	s := Species{name: name, geneCount: geneCount}
	if geneCount != 3 && geneCount != 4 {
		return Species{}, fmt.Errorf("Gene count is %d, must be 3 or 4", geneCount)
	}

	if geneCount == 3 && len(phenotypes) != 27 {
		return Species{}, fmt.Errorf("Got %d phenotypes, expected 27", len(phenotypes))
	}
	if geneCount == 4 && len(phenotypes) != 81 {
		return Species{}, fmt.Errorf("Got %d phenotypes, expected 81", len(phenotypes))
	}

	pts := map[Genotype]string{}
	for g, p := range phenotypes {
		g, err := s.ParseGenotype(g)
		if err != nil {
			return Species{}, err
		}
		pts[g] = p
	}
	s.phenotypes = pts
	return s, nil
}

func MustSpecies(name string, geneCount int, phenotypes map[string]string) Species {
	s, err := NewSpecies(name, geneCount, phenotypes)
	if err != nil {
		panic(err)
	}
	return s
}

func (s Species) ParseGenotype(genotype string) (Genotype, error) {
	var rslt Genotype

	if s.geneCount == 3 && len(genotype) != 6 {
		return 0, fmt.Errorf("genotype %q has wrong length (expected 6)", genotype)
	}
	if s.geneCount == 4 && len(genotype) != 8 {
		return 0, fmt.Errorf("genotype %q has wrong length (expected 8)", genotype)
	}

	for _, x := range []struct {
		gene   [3]string
		offset uint
	}{
		{gene0, 0},
		{gene1, 2},
		{gene2, 4},
		{gene3, 6},
	} {
		if len(genotype) == 6 && x.gene == gene3 {
			break
		}

		found := false
		for i, v := range x.gene {
			if v == genotype[x.offset:x.offset+2] {
				rslt |= Genotype(i << x.offset)
				found = true
				break
			}
		}
		if !found {
			return 0, fmt.Errorf("unparsable gene %q", genotype[x.offset:x.offset+2])
		}
	}
	return rslt, nil
}

func (s Species) RenderGenotype(g Genotype) string {
	switch s.geneCount {
	case 3:
		return fmt.Sprintf("%s%s%s", gene0[g.gene0()], gene1[g.gene1()], gene2[g.gene2()])
	case 4:
		return fmt.Sprintf("%s%s%s%s", gene0[g.gene0()], gene1[g.gene1()], gene2[g.gene2()], gene3[g.gene3()])
	default:
		panic(fmt.Sprintf("Gene count is %d, expect 3 or 4", s.geneCount))
	}
}

func (s Species) ToGeneticDistribution(g Genotype) GeneticDistribution {
	var gd GeneticDistribution
	switch s.geneCount {
	case 3:
		gd.dist = make([]uint64, threeGeneGenotypeCount)
	case 4:
		gd.dist = make([]uint64, fourGeneGenotypeCount)
	}
	gd.dist[g] = 1
	return gd
}

// Genotype represents a specific set of genes for a species, e.g. RrwwYY.
type Genotype uint8

// Internally, each two consecutive bits of a Genotype value represents a gene.
//  0 == 0b00 is dual-recessive (rr).
//  1 == 0b01 is dominant/recessive (Rr).
//  2 == 0b10 is dual-domninant (RR).
//  3 == 0b11 is unused.

func (g Genotype) gene0() uint8 { return uint8((g >> 0) & 0b11) }
func (g Genotype) gene1() uint8 { return uint8((g >> 2) & 0b11) }
func (g Genotype) gene2() uint8 { return uint8((g >> 4) & 0b11) }
func (g Genotype) gene3() uint8 { return uint8((g >> 6) & 0b11) }

const (
	threeGeneGenotypeCount = 64
	fourGeneGenotypeCount  = 256
)

// GeneticDistribution represents a probability distribution over all possible genotypes.
type GeneticDistribution struct {
	dist []uint64
}

func (gda GeneticDistribution) Breed(gdb GeneticDistribution) GeneticDistribution {
	rslt := GeneticDistribution{dist: make([]uint64, len(gda.dist))}

	if len(gda.dist) != len(gdb.dist) {
		panic(fmt.Sprintf("Mismatched genetic distributions (len %d != %d)", len(gda.dist), len(gdb.dist)))
	}
	var breedInto func(*GeneticDistribution, uint64, Genotype, Genotype)
	switch len(gda.dist) {
	case threeGeneGenotypeCount:
		breedInto = breedInto3
	case fourGeneGenotypeCount:
		breedInto = breedInto4
	default:
		panic(fmt.Sprintf("Bad genetic distribution length (%d, want 64 or 256)", len(gda.dist)))
	}

	// Breed each pair of possible genotypes into the result.
	for ga, pa := range gda.dist {
		if pa == 0 {
			continue
		}
		ga := Genotype(ga)
		for gb, pb := range gdb.dist {
			if pb == 0 {
				continue
			}
			gb := Genotype(gb)
			breedInto(&rslt, pa*pb, ga, gb)
		}
	}

	// Reduce the result.
	g := rslt.dist[0]
	for _, x := range rslt.dist[1:] {
		if g == 1 {
			break
		}
		g = gcd(g, x)
	}
	for i := range rslt.dist {
		rslt.dist[i] /= g
	}
	return rslt
}

func breedInto3(gd *GeneticDistribution, weight uint64, ga, gb Genotype) {
	wt0 := punnetSquareLookupTable[ga.gene0()][gb.gene0()]
	wt1 := punnetSquareLookupTable[ga.gene1()][gb.gene1()]
	wt2 := punnetSquareLookupTable[ga.gene2()][gb.gene2()]

	for g0, w0 := range wt0 {
		for g1, w1 := range wt1 {
			for g2, w2 := range wt2 {
				gd.dist[g0|(g1<<2)|(g2<<4)] += weight * w0 * w1 * w2
			}
		}
	}
}

func breedInto4(gd *GeneticDistribution, weight uint64, ga, gb Genotype) {
	wt0 := punnetSquareLookupTable[ga.gene0()][gb.gene0()]
	wt1 := punnetSquareLookupTable[ga.gene1()][gb.gene1()]
	wt2 := punnetSquareLookupTable[ga.gene2()][gb.gene2()]
	wt3 := punnetSquareLookupTable[ga.gene3()][gb.gene3()]

	for g0, w0 := range wt0 {
		for g1, w1 := range wt1 {
			for g2, w2 := range wt2 {
				for g3, w3 := range wt3 {
					gd.dist[g0|(g1<<2)|(g2<<4)|(g3<<6)] += weight * w0 * w1 * w2 * w3
				}
			}
		}
	}
}

var (
	// TODO: generate this lookup table from code, to decrease odds of error
	punnetSquareLookupTable = [3][3][3]uint64{
		// ga == 0 (rr)
		[3][3]uint64{
			[3]uint64{4, 0, 0},
			[3]uint64{2, 2, 0},
			[3]uint64{0, 4, 0},
		},

		// ga = 1 (Rr)
		[3][3]uint64{
			[3]uint64{2, 2, 0},
			[3]uint64{1, 2, 1},
			[3]uint64{0, 2, 2},
		},

		// ga = 2 (RR)
		[3][3]uint64{
			[3]uint64{0, 4, 0},
			[3]uint64{0, 2, 2},
			[3]uint64{0, 0, 4},
		},
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
