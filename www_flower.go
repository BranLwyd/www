package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
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

		gs, err := NewGenotypeSerdeFromExample(firstGenotypeParam)
		if err != nil {
			errStr = fmt.Sprintf("Couldn't parse first genotype: %v", err)
			goto writeResult
		}
		if gs.GeneCount() != s.geneCount {
			errStr = fmt.Sprintf("First genotype has wrong number of genes (%d, wanted %d)", gs.GeneCount(), s.geneCount)
			goto writeResult
		}

		ga, err := gs.DeserializeGenotype(firstGenotypeParam)
		if err != nil {
			errStr = fmt.Sprintf("Couldn't parse first genotype: %v", err)
			goto writeResult
		}
		gb, err := gs.DeserializeGenotype(secondGenotypeParam)
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
			results = append(results, fmt.Sprintf("%d: %s (%s)", p, gs.SerializeGenotype(g), s.phenotypes[g]))
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
		"cosmos": MustSpecies("Cosmos", map[string]string{
			"rryyss": "White",
			"rryySs": "White (seed)",
			"rryySS": "White",
			"rrYyss": "Yellow",
			"rrYySs": "Yellow",
			"rrYySS": "White",
			"rrYYss": "Yellow",
			"rrYYSs": "Yellow (seed)",
			"rrYYSS": "Yellow",
			"Rryyss": "Pink",
			"RryySs": "Pink",
			"RryySS": "Pink",
			"RrYyss": "Orange",
			"RrYySs": "Orange",
			"RrYySS": "Pink",
			"RrYYss": "Orange",
			"RrYYSs": "Orange",
			"RrYYSS": "Orange",
			"RRyyss": "Red (seed)",
			"RRyySs": "Red",
			"RRyySS": "Red",
			"RRYyss": "Orange",
			"RRYySs": "Orange",
			"RRYySS": "Red",
			"RRYYss": "Black",
			"RRYYSs": "Black",
			"RRYYSS": "Red",
		}),

		"hyacinths": MustSpecies("Hyacinths", map[string]string{
			"rryyWW": "White",
			"rryyWw": "White (seed)",
			"rryyww": "Blue",
			"rrYyWW": "Yellow",
			"rrYyWw": "Yellow",
			"rrYyww": "White",
			"rrYYWW": "Yellow (seed)",
			"rrYYWw": "Yellow",
			"rrYYww": "Yellow",
			"RryyWW": "Red",
			"RryyWw": "Pink",
			"Rryyww": "White",
			"RrYyWW": "Orange",
			"RrYyWw": "Yellow",
			"RrYyww": "Yellow",
			"RrYYWW": "Orange",
			"RrYYWw": "Yellow",
			"RrYYww": "Yellow",
			"RRyyWW": "Red",
			"RRyyWw": "Red (seed)",
			"RRyyww": "Red",
			"RRYyWW": "Blue",
			"RRYyWw": "Red",
			"RRYyww": "Red",
			"RRYYWW": "Purple",
			"RRYYWw": "Purple",
			"RRYYww": "Purple",
		}),

		"lilies": MustSpecies("Lilies", map[string]string{
			"rryyss": "White",
			"rryySs": "White",
			"rryySS": "White (seed)",
			"rrYyss": "Yellow",
			"rrYySs": "White",
			"rrYySS": "White",
			"rrYYss": "Yellow (seed)",
			"rrYYSs": "Yellow",
			"rrYYSS": "White",
			"Rryyss": "Red",
			"RryySs": "Pink",
			"RryySS": "White",
			"RrYyss": "Orange",
			"RrYySs": "Yellow",
			"RrYySS": "Yellow",
			"RrYYss": "Orange",
			"RrYYSs": "Yellow",
			"RrYYSS": "Yellow",
			"RRyyss": "Black",
			"RRyySs": "Red (seed)",
			"RRyySS": "Pink",
			"RRYyss": "Black",
			"RRYySs": "Red",
			"RRYySS": "Pink",
			"RRYYss": "Orange",
			"RRYYSs": "Orange",
			"RRYYSS": "White",
		}),

		"mums": MustSpecies("Mums", map[string]string{
			"rryyWW": "White",
			"rryyWw": "White (seed)",
			"rryyww": "Purple",
			"rrYyWW": "Yellow",
			"rrYyWw": "Yellow",
			"rrYyww": "White",
			"rrYYWW": "Yellow (seed)",
			"rrYYWw": "Yellow",
			"rrYYww": "Yellow",
			"RryyWW": "Pink",
			"RryyWw": "Pink",
			"Rryyww": "Pink",
			"RrYyWW": "Yellow",
			"RrYyWw": "Red",
			"RrYyww": "Pink",
			"RrYYWW": "Purple",
			"RrYYWw": "Purple",
			"RrYYww": "Purple",
			"RRyyWW": "Red (seed)",
			"RRyyWw": "Red",
			"RRyyww": "Red",
			"RRYyWW": "Purple",
			"RRYyWw": "Purple",
			"RRYyww": "Red",
			"RRYYWW": "Green",
			"RRYYWw": "Green",
			"RRYYww": "Red",
		}),

		"pansies": MustSpecies("Pansies", map[string]string{
			"rryyWW": "White",
			"rryyWw": "White (seed)",
			"rryyww": "Blue",
			"rrYyWW": "Yellow",
			"rrYyWw": "Yellow",
			"rrYyww": "Blue",
			"rrYYWW": "Yellow (seed)",
			"rrYYWw": "Yellow",
			"rrYYww": "Yellow",
			"RryyWW": "Red",
			"RryyWw": "Red",
			"Rryyww": "Blue",
			"RrYyWW": "Orange",
			"RrYyWw": "Orange",
			"RrYyww": "Orange",
			"RrYYWW": "Yellow",
			"RrYYWw": "Yellow",
			"RrYYww": "Yellow",
			"RRyyWW": "Red (seed)",
			"RRyyWw": "Red",
			"RRyyww": "Purple",
			"RRYyWW": "Red",
			"RRYyWw": "Red",
			"RRYyww": "Purple",
			"RRYYWW": "Orange",
			"RRYYWw": "Orange",
			"RRYYww": "Purple",
		}),

		"roses": MustSpecies("Roses", map[string]string{
			"rryyWWss": "White",
			"rryyWWSs": "White",
			"rryyWWSS": "White",
			"rryyWwss": "White (seed)",
			"rryyWwSs": "White",
			"rryyWwSS": "White",
			"rryywwss": "Purple",
			"rryywwSs": "Purple",
			"rryywwSS": "Purple",
			"rrYyWWss": "Yellow",
			"rrYyWWSs": "Yellow",
			"rrYyWWSS": "Yellow",
			"rrYyWwss": "White",
			"rrYyWwSs": "White",
			"rrYyWwSS": "White",
			"rrYywwss": "Purple",
			"rrYywwSs": "Purple",
			"rrYywwSS": "Purple",
			"rrYYWWss": "Yellow (seed)",
			"rrYYWWSs": "Yellow",
			"rrYYWWSS": "Yellow",
			"rrYYWwss": "Yellow",
			"rrYYWwSs": "Yellow",
			"rrYYWwSS": "Yellow",
			"rrYYwwss": "White",
			"rrYYwwSs": "White",
			"rrYYwwSS": "White",
			"RryyWWss": "Red",
			"RryyWWSs": "Pink",
			"RryyWWSS": "White",
			"RryyWwss": "Red",
			"RryyWwSs": "Pink",
			"RryyWwSS": "White",
			"Rryywwss": "Red",
			"RryywwSs": "Pink",
			"RryywwSS": "Purple",
			"RrYyWWss": "Orange",
			"RrYyWWSs": "Yellow",
			"RrYyWWSS": "Yellow",
			"RrYyWwss": "Red",
			"RrYyWwSs": "Pink",
			"RrYyWwSS": "White",
			"RrYywwss": "Red",
			"RrYywwSs": "Pink",
			"RrYywwSS": "Purple",
			"RrYYWWss": "Orange",
			"RrYYWWSs": "Yellow",
			"RrYYWWSS": "Yellow",
			"RrYYWwss": "Orange",
			"RrYYWwSs": "Yellow",
			"RrYYWwSS": "Yellow",
			"RrYYwwss": "Red",
			"RrYYwwSs": "Pink",
			"RrYYwwSS": "White",
			"RRyyWWss": "Black",
			"RRyyWWSs": "Red (seed)",
			"RRyyWWSS": "Pink",
			"RRyyWwss": "Black",
			"RRyyWwSs": "Red",
			"RRyyWwSS": "Pink",
			"RRyywwss": "Black",
			"RRyywwSs": "Red",
			"RRyywwSS": "Pink",
			"RRYyWWss": "Orange",
			"RRYyWWSs": "Orange",
			"RRYyWWSS": "Yellow",
			"RRYyWwss": "Red",
			"RRYyWwSs": "Red",
			"RRYyWwSS": "White",
			"RRYywwss": "Black",
			"RRYywwSs": "Red",
			"RRYywwSS": "Purple",
			"RRYYWWss": "Orange",
			"RRYYWWSs": "Orange",
			"RRYYWWSS": "Yellow",
			"RRYYWwss": "Orange",
			"RRYYWwSs": "Orange",
			"RRYYWwSS": "Yellow",
			"RRYYwwss": "Blue",
			"RRYYwwSs": "Red",
			"RRYYwwSS": "White",
		}),

		"tulips": MustSpecies("Tulips", map[string]string{
			"rryyss": "White",
			"rryySs": "White (seed)",
			"rryySS": "White",
			"rrYyss": "Yellow",
			"rrYySs": "Yellow",
			"rrYySS": "White",
			"rrYYss": "Yellow (seed)",
			"rrYYSs": "Yellow",
			"rrYYSS": "Yellow",
			"Rryyss": "Red",
			"RryySs": "Pink",
			"RryySS": "White",
			"RrYyss": "Orange",
			"RrYySs": "Yellow",
			"RrYySS": "Yellow",
			"RrYYss": "Orange",
			"RrYYSs": "Yellow",
			"RrYYSS": "Yellow",
			"RRyyss": "Black",
			"RRyySs": "Red (seed)",
			"RRyySS": "Red",
			"RRYyss": "Black",
			"RRYySs": "Red",
			"RRYySS": "Red",
			"RRYYss": "Purple",
			"RRYYSs": "Purple",
			"RRYYSS": "Purple",
		}),

		"windflowers": MustSpecies("Windflowers", map[string]string{
			"rrooWW": "White",
			"rrooWw": "White (seed)",
			"rrooww": "Blue",
			"rrOoWW": "Orange",
			"rrOoWw": "Orange",
			"rrOoww": "Blue",
			"rrOOWW": "Orange (seed)",
			"rrOOWw": "Orange",
			"rrOOww": "Orange",
			"RrooWW": "Red",
			"RrooWw": "Red",
			"Rrooww": "Blue",
			"RrOoWW": "Pink",
			"RrOoWw": "Pink",
			"RrOoww": "Pink",
			"RrOOWW": "Orange",
			"RrOOWw": "Orange",
			"RrOOww": "Orange",
			"RRooWW": "Red (seed)",
			"RRooWw": "Red",
			"RRooww": "Purple",
			"RROoWW": "Red",
			"RROoWw": "Red",
			"RROoww": "Purple",
			"RROOWW": "Pink",
			"RROOWw": "Pink",
			"RROOww": "Purple",
		}),
	}
)

