// Transform processes incoming records returning a aggregrate of individual user sales.
function Transform(recordsJSON){
	var processed = []
	
	var records = JSON.parse(recordsJSON)
	for(var index in records){
		var record = records[index]
		
		var totalSales = 0
		for(var saleIndex in record.sales){
			totalSales += record.sales[saleIndex]
		}
		
		processed.push({
			"user": record.name,
			"sales": totalSales,
		})
	}
	
	return JSON.stringify(processed)
}