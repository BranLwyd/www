# AC:NH Flower Breeding Tool

Enter genotypes in the format `RrYyWw` (or `RrYyWwSs` for roses), then click "Breed 'Em!" to see what you might get. You can also input genetic distributions: for example, `{1:RRYYWW, 2:rryyww}` represents a flower which has a 1/3 chance of being `RRYYWW` and a 2/3 chance of being `rryyww`.

For more information on the genotype notation used by this tool, please refer to [Paleh's Flower Genetics Guide](https://docs.google.com/document/d/1ARIQCUc5YVEd01D7jtJT9EEJF45m07NXhAm4fOpNvCs).

<form action="" method="get">
<p>
<label for="species">Species: </label>
<select name="species">
<option value="cosmos"{{if eq .Species "cosmos"}} selected{{end}}>Cosmos</option>
<option value="hyacinths"{{if eq .Species "hyacinths"}} selected{{end}}>Hyacinths</option>
<option value="lilies"{{if eq .Species "lilies"}} selected{{end}}>Lilies</option>
<option value="mums"{{if eq .Species "mums"}} selected{{end}}>Mums</option>
<option value="pansies"{{if eq .Species "pansies"}} selected{{end}}>Pansies</option>
<option value="roses"{{if eq .Species "roses"}} selected{{end}}>Roses</option>
<option value="tulips"{{if eq .Species "tulips"}} selected{{end}}>Tulips</option>
<option value="windflowers"{{if eq .Species "windflowers"}} selected{{end}}>Windflowers</option>
</select>
</p>
<p><label for="firstGenotype">First genotype: </label><input type="text" name="firstGenotype" id="firstGenotype" value="{{.FirstGenotype}}"></p>
<p><label for="secondGenotype">Second genotype: </label><input type="text" name="secondGenotype" id="secondGenotype" value="{{.SecondGenotype}}"></p>
<p><input type="submit" value="Breed 'em!"></p>
</form>

{{if .Error}}<div><b>{{.Error}}</b></div>{{end}}

{{if .Results}}
<div>
<table class="flower-table">
<tr>
<th>Odds</th>
<th>Genotype</th>
<th>Phenotype</th>
</tr>
{{$spc := .Species}}{{range .Results}}
<tr>
<td>{{.Odds}}</td>
<td><code>{{.Genotype}}</code></td>
<td><img src="{{pic $spc .Phenotype}}"> <span>{{.Phenotype}}</span></td>
</tr>
{{end}}
</table>
**Result:** `{{.ResultString}}`

(This can be used as an input genotype.)
</div>
{{end}}
</body>
</html>