// Species represents a specific species of flower, such as Windflower or Mum.
type Species struct {
	name       string              // a human-readable name for this species, e.g. "Windflower".
	geneCount  int                 // how many genes this flower has; assumed to be 3 or 4.
	phenotypes map[Genotype]string // phenotypes by genotype
}

func NewSpecies(name string, phenotypes map[string]string) (Species, error) {
	s := Species{name: name}
	gsInit := false
	var gs GenotypeSerde
	pts := map[Genotype]string{}
	for gStr, p := range phenotypes {
		if !gsInit {
			serde, err := NewGenotypeSerdeFromExample(gStr)
			if err != nil {
				return Species{}, fmt.Errorf("couldn't parse genotype %q: %v", gStr, err)
			}
			gs, gsInit = serde, true
			s.geneCount = gs.GeneCount()
		}

		g, err := gs.DeserializeGenotype(gStr)
		if err != nil {
			return Species{}, err
		}
		pts[g] = p
	}
	s.phenotypes = pts

	if s.geneCount == 3 && len(phenotypes) != 27 {
		return Species{}, fmt.Errorf("Got %d phenotypes, expected 27", len(phenotypes))
	}
	if s.geneCount == 4 && len(phenotypes) != 81 {
		return Species{}, fmt.Errorf("Got %d phenotypes, expected 81", len(phenotypes))
	}

	return s, nil
}

func MustSpecies(name string, phenotypes map[string]string) Species {
	s, err := NewSpecies(name, phenotypes)
	if err != nil {
		panic(err)
	}
	return s
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

type GenotypeSerde struct {
	gene0 [3]string // contents of these will be something like {"rr", "Rr", "RR"}
	gene1 [3]string
	gene2 [3]string
	gene3 [3]string // {"", "", ""} for 3-gene species
}

func (gs GenotypeSerde) GeneCount() int {
	if gs.gene3[0] == "" {
		return 3
	}
	return 4
}

func NewGenotypeSerdeFromExample(genotype string) (GenotypeSerde, error) {
	if len(genotype) != 6 && len(genotype) != 8 {
		return GenotypeSerde{}, fmt.Errorf("genotype %q has wrong length (expected 6 or 8)", genotype)
	}

	genesFrom := func(gene string) ([3]string, error) {
		lo, hi := strings.ToLower(gene[0:1]), strings.ToUpper(gene[0:1])
		genes := [3]string{lo + lo, hi + lo, hi + hi}
		if gene != genes[0] && gene != genes[1] && gene != genes[2] {
			return [3]string{}, fmt.Errorf("could not parse gene %q", gene)
		}
		return genes, nil
	}

	gene0, err := genesFrom(genotype[0:2])
	if err != nil {
		return GenotypeSerde{}, err
	}
	gene1, err := genesFrom(genotype[2:4])
	if err != nil {
		return GenotypeSerde{}, err
	}
	gene2, err := genesFrom(genotype[4:6])
	if err != nil {
		return GenotypeSerde{}, err
	}
	var gene3 [3]string
	if len(genotype) == 8 {
		gene3, err = genesFrom(genotype[6:8])
		if err != nil {
			return GenotypeSerde{}, err
		}
	}

	if gene0 == gene1 || gene0 == gene2 || gene0 == gene3 || gene1 == gene2 || gene1 == gene3 || gene2 == gene3 {
		return GenotypeSerde{}, fmt.Errorf("duplicate gene letters (%q, %q, %q, %q)", gene0[0], gene1[0], gene2[0], gene3[0])
	}

	return GenotypeSerde{gene0, gene1, gene2, gene3}, nil
}

func (gs GenotypeSerde) DeserializeGenotype(genotype string) (Genotype, error) {
	var rslt Genotype

	if gs.gene3[0] == "" && len(genotype) != 6 {
		return 0, fmt.Errorf("genotype %q has wrong length (expected 6)", genotype)
	}
	if gs.gene3[0] != "" && len(genotype) != 8 {
		return 0, fmt.Errorf("genotype %q has wrong length (expected 8)", genotype)
	}

	for _, x := range []struct {
		gene   [3]string
		offset uint
	}{
		{gs.gene0, 0},
		{gs.gene1, 2},
		{gs.gene2, 4},
		{gs.gene3, 6},
	} {
		if x.gene[0] == "" {
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

func (gs GenotypeSerde) SerializeGenotype(g Genotype) string {
	if gs.gene3[0] == "" {
		return fmt.Sprintf("%s%s%s", gs.gene0[g.gene0()], gs.gene1[g.gene1()], gs.gene2[g.gene2()])
	}
	return fmt.Sprintf("%s%s%s%s", gs.gene0[g.gene0()], gs.gene1[g.gene1()], gs.gene2[g.gene2()], gs.gene3[g.gene3()])
}

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
